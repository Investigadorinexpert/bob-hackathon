# Integración con Backend BOB

## Qué se hizo

Se integró el bot de WhatsApp con el backend BOB (Go + Gemini AI 2.5 Flash) para procesar mensajes inteligentemente con lead scoring automático.

## Cambios en `cmd/whserver/main.go`

### 1. Nueva función `callBOBBackend` (línea ~838)

```go
func callBOBBackend(fromPhone string, message string, logger jlog) string {
    sessionId := "wa-" + fromPhone

    payload := map[string]string{
        "sessionId": sessionId,
        "message":   message,
        "channel":   "whatsapp",
    }
    jsonData, _ := json.Marshal(payload)

    resp, err := http.Post(
        "http://localhost:3000/api/chat/message",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        logger.Warn("bob_backend_error", "err", err)
        return "Lo siento, hubo un error procesando tu mensaje."
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        logger.Warn("bob_backend_decode_error", "err", err)
        return "Error procesando la respuesta."
    }

    if reply, ok := result["reply"].(string); ok {
        // Log del lead score si viene
        if score, ok2 := result["leadScore"].(float64); ok2 {
            if category, ok3 := result["category"].(string); ok3 {
                logger.Info("bob_backend_reply",
                    "from", fromPhone,
                    "score", int(score),
                    "category", category,
                    "reply_len", len([]rune(reply)),
                )
            }
        }
        return reply
    }

    return "No se pudo obtener respuesta del sistema."
}
```

### 2. Modificación del callback `onFlush` (línea ~1054)

Se reemplazó el engine de reglas (`router.eng.Eval`) por la llamada al backend BOB:

```go
// ANTES:
res, handled := router.eng.Eval(context.Background(), env)

// AHORA:
reply := callBOBBackend(env.SenderJID, env.Text, logger)
```

## Cómo funciona

1. Usuario envía mensaje por WhatsApp
2. El bot recibe el mensaje y lo procesa
3. En el callback `onFlush`, llama a `callBOBBackend()`
4. El backend procesa con Gemini AI 2.5 Flash
5. Retorna respuesta + lead score
6. El bot envía la respuesta por WhatsApp

## Flujo de datos

```
WhatsApp User
    ↓
Bot Engine (whserver)
    ↓
callBOBBackend()
    ↓
POST http://localhost:3000/api/chat/message
{
  "sessionId": "wa-51999999999",
  "message": "Hola, busco un auto",
  "channel": "whatsapp"
}
    ↓
Backend BOB (Gemini AI)
    ↓
Response:
{
  "reply": "¡Hola! Claro, te ayudo...",
  "leadScore": 45,
  "category": "warm"
}
    ↓
WhatsApp User (recibe respuesta)
```

## Logs

El sistema ahora genera logs con información del lead scoring:

```
[INFO] bob_backend_reply from=51999999999 score=75 category=hot reply_len=156
```

## Testing

### 1. Asegúrate que el backend esté corriendo:
```bash
curl http://localhost:3000/health
```

### 2. Prueba la integración con curl:
```bash
curl -X POST http://localhost:3000/api/chat/message \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "wa-51999999999",
    "message": "Hola, busco un auto Toyota",
    "channel": "whatsapp"
  }'
```

### 3. Inicia el bot:
```bash
cd bot
go run cmd/whserver/main.go
```

### 4. Envía un mensaje por WhatsApp y verifica:
- Que el bot responde
- Que los logs muestran `bob_backend_reply`
- Que aparece el score y category

## Ventajas de esta integración

- **Inteligencia Gemini AI**: Respuestas contextuales y naturales
- **Lead Scoring Automático**: Califica leads 0-100 en tiempo real
- **Categorización**: hot/warm/cold automática
- **Centralizado**: Toda la lógica de IA en un solo lugar
- **Escalable**: Fácil agregar más canales (web, telegram, etc.)
- **Logs detallados**: Score y categoría en cada interacción

## Notas importantes

- El backend debe estar corriendo en `localhost:3000`
- Session IDs usan formato `wa-{numero_telefono}`
- Los mensajes se procesan con ventana de agregación (3 segundos por defecto)
- El lead scoring se actualiza en cada mensaje
- Los datos se persisten en `backend/data/sessions.json` y `backend/data/leads.json`

## Si necesitas modificar

### Cambiar URL del backend:
En `callBOBBackend()`, línea ~848:
```go
resp, err := http.Post(
    "http://TU_URL_AQUI/api/chat/message",  // <-- cambiar aquí
    "application/json",
    bytes.NewBuffer(jsonData),
)
```

### Agregar timeout personalizado:
Antes de la función `callBOBBackend()`, modificar el httpc:
```go
var httpc = &http.Client{Timeout: 10 * time.Second}  // cambiar timeout
```

### Deshabilitar integración (volver al engine antiguo):
En el callback `onFlush`, comentar la sección BOB y descomentar la sección del engine:
```go
// Comentar esto:
// reply := callBOBBackend(env.SenderJID, env.Text, logger)

// Descomentar esto:
// res, handled := router.eng.Eval(context.Background(), env)
```

---

**Integración completa y funcionando!**

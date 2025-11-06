# 📢 PARA EL EQUIPO - BOB Hackathon

## 🎯 Resumen Ejecutivo

**YA ESTÁ TODO LISTO Y FUNCIONANDO** ✅

- ✅ Backend con Gemini AI 2.5 Flash
- ✅ Frontend con chat interactivo
- ✅ Lead scoring automático
- ✅ Dashboard de leads en tiempo real
- ✅ API REST completa

---

## 🚀 URLs Activas

```
Frontend Web:  http://localhost:5173
Backend API:   http://localhost:3000
```

---

## 👥 División de Trabajo

### Kevin (Backend + Frontend) ✅ COMPLETADO
- ✅ Backend Node.js + Express
- ✅ Integración Gemini AI
- ✅ Sistema de scoring
- ✅ API REST completa
- ✅ Frontend React
- ✅ Chat widget
- ✅ Dashboard de leads

### Compañero (WhatsApp Integration) ⏳ PENDIENTE
**Solo necesitas agregar 20 líneas de código en tu bot de Go.**

---

## 🔌 INTEGRACIÓN WHATSAPP (Copy & Paste)

### En tu archivo `bot/cmd/whserver/main.go`:

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

// Agregar esta función
func callKevinBackend(fromPhone string, message string) string {
    // Session ID único por número de WhatsApp
    sessionId := fmt.Sprintf("wa-%s", fromPhone)

    // Preparar payload
    payload := map[string]string{
        "sessionId": sessionId,
        "message":   message,
        "channel":   "whatsapp",
    }

    jsonData, _ := json.Marshal(payload)

    // Llamar al backend de Kevin
    resp, err := http.Post(
        "http://localhost:3000/api/chat/message",
        "application/json",
        bytes.NewBuffer(jsonData),
    )

    if err != nil {
        log.Printf("Error calling backend: %v", err)
        return "Lo siento, hubo un error. Intenta de nuevo."
    }
    defer resp.Body.Close()

    // Parsear respuesta
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        log.Printf("Error decoding response: %v", err)
        return "Error procesando respuesta."
    }

    // Obtener reply
    if reply, ok := result["reply"].(string); ok {
        return reply
    }

    return "No se pudo obtener respuesta."
}

// Luego en tu handler de mensajes de WhatsApp:
func handleIncomingWhatsAppMessage(from string, body string) {
    // Llamar al backend de Kevin
    reply := callKevinBackend(from, body)

    // Enviar respuesta por WhatsApp
    sendWhatsAppMessage(from, reply)
}
```

### ¡ESO ES TODO!

---

## 🧪 Cómo Probar la Integración

### 1. Asegúrate que el backend de Kevin esté corriendo:
```bash
curl http://localhost:3000/health
```

Si responde `{"status":"ok",...}` → ✅ Backend listo

### 2. Prueba la integración con curl:
```bash
curl -X POST http://localhost:3000/api/chat/message \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "wa-51999999999",
    "message": "Hola, busco un auto",
    "channel": "whatsapp"
  }'
```

Deberías recibir:
```json
{
  "success": true,
  "sessionId": "wa-51999999999",
  "reply": "¡Hola! Claro, te ayudo...",
  "leadScore": 45,
  "category": "warm"
}
```

### 3. Integra en tu bot:
- Copia el código de arriba
- Llama `callKevinBackend(from, message)` cuando recibas un mensaje
- Envía el `reply` de vuelta por WhatsApp

---

## 📊 Lo que Obtienes del Backend

Cuando llamas al endpoint `/api/chat/message`, recibes:

```json
{
  "success": true,
  "sessionId": "wa-51999999999",
  "reply": "Respuesta inteligente del bot",
  "leadScore": 75,
  "category": "hot",
  "timestamp": "2025-11-06T02:00:00Z"
}
```

- `reply`: La respuesta que debes enviar por WhatsApp
- `leadScore`: Score del lead (0-100)
- `category`: Categoría (hot/warm/cold)

---

## 🎯 Para la Demo del Hackathon

### Flujo de Presentación:

**1. Mostrar Frontend Web (2 min)**
- Abrir http://localhost:5173
- Hacer conversación de ejemplo
- Mostrar cómo sube el score
- Cambiar a dashboard de leads

**2. Mostrar Integración WhatsApp (2 min)**
- Escanear QR del bot
- Enviar mensaje por WhatsApp
- Mostrar respuesta del bot
- Abrir dashboard web y mostrar el lead de WhatsApp

**3. Explicar Arquitectura (1 min)**
```
WhatsApp → Bot Go → Backend Kevin → Gemini AI
   ↓                                    ↓
Respuesta ← Bot Go ← Backend Kevin ← Respuesta
```

**4. Destacar Features (1 min)**
- ✅ Modular (REST API)
- ✅ Multi-canal (Web + WhatsApp)
- ✅ IA Moderna (Gemini 2.5 Flash)
- ✅ Lead scoring automático
- ✅ Dashboard en tiempo real

---

## ❓ Preguntas Frecuentes

### ¿Qué hace el backend de Kevin?
- Recibe mensajes (web o WhatsApp)
- Procesa con Gemini AI
- Calcula score del lead
- Devuelve respuesta inteligente

### ¿Qué debo hacer yo (WhatsApp)?
- Solo llamar al endpoint `/api/chat/message`
- Enviar la respuesta por WhatsApp
- **20 líneas de código máximo**

### ¿Necesito instalar algo más?
**NO**. El backend de Kevin ya tiene todo:
- Gemini AI configurado
- FAQs cargadas
- API BOB conectada
- Sistema de scoring listo

### ¿Y si quiero ver los leads?
- Frontend: http://localhost:5173 → Tab "Leads"
- API: `curl http://localhost:3000/api/leads`

---

## 🆘 Si Algo No Funciona

### Backend no responde:
```bash
# Verificar que esté corriendo
curl http://localhost:3000/health

# Si no responde, iniciar:
cd backend
npm start
```

### ¿Cómo sé si mi integración funciona?
```bash
# Probar con curl primero
curl -X POST http://localhost:3000/api/chat/message \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "wa-test",
    "message": "test",
    "channel": "whatsapp"
  }'
```

Si esto funciona → Tu integración funcionará.

---

## 📱 Contacto

**Kevin (Backend/Frontend)**
- Ya está todo listo de mi lado
- Si necesitas algo, avísame

**Compañero (WhatsApp)**
- Solo agrega el código de integración
- Prueba primero con curl
- Cualquier duda, pregunta

---

## ✅ Checklist Pre-Demo

- [ ] Backend corriendo (http://localhost:3000/health)
- [ ] Frontend corriendo (http://localhost:5173)
- [ ] Bot WhatsApp conectado
- [ ] Integración probada con curl
- [ ] Mensaje de prueba por WhatsApp funcionando
- [ ] Dashboard mostrando leads

---

**🎯 TODO LISTO. Solo falta integrar WhatsApp.**

**Tiempo estimado de integración: 30 minutos** ⏱️

---

**Hackathon BOB 2025** 🚀

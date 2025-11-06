# 🚀 QUICKSTART - BOB Chatbot

## ⚡ Inicio Rápido (2 comandos)

### 1️⃣ Iniciar Backend (Go)
```bash
cd backend
go run cmd/server/main.go
```
✅ Backend corriendo en: **http://localhost:3000**

### 2️⃣ Iniciar Frontend (Nueva terminal)
```bash
cd frontend
npm run dev
```
✅ Frontend corriendo en: **http://localhost:5173**

---

## 🌐 URLs Importantes

| Servicio | URL | Descripción |
|----------|-----|-------------|
| **Frontend** | http://localhost:5173 | Interfaz web del chatbot |
| **Backend API** | http://localhost:3000 | API REST |
| **Health Check** | http://localhost:3000/health | Estado del backend |
| **API Docs** | http://localhost:3000 | Documentación endpoints |

---

## 🧪 Probar el Sistema

### Opción 1: Interfaz Web (Recomendado)
1. Abre: **http://localhost:5173**
2. Escribe un mensaje en el chat
3. El bot responde con Gemini AI
4. Ve tu score en tiempo real
5. Cambia a pestaña "📊 Leads" para ver dashboard

### Opción 2: Con curl (Testing API)
```bash
# Enviar mensaje
curl -X POST http://localhost:3000/api/chat/message \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hola, busco un auto Toyota",
    "channel": "web"
  }'

# Ver leads
curl http://localhost:3000/api/leads

# Ver vehículos disponibles
curl 'http://localhost:3000/api/vehicles?limit=5'

# Ver FAQs
curl http://localhost:3000/api/faqs
```

---

## 📋 Endpoints API

### Chat
```bash
POST /api/chat/message          # Enviar mensaje
POST /api/chat/score            # Calcular score
GET  /api/chat/history/:id      # Ver historial
DELETE /api/chat/session/:id    # Limpiar sesión
```

### Leads
```bash
GET /api/leads                  # Listar leads
GET /api/leads/:sessionId       # Ver lead específico
GET /api/leads/stats            # Estadísticas
```

### Recursos
```bash
GET /api/faqs                   # Preguntas frecuentes
GET /api/vehicles               # Vehículos en subasta
GET /api/vehicles/:id           # Vehículo específico
```

---

## 🔧 Troubleshooting

### Backend no inicia
```bash
# Verificar que el puerto 3000 esté libre
lsof -ti:3000 | xargs kill -9

# Reinstalar dependencias
cd backend
rm -rf node_modules package-lock.json
npm install
npm start
```

### Frontend no inicia
```bash
# Verificar que el puerto 5173 esté libre
lsof -ti:5173 | xargs kill -9

# Reinstalar dependencias
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

### Gemini API no responde
1. Verifica que la key en `backend/.env` sea correcta
2. Chequea que tengas internet
3. Revisa logs del backend: `cd backend && tail -f logs/*.log`

---

## 💻 Desarrollo

### Backend
```bash
cd backend

# Iniciar en modo desarrollo
npm start

# Ver logs
tail -f logs/server.log

# Limpiar datos temporales
rm src/data/sessions.json src/data/leads.json
```

### Frontend
```bash
cd frontend

# Iniciar dev server
npm run dev

# Build para producción
npm run build

# Preview build
npm run preview
```

---

## 🔌 Integración WhatsApp (Para tu compañero)

Agregar en el bot de Go (`bot/cmd/whserver/main.go`):

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func handleWhatsAppMessage(from string, message string) {
    // Crear sessionId único
    sessionId := "wa-" + from

    // Preparar request
    payload := map[string]string{
        "sessionId": sessionId,
        "message":   message,
        "channel":   "whatsapp",
    }
    jsonData, _ := json.Marshal(payload)

    // Llamar a tu backend
    resp, err := http.Post(
        "http://localhost:3000/api/chat/message",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        log.Println("Error:", err)
        return
    }
    defer resp.Body.Close()

    // Parsear respuesta
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    // Obtener reply
    if reply, ok := result["reply"].(string); ok {
        // Enviar respuesta por WhatsApp
        sendWhatsAppMessage(from, reply)
    }
}
```

**Eso es TODO lo que necesita agregar.** El resto ya está funcionando.

---

## 📊 Features Disponibles

### ✅ Backend
- [x] Gemini 2.5 Flash AI
- [x] Lead scoring automático (0-100)
- [x] Session management
- [x] API BOB integrada (datos reales)
- [x] 62 FAQs cargadas
- [x] Soporte multi-canal (web + WhatsApp)

### ✅ Frontend
- [x] Chat widget interactivo
- [x] Lead scoring en tiempo real
- [x] Dashboard de leads
- [x] Filtros por categoría
- [x] Auto-refresh
- [x] Responsive design

---

## 🎯 Para el Hackathon

### Demo Flow:
1. **Mostrar Frontend**: http://localhost:5173
2. **Conversar con el bot**: Hacer 3-4 preguntas
3. **Ver scoring**: Mostrar cómo sube el score
4. **Dashboard**: Cambiar a pestaña Leads
5. **WhatsApp**: Mostrar integración (si está lista)
6. **API**: Mostrar endpoints en Postman

### Puntos Clave del Pitch:
- ✅ **Modular**: Backend REST API + Frontend + WhatsApp
- ✅ **IA Moderna**: Gemini 2.5 Flash (último modelo)
- ✅ **Scoring Automático**: Califica leads en tiempo real
- ✅ **Sin n8n**: Más simple, más rápido
- ✅ **Multi-canal**: Web + WhatsApp ready
- ✅ **Datos Reales**: Conectado a API de BOB

---

## 📝 Notas

- Backend guarda datos en `src/data/` (JSON temporal)
- Frontend usa proxy de Vite para llamar al backend
- Gemini key está en `backend/.env`
- FAQs y vehículos se cargan al iniciar el backend

---

**¿Problemas?** Revisa el README.md principal para más detalles.

**Hackathon BOB 2025** 🚀

# 🤖 BOB Chatbot - Hackathon 2025

Backend en **Go** con Gemini AI, lead scoring automático y arquitectura modular para BOB Subastas.

## 🎯 Features

- ✅ **Gemini AI Integration** - Respuestas inteligentes y contextuales
- ✅ **Lead Scoring Automático** - Calificación de leads en tiempo real
- ✅ **Session Management** - Conversaciones persistentes por sessionId
- ✅ **Multi-canal** - Soporte para Web y WhatsApp
- ✅ **API REST** - Endpoints documentados y fáciles de integrar
- ✅ **BOB API Integration** - Datos en tiempo real de subastas
- ✅ **FAQs Inteligentes** - Base de conocimiento con 62+ preguntas frecuentes

## 🏗️ Arquitectura

```
┌─────────────┐         ┌──────────────┐
│  Frontend   │         │   WhatsApp   │
│    Web      │         │     Bot      │
└──────┬──────┘         └───────┬──────┘
       │                        │
       └────────┬───────────────┘
                │
         ┌──────▼──────┐
         │   Backend   │
         │  (Node.js)  │
         └──────┬──────┘
                │
     ┌──────────┼──────────┐
     │          │          │
┌────▼───┐ ┌───▼────┐ ┌──▼─────┐
│Gemini  │ │BOB API │ │ FAQs   │
│  AI    │ │        │ │  DB    │
└────────┘ └────────┘ └────────┘
```

## 🚀 Quick Start

> **📌 IMPORTANTE**: Lee el [QUICKSTART.md](QUICKSTART.md) para instrucciones detalladas paso a paso.

### ⚡ Inicio Rápido (2 comandos)

**Terminal 1 - Backend:**
```bash
cd backend
npm install    # Solo primera vez
npm start
```
✅ Backend: http://localhost:3000

**Terminal 2 - Frontend:**
```bash
cd frontend
npm install    # Solo primera vez
npm run dev
```
✅ Frontend: http://localhost:5173

### 🎯 Probar Ahora
Abre tu navegador en: **http://localhost:5173**

## 📚 API Endpoints

### Chat

#### POST /api/chat/message
Envía un mensaje y recibe respuesta con IA

```json
// Request
{
  "sessionId": "web-uuid-123",  // opcional, se genera automáticamente
  "message": "Hola, busco un auto",
  "channel": "web" | "whatsapp"
}

// Response
{
  "success": true,
  "sessionId": "web-uuid-123",
  "reply": "¡Hola! Claro, te ayudo a encontrar un auto...",
  "leadScore": 45,
  "category": "warm",
  "timestamp": "2025-11-05T20:30:00Z"
}
```

#### POST /api/chat/score
Calcula el score de un lead

```json
// Request
{
  "sessionId": "web-uuid-123"
}

// Response
{
  "success": true,
  "score": 85,
  "category": "hot",
  "reasons": ["Necesidad urgente", "Presupuesto definido"],
  "urgency": "high",
  "budget": "defined",
  "businessType": "company"
}
```

#### GET /api/chat/history/:sessionId
Obtiene historial de conversación

### Leads

#### GET /api/leads
Lista todos los leads

```
Query params:
- category: hot | warm | cold
- channel: web | whatsapp
```

#### GET /api/leads/:sessionId
Obtiene un lead específico

#### GET /api/leads/stats
Estadísticas de leads

### Resources

#### GET /api/faqs
Obtiene FAQs

```
Query params:
- search: texto a buscar
- categoria: filtrar por categoría
- empresa: filtrar por empresa
```

#### GET /api/vehicles
Obtiene vehículos en subasta

```
Query params:
- marca: filtrar por marca
- modelo: filtrar por modelo
- precio_min: precio mínimo
- precio_max: precio máximo
- tipo_subasta: En vivo | Dinámica | Sobre cerrado
- limit: número de resultados (default: 10)
```

## 🔥 Testing Rápido

### Con curl:

```bash
# Health check
curl http://localhost:3000/health

# Enviar mensaje
curl -X POST http://localhost:3000/api/chat/message \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hola, busco un auto Toyota",
    "channel": "web"
  }'

# Ver leads
curl http://localhost:3000/api/leads

# Ver FAQs
curl http://localhost:3000/api/faqs

# Ver vehículos
curl http://localhost:3000/api/vehicles?limit=5
```

### Con Postman:

Importa esta colección:

```json
{
  "info": {
    "name": "BOB Chatbot API"
  },
  "item": [
    {
      "name": "Send Message",
      "request": {
        "method": "POST",
        "url": "http://localhost:3000/api/chat/message",
        "body": {
          "mode": "raw",
          "raw": "{\"message\": \"Hola, busco un auto\", \"channel\": \"web\"}"
        }
      }
    }
  ]
}
```

## 🔌 Integración con WhatsApp Bot

Tu compañero solo necesita hacer esto en su bot de Go:

```go
func handleWhatsAppMessage(from string, message string) {
    // Crear sessionId único por número
    sessionId := "wa-" + from

    // Llamar a tu backend
    body := map[string]string{
        "sessionId": sessionId,
        "message": message,
        "channel": "whatsapp",
    }

    resp, _ := http.Post(
        "http://localhost:3000/api/chat/message",
        "application/json",
        toJSON(body)
    )

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    // Enviar respuesta por WhatsApp
    reply := result["reply"].(string)
    sendWhatsAppMessage(from, reply)
}
```

## 📊 Sistema de Scoring

El sistema califica leads automáticamente:

- **🔥 Hot (80-100)**: Lead caliente, enviar a comercial inmediatamente
  - Necesidad urgente y específica
  - Presupuesto definido
  - Empresa/negocio

- **😐 Warm (50-79)**: Lead tibio, hacer seguimiento
  - Interés genuino
  - Explorando opciones
  - Necesidad real pero no urgente

- **❄️ Cold (0-49)**: Lead frío, base de datos
  - Solo curiosidad
  - Preguntas muy generales
  - Sin necesidad clara

## 🛠️ Tecnologías

- **Go 1.21+** - Backend de alto rendimiento
- **Gin** - Framework web (similar a Express)
- **Gemini AI 2.0 Flash** (Google)
- **UUID** para session management
- **CORS** habilitado
- **JSON** como base de datos temporal

## 📁 Estructura del Proyecto

```
backend/                  # Backend Go
├── cmd/
│   └── server/
│       └── main.go      # Servidor principal
├── internal/
│   ├── config/          # Configuración
│   ├── controllers/     # Chat & Lead controllers
│   ├── services/        # Gemini, Session, BOB API, FAQs
│   └── models/          # Estructuras de datos
├── data/                # FAQs, vehículos, sesiones, leads
├── .env                 # Configuración
├── go.mod               # Dependencias
└── README.md            # Documentación
```

## 🎯 Estado del Proyecto

1. ✅ Backend Go funcionando
2. ✅ Frontend React completo
3. ✅ Dashboard de leads
4. ✅ Gemini AI 2.0 Flash integrado
5. ⏳ Integración WhatsApp (pendiente - tu compañero)
6. ⏳ Deploy a producción

## 👥 Equipo

- **Kevin Navarro** - Backend + Frontend Web
- **[Compañero]** - Integración WhatsApp

## 📝 Notas

- **Backend en Go**: Más rápido y eficiente que Node.js
- **NO usa n8n**: Arquitectura más simple y directa
- **Session IDs** únicos por canal: `web-uuid` o `wa-{numero}`
- **FAQs y vehículos** se cargan en memoria al iniciar
- **Cache de API BOB**: 5 minutos
- **Gemini 2.0 Flash**: Modelo más reciente y rápido

## 🚀 Ventajas de Go

- **5-10x más rápido** que Node.js
- **Concurrencia nativa** con goroutines
- **Tipado fuerte** (menos errores en runtime)
- **Single binary** para deployment fácil
- **Menor uso de memoria**

---

**Hackathon BOB 2025** 🚀

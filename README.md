# BOB Chatbot - Hackathon 2025

Sistema de chatbot inteligente con lead scoring automático para BOB Subastas. Backend en Go con Gemini AI 2.5 Flash.

## Inicio Rápido

### Backend
```bash
cd backend
go run cmd/server/main.go
```
Backend en: http://localhost:3000

### Frontend
```bash
cd frontend
npm install  # solo primera vez
npm run dev
```
Frontend en: http://localhost:5173

### Verificar
```bash
curl http://localhost:3000/health
```

## Features

- **Gemini AI 2.5 Flash** - Respuestas inteligentes contextuales
- **Lead Scoring Automático** - Calificación 0-100 en tiempo real
- **Multi-canal** - Web + WhatsApp
- **Session Management** - Conversaciones persistentes
- **BOB API Integration** - Datos reales de subastas
- **62+ FAQs** - Base de conocimiento completa

## Arquitectura

```
Frontend Web ────┐
                 ├──> Backend Go ──┬──> Gemini AI 2.5
WhatsApp Bot ────┘                 ├──> BOB API
                                   └──> FAQs DB
```

## API Endpoints

### Chat
```bash
# Enviar mensaje
POST /api/chat/message
{
  "message": "Busco un auto Toyota",
  "channel": "web",
  "sessionId": "opcional"
}

# Calcular score
POST /api/chat/score
{ "sessionId": "web-123" }

# Ver historial
GET /api/chat/history/:sessionId
```

### Leads
```bash
# Listar leads
GET /api/leads?category=hot&channel=web

# Lead específico
GET /api/leads/:sessionId

# Estadísticas
GET /api/leads/stats
```

### Recursos
```bash
# FAQs
GET /api/faqs?search=subasta

# Vehículos
GET /api/vehicles?marca=Toyota&limit=10
```

## Testing

```bash
# Test completo
curl -X POST http://localhost:3000/api/chat/message \
  -H "Content-Type: application/json" \
  -d '{"message": "Hola", "channel": "web"}'

# Ver leads
curl http://localhost:3000/api/leads

# Ver vehículos
curl http://localhost:3000/api/vehicles?limit=5
```

## Lead Scoring

- **Hot (80-100)**: Necesidad urgente, presupuesto definido, empresa
- **Warm (50-79)**: Interés genuino, explorando opciones
- **Cold (0-49)**: Curiosidad, preguntas generales

## Integración WhatsApp

En tu bot de Go, agrega esto en `bot/cmd/whserver/main.go`:

```go
func handleWhatsAppMessage(from string, message string) {
    sessionId := "wa-" + from

    payload := map[string]string{
        "sessionId": sessionId,
        "message":   message,
        "channel":   "whatsapp",
    }
    jsonData, _ := json.Marshal(payload)

    resp, _ := http.Post(
        "http://localhost:3000/api/chat/message",
        "application/json",
        bytes.NewBuffer(jsonData),
    )

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    reply := result["reply"].(string)
    sendWhatsAppMessage(from, reply)
}
```

## Configuración

Copia `.env.example` a `.env`:

```bash
cd backend
cp .env.example .env
```

Edita `backend/.env`:
```env
GEMINI_API_KEY=tu_api_key_aqui
GEMINI_MODEL=gemini-2.5-flash
PORT=3000
BOB_API_BASE_URL=https://apiv3.somosbob.com/v3
```

## Estructura del Proyecto

```
backend/
├── cmd/server/main.go      # Servidor principal
├── internal/
│   ├── config/             # Configuración
│   ├── controllers/        # Chat & Leads
│   ├── services/           # Gemini, Session, BOB API
│   └── models/             # Estructuras de datos
├── data/                   # FAQs, vehículos, sesiones
├── .env                    # Configuración (no en git)
└── go.mod

frontend/
├── src/
│   ├── components/         # ChatWidget, Dashboard
│   └── App.jsx
└── vite.config.js

bot/                        # Bot WhatsApp (separado)
```

## Troubleshooting

### Backend no inicia
```bash
lsof -ti:3000 | xargs kill -9
cd backend
go run cmd/server/main.go
```

### Frontend no inicia
```bash
lsof -ti:5173 | xargs kill -9
cd frontend
npm run dev
```

### Gemini no responde
- Verificar API key en `backend/.env`
- Verificar modelo: debe ser `gemini-2.5-flash`
- Verificar conexión a internet

## Stack Tecnológico

- **Go 1.21+** - Backend (5-10x más rápido que Node.js)
- **Gin** - Framework web
- **Gemini AI 2.5 Flash** - IA conversacional
- **React + Vite** - Frontend
- **UUID** - Session management
- **JSON** - Almacenamiento temporal

## Demo Hackathon

### Flow de presentación:
1. Mostrar frontend web (http://localhost:5173)
2. Chatear con el bot (3-4 preguntas)
3. Ver scoring en tiempo real
4. Dashboard de leads
5. Integración WhatsApp (si está lista)
6. Mostrar API endpoints

### Puntos clave:
- Modular (REST API)
- Multi-canal (Web + WhatsApp)
- IA moderna (Gemini 2.5 Flash)
- Scoring automático
- Sin n8n (más simple)
- Datos reales (BOB API)

## Estado del Proyecto

- [x] Backend Go funcionando
- [x] Frontend React completo
- [x] Dashboard de leads
- [x] Gemini AI 2.5 Flash integrado
- [x] API BOB conectada
- [x] 62 FAQs cargadas
- [ ] Integración WhatsApp
- [ ] Deploy a producción

## Equipo

- **Kevin Navarro** - Backend Go + Frontend React
- **Compañero** - Bot WhatsApp

## Notas

- Backend guarda datos en `backend/data/*.json`
- Session IDs: `web-uuid` o `wa-{numero}`
- Cache BOB API: 5 minutos
- FAQs y vehículos se cargan al iniciar

---

**Hackathon BOB 2025**

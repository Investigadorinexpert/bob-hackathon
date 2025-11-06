# Backend Go - BOB Chatbot

## Stack Tecnológico

- **Go 1.21+**
- **Gin** - Framework web
- **Gemini AI** (Google) - IA conversacional
- **UUID** - Generación de IDs
- **CORS** - Habilitado para frontend

## Estructura del Proyecto

```
backend-go/
├── cmd/
│   └── server/
│       └── main.go           # Servidor principal
├── internal/
│   ├── config/
│   │   └── config.go         # Configuración
│   ├── controllers/
│   │   ├── chat_controller.go
│   │   └── lead_controller.go
│   ├── services/
│   │   ├── gemini_service.go
│   │   ├── session_service.go
│   │   ├── bob_api_service.go
│   │   └── faq_service.go
│   └── models/
│       └── models.go         # Estructuras de datos
├── data/
│   ├── faqs.csv             # FAQs (62 preguntas)
│   ├── vehicles.csv         # Vehículos
│   ├── sessions.json        # Sesiones (auto-generado)
│   └── leads.json          # Leads (auto-generado)
├── .env                     # Configuración (Gemini API Key)
├── go.mod
└── go.sum
```

## Inicio Rápido

### 1. Descargar dependencias

```bash
cd backend-go
go mod tidy
```

### 2. Configurar .env

El archivo `.env` ya está configurado con:

```env
GEMINI_API_KEY=AIzaSyAwPmY89hvvTek-o4CT5Svn4mjeoV1B8pg
GEMINI_MODEL=gemini-2.0-flash-exp
PORT=3000
BOB_API_BASE_URL=https://apiv3.somosbob.com/v3
```

### 3. Ejecutar el servidor

```bash
go run cmd/server/main.go
```

El servidor iniciará en `http://localhost:3000`

## Endpoints API

### Health Check
```bash
GET /health
```

### Chat
```bash
POST /api/chat/message
POST /api/chat/score
GET  /api/chat/history/:sessionId
DELETE /api/chat/session/:sessionId
```

### Leads
```bash
GET /api/leads
GET /api/leads/:sessionId
GET /api/leads/stats
```

### Recursos
```bash
GET /api/faqs
GET /api/vehicles
GET /api/vehicles/:id
```

## Testing

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
```

## Características

**Gemini 2.0 Flash** - Modelo más reciente
**Lead Scoring** - Automático con IA
**Multi-canal** - Web + WhatsApp
**Session Management** - Persistente
**BOB API** - Integración con datos reales
**FAQs** - Base de conocimiento
**CORS** - Habilitado para frontend

## Diferencias con la versión Node.js

1. **Rendimiento**: Go es compilado, más rápido
2. **Concurrencia**: Goroutines nativas
3. **Tipado fuerte**: Menos errores en runtime
4. **Single binary**: Fácil deployment
5. **Menor uso de memoria**: Más eficiente

## Compatible con Frontend

El frontend React existente (`http://localhost:5173`) funciona perfectamente con este backend. Las APIs son 100% compatibles.

## Deployment

### Compilar para producción

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bob-server cmd/server/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o bob-server.exe cmd/server/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o bob-server cmd/server/main.go
```

Luego solo ejecutar el binario:
```bash
./bob-server
```

## Notas

- Todas las dependencias se descargan con `go mod tidy`
- El servidor crea automáticamente `data/sessions.json` y `data/leads.json`
- Los FAQs y vehículos se cargan en memoria al inicio
- Cache de BOB API: 5 minutos

---

**Hackathon BOB 2025** 🚀

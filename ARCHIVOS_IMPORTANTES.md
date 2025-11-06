# 📁 ARCHIVOS IMPORTANTES - BOB Hackathon

## 📖 Documentación (LEE PRIMERO)

| Archivo | Descripción | Para Quién |
|---------|-------------|------------|
| **README.md** | Documentación completa del proyecto | Todos |
| **QUICKSTART.md** | Guía de inicio rápido (URLs, comandos) | Todos |
| **PARA_EL_EQUIPO.md** | Instrucciones para integración WhatsApp | Compañero |
| **ARCHIVOS_IMPORTANTES.md** | Este archivo (índice de todo) | Todos |

---

## 🔧 Configuración

| Archivo | Ubicación | Descripción |
|---------|-----------|-------------|
| `.env` | `backend/.env` | **Gemini API Key aquí** |
| `package.json` | `backend/package.json` | Dependencias backend |
| `package.json` | `frontend/package.json` | Dependencias frontend |

---

## 🎯 Backend (Node.js + Express + Gemini)

### Core Services
```
backend/src/services/
├── geminiService.js       # 🧠 Gemini AI integration
├── sessionService.js      # 💾 Session management
├── bobApiService.js       # 🔗 BOB API client
└── faqService.js          # ❓ FAQs management
```

### Controllers
```
backend/src/controllers/
├── chatController.js      # 💬 Chat endpoints
└── leadController.js      # 📊 Leads endpoints
```

### Routes
```
backend/src/routes/
├── chat.routes.js         # /api/chat/*
├── leads.routes.js        # /api/leads/*
└── resources.routes.js    # /api/faqs, /api/vehicles
```

### Data
```
backend/src/data/
├── faqs.csv              # 📋 62 Preguntas frecuentes
├── vehicles.csv          # 🚗 Vehículos en subasta
├── prompts.js            # 📝 System prompts para Gemini
├── sessions.json         # 💾 Sesiones activas (auto-generado)
└── leads.json            # 📊 Leads guardados (auto-generado)
```

### Main Server
```
backend/src/server.js     # 🚀 Express server principal
```

---

## 🎨 Frontend (React + Vite)

### Componentes
```
frontend/src/components/
├── ChatWidget.jsx        # 💬 Chat interactivo
├── ChatWidget.css        # 🎨 Estilos del chat
├── LeadsDashboard.jsx    # 📊 Dashboard de leads
└── LeadsDashboard.css    # 🎨 Estilos del dashboard
```

### App Principal
```
frontend/src/
├── App.jsx               # 🏠 Componente principal
├── App.css               # 🎨 Estilos principales
├── main.jsx              # 🚀 Entry point
└── index.css             # 🎨 Estilos globales
```

### Configuración
```
frontend/
├── vite.config.js        # ⚙️ Config Vite + proxy
└── index.html            # 📄 HTML principal
```

---

## 🔑 Archivos Críticos (NO TOCAR sin saber)

| Archivo | Ubicación | Por Qué es Crítico |
|---------|-----------|-------------------|
| `.env` | `backend/.env` | Contiene Gemini API Key |
| `geminiService.js` | `backend/src/services/` | Integración con IA |
| `server.js` | `backend/src/server.js` | Servidor principal |
| `ChatWidget.jsx` | `frontend/src/components/` | UI del chat |

---

## 📝 Archivos que SÍ puedes modificar

### Para cambiar prompts/comportamiento del bot:
```
backend/src/data/prompts.js
```

### Para agregar/editar FAQs:
```
backend/src/data/faqs.csv
```

### Para cambiar colores/estilos:
```
frontend/src/App.css
frontend/src/components/ChatWidget.css
frontend/src/components/LeadsDashboard.css
```

### Para agregar nuevos endpoints:
```
backend/src/controllers/
backend/src/routes/
```

---

## 🚀 Comandos Importantes

### Backend
```bash
cd backend
npm start              # Iniciar servidor
npm install            # Instalar dependencias
```

### Frontend
```bash
cd frontend
npm run dev            # Iniciar dev server
npm run build          # Build para producción
npm install            # Instalar dependencias
```

---

## 🔗 URLs Importantes

| Servicio | URL | Estado |
|----------|-----|--------|
| Frontend | http://localhost:5173 | ✅ Funcionando |
| Backend | http://localhost:3000 | ✅ Funcionando |
| Health Check | http://localhost:3000/health | ✅ Disponible |
| API Docs | http://localhost:3000 | ✅ Disponible |

---

## 📊 Estructura del Proyecto

```
bob-hackathon/
│
├── 📖 README.md                  # Documentación principal
├── 📖 QUICKSTART.md              # Guía rápida
├── 📖 PARA_EL_EQUIPO.md          # Instrucciones equipo
├── 📖 ARCHIVOS_IMPORTANTES.md    # Este archivo
│
├── 📁 backend/                   # Node.js + Express + Gemini
│   ├── src/
│   │   ├── controllers/          # Lógica de endpoints
│   │   ├── services/             # Servicios (AI, Session, APIs)
│   │   ├── routes/               # Definición de rutas
│   │   ├── data/                 # FAQs, Prompts, Data
│   │   └── server.js             # Servidor principal
│   ├── .env                      # ⚠️ API Keys (NO COMMITEAR)
│   ├── .gitignore
│   └── package.json
│
└── 📁 frontend/                  # React + Vite
    ├── src/
    │   ├── components/           # ChatWidget, Dashboard
    │   ├── App.jsx               # Componente principal
    │   └── main.jsx              # Entry point
    ├── index.html
    ├── vite.config.js
    └── package.json
```

---

## ✅ Checklist de Verificación

### Antes de Demo:
- [ ] Backend corriendo → `curl http://localhost:3000/health`
- [ ] Frontend corriendo → Abrir http://localhost:5173
- [ ] Chat funciona → Enviar mensaje de prueba
- [ ] Dashboard muestra leads → Ver pestaña "Leads"
- [ ] Gemini responde → Verificar respuestas del bot
- [ ] Integración WhatsApp lista → Probar mensaje

---

## 🆘 En Caso de Emergencia

### Si algo no funciona:

**1. Backend no inicia:**
```bash
cd backend
rm -rf node_modules
npm install
npm start
```

**2. Frontend no inicia:**
```bash
cd frontend
rm -rf node_modules
npm install
npm run dev
```

**3. Gemini no responde:**
- Verificar `backend/.env` tiene la key correcta
- Revisar modelo: debe ser `gemini-2.5-flash`
- Verificar internet

**4. Frontend no conecta con backend:**
- Verificar `frontend/vite.config.js` tiene proxy correcto
- Backend debe estar en puerto 3000
- Frontend debe estar en puerto 5173

---

## 📞 Ayuda

**Kevin (Backend/Frontend)**
- Gemini AI
- Scoring system
- API REST
- Frontend React

**Compañero (WhatsApp)**
- Bot de Go
- Integración WhatsApp
- Solo necesita llamar a `http://localhost:3000/api/chat/message`

---

## 🎯 Próximos Pasos

1. ✅ **Verificar que todo funciona** → QUICKSTART.md
2. ⏳ **Integrar WhatsApp** → PARA_EL_EQUIPO.md
3. ⏳ **Practicar demo** → Flujo de presentación
4. ⏳ **Preparar pitch** → Destacar features clave

---

**Todo está listo. ¡A ganar el hackathon!** 🏆🚀

**Hackathon BOB 2025**

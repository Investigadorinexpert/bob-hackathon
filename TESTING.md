# entorno de testing para claude

este entorno permite testing exhaustivo del sistema sin afectar datos de produccion.

## setup rapido

```bash
./setup_test_env.sh
```

esto crea:
- `backend/data-test/` - directorio de datos aislado
- `backend/.env.test` - configuracion para test (puerto 3001)
- `backend/data-test/sessions.json` - sesiones vacias
- `backend/data-test/leads.json` - leads vacios
- `backend/data-test/faqs.csv` - copia de faqs produccion

## iniciar servidor de test

```bash
cd backend
export $(cat .env.test | xargs)
go run cmd/server/main.go
```

el servidor corre en:
- puerto: 3001 (no conflicto con produccion en 3000)
- data: `data-test/` (aislado de `data/`)

## script de explotacion completa

`exploit_system.py` permite:

### opcion 1: stress test + analisis completo

simula 6 clientes diferentes:
1. comprador urgente (hot lead esperado)
2. tire-patadas (cold/discarded esperado)
3. interesado moderado (warm lead esperado)
4. spam/ambiguo
5. solo faqs
6. buscador vehiculos

luego analiza:
- sessions.json raw
- leads.json raw
- stats via api
- cada lead con 7 dimensiones detalladas
- historial completo de conversaciones
- scoring, boosts, penalizaciones

### opcion 2: analisis de data actual

lee y muestra:
- todos los archivos de datos
- todos los leads via api
- todas las sesiones
- estadisticas completas

### opcion 3: cliente custom

permite simular un cliente paso a paso y ver:
- respuestas del bot
- routing de agentes
- scoring progresivo
- lead final generado

## ejecucion

```bash
python3 exploit_system.py
```

## que puede ver claude

con este setup, claude puede:

### datos raw
- leer `backend/data-test/sessions.json` directamente
- leer `backend/data-test/leads.json` directamente
- ver toda la metadata, timestamps, contexto

### pipeline completo
- mensaje del usuario
- respuesta del orchestrator
- routing a agente (faq/auction/scoring)
- respuesta final
- score calculado
- categoria asignada

### scoring detallado
- las 7 dimensiones con puntajes
- boosts aplicados (+3 a +7)
- penalizaciones aplicadas (-2 a -6)
- score total
- clasificacion (hot/warm/cold/discarded)

### metricas
- total leads por categoria
- distribucion por canal
- score promedio
- total mensajes por sesion

### conversaciones completas
- historial mensaje por mensaje
- role (user/assistant)
- timestamp de cada mensaje
- contenido completo

## ventajas

- no afecta datos de produccion
- testing agresivo permitido
- multiples clientes simultaneos
- reset facil (borrar data-test y re-run setup)
- puerto diferente (3001) no interfiere con 3000

## reset del entorno

```bash
rm -rf backend/data-test
./setup_test_env.sh
```

## configuracion

`DATA_DIR` en `.env` controla donde se guardan datos:
- produccion: `DATA_DIR=data` (puerto 3000)
- testing: `DATA_DIR=data-test` (puerto 3001)

esto hace el backend 100% aislable para testing sin tocar codigo.

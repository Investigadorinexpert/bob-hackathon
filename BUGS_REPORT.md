# reporte de bugs - bob chatbot

analisis exhaustivo del sistema multiagente. encontrados 7 bugs criticos.

## bug 1: tiempos de respuesta inaceptables ⚠️ CRITICO

**ubicacion**: todo el sistema
**severidad**: ALTA - UX completamente rota

**problema**:
- primeros mensajes: 4-7 segundos
- mensajes subsecuentes: 22-30 segundos
- timeout acumulativo por multiples llamadas a gemini

**evidencia**:
```
[GIN] 00:30:05 | 200 |  6.789452709s | POST /api/chat/message
[GIN] 00:30:50 | 200 | 22.817683833s | POST /api/chat/message
[GIN] 00:31:19 | 200 | 28.207856834s | POST /api/chat/message
[GIN] 00:32:16 | 200 |     30.50093s | POST /api/chat/message
```

**impacto**: experiencia de usuario completamente inaceptable. ningún usuario esperará 30 segundos por respuesta.

**solucion sugerida**:
- agregar timeout a http clients (bob api, gemini)
- implementar streaming de respuestas
- cachear respuestas de gemini para patrones comunes
- considerar modelo más rápido (gemini-flash en vez de gemini-2.5-flash)

---

## bug 2: inconsistencia mayusculas/minusculas en categorias ⚠️ MEDIO

**ubicacion**: `backend/internal/agents/scoring_agent.go:336`

**problema**:
```go
// scoring_agent.go:336
msg += fmt.Sprintf("**Lead Score: %d/100** - Categoría: %s\n\n",
    categoryEmoji, data.TotalScore, strings.ToUpper(data.Category))
                                    ^^^^^^^^^^^^^^^
```

el mensaje usa `strings.ToUpper()` pero los logs usan el valor original:

**evidencia**:
```
Score: 59 (cold)   <- minúsculas
Score: 63 (COLD)   <- MAYÚSCULAS (mensaje al usuario)
Score: 85 (hot)    <- minúsculas
```

**impacto**: inconsistencia visual en UI y logs, confusion del usuario

**solucion**: usar consistentemente minúsculas o mayúsculas en todo el sistema

---

## bug 3: scoring erratico con saltos brutales ⚠️ ALTO

**ubicacion**: `backend/internal/agents/scoring_agent.go`

**problema**: scoring no es progresivo ni estable

**evidencia**:
```
mensaje 6:  score 59 (cold)
mensaje 8:  score 63 (cold)  +4 puntos
mensaje 10: score 85 (hot)   +22 puntos ← SALTO BRUTAL
mensaje 12: score 89 (hot)   +4 puntos
mensaje 14: score 100 (hot)  +11 puntos
```

de 63 a 85 en 2 mensajes (+22 puntos) es un salto muy grande que indica:
1. gemini no está siendo consistente en su scoring
2. no hay validación de que el scoring sea progresivo
3. conversaciones similares darán scores muy diferentes

**impacto**: leads clasificados incorrectamente, falta de confiabilidad

**solucion sugerida**:
- implementar validación que el score no cambie más de ±15 puntos entre mensajes
- guardar scoring previo y usarlo como contexto para gemini
- agregar smoothing: `new_score = 0.7 * prev_score + 0.3 * calculated_score`

---

## bug 4: auction agent completamente roto 🔥 CRITICO

**ubicacion**: `backend/internal/agents/auction_agent.go:39-44`

**problema**: TODAS las respuestas del auction agent fallan

**codigo problematico**:
```go
// auction_agent.go:39
vehicles, err := a.bobAPIService.GetSublots(false)
if err != nil {
    return &AgentOutput{
        Response: "Lo siento, tuve un problema consultando las subastas disponibles. ¿Podrías intentar de nuevo?",
    }, nil  // ← ERROR SILENCIADO
}
```

**evidencia** (sessions.json):
```json
{
  "role": "user",
  "content": "hola, necesito comprar un auto URGENTE para mi empresa"
},
{
  "role": "assistant",
  "content": "Lo siento, tuve un problema consultando las subastas disponibles. ¿Podrías intentar de nuevo?"
}
```

TODAS las 8 respuestas del auction agent son este error.

**causa raiz**: `bob_api_service.go:50`
```go
resp, err := http.Get(url)  // sin timeout configurado
```

la api de bob responde pero muy lento o devuelve demasiados datos, causando timeout implícito del http client de go.

**test manual**:
```bash
curl "https://apiv3.somosbob.com/v3/sublots/details"
# responde 200 OK pero se queda descargando datos indefinidamente
```

**impacto**:
- funcionalidad principal del chatbot COMPLETAMENTE ROTA
- usuarios no pueden buscar vehículos (razon de ser del bot)
- sistema INUTILIZABLE para producción

**solucion urgente**:
```go
// bob_api_service.go
client := &http.Client{
    Timeout: 10 * time.Second,
}
resp, err := client.Get(url)

// auction_agent.go - NO silenciar error
if err != nil {
    return nil, fmt.Errorf("error obteniendo vehiculos: %w", err)
}
```

---

## bug 5: score imposible de 100/100 ⚠️ ALTO

**ubicacion**: `backend/internal/agents/scoring_agent.go:298`

**problema**: scoring acepta cualquier valor de gemini sin validación

**codigo**:
```go
// scoring_agent.go:298
return &models.ScoringData{
    TotalScore:         scoring.TotalScore,  // ← CONFIA CIEGAMENTE EN GEMINI
    Category:           scoring.Category,
    DimensionScores:    dimensionScores,
    Boosts:             scoring.Boosts,
    Penalizaciones:     scoring.Penalizaciones,
    ...
}
```

**evidencia**:
```json
{
  "sessionId": "whatsapp-6501...",
  "score": 100,
  "category": "hot"
}
```

score perfecto de 100/100 es estadísticamente imposible dado que:
- 7 dimensiones suman max 100 puntos (10+15+25+15+10+10+15)
- boosts max +7
- para llegar a 100 necesitas score perfecto en TODAS las dimensiones

el prompt dice (línea 204):
```
4. El totalScore debe ser la suma de todas las dimensiones + boosts - penalizaciones
```

pero el código NO valida esto.

**impacto**: scores irreales, leads mal clasificados

**solucion**:
```go
// validar que totalScore == suma de dimensiones + boosts - penalties
calculatedScore := 0
for _, score := range dimensionScores {
    calculatedScore += score
}
// agregar boosts
// restar penalizaciones

if scoring.TotalScore != calculatedScore {
    log.Printf("Warning: Gemini score %d != calculated %d",
        scoring.TotalScore, calculatedScore)
    scoring.TotalScore = calculatedScore  // usar el calculado
}
```

---

## bug 6: stats ignoran leads "discarded" ⚠️ MEDIO

**ubicacion**: `backend/internal/services/session_service.go:196-203`

**codigo problematico**:
```go
switch lead.Category {
case "hot":
    stats.Hot++
case "warm":
    stats.Warm++
case "cold":
    stats.Cold++
}
// NO HAY case "discarded" !!!
```

**evidencia**:
leads.json:
```json
{
  "sessionId": "web-25b968f5...",
  "score": 7,
  "category": "discarded"  // ← existe
}
```

stats api:
```json
{
  "total": 2,
  "hot": 1,
  "warm": 0,
  "cold": 0  // ← discarded no aparece
}
```

**impacto**:
- métricas incorrectas
- leads descartados invisibles en dashboard
- total != hot + warm + cold + discarded

**solucion**:
```go
stats := &models.LeadStats{
    Total:     len(s.leads),
    Hot:       0,
    Warm:      0,
    Cold:      0,
    Discarded: 0,  // ← AGREGAR
    ...
}

switch lead.Category {
case "hot":
    stats.Hot++
case "warm":
    stats.Warm++
case "cold":
    stats.Cold++
case "discarded":
    stats.Discarded++  // ← AGREGAR
}
```

---

## bug 7: http client sin timeout ⚠️ ALTO

**ubicacion**: `backend/internal/services/bob_api_service.go:50`

**problema**: cliente http sin timeout configurado

**codigo**:
```go
// bob_api_service.go:50
url := fmt.Sprintf("%s/sublots/details", b.baseURL)
resp, err := http.Get(url)  // usa DefaultClient sin timeout
```

**impacto**:
- requests pueden colgarse indefinidamente
- goroutines leak si api no responde
- timeout implícito de OS (puede ser 2+ minutos)
- en producción esto causará memory leaks

**solucion**:
```go
client := &http.Client{
    Timeout: 10 * time.Second,
}
resp, err := client.Get(url)
```

---

## problemas adicionales encontrados (no bugs pero importantes)

### 1. sin rate limiting
gemini api probablemente tiene rate limits. sistema puede ser baneado por uso excesivo.

### 2. sin retries
si gemini falla, el request muere. debería reintentar 2-3 veces.

### 3. sin circuit breaker
si bob api está caída, seguirá intentando en cada request. debería implementar circuit breaker pattern.

### 4. logs en español e inglés mezclados
```
2025/11/06 00:30:02 🔀 Ruteando a Auction Agent
2025/11/06 00:30:19 🔀 Ruteando a Auction Agent
```
vs
```
[GIN] POST /api/chat/message
```

### 5. sin observabilidad
no hay métricas, tracing, ni monitoring. imposible debuggear en producción.

---

## prioridad de fixes

### p0 - bloqueantes (no funciona en producción):
1. **bug 4**: auction agent roto (sistema inutilizable)
2. **bug 1**: tiempos de respuesta (ux inaceptable)

### p1 - criticos (datos incorrectos):
3. **bug 5**: scores imposibles (confiabilidad)
4. **bug 3**: scoring erratico (clasificación incorrecta)

### p2 - importantes (calidad):
5. **bug 7**: http timeout (estabilidad)
6. **bug 6**: stats incorrectas (métricas)

### p3 - menores (polish):
7. **bug 2**: inconsistencia mayúsculas (cosmético)

---

## testing realizado

se ejecutó stress test con 6 clientes simultáneos:
- 1 comprador urgente (esperado: hot) → resultado: hot con score 100 (sospechoso)
- 1 tire-patadas (esperado: discarded) → resultado: discarded con score 7 (correcto)
- 4 clientes adicionales (interrumpido por timeout)

datos generados:
- 2 sessions completas
- 2 leads (1 hot, 1 discarded)
- ~16 mensajes procesados
- tiempo total: ~3 minutos para 2 conversaciones (inaceptable)

## archivos afectados

```
backend/internal/agents/auction_agent.go        - bug 4
backend/internal/agents/scoring_agent.go        - bugs 2, 3, 5
backend/internal/services/bob_api_service.go    - bugs 4, 7
backend/internal/services/session_service.go    - bug 6
```

## conclusión

el sistema tiene bugs críticos que lo hacen **NO APTO PARA PRODUCCIÓN**:

1. auction agent completamente roto (funcionalidad principal)
2. tiempos de respuesta inaceptables (30 segundos)
3. scoring no confiable (saltos de 22 puntos, scores imposibles)
4. métricas incorrectas (discarded invisible)

**recomendación**: NO DEPLOYAR hasta fix de p0 y p1.

**esfuerzo estimado de fixes**:
- p0: 4-6 horas
- p1: 6-8 horas
- total: 10-14 horas para tener sistema production-ready

---

fecha: 2025-11-06
entorno: test aislado (puerto 3001, data-test/)

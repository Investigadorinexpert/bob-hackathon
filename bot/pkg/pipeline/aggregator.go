// /home/investigadorinexperto/Negocios/bot/pkg/pipeline/aggregator.go
package pipeline

import (
	"sync"
	"time"
)

// Llama onFlush(chat, count) cuando cierra la ventana.
// Ventana DESLIZANTE: cada Add o Touch reinicia el timer.
// - Add(chat): incrementa el contador y reinicia la ventana.
// - Touch(chat): NO incrementa el contador, pero reinicia la ventana.
// Uso de "generación" para evitar carreras con timers viejos.
type Aggregator struct {
	mu      sync.Mutex
	perChat map[string]*batch
	window  time.Duration
	onFlush func(chat string, count int)
	onReset func(chat string, reason string, count int, window time.Duration)
}

type batch struct {
	count int
	timer *time.Timer
	gen   uint64
}

// NewAggregator crea un agregador de ventana deslizante.
// onFlush(chat, count) se invoca al cerrar la ventana con count>0.
// onReset(chat, reason, count, window) se invoca en:
//   - "start": cuando se crea el batch para un chat
//   - "message": después de Add (count ya incrementado)
//   - "typing": después de Touch (count NO cambia)
func NewAggregator(
	window time.Duration,
	onFlush func(chat string, count int),
	onReset func(chat string, reason string, count int, window time.Duration),
) *Aggregator {
	return &Aggregator{
		perChat: make(map[string]*batch),
		window:  window,
		onFlush: onFlush,
		onReset: onReset,
	}
}

// Add incrementa el conteo y reinicia la ventana.
func (a *Aggregator) Add(chat string) {
	if chat == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	b := a.ensureBatchLocked(chat)
	b.count++
	a.resetTimerLocked(chat, b)

	if a.onReset != nil {
		a.onReset(chat, "message", b.count, a.window)
	}
}

// Touch reinicia la ventana sin incrementar el conteo.
// Útil para eventos como "usuario está escribiendo".
func (a *Aggregator) Touch(chat string) {
	if chat == "" {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	b := a.ensureBatchLocked(chat)
	// NO incrementa b.count
	a.resetTimerLocked(chat, b)

	if a.onReset != nil {
		a.onReset(chat, "typing", b.count, a.window)
	}
}

// TouchTyping es un alias semántico de Touch para eventos "usuario está escribiendo".
func (a *Aggregator) TouchTyping(chat string) {
	a.Touch(chat)
}

// ensureBatchLocked obtiene o crea el batch del chat.
// Debe llamarse con el candado tomado.
func (a *Aggregator) ensureBatchLocked(chat string) *batch {
	if b, ok := a.perChat[chat]; ok {
		return b
	}
	// Primera generación
	b := &batch{count: 0, gen: 1}
	gen := b.gen
	b.timer = time.AfterFunc(a.window, func() {
		a.flushGen(chat, gen)
	})
	a.perChat[chat] = b

	// Log de inicio de ventana
	if a.onReset != nil {
		a.onReset(chat, "start", b.count, a.window)
	}
	return b
}

// resetTimerLocked reinicia el timer de un batch en una ventana deslizante,
// manejando la carrera en que el timer pudo estar por disparar.
// Debe llamarse con el candado tomado.
func (a *Aggregator) resetTimerLocked(chat string, b *batch) {
	if b.timer == nil {
		b.gen++
		gen := b.gen
		b.timer = time.AfterFunc(a.window, func() {
			a.flushGen(chat, gen)
		})
		return
	}
	if b.timer.Stop() {
		// Si alcanzamos a detenerlo, reseteamos sobre el mismo timer.
		b.timer.Reset(a.window)
		return
	}
	// El timer pudo haber disparado o estar ejecutándose:
	// creamos una NUEVA generación y un nuevo timer.
	b.gen++
	gen := b.gen
	b.timer = time.AfterFunc(a.window, func() {
		a.flushGen(chat, gen)
	})
}

// flushGen sólo hace flush si la generación del timer coincide con la generación activa.
func (a *Aggregator) flushGen(chat string, gen uint64) {
	a.mu.Lock()
	b, ok := a.perChat[chat]
	if !ok {
		a.mu.Unlock()
		return
	}
	if b.gen != gen {
		a.mu.Unlock()
		return
	}
	delete(a.perChat, chat)
	count := b.count
	a.mu.Unlock()

	if count > 0 && a.onFlush != nil {
		a.onFlush(chat, count)
	}
}

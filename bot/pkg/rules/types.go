package rules

import (
	"context"
	"regexp"
	"strings"
	"time"
)

type Envelope struct {
	EventType string
	ChatJID   string
	SenderJID string
	ChatName  string
	MessageID string
	Text      string
	At        time.Time
	// + lo que necesites
}

type ActionResult struct {
	Handled bool
	Reply   string            // respuesta de texto
	Meta    map[string]string // datos extra para logging o side-effects
}

type Action func(ctx context.Context, env Envelope) (ActionResult, error)
type Predicate func(env Envelope) bool

type Rule struct {
	Name      string
	WhenAll   []Predicate
	Then      Action
	Priority  int // mayor primero
	StopChain bool
}

// Helpers de predicado
func Command(cmd string) Predicate {
	prefix := "/" + strings.ToLower(strings.TrimPrefix(cmd, "/"))
	return func(env Envelope) bool {
		t := strings.TrimSpace(strings.ToLower(env.Text))
		return t == prefix || strings.HasPrefix(t, prefix+" ")
	}
}
func Contains(substr string) Predicate {
	substr = strings.ToLower(substr)
	return func(env Envelope) bool {
		return strings.Contains(strings.ToLower(env.Text), substr)
	}
}
func Regex(rx *regexp.Regexp) Predicate {
	return func(env Envelope) bool { return rx.MatchString(env.Text) }
}
func OnDM() Predicate {
	return func(env Envelope) bool { return !strings.HasSuffix(env.ChatJID, "@g.us") }
}
func OnGroup() Predicate {
	return func(env Envelope) bool { return strings.HasSuffix(env.ChatJID, "@g.us") }
}

// Helpers de acción
func Reply(text string) Action {
	return func(ctx context.Context, env Envelope) (ActionResult, error) {
		return ActionResult{Handled: true, Reply: text}, nil
	}
}

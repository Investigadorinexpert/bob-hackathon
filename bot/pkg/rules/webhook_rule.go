package rules

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*** Perfil mínimo compatible con tu snapshot ***/
type profileBlock struct {
	Spam      bool      `json:"spam"`
	Malicious bool      `json:"malicious"`
	Permanent bool      `json:"permanent"`
	Until     time.Time `json:"until"`
}
type mediaEntry struct {
	Direction string    `json:"direction"`
	Type      string    `json:"type"`
	Mimetype  string    `json:"mimetype"`
	At        time.Time `json:"at"`
}
type profileMedia struct {
	In  []mediaEntry `json:"in"`
	Out []mediaEntry `json:"out"`
}
type ProfileSnapshot struct {
	SenderJID string            `json:"sender_jid"`
	Lang      string            `json:"lang"`
	Tier      string            `json:"tier"`
	Tags      map[string]string `json:"tags"`
	Media     profileMedia      `json:"media"`
	Block     profileBlock      `json:"block"`
}

/*** Utils ***/
func sanitizePathPart(s string) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "unknown"
	}
	return b.String()
}
func loadProfileSnapshot(chatJID string) (*ProfileSnapshot, error) {
	path := filepath.Join("outbox", "profiles", sanitizePathPart(chatJID)+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p ProfileSnapshot
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
func profileAllowsWebhook(p *ProfileSnapshot) bool {
	if p == nil {
		return false
	}
	// regla base: si cualquiera de los bloqueos está activo, NO mandamos
	if p.Block.Spam || p.Block.Malicious || p.Block.Permanent {
		return false
	}
	// aquí puedes meter más gates: tier/lang/tags/lo que quieras
	return true
}
func pickWebhookURL(p *ProfileSnapshot) string {
	// Prioridad: por tier -> por lang -> default
	if p != nil {
		if u := strings.TrimSpace(os.Getenv("N8N_" + strings.ToUpper(p.Tier) + "_WEBHOOK_URL")); u != "" {
			return u
		}
		if u := strings.TrimSpace(os.Getenv("N8N_LANG_" + strings.ToUpper(p.Lang) + "_WEBHOOK_URL")); u != "" {
			return u
		}
	}
	return strings.TrimSpace(os.Getenv("N8N_WEBHOOK_URL"))
}

/*** Payload de salida hacia n8n ***/
type webhookReq struct {
	ChatJID    string           `json:"chat_jid"`
	Aggregated string           `json:"aggregated"` // último texto del batch (env.Text)
	Profile    *ProfileSnapshot `json:"profile"`    // incluye media.in/out recientes (tu “acumulado” práctico)
}
type webhookResp struct {
	Reply string `json:"reply"`
}

/*** Predicate + Action públicos para usar en Builtin() ***/

// Predicate: verifica si el perfil permite mandar a webhook
func ProfileAllowsWebhook() Predicate {
	return func(env Envelope) bool {
		p, err := loadProfileSnapshot(env.ChatJID)
		if err != nil {
			return false
		}
		return profileAllowsWebhook(p)
	}
}

// Action: hace POST a n8n y devuelve reply si hay
func SendToWebhook() Action {
	return func(ctx context.Context, env Envelope) (ActionResult, error) {
		p, err := loadProfileSnapshot(env.ChatJID)
		if err != nil || !profileAllowsWebhook(p) {
			return ActionResult{Handled: false}, nil
		}
		url := pickWebhookURL(p)
		if url == "" {
			return ActionResult{Handled: false}, nil
		}

		reqBody := webhookReq{
			ChatJID:    env.ChatJID,
			Aggregated: strings.TrimSpace(env.Text),
			Profile:    p,
		}
		b, _ := json.Marshal(reqBody)

		to := 5 * time.Second
		if ms := strings.TrimSpace(os.Getenv("N8N_TIMEOUT_MS")); ms != "" {
			// “1234” -> 1234ms
			if d, err := time.ParseDuration(ms + "ms"); err == nil && d > 0 {
				to = d
			}
		}
		httpc := &http.Client{Timeout: to}
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		if tok := strings.TrimSpace(os.Getenv("N8N_AUTH_TOKEN")); tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}

		resp, err := httpc.Do(req)
		if err != nil {
			return ActionResult{Handled: false}, nil
		}
		defer resp.Body.Close()
		if resp.StatusCode/100 != 2 {
			io.Copy(io.Discard, resp.Body)
			return ActionResult{Handled: false}, nil
		}

		var out webhookResp
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return ActionResult{Handled: false}, nil
		}
		reply := strings.TrimSpace(out.Reply)
		if reply == "" {
			return ActionResult{Handled: false}, nil
		}
		return ActionResult{Handled: true, Reply: reply}, nil
	}
}

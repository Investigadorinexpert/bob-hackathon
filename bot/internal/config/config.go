// internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	// ===== Engine/Bot =====
	DBPath                string
	MsgDBPath             string
	HTTPPort              int
	EnableStatus          bool
	BackupEvery           time.Duration
	MaxConnAttempts       int
	ReconnectBaseDelay    time.Duration
	SendPresenceAvailable bool // nuevo

	// ===== Forward (Folder + Webhook) =====
	ForwardMode      string
	Outbox           string
	ContextDepth     int
	ForwardExtraJSON map[string]any

	WebhookEnabled bool
	WebhookURL     string
	WebhookSecret  string
	WebhookHeaders map[string]string

	// ===== Typing =====
	TypingEnabled    bool
	TypingDebounce   time.Duration
	TypingPauseAfter time.Duration
	TypingMedia      string

	// ===== Rules =====
	RulesMode     string
	RulesJSONPath string

	// ===== Server (whserver) =====
	ServerAddr             string
	ServerRequireSig       bool
	ServerBodyLimit        int64
	ServerTSSkew           time.Duration
	ServerLogJSON          bool
	ServerDedupeWindow     time.Duration
	ServerUseTimestamp     bool
	ServerAllowNoSecretDev bool

	// Punteros REST del server hacia el engine
	ServerEngineSendURL     string // WH_ENGINE_SEND_URL
	ServerEngineTypingURL   string // WH_ENGINE_TYPING_URL
	ServerEngineMarkReadURL string // WH_ENGINE_MARKREAD_URL   <-- NUEVO

	// ===== Reply typing wait (tunable por .env) =====
	ReplyBaseWait  time.Duration
	ReplyPerCharMs int
	ReplyJitterMs  int
	ReplyMaxWait   time.Duration
	PreReplyDelay  time.Duration
	AggWindow      time.Duration
}

// ---------- helpers ----------
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
func getenvBool01(k string, def bool) bool {
	switch os.Getenv(k) {
	case "1":
		return true
	case "0":
		return false
	default:
		return def
	}
}
func getenvDur(k, def string) time.Duration {
	raw := getenv(k, def)
	d, err := time.ParseDuration(raw)
	if err != nil {
		d, _ = time.ParseDuration(def)
	}
	return d
}

func Load() *AppConfig {
	// Headers extra para webhook
	var hdrs map[string]string
	_ = json.Unmarshal([]byte(getenv("WH_WEBHOOK_HEADERS_JSON", "{}")), &hdrs)

	// Extra params para adjuntar en el envelope
	var extra map[string]any
	_ = json.Unmarshal([]byte(getenv("WH_FORWARD_EXTRA_JSON", `{"tenant":"acme","lang":"es"}`)), &extra)

	base := getenv("WH_ENGINE_BASE_URL", "http://localhost:8080")

	return &AppConfig{
		// ===== Engine/Bot =====
		DBPath:                getenv("WH_DB_PATH", "data/session.db"),
		MsgDBPath:             getenv("WH_MSG_DB_PATH", "data/messages.db"),
		HTTPPort:              getenvInt("WH_HTTP_PORT", 8080),
		EnableStatus:          getenvBool01("WH_ENABLE_STATUS", true),
		BackupEvery:           getenvDur("WH_BACKUP_EVERY", "30m"),
		MaxConnAttempts:       getenvInt("WH_MAX_CONN_ATTEMPTS", 5),
		ReconnectBaseDelay:    getenvDur("WH_RECONNECT_BASE_DELAY", "2s"),
		SendPresenceAvailable: getenvBool01("WH_SEND_PRESENCE_AVAILABLE", true),

		// ===== Forward (Folder + Webhook) =====
		ForwardMode:      getenv("WH_FORWARD_MODE", "folder"),
		Outbox:           getenv("WH_OUTBOX", "outbox"),
		ContextDepth:     getenvInt("WH_CONTEXT_DEPTH", 3),
		ForwardExtraJSON: extra,

		WebhookEnabled: getenvBool01("WH_WEBHOOK_ENABLED", true),
		WebhookURL:     getenv("WH_WEBHOOK_URL", "http://127.0.0.1:9000/wh"),
		WebhookSecret:  getenv("WH_WEBHOOK_SECRET", ""),
		WebhookHeaders: hdrs,

		// ===== Typing =====
		TypingEnabled:    getenvBool01("WH_TYPING_ENABLED", true),
		TypingDebounce:   getenvDur("WH_TYPING_DEBOUNCE", "800ms"),
		TypingPauseAfter: getenvDur("WH_TYPING_PAUSE_AFTER", "3s"),
		TypingMedia:      getenv("WH_TYPING_MEDIA", "text"),

		// ===== Rules =====
		RulesMode:     getenv("WH_RULES_MODE", "code"),
		RulesJSONPath: getenv("WH_RULES_JSON_PATH", "rules.json"),

		// ===== Server (whserver) =====
		ServerAddr:             getenv("WH_SERVER_ADDR", ":9000"),
		ServerRequireSig:       getenvBool01("WH_REQUIRE_SIG", true),
		ServerBodyLimit:        int64(getenvInt("WH_BODY_LIMIT", 2<<20)),
		ServerTSSkew:           getenvDur("WH_TS_SKEW", "2m"),
		ServerLogJSON:          getenvBool01("WH_LOG_JSON", false),
		ServerDedupeWindow:     getenvDur("WH_DEDUPE_WINDOW", "10m"),
		ServerUseTimestamp:     getenvBool01("WH_USE_TIMESTAMP", false),
		ServerAllowNoSecretDev: getenvBool01("WH_ALLOW_NO_SECRET_DEV", true),

		// Punteros al engine
		ServerEngineSendURL:     getenv("WH_ENGINE_SEND_URL", base+"/api/send"),
		ServerEngineTypingURL:   getenv("WH_ENGINE_TYPING_URL", base+"/api/typing"),
		ServerEngineMarkReadURL: getenv("WH_ENGINE_MARKREAD_URL", base+"/api/markread"),

		// ===== Reply typing wait =====
		ReplyBaseWait:  getenvDur("WH_REPLY_BASE_WAIT", "400ms"),
		ReplyPerCharMs: getenvInt("WH_REPLY_PER_CHAR_MS", 35),
		ReplyJitterMs:  getenvInt("WH_REPLY_JITTER_MS", 400),
		ReplyMaxWait:   getenvDur("WH_REPLY_MAX_WAIT", "4s"),
		PreReplyDelay:  getenvDur("WH_PRE_REPLY_DELAY", "900ms"),
		AggWindow:      getenvDur("WH_AGGREGATOR_WINDOW", "2s"),
	}
}

package rules

import "regexp"

func Builtin() []Rule {
	rxEmail := regexp.MustCompile(`...`)

	// 🔝 Regla de webhook basada en perfil (más prioridad)
	hookRule := Rule{
		Name:     "profile→webhook",
		Priority: 1000, // mayor que las demás
		WhenAll:  []Predicate{ProfileAllowsWebhook()},
		Then:     SendToWebhook(),
	}

	return []Rule{
		hookRule,
		{
			Name:     "help",
			Priority: 100,
			WhenAll:  []Predicate{Command("help")},
			Then:     Reply("Comandos: /help, /ping, /demo"),
		},
		{
			Name:     "ping",
			Priority: 90,
			WhenAll:  []Predicate{Command("ping")},
			Then:     Reply("pong"),
		},
		{
			Name:     "demo group",
			Priority: 70,
			WhenAll:  []Predicate{OnGroup(), Command("demo")},
			Then:     Reply("Demo en grupo 👥"),
		},
		{
			Name:     "email detector",
			Priority: 60,
			WhenAll:  []Predicate{Regex(rxEmail)},
			Then:     Reply("Pillado un email, gracias!"),
		},
	}
}

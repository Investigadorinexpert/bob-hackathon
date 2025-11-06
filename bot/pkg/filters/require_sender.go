package filters

type RequireSender struct{}

func (RequireSender) Apply(e EnvView) bool { return e.SenderJID != "" }

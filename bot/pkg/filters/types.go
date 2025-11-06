package filters

type EnvView struct {
	Direction string
	SenderJID string
	ChatJID   string
}

type Filter interface{ Apply(e EnvView) bool }

type Chain struct{ Filters []Filter }

func (c Chain) Pass(e EnvView) bool {
	for _, f := range c.Filters {
		if !f.Apply(e) {
			return false
		}
	}
	return true
}

package filters

type NotOut struct{}

func (NotOut) Apply(e EnvView) bool { return e.Direction != "out" }

package soap

type Port struct {
	actions map[string]*Action
}

func NewPort() *Port {
	p := &Port{
		actions: make(map[string]*Action, 20),
	}

	return p
}

func (p *Port) AddAction(name string, a *Action) {
	p.actions[name] = a
}

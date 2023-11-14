package soap

type Port struct {
	actions map[string]IAction
}

func NewPort() *Port {
	p := &Port{
		actions: make(map[string]IAction, 20),
	}

	return p
}

func (p *Port) AddAction(name string, a IAction) {
	p.actions[name] = a
}

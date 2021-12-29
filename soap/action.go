package soap

type (
	ActionOpt     func(a *Action)
	ActionHandler func(in interface{}) (interface{}, error)
)

type Action struct {
	in            interface{}
	out           interface{}
	documentation string
	handler       ActionHandler
}

func WithDocumentation(s string) ActionOpt {
	return func(a *Action) {
		a.documentation = s
	}
}

func NewAction(in, out interface{}, handler ActionHandler, opts ...ActionOpt) *Action {
	a := &Action{
		in:      in,
		out:     out,
		handler: handler,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Action) IsOneWay() bool {
	return a.out == nil
}

func (a *Action) Run(in interface{}) (interface{}, error) {
	return a.handler(in)
}

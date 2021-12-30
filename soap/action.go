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

// NewAction create a new SOAP action
//  - in is the type of request. should be a struct or pointer to struct
//  - out is the type of response. same as in
// The fields of in & out will also be used to generate WSDL definitions
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

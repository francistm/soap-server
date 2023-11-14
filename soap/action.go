package soap

import (
	"context"
	"fmt"
	"reflect"
)

type (
	ActionOpt[I, O any]     func(a *Action[I, O])
	ActionHandler[I, O any] func(ctx context.Context, in I) (O, error)

	NilOut struct{}
)

type IAction interface {
	KindIn() any
	KindOut() any
	Run(context.Context, any) (any, error)
}

type Action[I, O any] struct {
	documentation string
	handler       ActionHandler[I, O]
}

// NewAction create a new SOAP action
// - in is the type of request. should be a struct or pointer to struct
// - out is the type of response. same as in
// The fields of in & out will also be used to generate WSDL definitions
func NewAction[I, O any](handler ActionHandler[I, O], opts ...ActionOpt[I, O]) *Action[I, O] {
	a := &Action[I, O]{
		handler: handler,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

func (a *Action[I, O]) Run(ctx context.Context, in any) (any, error) {
	typedIn, ok := in.(I)

	if !ok {
		return nil, fmt.Errorf("want request as %T, got %T", typedIn, in)
	}

	return a.handler(ctx, typedIn)
}

func (a *Action[I, O]) KindIn() any {
	var zeroIn I

	return zeroIn
}

func (a *Action[I, O]) KindOut() any {
	var (
		zeroOut        O
		zeroOutTypeRef = reflect.TypeOf(zeroOut)
	)

	for zeroOutTypeRef.Kind() == reflect.Ptr {
		zeroOutTypeRef = zeroOutTypeRef.Elem()
	}

	if zeroOutTypeRef.String() == "soap.NilOut" {
		return nil
	}

	return zeroOut
}

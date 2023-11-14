package soap

func WithDocumentation[I, O any](s string) ActionOpt[I, O] {
	return func(a *Action[I, O]) {
		a.documentation = s
	}
}

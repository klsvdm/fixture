package fixtures

type options struct {
	editor func(value any)
}

type Option func(*options)

func WithEditor[T any](editor func(value *T)) Option {
	return func(o *options) {
		o.editor = func(value any) {
			if v, ok := value.(*T); ok {
				editor(v)
			}
		}
	}
}

func applyOptions(opts []Option) *options {
	result := &options{}

	for _, option := range opts {
		option(result)
	}

	return result
}

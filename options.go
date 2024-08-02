package fixtures

type options[T any] struct {
	editor func(value *T)
}

type Option[T any] func(*options[T])

func WithEditor[T any](editor func(value *T)) Option[T] {
	return func(o *options[T]) {
		o.editor = editor
	}
}

func applyOptions[T any](opts []Option[T]) *options[T] {
	result := &options[T]{}

	for _, option := range opts {
		option(result)
	}

	return result
}

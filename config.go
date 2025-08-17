package stack

type Option[T any] func(*stack[T])

const (
	UnlimitedCapacity = -1
)

func WithCapacity[T any](cap int) Option[T] {
	return func(s *stack[T]) {
		if cap < UnlimitedCapacity {
			panic("cannot specify arbitrary negative capacity")
		}
		s.capacity = cap
	}
}

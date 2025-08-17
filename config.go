package stack

// Option represents a configuration function that can be applied to a stack during creation.
// Options follow the functional options pattern for flexible and extensible configuration.
type Option[T any] func(*stack[T])

const (
	// UnlimitedCapacity indicates that the stack should have no size limit.
	// This is the default capacity when no WithCapacity option is provided.
	UnlimitedCapacity = -1
)

// WithCapacity returns an option that sets the maximum capacity of the stack.
//
// The capacity must be >= 0 or equal to UnlimitedCapacity (-1).
// Any other negative value will cause a panic.
//
// Parameters:
//   - cap: The maximum number of items the stack can hold
//   - Use 0 for a stack that cannot hold any items
//   - Use any positive integer for a fixed capacity
//   - Use UnlimitedCapacity (-1) for unlimited capacity
//
// Example:
//
//	s := stack.New[int](stack.WithCapacity[int](100))  // Max 100 items
//	s := stack.New[int](stack.WithCapacity[int](0))    // No items allowed
//	s := stack.New[int](stack.WithCapacity[int](stack.UnlimitedCapacity)) // No limit
//
// Panics if cap < UnlimitedCapacity (i.e., cap < -1).
func WithCapacity[T any](cap int) Option[T] {
	return func(s *stack[T]) {
		if cap < UnlimitedCapacity {
			panic("cannot specify arbitrary negative capacity")
		}
		s.capacity = cap
	}
}

package stack

import "errors"

var (
	// ErrOverflow is returned when attempting to push an item to a stack
	// that has reached its maximum capacity.
	//
	// This error occurs when:
	//   - The stack was created with WithCapacity option
	//   - The current size equals the specified capacity
	//   - Push() is called on the full stack
	//
	// Example:
	//
	//	s := stack.New[int](stack.WithCapacity[int](2))
	//	s.Push(1) // OK
	//	s.Push(2) // OK
	//	err := s.Push(3) // Returns ErrOverflow
	//	if errors.Is(err, stack.ErrOverflow) {
	//		fmt.Println("Stack is full")
	//	}
	ErrOverflow = errors.New("stack overflow")

	// ErrUnderflow is returned when attempting to pop or peek at an empty stack.
	//
	// This error occurs when:
	//   - Pop() is called on an empty stack
	//   - Peek() is called on an empty stack
	//
	// When this error is returned, the operation returns the zero value for type T.
	//
	// Example:
	//
	//	s := stack.New[int]()
	//	val, err := s.Pop() // Returns 0, ErrUnderflow
	//	if errors.Is(err, stack.ErrUnderflow) {
	//		fmt.Println("Stack is empty")
	//	}
	ErrUnderflow = errors.New("stack underflow")
)

// Package stack provides a thread-safe, generic stack implementation with configurable capacity limits.
//
// The stack follows LIFO (Last-In-First-Out) semantics and supports any type through Go generics.
// All operations are safe for concurrent use across multiple goroutines.
//
// Example usage:
//
//	s := stack.New[int]()
//	s.Push(1)
//	s.Push(2)
//	val, err := s.Pop() // returns 2, nil
package stack

import (
	"sync"
)

// Stack defines the interface for a generic stack data structure.
// All operations are thread-safe and support any type T.
type Stack[T any] interface {
	// Push adds an item to the top of the stack.
	// Returns ErrOverflow if the stack is at capacity.
	Push(val T) error

	// Pop removes and returns the top item from the stack.
	// Returns ErrUnderflow if the stack is empty.
	Pop() (T, error)

	// Size returns the current number of items in the stack.
	Size() int

	// Peek returns the top item without removing it from the stack.
	// Returns ErrUnderflow if the stack is empty.
	Peek() (T, error)
}

// New creates a new stack with the specified options.
// If no options are provided, creates an unlimited capacity stack.
//
// Example:
//
//	s := stack.New[int]()                           // Unlimited capacity
//	s := stack.New[int](stack.WithCapacity[int](10)) // Capacity of 10
func New[T any](opts ...Option[T]) Stack[T] {
	return newStack(opts...)
}

type stack[T any] struct {
	mu       sync.RWMutex
	capacity int
	items    []T
}

func newStack[T any](opts ...Option[T]) *stack[T] {
	s := &stack[T]{
		capacity: UnlimitedCapacity,
	}
	for _, opt := range opts {
		opt(s)
	}

	s.items = make([]T, 0)

	return s
}

func (s *stack[T]) Push(val T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.capacity >= 0 && len(s.items)+1 > s.capacity {
		return ErrOverflow
	}

	s.items = append(s.items, val)

	return nil
}

func (s *stack[T]) Pop() (T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sz := len(s.items)
	if sz == 0 {
		var zero T
		return zero, ErrUnderflow
	}

	idx := sz - 1

	result := s.items[idx]
	s.items = s.items[:idx]

	return result, nil
}

func (s *stack[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

func (s *stack[T]) Peek() (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sz := len(s.items)
	if sz == 0 {
		var zero T
		return zero, ErrUnderflow
	}

	idx := sz - 1

	return s.items[idx], nil
}

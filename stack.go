package stack

import (
	"sync"
)

type Stack[T any] interface {
	Push(val T) error
	Pop() (T, error)
	Size() int
	Peek() (T, error)
}

type stack[T any] struct {
	mu       sync.RWMutex
	capacity int
	items    []T
}

func New[T any](opts ...Option[T]) Stack[T] {
	return newStack(opts...)
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

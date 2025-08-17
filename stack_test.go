package stack

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New[int]()
	if s == nil {
		t.Fatal("New() returned nil")
	}

	if size := s.Size(); size != 0 {
		t.Errorf("New stack size = %d, want 0", size)
	}
}

func TestNewWithCapacity(t *testing.T) {
	s := New[int](WithCapacity[int](5))
	if s == nil {
		t.Fatal("New() with capacity returned nil")
	}

	if size := s.Size(); size != 0 {
		t.Errorf("New stack with capacity size = %d, want 0", size)
	}
}

func TestPushPop(t *testing.T) {
	s := New[int]()

	// Test push
	err := s.Push(42)
	if err != nil {
		t.Errorf("Push(42) error = %v, want nil", err)
	}

	if size := s.Size(); size != 1 {
		t.Errorf("Size after push = %d, want 1", size)
	}

	// Test pop
	val, err := s.Pop()
	if err != nil {
		t.Errorf("Pop() error = %v, want nil", err)
	}

	if val != 42 {
		t.Errorf("Pop() = %d, want 42", val)
	}

	if size := s.Size(); size != 0 {
		t.Errorf("Size after pop = %d, want 0", size)
	}
}

func TestPushPopMultiple(t *testing.T) {
	s := New[int]()
	values := []int{1, 2, 3, 4, 5}

	// Push all values
	for _, v := range values {
		err := s.Push(v)
		if err != nil {
			t.Errorf("Push(%d) error = %v, want nil", v, err)
		}
	}

	if size := s.Size(); size != len(values) {
		t.Errorf("Size after pushes = %d, want %d", size, len(values))
	}

	// Pop all values (should be in reverse order)
	for i := len(values) - 1; i >= 0; i-- {
		val, err := s.Pop()
		if err != nil {
			t.Errorf("Pop() error = %v, want nil", err)
		}
		if val != values[i] {
			t.Errorf("Pop() = %d, want %d", val, values[i])
		}
	}

	if size := s.Size(); size != 0 {
		t.Errorf("Size after all pops = %d, want 0", size)
	}
}

func TestPeek(t *testing.T) {
	s := New[string]()

	// Test peek on empty stack
	_, err := s.Peek()
	if !errors.Is(err, ErrUnderflow) {
		t.Errorf("Peek() on empty stack error = %v, want ErrUnderflow", err)
	}

	// Push and peek
	err = s.Push("hello")
	if err != nil {
		t.Errorf("Push() error = %v, want nil", err)
	}

	val, err := s.Peek()
	if err != nil {
		t.Errorf("Peek() error = %v, want nil", err)
	}
	if val != "hello" {
		t.Errorf("Peek() = %q, want %q", val, "hello")
	}

	// Size should remain the same after peek
	if size := s.Size(); size != 1 {
		t.Errorf("Size after peek = %d, want 1", size)
	}

	// Push another value
	err = s.Push("world")
	if err != nil {
		t.Errorf("Push() error = %v, want nil", err)
	}

	// Peek should return the last pushed value
	val, err = s.Peek()
	if err != nil {
		t.Errorf("Peek() error = %v, want nil", err)
	}
	if val != "world" {
		t.Errorf("Peek() = %q, want %q", val, "world")
	}
}

func TestPopUnderflow(t *testing.T) {
	s := New[int]()

	val, err := s.Pop()
	if !errors.Is(err, ErrUnderflow) {
		t.Errorf("Pop() on empty stack error = %v, want ErrUnderflow", err)
	}

	// Check that zero value is returned
	if val != 0 {
		t.Errorf("Pop() on empty stack value = %d, want 0 (zero value)", val)
	}
}

func TestCapacityOverflow(t *testing.T) {
	s := New[int](WithCapacity[int](2))

	// Push up to capacity
	err := s.Push(1)
	if err != nil {
		t.Errorf("Push(1) error = %v, want nil", err)
	}

	err = s.Push(2)
	if err != nil {
		t.Errorf("Push(2) error = %v, want nil", err)
	}

	// Try to exceed capacity
	err = s.Push(3)
	if !errors.Is(err, ErrOverflow) {
		t.Errorf("Push(3) exceeding capacity error = %v, want ErrOverflow", err)
	}

	// Size should still be 2
	if size := s.Size(); size != 2 {
		t.Errorf("Size after overflow attempt = %d, want 2", size)
	}
}

func TestUnlimitedCapacity(t *testing.T) {
	s := New[int]()

	// Push many items
	for i := 0; i < 1000; i++ {
		err := s.Push(i)
		if err != nil {
			t.Errorf("Push(%d) with unlimited capacity error = %v, want nil", i, err)
		}
	}

	if size := s.Size(); size != 1000 {
		t.Errorf("Size with unlimited capacity = %d, want 1000", size)
	}
}

func TestDifferentTypes(t *testing.T) {
	t.Run("string stack", func(t *testing.T) {
		s := New[string]()
		s.Push("first")
		s.Push("second")

		val, _ := s.Pop()
		if val != "second" {
			t.Errorf("Pop() = %q, want %q", val, "second")
		}
	})

	t.Run("struct stack", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		s := New[Person]()
		person := Person{Name: "Alice", Age: 30}
		s.Push(person)

		val, err := s.Pop()
		if err != nil {
			t.Errorf("Pop() error = %v, want nil", err)
		}
		if val != person {
			t.Errorf("Pop() = %+v, want %+v", val, person)
		}
	})
}

func TestConcurrency(t *testing.T) {
	s := New[int]()
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup

	// Start multiple goroutines pushing values
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				s.Push(start*numOperations + j)
			}
		}(i)
	}

	wg.Wait()

	// Check that all items were pushed
	expectedSize := numGoroutines * numOperations
	if size := s.Size(); size != expectedSize {
		t.Errorf("Size after concurrent pushes = %d, want %d", size, expectedSize)
	}

	// Start multiple goroutines popping values
	results := make(chan int, expectedSize)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				if val, err := s.Pop(); err == nil {
					results <- val
				}
			}
		}()
	}

	wg.Wait()
	close(results)

	// Count results
	count := 0
	for range results {
		count++
	}

	if count != expectedSize {
		t.Errorf("Popped %d items, want %d", count, expectedSize)
	}

	// Stack should be empty
	if size := s.Size(); size != 0 {
		t.Errorf("Size after concurrent pops = %d, want 0", size)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("zero capacity", func(t *testing.T) {
		s := New[int](WithCapacity[int](0))
		err := s.Push(1)
		if !errors.Is(err, ErrOverflow) {
			t.Errorf("Push() with zero capacity error = %v, want ErrOverflow", err)
		}
	})

	t.Run("negative capacity (should panic)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("WithCapacity(-5) should panic, but it didn't")
			}
		}()

		// This should panic
		New[int](WithCapacity[int](-5))

		// If we get here, the test should fail
		t.Error("Expected panic did not occur")
	})

	t.Run("multiple options", func(t *testing.T) {
		// Last option should win
		s := New[int](
			WithCapacity[int](5),
			WithCapacity[int](3),
		)

		// Should only be able to push 3 items
		for i := 0; i < 3; i++ {
			err := s.Push(i)
			if err != nil {
				t.Errorf("Push(%d) error = %v, want nil", i, err)
			}
		}

		err := s.Push(3)
		if !errors.Is(err, ErrOverflow) {
			t.Errorf("Push(3) exceeding final capacity error = %v, want ErrOverflow", err)
		}
	})
}

// Benchmark tests
func BenchmarkPush(b *testing.B) {
	s := New[int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	s := New[int]()
	// Pre-populate stack
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Pop()
	}
}

func BenchmarkPeek(b *testing.B) {
	s := New[int]()
	s.Push(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Peek()
	}
}

func BenchmarkSize(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Size()
	}
}

func TestRaceConditions(t *testing.T) {
	t.Run("concurrent push/pop/peek/size", func(t *testing.T) {
		s := New[int]()
		const numGoroutines = 50
		const numOperations = 200

		var wg sync.WaitGroup

		// Concurrent pushes
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(start int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					s.Push(start*numOperations + j)
				}
			}(i)
		}

		// Concurrent pops
		for i := 0; i < numGoroutines/2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations/2; j++ {
					s.Pop() // Ignore errors for this test
				}
			}()
		}

		// Concurrent peeks
		for i := 0; i < numGoroutines/4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					s.Peek() // Ignore errors
				}
			}()
		}

		// Concurrent size checks
		for i := 0; i < numGoroutines/4; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					s.Size()
				}
			}()
		}

		wg.Wait()

		// Just verify no panics occurred and stack is in valid state
		size := s.Size()
		if size < 0 {
			t.Errorf("Invalid size after concurrent operations: %d", size)
		}
	})

	t.Run("rapid push/pop cycles", func(t *testing.T) {
		s := New[int]()
		const numGoroutines = 100
		const cycles = 1000

		var wg sync.WaitGroup

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < cycles; j++ {
					// Push then immediately try to pop
					s.Push(id*cycles + j)
					s.Pop() // May succeed or fail, both are valid
				}
			}(i)
		}

		wg.Wait()

		// Stack should be in a valid state
		size := s.Size()
		if size < 0 {
			t.Errorf("Invalid size after rapid cycles: %d", size)
		}
	})

	t.Run("capacity-limited concurrent access", func(t *testing.T) {
		s := New[int](WithCapacity[int](10))
		const numGoroutines = 50
		const attempts = 100

		var wg sync.WaitGroup
		var pushCount, popCount int64

		// Many goroutines trying to push
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < attempts; j++ {
					err := s.Push(id*attempts + j)
					if err == nil {
						atomic.AddInt64(&pushCount, 1)
					}
				}
			}(i)
		}

		// Some goroutines trying to pop
		for i := 0; i < numGoroutines/2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < attempts/2; j++ {
					_, err := s.Pop()
					if err == nil {
						atomic.AddInt64(&popCount, 1)
					}
				}
			}()
		}

		wg.Wait()

		// Verify consistency
		finalSize := int64(s.Size())
		expectedSize := pushCount - popCount

		if finalSize != expectedSize {
			t.Errorf("Size inconsistency: got %d, expected %d (pushed: %d, popped: %d)",
				finalSize, expectedSize, pushCount, popCount)
		}

		// Size should never exceed capacity
		if finalSize > 10 {
			t.Errorf("Size exceeded capacity: %d > 10", finalSize)
		}
	})
}

func TestStressTest(t *testing.T) {

	s := New[int]()
	const duration = 1000 * time.Millisecond
	const numWorkers = 20

	var wg sync.WaitGroup
	start := time.Now()

	// Continuous operations for the duration
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			counter := 0
			for time.Since(start) < duration {
				switch counter % 4 {
				case 0:
					s.Push(workerID*10000 + counter)
				case 1:
					s.Pop()
				case 2:
					s.Peek()
				case 3:
					s.Size()
				}
				counter++
			}
		}(i)
	}

	wg.Wait()

	// Just verify the stack is in a valid state
	size := s.Size()
	if size < 0 {
		t.Errorf("Invalid final size: %d", size)
	}
}

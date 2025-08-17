# Go Stack

A thread-safe, generic LIFO stack implementation for Go with configurable capacity limits.

[![Go Reference](https://pkg.go.dev/badge/github.com/mghyo/go-stack.svg)](https://pkg.go.dev/github.com/mghyo/go-stack)
[![Go Report Card](https://goreportcard.com/badge/github.com/mghyo/go-stack)](https://goreportcard.com/report/github.com/mghyo/go-stack)

## Features

- **Generic**: Works with any type using Go's type parameters
- **Thread-safe**: Safe for concurrent use across multiple goroutines
- **LIFO semantics**: Last-In-First-Out behavior
- **Configurable capacity**: Set maximum size or use unlimited capacity
- **Zero dependencies**: Uses only Go standard library

## Installation

```bash
go get github.com/mghyo/go-stack
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/mghyo/go-stack"
)

func main() {
    // Create a new stack
    s := stack.New[int]()
    
    // Add items (LIFO)
    s.Push(1)
    s.Push(2)
    s.Push(3)
    
    // Remove items in reverse order
    for s.Size() > 0 {
        val, _ := s.Pop()
        fmt.Println(val) // Prints: 3, 2, 1
    }
}
```

## Usage

### Basic Operations

```go
s := stack.New[string]()

// Add items to top of stack
s.Push("first")
s.Push("second")

// Check top item without removing
top, err := s.Peek() // Returns "second"

// Remove items from top
val, err := s.Pop() // Returns "second"
val, err = s.Pop()  // Returns "first"

// Check size
fmt.Println(s.Size()) // 0
```

### Capacity-Limited Stack

```go
// Create stack with max capacity of 3
s := stack.New[int](stack.WithCapacity[int](3))

s.Push(1) // OK
s.Push(2) // OK  
s.Push(3) // OK
err := s.Push(4) // Returns stack.ErrOverflow
```

### Error Handling

```go
s := stack.New[int]()

// Empty stack operations return ErrUnderflow
val, err := s.Pop()
if errors.Is(err, stack.ErrUnderflow) {
    fmt.Println("Stack is empty")
}

val, err = s.Peek()
if errors.Is(err, stack.ErrUnderflow) {
    fmt.Println("Stack is empty")
}
```

## API Reference

### Types

```go
type Stack[T any] interface {
    Push(val T) error      // Add item to top
    Pop() (T, error)       // Remove item from top  
    Size() int             // Current number of items
    Peek() (T, error)      // View top item without removing
}
```

### Functions

```go
// Create new stack
func New[T any](opts ...Option[T]) Stack[T]

// Set maximum capacity (-1 for unlimited)
func WithCapacity[T any](cap int) Option[T]
```

### Constants & Errors

```go
const UnlimitedCapacity = -1

var ErrOverflow = errors.New("stack overflow")   // Stack is full
var ErrUnderflow = errors.New("stack underflow") // Stack is empty
```

## Performance

- **Push**: O(1) amortized
- **Pop**: O(1)
- **Peek**: O(1)
- **Size**: O(1)

All operations are highly efficient with optimal time complexity.

## Thread Safety

All operations are thread-safe and can be used concurrently:

```go
s := stack.New[int]()
var wg sync.WaitGroup

// Multiple pushers
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(val int) {
        defer wg.Done()
        s.Push(val)
    }(i)
}

// Multiple poppers  
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        s.Pop()
    }()
}

wg.Wait()
```

## Testing

```bash
go test              # Run tests
go test -race        # Run with race detection
go test -bench=.     # Run benchmarks
go test -cover       # Check coverage
```

## Requirements

- Go 1.18+ (for generics support)

## License

MIT License - see [LICENSE](LICENSE) file for details.
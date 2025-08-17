package stack

import "errors"

var (
	ErrOverflow  = errors.New("stack overflow")
	ErrUnderflow = errors.New("stack underflow")
)

package errors

import (
	"errors"
)

// TODO: Support validation with multiple errors.
type Errors []error

var (
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrEmptyKey        = errors.New("empty key")
)

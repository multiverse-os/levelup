package leveldb

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////
var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrEmptyKey        = errors.New("key cannot be empty")
	ErrIndexOutOfRange = errors.New("iterator index out of range")
	ErrClosed          = errors.New("closed")
	ErrCorrupt         = errors.New("corrupt")
)

package backend

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrEmptyKey         = errors.New("empty key")
	ErrReadOnly         = errors.New("read-only mode")
	ErrSnapshotReleased = errors.New("snapshot released")
	ErrIterReleased     = errors.New("iterator released")
	ErrClosed           = errors.New("closed")
)

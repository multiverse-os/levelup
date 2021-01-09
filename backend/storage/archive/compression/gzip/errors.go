package gzip

import (
	"errors"
)

var (
	ErrChecksum = errors.New("gzip: invalid checksum")
	ErrHeader   = errors.New("gzip: invalid header")
)

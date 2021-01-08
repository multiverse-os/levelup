package model

import (
	"time"

	codec "github.com/multiverse-os/codec"
)

type Record struct {
	CollectionId uint32

	Id []byte

	Codec codec.Codec

	CreatedAt time.Time
	UpdatedAt time.Time

	Index int

	Key  []byte
	Data []byte

	Hooks *Hooks

	VersionHistory []*Record
}

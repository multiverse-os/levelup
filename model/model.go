package model

import (
	"time"

	codec "github.com/multiverse-os/levelup/data/codec"
	id "github.com/multiverse-os/levelup/id"
)

type Model interface {
	// Hooks

	// Validations

	// Updating timestamps

	// encoding and compression
}

type Record struct {
	Id id.Id

	Codec codec.Codec

	CreatedAt time.Time
	UpdatedAt time.Time

	Index int

	Key  []byte
	Data []byte

	Hooks *Hooks

	VersionHistory []*Record
}

func New() Model {
	return Record{
		Id:        id.New(),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

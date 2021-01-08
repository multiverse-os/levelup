package model

import (
	"time"

	"github.com/multiverse-os/levelup/id"
)

type Type int

const (
	Standard Type = iota
	Temporary
	Immutable
	Verionable
)

type Action func() error

type Model struct {
	Type Type
	Timestamps
	Hooks
}

////////////////////////////////////////////////////////////////////////////////
type Document struct {
	Model

	Id           id.Id
	CollectionId uint32

	Key  []byte
	Data []byte

	Value interface{}
}

type Record struct {
	Model

	Id           id.Id
	CollectionId uint32

	Key  []byte
	Data []byte
}

////////////////////////////////////////////////////////////////////////////////
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

////////////////////////////////////////////////////////////////////////////////
type Version struct {
	ChangedAt time.Time
	Versions  []*Version
}

////////////////////////////////////////////////////////////////////////////////
type Cache struct {
	TTL time.Duration
}

////////////////////////////////////////////////////////////////////////////////
type Hooks struct {
	BeforeGet []*Action
	AfterGet  []*Action

	BeforeSet []*Action
	AfterSet  []*Action

	BeforeDelete []*Action
	AfterDelete  []*Action
}

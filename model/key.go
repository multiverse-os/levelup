package model

import (
	"fmt"

	id "github.com/multiverse-os/levelup/id"

	util "github.com/syndtr/goleveldb/leveldb/util"
)

type RecordType int

const (
	Standard RecordType = iota
	Immutable
	Temporary
	Versioning
	All
)

func (self RecordType) String() string {
	switch self {
	case Standard:
		return "s"
	case Immutable:
		return "i"
	case Temporary:
		return "t"
	case Versioning:
		return "v"
	default: // All
		return ""
	}
}

type key struct {
	Collection *Collection

	Type RecordType

	Prefix []byte
	Name   string

	Parent    *key
	ChildKeys []*key
}

func Key(keyString string) key {
	if EmptyString(keyString) {
		return key{
			Type:   All,
			Prefix: []byte(""),
		}
	} else {
		return key{
			Type:   Standard,
			Name:   keyString,
			Prefix: []byte(id.New().Short()),
		}
	}
}

func (self key) Bytes() []byte {
	if self.Type == All {
		return []byte("")
	} else {
		return []byte(fmt.Sprintf("%s:%s:%s", self.Type.String(), self.Prefix, self.Name))
	}
}

func (self key) String() string {
	if self.Type == All {
		return ""
	} else {
		return fmt.Sprintf("%s:%s:%s", self.Type.String(), self.Prefix, self.Name)
	}
}

func (self key) Range() *util.Range {
	return util.BytesPrefix(self.Bytes())
}

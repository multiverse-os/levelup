package levelup

import (
	"fmt"
	"regexp"

	id "github.com/multiverse-os/levelup/id"
)

////////////////////////////////////////////////////////////////////////////////
type Key struct {
	Database   *Database
	Collection *Collection

	Type KeyType

	Subset id.Id
	Record id.Id
}

func IsValidKey(k string) (bool, error) {
	return regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{1,255}$`, k)
}

////////////////////////////////////////////////////////////////////////////////
type KeyType int

const (
	bytesKey KeyType = iota
	specialKey
	indexKey
	documentKey
	immutableKey
	cacheKey
	allKey
)

func (self KeyType) String() string {
	switch self {
	case documentKey:
		return "d"
	case indexKey:
		return "x"
	case immutableKey:
		return "i"
	case cacheKey:
		return "c"
	case specialKey:
		return "s"
	case bytesKey:
		return "b"
	default: // Standard
		return ""
	}
}

func (self KeyType) Bytes() []byte { return []byte(self.String()) }

////////////////////////////////////////////////////////////////////////////////
//func Key(t KeyType) Key {
//	return Key{
//		Type: standardKey,
//	}
//}
////////////////////////////////////////////////////////////////////////////////

func Document(subsetName string) Key {
	return Key{
		Type:   documentKey,
		Subset: id.Hash(subsetName),
	}
}

func ValueBucket(subsetName string) Key {
	return Key{
		Type:   bytesKey,
		Subset: id.Hash(subsetName),
	}
}

func Special(subsetName string) Key {
	return Key{
		Type:   specialKey,
		Subset: id.Hash(subsetName),
	}
}

func Index(subsetName string) Key {
	return Key{
		Type:   indexKey,
		Subset: id.Hash(subsetName),
	}
}

func Cache(subsetName string) Key {
	return Key{
		Type:   cacheKey,
		Subset: id.Hash(subsetName),
	}
}

func Immutable(subsetName string) Key {
	return Key{
		Type:   immutableKey,
		Subset: id.Hash(subsetName),
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self Key) Key(recordName string) Key {
	self.Record = id.Hash(recordName)
	return self
}

////////////////////////////////////////////////////////////////////////////////
func (self Key) Bytes() []byte {
	switch self.Type {
	case documentKey:
		return id.Hash(fmt.Sprintf("%s.%s.%s", self.Type.Bytes(), self.Subset.Bytes(), self.Record.Bytes())).Bytes()
	case bytesKey:
		return id.Hash(fmt.Sprintf("%s.%s.%s", self.Type.Bytes(), self.Subset.Bytes(), self.Record.Bytes())).Bytes()
	case immutableKey:
		return []byte(fmt.Sprintf("%s!%s.%s", self.Type.Bytes(), self.Subset.Bytes(), self.Record.Bytes()))
	case cacheKey:
		return []byte(fmt.Sprintf("%s?%s.%s", self.Type.Bytes(), self.Subset.Bytes(), self.Record.Bytes()))
	case specialKey:
		return []byte(fmt.Sprintf("%s$%s.%s", self.Type.Bytes(), self.Subset.Bytes(), self.Record.Bytes()))
	default: // allKey
		return []byte("")
	}
}

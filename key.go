package levelup

import (
	"fmt"

	id "github.com/multiverse-os/levelup/id"
)

////////////////////////////////////////////////////////////////////////////////
type key struct {
	Type KeyType

	Subset id.Id
	Record id.Id
}

////////////////////////////////////////////////////////////////////////////////
type KeyType int

const (
	bytesKey KeyType = iota
	specialKey
	documentKey
	immutableKey
	cacheKey
	allKey
)

func (self KeyType) String() string {
	switch self {
	case documentKey:
		return "d"
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
//func Key(t KeyType) key {
//	return key{
//		Type: standardKey,
//	}
//}
////////////////////////////////////////////////////////////////////////////////

func Document(subsetName string) key {
	return key{
		Type:   documentKey,
		Subset: id.Hash(subsetName),
	}
}

func ValueBucket(subsetName string) key {
	return key{
		Type:   bytesKey,
		Subset: id.Hash(subsetName),
	}
}

func Special(subsetName string) key {
	return key{
		Type:   specialKey,
		Subset: id.Hash(subsetName),
	}
}

func Cache(subsetName string) key {
	return key{
		Type:   cacheKey,
		Subset: id.Hash(subsetName),
	}
}

func Immutable(subsetName string) key {
	return key{
		Type:   immutableKey,
		Subset: id.Hash(subsetName),
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self key) Key(recordName string) key {
	self.Record = id.Hash(recordName)
	return self
}

////////////////////////////////////////////////////////////////////////////////
func (self key) Bytes() []byte {
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

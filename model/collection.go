package levelup

import (
	"fmt"
	"reflect"
	"time"
)

type CollectionType int

const (
	Persisting CollectionType = iota
	Versioning
	Temporary
	Memory
)

type Collection struct {
	Database *Database

	Type         CollectionType
	Name         string
	NameChecksum []byte

	Records      map[int64]*Record
	Size         int
	LastUpdateAt time.Time

	//Relationships []*collection.Relationships
}

func (self *Database) NewCollection(name string) *Collection {
	collection := &Collection{
		Database:     self,
		Name:         name,
		NameChecksum: []byte(self.Codec.Checksum([]byte(name))),
	}
	self.Collections = append(self.Collections, collection)
	return collection
}

func (self *Database) Collection(name string) *Collection {
	for _, collection := range self.Collections {
		if reflect.DeepEqual(collection.NameChecksum, self.Codec.Checksum([]byte(name))) {
			fmt.Println("found collection!:", name)
			return collection
		}
	}
	return nil
}

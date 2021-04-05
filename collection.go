package levelup

import (
	"time"

	id "github.com/multiverse-os/levelup/id"
)

type Collection struct {
	Database Database

	Id id.Id

	Type KeyType
	Name string

	Records      map[uint32]*Record
	Size         int
	LastUpdateAt time.Time

	//Relationships []*collection.Relationships
}

func (self Collection) Key(recordKey string) Key { return Document(self.Name).Key(recordKey) }

////////////////////////////////////////////////////////////////////////////////
// TODO: Direct get/put/delete access on collection

func (self *Collection) Put(key string, value []byte) error {
	return self.Database.Put(self.Key(key), value)
}

func (self *Collection) Get(key string) ([]byte, error) {
	return self.Database.Get(self.Key(key))
}

func (self *Collection) Delete(key string) error {
	return self.Database.Delete(self.Key(key))
}

package levelup

// TODO[IMPORTANT]: Basically make everything private, then use the database,
// and slowly only make the things public that we are using.
// TODO: Do benchmarking and if we receive difference using bool and uint8 on
// our enumerators.

import (
	"sync"

	"time"

	codec "github.com/multiverse-os/codec"
	checksum "github.com/multiverse-os/codec/checksum"
	encoding "github.com/multiverse-os/codec/encoding"
	backend "github.com/multiverse-os/levelup/backend"
	id "github.com/multiverse-os/levelup/id"
	threadsafe "github.com/multiverse-os/levelup/threadsafe"
)

// TODO: Combine this into our events object
type Storage interface {
	Get(key Key)

	// TODO: Here we will setup our ability to define Hooks:
	//        (BeforeGet, AfterGet)
	//        (BeforeSet, AfterSet)
	//        (BeforeDelete, AfterDelete)

	// TODO: This will also allow us to maintain multiple databases, we can load
	// a single database specifically for writing. Then take have a saprately
	// stored location for a immutable database.
	OnReject() func(record *Record)
}

type Database struct {
	sync.RWMutex
	Throttle time.Duration // Disk write throttling
	Active   threadsafe.Boolean
	Events   Events

	// TODO: Now implement a buffered channel to control and throttle writes

	storage backend.Database
	Codec   codec.Codec

	Collections map[uint32]*Collection
	Records     map[uint32]*Record
}

// TODO: A database Loader to get everything read from the logs into cache and
// get ready for writing is separate from the Writer function under events

////////////////////////////////////////////////////////////////////////////////
func (self Database) NewCollection(name string) *Collection {
	collection := &Collection{
		Id:       id.Hash(name),
		Active:   threadsafe.Switch().Off(),
		Database: self,
		Name:     name,
		Records:  make(map[uint32]*Record),
	}
	self.Collections[collection.Id.UInt32()] = collection
	return collection
}

func (self *Database) Collection(name string) *Collection {
	return self.Collections[id.Hash(name).UInt32()]
}

////////////////////////////////////////////////////////////////////////////////
func Open(path string) (*Database, error) {
	if db, err := backend.FileStorage(path); err != nil {
		return nil, err
	} else {

		database := &Database{
			storage:     db,
			Codec:       codec.EncodingFormat(encoding.BSON).ChecksumAlgorithm(checksum.XXH64),
			Collections: make(map[uint32]*Collection),
			Records:     make(map[uint32]*Record),
			Events:      []*Event{},
		}

		// TODO: HERE IS WHERE WE WOULD NORMALLY DO OUR
		//   GO FUNC BASED UPDATE LOOP ON A BUFFERED CHANNEL
		return database, nil
	}
}

func (self *Database) Close() error {
	err := self.storage.Close()
	self.storage.Close()
	self = nil
	return err
}

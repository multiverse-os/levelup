package levelup

import (
	"sync"

	backend "github.com/multiverse-os/levelup/backend"
	history "github.com/multiverse-os/levelup/history"
	id "github.com/multiverse-os/levelup/id"
	model "github.com/multiverse-os/levelup/model"

	codec "github.com/multiverse-os/codec"
	checksum "github.com/multiverse-os/codec/checksum"
	encoding "github.com/multiverse-os/codec/encoding"
)

// TODO: SPECIAL COLLECTIONS BASED ON SPECIALIZED PREFIXES

//       USE THESE TO LOAD the DATABASE FROM FILES.
//       _SETTINGS_ will have database configuration, a singlton
//       _LOG_ will store logs allowing stepping fowards and backwards

// NOTE: We will only interact with the **Storage** object to preform writes,
//       and initialization. Every other READ should occur from the cache
//       provided by the maps which IDs will be xxhash of the name.

//       Ephemeral cache databases will be loaded in this way. And we will use
//       locking and buffered channels to handle writes and emits to handle
//       updating the cache. This in addition to transactions should provide us
//       with a solid thread-safe atomic system that supports heavy writes and
//       heavy reads. This is just a first draft, next draft will likely use a
//       different KVdb we have been working on.
type Database struct {
	storage backend.Storage

	Access sync.RWMutex
	Codec  codec.Codec

	Collections map[uint32]*Collection
	Records     map[int64]*model.Record
	Logs        []*history.Log
}

////////////////////////////////////////////////////////////////////////////////
func (self *Database) NewCollection(name string) *Collection {
	collection := &Collection{
		Id:       id.Hash(name),
		Database: self,
		Name:     name,
	}
	self.Collections[collection.Id.UInt32()] = collection
	return collection
}

////////////////////////////////////////////////////////////////////////////////
func Open(path string) (*Database, error) {
	if db, err := backend.DatabaseFile(path); err != nil {
		return nil, err
	} else {
		return &Database{
			// Private
			// NOTE: We want this to be private so we can control access at a
			// bottleneck letting us base the GET/PUT/DELETE functionality on the key
			// used in the request.
			storage: db,
			// Public
			Codec:       codec.EncodingFormat(encoding.BSON).ChecksumAlgorithm(checksum.XXH64),
			Collections: make(map[uint32]*Collection),
			Records:     make(map[int64]*model.Record),
			Logs:        []*history.Log{},
		}, nil

		// TODO: Scan existing database, and load it into READ CACHE.
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self *Database) Close() {
	self.storage.Close()
}

////////////////////////////////////////////////////////////////////////////////
func (self *Database) Collection(name string) *Collection {
	return self.Collections[id.Hash(name).UInt32()]
}

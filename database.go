package levelup

import (
	"sync"

	leveldb "github.com/multiverse-os/levelup/backend"

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
	Storage *leveldb.Storage
	Access  sync.RWMutex
	Codec   codec.Codec

	Collections     map[int64]*model.Collection
	CollectionCount int

	Records     map[int64]*model.Record
	RecordCount int

	Logs     []*Log
	LogCount int
}

////////////////////////////////////////////////////////////////////////////////
func Open(path string) (*Database, error) {
	if db, err := leveldb.Open(leveldb.FileStore, path); err != nil {
		return nil, err
	} else {
		return &Database{
			Storage: db,
			Access:  make(sync.RWMutex),
			Codec:   codec.EncodingFormat(encoding.BSON).ChecksumAlgorithm(checksum.XXH64),

			Collections:     make(map[int64]*model.Collection),
			CollectionCount: 0,

			Records:     make(map[int64]*model.Record),
			RecordCount: 0,

			Logs:     make([]*Log),
			LogCount: 0,
		}, nil

		// TODO: Scan existing database, and load it into READ CACHE.
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self *Database) Close() {
	self.Storage.Close()
}

////////////////////////////////////////////////////////////////////////////////

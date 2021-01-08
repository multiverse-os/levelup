package levelup

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	backend "github.com/multiverse-os/levelup/backend"
	history "github.com/multiverse-os/levelup/history"
	id "github.com/multiverse-os/levelup/id"
	model "github.com/multiverse-os/levelup/model"

	codec "github.com/multiverse-os/codec"
	checksum "github.com/multiverse-os/codec/checksum"
	encoding "github.com/multiverse-os/codec/encoding"
)

type Database struct {
	storage backend.Storage

	Codec codec.Codec

	throttle time.Duration // Disk throttling to prevent a heavy upgrade from hogging resources
	active   uint32        // Flag whether the event loop was started

	update    chan struct{}   // Notification channel that headers should be processed
	quit      chan chan error // Quit channel to tear down running goroutines
	ctx       context.Context
	ctxCancel func()

	// TODO: This mutex is primarily used, hwen updating the caches, or grabbing
	// reads. This should allow our database to avoid race conditions. But this
	// will obviously need heavy heavy testing.
	Access      sync.RWMutex
	Collections map[uint32]*Collection
	Records     map[int64]*model.Record
	// NOTE: These are intended to be more than just logs, but rather a history of
	// every WRITE ACTION (_CREATE, _UPDATE, _DELETE), so that can specifically
	// rewind the database.
	ActionLogs []*history.ActionLog
}

// TODO: Add the ability to rewind the database using action logs, to step
// backwards.

// TODO: Use a buffered channel to submit all writes, so that that only a single
// write operation is capable of happening at any time. USE TRANSACTIONS OR ADD
// TRANSACTION SUPPORT.
// READ operations should NEVER hit the actual database, we load the the
// database from the file into memory, and reads only interact with this
// memory database. Changes to the cache where reads come from, happen after
// a buffered write comes down the channel, succeeds, then causes the update to
// the cache.
// This means that the only interactions with the database are write actions--
// controlled by a buffered channel and a mutex.
// This should allow for A LOT of write actions to occur, since we are not
// doing any reads from the disk io, we are doing all the reads from the
// memory IO.

// After this is successful, we want to abstract ontop of the cached memory
// database a graph database built using buckets, or maps, of relationships.
// These relationships are stored separatedly in a differnet database file.
// These relationships, do not change any data in the objects, they just store
// data baout how obgjects are related.

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

func (self *Database) Collection(name string) *Collection {
	return self.Collections[id.Hash(name).UInt32()]
}

////////////////////////////////////////////////////////////////////////////////
func Open(path string) (*Database, error) {
	if db, err := backend.DatabaseFile(path); err != nil {
		return nil, err
	} else {

		database := &Database{
			// Private
			// NOTE: We want this to be private so we can control access at a
			// bottleneck letting us base the GET/PUT/DELETE functionality on the key
			// used in the request.
			storage: db,

			update: make(chan struct{}, 1),
			quit:   make(chan chan error),

			// Public
			Codec:       codec.EncodingFormat(encoding.BSON).ChecksumAlgorithm(checksum.XXH64),
			Collections: make(map[uint32]*Collection),
			Records:     make(map[int64]*model.Record),
			ActionLogs:  []*history.ActionLog{},
		}
		database.ctx, database.ctxCancel = context.WithCancel(context.Background())

		//go database.UpdateLoop()
		return database, nil

		// TODO: Scan existing database, and load it into READ CACHE.
	}
}

func (self *Database) Close() error {
	var errs []error

	self.ctxCancel()

	//// Tear down the primary update loop
	errc := make(chan error)
	self.quit <- errc
	if err := <-errc; err != nil {
		errs = append(errs, err)
	}
	//// If needed, tear down the secondary event loop
	if atomic.LoadUint32(&self.active) != 0 {
		self.quit <- errc
		if err := <-errc; err != nil {
			errs = append(errs, err)
		}
	}
	//// Return any failures
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return fmt.Errorf("%v", errs)
	}

	err := self.storage.Close()
	self.storage = nil
	self.storage.Close()
	return err
}

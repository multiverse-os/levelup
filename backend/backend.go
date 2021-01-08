package backend

import (
	leveldb "github.com/multiverse-os/levelup/backend/leveldb"

	iterator "github.com/syndtr/goleveldb/leveldb/iterator"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
	util "github.com/syndtr/goleveldb/leveldb/util"
)

type StorageType int

const (
	File StorageType = iota
	Memory
)

////////////////////////////////////////////////////////////////////////////////
type Backend interface {
	// Information & Statistics
	Name() string
	Storage() string

	//Collections() []*collection.Collection

	// Cleanup
	Close() error

	// KeyValue
	Set(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error

	//All() []*record.Record
	//Paginate(page, perPage int) []*record.Record
	//Collection(name string) []*record.Record

	//Find(key []byte) []*record.Record
	//Where(field string, comparison Comparison, value interface{}) []*record.Record
	//

	// Iterators
	//Find(key []byte) []byte

	// LEVELDB Specific
	NewIterator(iteratorRange *util.Range, options *opt.ReadOptions) iterator.Iterator
}

func Open(storageType StorageType, path string) (Backend, error) {
	switch storageType {
	case File:
		return leveldb.FileStorage(path)
	case Memory:
		return leveldb.MemoryStorage()
	default:
		panic("invalid database")
	}
}

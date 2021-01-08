package backend

import (
	leveldb "github.com/multiverse-os/levelup/backend/leveldb"

	iterator "github.com/syndtr/goleveldb/leveldb/iterator"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
	util "github.com/syndtr/goleveldb/leveldb/util"
)

////////////////////////////////////////////////////////////////////////////////
type Storage interface {
	// Information & Statistics
	Name() string

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

func DatabaseFile(path string) (Storage, error) { return leveldb.FileStorage(path) }

//func Memory() (Storage, error)                  { return Storage{}, nil }

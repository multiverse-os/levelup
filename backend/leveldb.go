package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	leveldb "github.com/syndtr/goleveldb/leveldb"
	errors "github.com/syndtr/goleveldb/leveldb/errors"
	filter "github.com/syndtr/goleveldb/leveldb/filter"
	iterator "github.com/syndtr/goleveldb/leveldb/iterator"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
	util "github.com/syndtr/goleveldb/leveldb/util"
)

////////////////////////////////////////////////////////////////////////////////
type Type int

const (
	File Type = iota
	Archive
	Memory
)

type Storage struct {
	Type  Type
	Path  string
	Files []string
}

////////////////////////////////////////////////////////////////////////////////
type Database struct {
	sync.Mutex
	Type Type

	Store *leveldb.DB
	Batch *leveldb.Batch

	PersistentStorage Storage

	Options Options
}

////////////////////////////////////////////////////////////////////////////////
type Options struct {
	ReadOnly  bool
	WriteOnly bool

	Runtime *opt.Options
	Read    *opt.ReadOptions
	Write   *opt.WriteOptions
}

////////////////////////////////////////////////////////////////////////////////
// NOTE: In our implementation we are going to keep all the database files in
//       in a *.tar archive. Then by doing modifications to the tar file, we
//       will make our changes.
//       This new storage technique will include diff patches, regular snapshots
//       after x amount of data has been changed, and several iterations of
//       the database within the tar file. We could also use this to have a
//       a read-only version and a write version in the same tar.
//       Writes could be creating the diff-files, then the diff patches can be
//       applied after x amount have been built up.
//
//       **[Key concepts]**
//
//         1) Guarantee atomicity
//
//         2) Archive to keep all files, replicated versions (for speed, think
//            RAID), iterations, diffs, in a single file.
//
//         3) Add another level of compression possibly at rest.
//
//         4) At rest encryption with root-key, and re-encrypt. By having
//            replications, it may be possible to decrypt parts of the db,
//            either way it could be encrypted to the session key.
//
//         5) Store all public keys, information for backing up, git data, etc.
//
//
////////////////////////////////////////////////////////////////////////////////

//type FileStorage struct {
//	Path string
//
//	// TODO: tar.zstd !
//	// TODO: Store cryptokeysm identify owners
//
//	// TODO: Digest file=> signed could have things like who can write,
//	//                     what is the root key, recovery key,
//	//                     active session key (when it expires)
//
//	//
//
//	Archive *tar.Writer
//	Files   []string
//}

//func (self FileStorage) ParseFiles() {
//
//}
//
//func (self FileStorage) CreateArchive() ([]byte, error) {
//
//}
//
//func (self FileStorage) AddFile(path string) error {
//	file, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//	if stat, err := file.Stat(); err == nil {
//		// now lets create the header as needed for this file within the tarball
//		header := new(tar.Header)
//		header.Name = path
//		header.Size = stat.Size()
//		header.Mode = int64(stat.Mode())
//		header.ModTime = stat.ModTime()
//		// write the header to the tarball archive
//		if err := tw.WriteHeader(header); err != nil {
//			return err
//		}
//		// copy the file data to the tarball
//		if _, err := io.Copy(tw, file); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (self FileStorage) DeleteFiles() error {
//	return os.RemoveAll(self.Path)
//}

////////////////////////////////////////////////////////////////////////////////
func (Database) Name() string {
	return "leveldb"
}

////////////////////////////////////////////////////////////////////////////////
func MemoryStorage() (database Database, err error) {
	return Database{}, fmt.Errorf("memory databse not yet implemented")
}

func FileStorage(path string) (database Database, err error) {
	database = Database{
		Options: Options{
			Runtime: &opt.Options{
				CompactionTableSize:           (4 * opt.MiB),
				CompactionTotalSize:           (128 * opt.MiB),
				WriteBuffer:                   (8 * opt.MiB),
				Compression:                   opt.SnappyCompression,
				ReadOnly:                      false,
				Filter:                        filter.NewBloomFilter(10),
				CompactionTotalSizeMultiplier: 20,
				BlockCacheCapacity:            0,
				WriteL0SlowdownTrigger:        16,
				WriteL0PauseTrigger:           64,
				// NOTE: This affects how we build code in this file.
				ErrorIfMissing: false,
				ErrorIfExist:   false,
				// TODO: Can use this to specify a virtual filesystem for greater
				//       security.
				// Comparator: nil,
				// Logger: nil,
				// FileSystem: opts.getFileSystem(),
			},
			Read: &opt.ReadOptions{
				DontFillCache: true,
			},
			Write: &opt.WriteOptions{
				Sync: true,
			},
		},
	}

	_ = os.MkdirAll(filepath.Clean(path), 0750)

	//////////////////////////////////////////////////////////////////////////////
	// Repairs Database                                                         //
	// rdb, err := leveldb.RecoverFile(db.path, db.opts);                       //
	//////////////////////////////////////////////////////////////////////////////
	database.Store, err = leveldb.OpenFile(path, database.Options.Runtime)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		database.Store, err = leveldb.RecoverFile(path, nil)
	}

	// TODO: Read up more on finalizers
	// runtime.SetFinalizer(database, (Database).finalize)
	return database, err
}

//func (self Database) finalize() {
//	go self.Store.Close()
//}

// SetReadOnly makes DB read-only. It will stay read-only until reopened.
func (self Database) SetReadOnly() error {
	return self.Store.SetReadOnly()
}

////////////////////////////////////////////////////////////////////////////////
//type EncryptedDB struct {
//	*leveldb.DB
//	scloser io.Closer
//}
//
//func (e *EncryptedDB) Close() {
//	e.DB.Close()
//	e.scloser.Close()
//}
//
//func OpenAESEncryptedFile(path string, key []byte, opt *opt.Options) (db *EncryptedDB, err error) {
//	stor, err := aesgcm.OpenEncryptedFile(path, key, opt.GetReadOnly())
//	if err != nil {
//		return
//	}
//	ldb, err := leveldb.Open(stor, opt)
//	if err != nil {
//		stor.Close()
//	} else {
//		db = &EncryptedDB{
//			DB:      ldb,
//			scloser: stor,
//		}
//	}
//	return
//}

////////////////////////////////////////////////////////////////////////////////
func (self Database) Close() error {
	return self.Store.Close()
}

////////////////////////////////////////////////////////////////////////////////
func (self Database) CompactDatastore() error {
	return self.Store.CompactRange(util.Range{Limit: nil, Start: nil})
}

////////////////////////////////////////////////////////////////////////////////
func (self Database) Has(key []byte) bool {
	data, err := self.Store.Get(key, self.Options.Read)
	if err != nil {
		return false
	} else {
		return len(data) != 0
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self Database) Get(key []byte) ([]byte, error) {
	return self.Store.Get(key, nil)
}

func (self Database) Set(key []byte, value []byte) error {
	if len(key) == 0 {
		return ErrEmptyKey
	} else if len(value) == 0 {
		return self.Store.Delete(key, nil)
	} else {
		return self.Store.Put(key, value, self.Options.Write)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self Database) Delete(key []byte) error {
	return self.Store.Delete(key, nil)
}

////////////////////////////////////////////////////////////////////////////////
func (self Database) writeBatch(batch *leveldb.Batch) error {
	if err := self.Store.Write(batch, self.Options.Write); err != nil {
		return err
	} else {
		return nil
	}
}

func (self Database) WriteBatch(data map[string][]byte) error {
	batch := &leveldb.Batch{}
	for key, value := range data {
		keyBytes := []byte(key)
		if value != nil {
			batch.Delete(keyBytes)
		} else {
			batch.Put(keyBytes, value)
		}
	}
	return self.writeBatch(batch)
}

////////////////////////////////////////////////////////////////////////////////

func (self Database) Snapshot() (map[string][]byte, error) {
	self.Lock()
	defer self.Unlock()

	snap, err := self.Store.GetSnapshot()
	if err != nil {
		return nil, fmt.Errorf("Error while taking snapshot:" + err.Error())
	}

	data := make(map[string][]byte)
	iter := snap.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		data[string(key)] = append([]byte{}, val...)
	}
	iter.Release()
	if iter.Error() != nil {
		return nil, iter.Error()
	}
	return data, nil
}

//func (self Database) Snapshot() *Snapshot {
//	ss := self.Store.NewSnapshot()
//	if ss == nil {
//		return &Snapshot{err: ErrDBClosed}
//	}
//	return newSnapshot(ss)
//}

////////////////////////////////////////////////////////////////////////////////
func (self Database) NewIterator(iteratorRange *util.Range, options *opt.ReadOptions) iterator.Iterator {
	return self.Store.NewIterator(iteratorRange, options)
}

func (self Database) All() (records [][]byte, err error) {
	iterator := self.Store.NewIterator(nil, nil)
	for iterator.Next() {
		records = append(records, iterator.Value())
	}
	iterator.Release()
	return records, iterator.Error()
}

func (self Database) AllKeys() (keys [][]byte, err error) {
	iterator := self.Store.NewIterator(nil, nil)
	for iterator.Next() {
		keys = append(keys, iterator.Key())
	}
	iterator.Release()
	return keys, iterator.Error()

}

func (self Database) Paginate(start, limit []byte) (records [][]byte, err error) {
	iterator := self.Store.NewIterator(&util.Range{Start: start, Limit: limit}, nil)
	for iterator.Next() {
		records = append(records, iterator.Value())
	}
	iterator.Release()
	return records, iterator.Error()
}

func (self Database) Collection(iteratorRange *util.Range) (records [][]byte, err error) {
	iterator := self.Store.NewIterator(iteratorRange, nil)
	for iterator.Next() {
		records = append(records, iterator.Value())
	}
	iterator.Release()
	return records, iterator.Error()
}

func (self Database) CollectionKeys(iteratorRange *util.Range) (keys [][]byte, err error) {
	iterator := self.Store.NewIterator(iteratorRange, nil)
	for iterator.Next() {
		keys = append(keys, iterator.Value())
	}
	iterator.Release()
	return keys, iterator.Error()
}

////////////////////////////////////////////////////////////////////////////////

// All returns an iterator catching all keys in db.
//func (self Database) All() iterator.Iterator {
//	return self.Store.All(self.Options.Read)
//}
//
//func (self Database) Find(start []byte) iterator.Iterator {
//	return self.Store.Find(start, self.Options.Read)
//}
//
//func (self Database) Range(start, limit []byte) iterator.Iterator {
//	return self.Store.Range(start, limit, self.Options.Read)
//}
//
//// Prefix returns an iterator catching all keys having prefix as prefix.
//func (self Database) Prefix(prefix []byte) iterator.Iterator {
//	return self.Store.Prefix(prefix, self.Options.Read)
//}
//
//func (self Database) CompactRange(start, limit []byte) error {
//	return self.Store.CompactRange(start, limit)
//}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////

// TODO: need to be able to pass in multiple types of action
//func (self *Database) BatchPut(kv ...[]byte) error {
//	batch := new(leveldb.Batch)
//
//	for _, keyValue := range kv {
//		batch.Put(keyValue[0], keyValue[1])
//	}
//
//	return self.KV.Write(batch, nil)
//}

// TODO: Batch should accept functions then run all those functioons then batch
// close.
//func (self *Database) BatchPut(func)
//batch := new(leveldb.Batch)
//batch.Put([]byte("foo"), []byte("value"))
//batch.Put([]byte("bar"), []byte("another value"))
//batch.Delete([]byte("baz"))
//err = db.Write(batch, nil)

////////////////////////////////////////////////////////////////////////////////

package leveldb

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	codec "github.com/multiverse-os/codec"
	errors "github.com/multiverse-os/levelup/errors"

	leveldb "github.com/syndtr/goleveldb/leveldb"
	filter "github.com/syndtr/goleveldb/leveldb/filter"
	iterator "github.com/syndtr/goleveldb/leveldb/iterator"
	opt "github.com/syndtr/goleveldb/leveldb/opt"
	util "github.com/syndtr/goleveldb/leveldb/util"
)

////////////////////////////////////////////////////////////////////////////////
type StorageType int

const (
	File StorageType = iota
	Memory
)

func (self StorageType) String() string {
	switch self {
	case File:
		return "file://"
	case Memory:
		return "memory://"
	default:
		return ""
	}
}

// NOTE: Used to sign updates, and encrypt the database when at rest. The
// database can use this in combination with a Model.Field to encrypt the
// entire model. This allows user data to be encrypted with the combination
// of a system key and a user key when at rest. The system key ideally will
// remain cold, and never be on the production server. It simply will send a
// session certificate over to the server which is valid for the process
// lifetime. This is done by decrypting the data and renecrypting it to the
// session key, and when closing, decryping from the session key and
// re-encrypting back to the root key. This allows the data to be stored at
// rest with a key that is never present on the production server.

//func GenerateSessionKey(rootKey []byte) []byte {
//	// TODO: Use codec.Crypto argon2
//	return argon2.Key([]byte(rootKey), salt, 3, 32*1024, 4, 32)
//}

type Options struct {
	Storage   StorageType
	Path      string
	ReadOnly  bool
	WriteOnly bool

	Runtime *opt.Options
	Read    *opt.ReadOptions
	Write   *opt.WriteOptions
}

type Database struct {
	Store *leveldb.DB

	Options Options

	Codec codec.Codec

	//RootKey    []byte
	//SessionKey []byte

	Access   sync.Mutex
	throttle time.Duration // Disk throttling to prevent a heavy upgrade from hogging resources
	log      log.Logger

	active    uint32          // Flag whether the event loop was started
	update    chan struct{}   // Notification channel that headers should be processed
	quit      chan chan error // Quit channel to tear down running goroutines
	ctx       context.Context
	ctxCancel func()
	// TODO: Ability to Subscribe to changes
	// TODO: Hooks on GET, SET, DELETE actions
}

////////////////////////////////////////////////////////////////////////////////
func (Database) Name() string {
	return "leveldb"
}

func (self Database) Storage() string {
	return self.Options.Storage.String()
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
				// MaxFileSize:                   (8 * opt.MiB),
				// MaxOpenFiles:                  1,
				// CreateIfMissing:               opt.CreateIfMissing,
				// ErrorIfExists:                 opt.ErrorIfExists,
				// Level0StopWriteFiles: 0,
				// Level0SlowdownWriteFiles: 0,
				// Level0CompactionFiles: 0,
				// IterationBytesPerSampleSeek: 0,
				// MinimalAllowedOverlapSeeks: 0,
				// CompactionBytesPerSeek: 0,
				// CompactionConcurrency: 0,
				// BlockCompressionRatio: 0.00,
				// BlockRestartInterval: 0,
				// BlockSize: 0,
				// MaxGrandparentOverlapBytes: 0,
				// Comparator: nil,
				// Logger: nil,
				// TODO: Can use this to specify a virtual filesystem for greater
				//       security.
				// FileSystem: opts.getFileSystem(),
			},
			Read: &opt.ReadOptions{
				DontFillCache: true,
				// TODO: Why is this erroring?
				//VerifyChecksums: false,
			},
			Write: &opt.WriteOptions{
				Sync: true,
			},
		},
	}

	database.Store, err = leveldb.OpenFile(path, database.Options.Runtime)
	database.update = make(chan struct{}, 1)
	database.quit = make(chan chan error)
	database.ctx, database.ctxCancel = context.WithCancel(context.Background())

	//go database.UpdateLoop()

	// TODO: Read up more on finalizers
	// runtime.SetFinalizer(database, (Database).finalize)
	return database, err
}

// SetReadOnly makes DB read-only. It will stay read-only until reopened.
func (self Database) SetReadOnly() error {
	return self.Store.SetReadOnly()
}

func (self Database) finalize() {
	go self.Store.Close()
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
func (self Database) WriteLoop() {
	var (
		updating bool
		updated  time.Time
		err      error
	)

	for {
		select {
		case errc := <-self.quit:
			// Chain indexer terminating, report no failure and abort
			errc <- nil
			return

		case <-self.update:
			// Section headers completed (or rolled back), update the index
			self.Access.Lock()
			if time.Since(updated) > 8*time.Second {
				updated = time.Now()
			}
			self.Access.Unlock()
			if err != nil {
				select {
				case <-self.ctx.Done():
					<-self.quit <- nil
					return
				default:
				}
			}
			self.Access.Lock()

			// If processing succeeded and no reorgs occurred, mark the section completed
			if err == nil {
				if updating {
					updating = false
				}
			} else {
				// If processing failed, don't retry until further notification
			}
			// If there are still further sections to process, reschedule
			time.AfterFunc(self.throttle, func() {
				select {
				case self.update <- struct{}{}:
				default:
				}
			})
		}
		self.Access.Unlock()
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self Database) Close() error {
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
	switch {
	case len(errs) == 0:
		return nil
	case len(errs) == 1:
		return errs[0]
	default:
		return fmt.Errorf("%v", errs)
	}

	runtime.SetFinalizer(self, nil)
	err := self.Store.Close()
	self.Store = nil
	return err
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
		return errors.ErrEmptyKey
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
	self.Access.Lock()
	defer self.Access.Unlock()

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
//		return &Snapshot{err: errors.ErrDBClosed}
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

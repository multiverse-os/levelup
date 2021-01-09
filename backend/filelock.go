package backend

//import (
//	"fmt"
//	"syscall"
//
//	"github.com/pkg/errors"
//	leveldb "github.com/syndtr/goleveldb/leveldb"
//	opt "github.com/syndtr/goleveldb/leveldb/opt"
//)
//
//// EXAMPLE FROM
//// https://github.com/vuongdh/fabric/blob/b89e02b20806018d27bd6d75c86c12ad8f602249/common/ledger/util/leveldbhelper/leveldb_helper.go
//
//// FileLock encapsulate the DB that holds the file lock.
//// As the FileLock to be used by a single process/goroutine,
//// there is no need for the semaphore to synchronize the
//// FileLock usage.
//type FileLock struct {
//	db       *leveldb.DB
//	filePath string
//}
//
//// NewFileLock returns a new file based lock manager.
//func NewFileLock(filePath string) *FileLock {
//	return &FileLock{
//		filePath: filePath,
//	}
//}
//
//// Lock acquire a file lock. We achieve this by opening
//// a db for the given filePath. Internally, leveldb acquires a
//// file lock while opening a db. If the db is opened again by the same or
//// another process, error would be returned. When the db is closed
//// or the owner process dies, the lock would be released and hence
//// the other process can open the db. We exploit this leveldb
//// functionality to acquire and release file lock as the leveldb
//// supports this for Windows, Solaris, and Unix.
//func (f *FileLock) Lock() error {
//	dbOpts := &opt.Options{}
//	var err error
//	var dirEmpty bool
//	if dirEmpty, err = fileutil.CreateDirIfMissing(f.filePath); err != nil {
//		panic(fmt.Sprintf("Error creating dir if missing: %s", err))
//	}
//	dbOpts.ErrorIfMissing = !dirEmpty
//	f.db, err = leveldb.OpenFile(f.filePath, dbOpts)
//	if err != nil && err == syscall.EAGAIN {
//		return errors.Errorf("lock is already acquired on file %s", f.filePath)
//	}
//	if err != nil {
//		panic(fmt.Sprintf("Error acquiring lock on file %s: %s", f.filePath, err))
//	}
//	return nil
//}
//
//// Unlock releases a previously acquired lock. We achieve this by closing
//// the previously opened db. FileUnlock can be called multiple times.
//func (f *FileLock) Unlock() {
//	if f.db == nil {
//		return
//	}
//	if err := f.db.Close(); err != nil {
//		logger.Warningf("unable to release the lock on file %s: %s", f.filePath, err)
//		return
//	}
//	f.db = nil
//}

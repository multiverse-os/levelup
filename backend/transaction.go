package backend

// EXAMPLE OF DOING TRANSACTIONS
// https://github.com/ljr78/syncthing/blob/f0e33d052a03f570cb854ff528e9680864c6d23b/lib/db/backend/leveldb_backend.go<Paste>

// AND SIMPLER
// https://github.com/rkfg/ns2query/blob/9f1c5e381c10bb1f13ae5feb6b82c5e0b67937e3/leveldb.go

//func (b *Database) NewReadTransaction() (ReadTransaction, error) {
//	return b.Snapshot()
//}
//
//func (b *Database) newSnapshot() (leveldbSnapshot, error) {
//	rel, err := newReleaser(b.closeWG)
//	if err != nil {
//		return leveldbSnapshot{}, err
//	}
//	snap, err := b.ldb.GetSnapshot()
//	if err != nil {
//		rel.Release()
//		return leveldbSnapshot{}, wrapLeveldbErr(err)
//	}
//	return leveldbSnapshot{
//		snap: snap,
//		rel:  rel,
//	}, nil
//}
//
//func (b *leveldbBackend) NewWriteTransaction() (WriteTransaction, error) {
//	rel, err := newReleaser(b.closeWG)
//	if err != nil {
//		return nil, err
//	}
//	snap, err := b.newSnapshot()
//	if err != nil {
//		rel.Release()
//		return nil, err // already wrapped
//	}
//	return &leveldbTransaction{
//		leveldbSnapshot: snap,
//		ldb:             b.ldb,
//		batch:           new(leveldb.Batch),
//		rel:             rel,
//	}, nil
//}

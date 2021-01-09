package backend

// DiskUsage returns the current disk size used by this levelDB.
// For in-mem datastores, it will return 0.
//func (d *Datastore) DiskUsage() (uint64, error) {
//	if d.path == "" { // in-mem
//		return 0, nil
//	}
//
//	var du uint64
//
//	err := filepath.Walk(d.path, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		du += uint64(info.Size())
//		return nil
//	})
//
//	if err != nil {
//		return 0, err
//	}
//
//	return du, nil
//}
//

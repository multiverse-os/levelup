package levelup

////////////////////////////////////////////////////////////////////////////////
// TODO: Here we will take the key and modify the set get or delete based on,
//       the key type, we need to ensure all put/get/delete goes through here
//       so that the the bottleneck is here, and the key checking cna happen
//       here. This will enable us to support the variety of key types we want
//       such as CACHE, IMMUTABLE, DOCUMENT, BYTES, SPECIAL keys.
func (self *Database) Put(k key, value []byte) error {
	return self.storage.Set(k.Bytes(), value)
}

func (self *Database) Get(k key) ([]byte, error) {
	return self.storage.Get(k.Bytes())
}

func (self *Database) Delete(k key) error {
	return self.storage.Delete(k.Bytes())
}

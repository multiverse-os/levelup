package levelup

////////////////////////////////////////////////////////////////////////////////
// TODO: Here we will take the key and modify the set get or delete based on,
//       the key type, we need to ensure all put/get/delete goes through here
//       so that the the bottleneck is here, and the key checking cna happen
//       here. This will enable us to support the variety of key types we want
//       such as CACHE, IMMUTABLE, DOCUMENT, BYTES, SPECIAL keys.
func (self *Database) Put(k Key, value []byte) error {
	self.Lock()
	defer self.Unlock()
	return self.storage.Set(k.Bytes(), value)
}

func (self *Database) Get(k Key) ([]byte, error) {
	self.Lock()
	defer self.Unlock()
	return self.storage.Get(k.Bytes())
}

func (self *Database) Delete(k Key) error {
	self.Lock()
	defer self.Unlock()
	return self.storage.Delete(k.Bytes())
}

////////////////////////////////////////////////////////////////////////////////
func (self *Database) PutObject(k Key, value interface{}) error {
	self.Lock()
	defer self.Unlock()
	data, err := self.Codec.Encode(value)
	//data = self.Codec.Compress(data)
	if err != nil {
		return err
	}
	return self.storage.Set(k.Bytes(), data)
}

func (self *Database) GetObject(key Key, value interface{}) error {
	self.Lock()
	defer self.Unlock()
	data, err := self.storage.Get(key.Bytes())
	if err != nil {
		return err
	}

	err = self.Codec.Decode(data, &value)
	if err != nil {
		return err
	}
	return nil
}

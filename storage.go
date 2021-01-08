package levelup

////////////////////////////////////////////////////////////////////////////////
// TODO: Here we will take the key and modify the set get or delete based on,
//       the key type, we need to ensure all put/get/delete goes through here
//       so that the the bottleneck is here, and the key checking cna happen
//       here. This will enable us to support the variety of key types we want
//       such as CACHE, IMMUTABLE, DOCUMENT, BYTES, SPECIAL keys.
func (self *Database) Put(k key, value []byte) error {
	self.Access.Lock()
	defer self.Access.Unlock()
	return self.storage.Set(k.Bytes(), value)
}

func (self *Database) Get(k key) ([]byte, error) {
	self.Access.Lock()
	defer self.Access.Unlock()
	return self.storage.Get(k.Bytes())
}

func (self *Database) Delete(k key) error {
	self.Access.Lock()
	defer self.Access.Unlock()
	return self.storage.Delete(k.Bytes())
}

////////////////////////////////////////////////////////////////////////////////
func (self *Database) PutObject(k key, value interface{}) error {
	self.Access.Lock()
	defer self.Access.Unlock()
	data, err := self.Codec.Encode(value)
	//data = self.Codec.Compress(data)
	if err != nil {
		return err
	}
	return self.storage.Set(k.Bytes(), data)
}

func (self *Database) GetObject(key key, value interface{}) error {
	self.Access.Lock()
	defer self.Access.Unlock()
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

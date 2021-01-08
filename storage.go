package levelup

import (
	validations "github.com/multiverse-os/levelup/model/validations"
)

////////////////////////////////////////////////////////////////////////////////
func (self *Database) Put(key string, value []byte) error {
	if validations.NotEmpty(key) {
		return self.Storage.Set(Key(key).Bytes(), value)
	}
}

func (self *Database) Get(key string) ([]byte, error) {
	if validations.NotEmpty(key) {
		return self.Storage.Get(Key(key).Bytes())
	}
}

func (self *Database) Delete(key string) error {
	if validations.NotEmpty(key) {
		return self.Storage.Delete(Key(key).Bytes())
	}
}

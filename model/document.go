package model

type Document struct {
	Collection *Collection

	Id    []byte
	Name  string
	Value interface{}
	Data  []byte

	// TODO: Require this public key to sign any changes to the document
	//EditorsKeys []*keypair.PublicKey
	//ViewersKeys []*keypair.PublicKey

	// Stats [ Track to automatically optimize the database ]
	//WritesPerHour int
	//ReadsPerHour  int

	// Cache
	//TTL time.Time

	// Versioning
	//Versions []*Document

	// Tree
	//Parent         *Document
	//ChildDocuments []*Document
}

////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////
//func (self *Database) PutObject(key string, value interface{}) error {
//	if validations.IsEmpty(key) {
//		return errors.ErrEmptyKey
//	} else {
//		data, err := self.Codec.Encode(value)
//		//data = self.Codec.Compress(data)
//		if validations.NoNil(err) {
//			return err
//		}
//		return self.Put(key, data)
//	}
//	return nil
//}
//
//func (self *Database) GetObject(key string, value interface{}) error {
//	if validations.IsEmpty(key) {
//		return errors.ErrEmptyKey
//	} else {
//		data, err := self.Get(key)
//		if err != nil {
//			return err
//		}
//
//		err = self.Codec.Decode(data, &value)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//}

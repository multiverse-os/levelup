package levelup

import (
	"time"
)

type Record struct {
	Database   *Database
	Collection *Collection

	CollectionId uint32

	// TODO: Could have codec be per database, per collection; then whichever is
	// slected is linked here
	//Codec *Codec

	CreatedAt     time.Time
	LastUpdatedAt time.Time

	//Encoding    string
	//Compression string

	Key   Key
	Value []byte

	CompressedValue   []byte
	UncompressedValue []byte

	Document interface{}
}

type value struct {
	Uncompressed []byte
	Compressed   []byte
	Encrypted    []byte
	Decrypted    []byte
	Signed       []byte
	Encoded      []byte
	Decoded      interface{}
}

func (self Key) NewRecord(input []byte) Record {
	return Record{
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Key:           self,
		Value:         input,
	}
}

func (self Record) Encode() []byte {
	encodedData, _ := self.Database.Codec.Encode(self.Value)
	return encodedData
}

func (self Record) Decode() interface{} {
	self.Database.Codec.Decode(self.Value, &self.Document)
	return self.Document
}

func (self Record) Compress() []byte {
	compressedData, _ := self.Database.Codec.Compress(self.Value)
	return compressedData
}

func (self Record) Uncompress() []byte {
	uncompressedData, _ := self.Database.Codec.Uncompress(self.Value)
	return uncompressedData
}

// TODO: Perhaps instead of using the codec, we can simply provide here the
//       cryptosystem (asymmetric vs. symmetric), and the algorith, and either a
//       seed to generate the key, or the key itself. This could be stored
//       inside the database, or the collection and be used solely or in
//       combination with a key generated from a field (like a password field).
func (self Record) Encrypt() []byte { return []byte{} }

func (self Record) Decrypt() []byte { return []byte{} }

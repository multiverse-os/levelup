package levelup

import (
	"time"
)

type ValueType int

const (
	intType ValueType = iota
	byteType
	runeType
	floatType
	stringType
	uintType
	intSliceType
	byteSliceType
	stringSliceType
	floatSliceType
	uintSliceType
)

type Record struct {
	Database   Database
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

////////////////////////////////////////////////////////////////////////////////
type value struct {
	Type ValueType

	value        []byte
	untypedValue interface{}
	// Compression
	Uncompressed []byte
	Compressed   []byte
	// Cryptographic
	Encrypted []byte
	Decrypted []byte
	Signed    []byte
	// Encoding
	Encoded []byte
	Decoded interface{}
}

func (self value) Int() int       { return self.untypedValue.(int) }
func (self value) Byte() byte     { return self.untypedValue.(byte) }
func (self value) Rune() rune     { return self.untypedValue.(rune) }
func (self value) Float() float64 { return self.untypedValue.(float64) }
func (self value) String() string { return self.untypedValue.(string) }
func (self value) UInt() uint     { return self.untypedValue.(uint) }

func (self value) IntSlice() []int       { return self.untypedValue.([]int) }
func (self value) ByteSlice() []byte     { return self.untypedValue.([]byte) }
func (self value) StringSlice() []string { return self.untypedValue.([]string) }
func (self value) FloatSlice() []float64 { return self.untypedValue.([]float64) }
func (self value) UIntSlice() []uint     { return self.untypedValue.([]uint) }

////////////////////////////////////////////////////////////////////////////////
// NOTE: Now we simply use Value(input []byte) to work with (key/value) level
//       of the database, or we use Document(input []byte). Then these will
//       handle encoding into and out of objects, instead of interacting with
//       two types of Put commands.

func (self Key) Value(input []byte) *Record {
	return &Record{
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Key:           self,
		Value:         input,
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self *Record) Encode() []byte {
	encodedData, _ := self.Database.Codec.Encode(self.Value)
	return encodedData
}

func (self *Record) Decode() interface{} {
	self.Database.Codec.Decode(self.Value, &self.Document)
	return self.Document
}

func (self *Record) Compress() []byte {
	compressedData, _ := self.Database.Codec.Compress(self.Value)
	return compressedData
}

func (self *Record) Uncompress() []byte {
	uncompressedData, _ := self.Database.Codec.Uncompress(self.Value)
	return uncompressedData
}

// TODO: Perhaps instead of using the codec, we can simply provide here the
//       cryptosystem (asymmetric vs. symmetric), and the algorith, and either a
//       seed to generate the key, or the key itself. This could be stored
//       inside the database, or the collection and be used solely or in
//       combination with a key generated from a field (like a password field).
func (self *Record) Encrypt() []byte { return []byte{} }

func (self *Record) Decrypt() []byte { return []byte{} }

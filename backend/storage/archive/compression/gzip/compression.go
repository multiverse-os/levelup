package gzip

import (
	"bytes"
	"io/ioutil"
)

// NOTE: Parallelization of gzip boosts the preformance by:
// 1573.23MB/s 	101.9x 	4902285651 	+2.53%

// https://github.com/klauspost/pgzip/

// GZIP.SetConcurrency(BlockSize, Blocks)
type Compression struct {
	ChunkSize int
	Parallel  int
}

func Load() Compression {
	return Compression{
		//ChunkSize: 512000, // 5 Megabytes
		ChunkSize: 102400, // 1 Megabyte
		Parallel:  4,      // 4 Parallel Processes
	}
}

// Compress([]byte) []byte
// Compress(interface {}) ([]byte, error)
func (self Compression) Compress(input []byte) ([]byte, error) {
	var buffer bytes.Buffer
	gzipWriter := NewWriter(&buffer)
	gzipWriter.SetConcurrency(self.ChunkSize, self.Parallel)
	gzipWriter.Write(input)
	gzipWriter.Close()
	return buffer.Bytes(), nil
}

// Uncompress([]byte) ([]byte, error)
// Uncompress([]byte, interface {}) error
func (Compression) Uncompress(input []byte) ([]byte, error) {
	gzipReader, _ := NewReader(bytes.NewReader(input))
	return ioutil.ReadAll(gzipReader)
}

func (Compression) String() string {
	return "gzip"
}

package compression

import (
	gzip "github.com/multiverse-os/codec/compression/gzip"
	none "github.com/multiverse-os/codec/compression/none"
	zstd "github.com/multiverse-os/codec/compression/zstd"
)

////////////////////////////////////////////////////////////////////////////////
type Type int

const (
	None Type = iota
	Gzip
	Snappy
	Zstd
)

// NOTE: Decompress means to relieve the pressure or compression on something,
//       whereas uncompress means to restore a compressed file.
type Codec interface {
	Compress(input []byte) ([]byte, error)
	Uncompress(data []byte) ([]byte, error)
	String() string
}

func Algorithm(algorithm Type) Codec {
	switch algorithm {
	case Gzip:
		return gzip.Compression{}
	case Zstd:
		return zstd.Compression{}
	default: // None
		return none.Compression{}
	}

}

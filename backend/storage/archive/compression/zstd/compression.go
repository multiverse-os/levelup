package zstd

import (
	zstd "github.com/klauspost/compress/zstd"
)

// TODO: Eventually options will go here
type Compression struct{}

func (Compression) Compress(input []byte) ([]byte, error) {
	encoder, _ := zstd.NewWriter(nil)
	return encoder.EncodeAll(input, make([]byte, 0, len(input))), nil
}

func (Compression) Uncompress(input []byte) ([]byte, error) {
	decoder, _ := zstd.NewReader(nil)
	return decoder.DecodeAll(input, nil)
}

func (Compression) String() string {
	return "zstd"
}

////////////////////////////////////////////////////////////////////////////////
// Streaming
//
//func Decompress(in io.Reader, out io.Writer) error {
//    d, err := zstd.NewReader(input)
//    if err != nil {
//        return err
//    }
//    defer d.Close()
//
//    // Copy content...
//    _, err := io.Copy(out, d)
//    return err
//}

package none

type Compression struct{}

func (Compression) Compress(input []byte) ([]byte, error) {
	return input, nil
}

func (Compression) Uncompress(input []byte) ([]byte, error) {
	return input, nil
}

func (Compression) String() string {
	return "none"
}

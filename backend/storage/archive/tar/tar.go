package tar

type FileDescriptor int

const (
	Memory FileDescriptor = iota
	File
)

////////////////////////////////////////////////////////////////////////////////
type FileArchive struct {
	FD FileDescriptor
}

func Archive(path string) FileArchive {
	return FileArchive{}
}

func (self FileArchive) Compress() ([]byte, error) {
	return []byte{}, nil
}

func (self FileArchive) Uncompress() ([]byte, error) {
	return []byte{}, nil
}

func (self FileArchive) String() string { return "tar" }

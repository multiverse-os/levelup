package archive

import (
	tar "github.com/multiverse-os/levelup/backend/archive/tar"
)

type Type int

const (
	TarArchive Type = iota
)

type Archive interface {
	String() string

	Compress() ([]byte, error)
	Uncompress() ([]byte, error)

	// Directory adds a new directory entry to the archive and sets the
	// directory for subsequent calls to Header.
	//Directory(name string) error

	//// Header adds a new file to the archive. The file is added
	//// to the directory
	//// set by Directory. The content of the file must be
	//// written to the returned
	//// writer.
	//Header(os.FileInfo) (io.Writer, error)

	//// Close flushes the archive and closes the
	//// underlying file.
	//Close() error
}

func new(t Type, path string) Archive {
	switch t {
	default: // Tar
		return tar.Archive("./test-data")
	}
}

func Tar(path string) Archive {
	return new(TarArchive, path)
}

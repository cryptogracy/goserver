package filesystem

import (
	"io"
	"os"
)

type Filesystem interface {
	Open(name string) (File, error)
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
}

// OsFS implements fileSystem using the local disk.
type OsFS struct{}

func (OsFS) Open(name string) (File, error) {
	return os.Open(name)
}

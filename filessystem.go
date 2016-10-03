package main

import (
	"io"
	"os"
)

type fileSystem interface {
	Open(name string) (file, error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (file, error) {
	return os.Open(name)
}

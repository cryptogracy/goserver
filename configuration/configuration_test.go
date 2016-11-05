package configuration

import (
	"io"
	"strings"
	"testing"

	"github.com/cryptogracy/goserver/filesystem"
)

type fakeFileSystem struct {
	content string
}

func (ffs fakeFileSystem) Open(name string) (filesystem.File, error) {
	reader := strings.NewReader(ffs.content)
	return fakeFile{reader, nil, nil}, nil
}

type fakeFile struct {
	reader io.Reader
	io.Seeker
	io.ReaderAt
}

func (ff fakeFile) Read(b []byte) (n int, err error) {
	return ff.reader.Read(b)
}

func (ff fakeFile) Close() error {
	return nil
}

func TestInit(t *testing.T) {
	fs = fakeFileSystem{`
address: 1.2.3.4:1234
static: static_test
dir: dir_test
tempdir: tempdir_test
database: database_test`}
	Init()
	if Config.Static != "static_test" {
		t.Error("Address does not match")
	}
	if Config.Address != "1.2.3.4:1234" {
		t.Error("Address does not match")
	}
	if Config.Dir != "dir_test" {
		t.Error("Address does not match")
	}
	if Config.Tempdir != "tempdir_test" {
		t.Error("Address does not match")
	}
	if Config.Database != "database_test" {
		t.Error("Address does not match")
	}
}

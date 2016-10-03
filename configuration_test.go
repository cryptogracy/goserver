package main

import (
	"io"
	"strings"
	"testing"
)

type fakeFileSystem struct {
	content string
}

func (ffs fakeFileSystem) Open(name string) (file, error) {
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

func TestReadConfiguationDefault(t *testing.T) {
	fs = fakeFileSystem{}
	config := readConfiguration()
	if config.Static != "ui" {
		t.Error("Address does not match")
	}
	if config.Address != "127.0.0.1:8000" {
		t.Error("Address does not match")
	}
	if config.Dir != "cache" {
		t.Error("Address does not match")
	}
	if config.Tempdir != "tmp" {
		t.Error("Address does not match")
	}
	if config.Database != "goserver.db" {
		t.Error("Address does not match")
	}
}

func TestReadConfiguationSet(t *testing.T) {
	fs = fakeFileSystem{`
address: 1.2.3.4:1234
static: static_test
dir: dir_test
tempdir: tempdir_test
database: database_test`}
	config := readConfiguration()
	if config.Static != "static_test" {
		t.Error("Address does not match")
	}
	if config.Address != "1.2.3.4:1234" {
		t.Error("Address does not match")
	}
	if config.Dir != "dir_test" {
		t.Error("Address does not match")
	}
	if config.Tempdir != "tempdir_test" {
		t.Error("Address does not match")
	}
	if config.Database != "database_test" {
		t.Error("Address does not match")
	}
}

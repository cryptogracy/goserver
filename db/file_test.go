package db

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/cryptogracy/goserver/configuration"
)

type testReader struct {
	content string
}

const hash string = "1ae597b48b2e9befd48d39e7c98b14c0bd5b4320f5214a7773f5a3ef321571d318b6b2450d55372269a12549a38f6683d421212fd59adc8873711696cefe09f4"

func TestAddFile(t *testing.T) {
	defer os.RemoveAll(beforeAdd(t))

	// Initialize in memory database
	Init(":memory:")
	defer db.Close()

	before := time.Now().UTC().Add(time.Duration(10 * time.Minute))
	if err := AddFile(hash, 10, strings.NewReader("This is the content")); err != nil {
		t.Error("Could not add file", err)
	}
	after := time.Now().UTC().Add(time.Duration(10 * time.Minute))

	var file File
	if info := db.First(&file, &File{Hash: hash}); info.Error != nil {
		t.Error("Unable to get file", info.Error)
	}

	if !(before.Before(file.Death) && after.After(file.Death)) {
		t.Errorf(`Wrong time.
  Before: '%v'
  File:   '%v'
  After:  '%v'`, before, file.Death, after)
	}
}

func TestAddFileTwice(t *testing.T) {
	defer os.RemoveAll(beforeAdd(t))

	// Initialize in memory database
	Init(":memory:")
	defer db.Close()

	if err := AddFile(hash, 10, strings.NewReader("This is the content")); err != nil {
		t.Fatal("Unable to add", err)
	}

	if err := AddFile(hash, 10, strings.NewReader("This is the content")); err != ErrFileExist {
		t.Error("Tried to Add same file twice, but wrong error", err)
	}
}

func TestAddFileWrongHash(t *testing.T) {
	defer os.RemoveAll(beforeAdd(t))

	// Initialize in memory database
	Init(":memory:")
	defer db.Close()

	if err := AddFile(hash, 10, strings.NewReader("This is the content")); err != nil {
		t.Fatal("Unable to add", err)
	}

	if err := AddFile(hash, 10, strings.NewReader("This is the content")); err != ErrFileExist {
		t.Error("Tried to Add same file twice, but wrong error", err)
	}
}

func TestRemoveFiles(t *testing.T) {
	defer os.RemoveAll(beforeAdd(t))

	// Initialize in memory database
	Init(":memory:")
	defer db.Close()

	fn := func(lifespan int, t *testing.T, del int64) {
		if err := AddFile(hash, lifespan, strings.NewReader("This is the content")); err != nil {
			t.Fatal("Unable to add", err)
		}

		if affected, err := removeFiles(); err != nil {
			t.Error("Unable to remove:", err)
		} else if affected != del {
			t.Errorf("Wrong number of removed files. %v != %v", affected, del)
		} else {
			t.Log("asd")
		}
	}

	fn(-5, t, 1)
	fn(0, t, 1)
	fn(5, t, 0)

}

func beforeAdd(t *testing.T) (resDir string) {
	// Create temporary directories
	if dir, err := ioutil.TempDir("", ""); err == nil {
		resDir = dir
	} else {
		t.Error("Unable to create temporary directory", err)
	}
	configuration.Config.Tempdir = path.Join(resDir, "tmp")
	if err := os.Mkdir(configuration.Config.Tempdir, 0777); err != nil {
		t.Error("Unable to create temporary directory", err)
	}
	configuration.Config.Dir = path.Join(resDir, "cache")
	if err := os.Mkdir(configuration.Config.Dir, 0777); err != nil {
		t.Error("Unable to create temporary directory", err)
	}
	return
}

func TestIsHash(t *testing.T) {
	if !isHash(hash, strings.NewReader("This is the content")) {
		t.Error("Correct hash, but isHash returns false")
	}
	if isHash(hash, strings.NewReader("This is the wrong content")) {
		t.Error("Incorrect hash, but isHash returns true")
	}
}

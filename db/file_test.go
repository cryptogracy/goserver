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

func TestAddFile(t *testing.T) {
	// Create temporary directories
	if dir, err := ioutil.TempDir("", ""); err == nil {
		configuration.Config.Tempdir = path.Join(dir, "tmp")
		configuration.Config.Dir = path.Join(dir, "cache")
		defer os.RemoveAll(dir)
	} else {
		t.Fatal("Unable to create temporary directory", err)
	}
	if err := os.Mkdir(configuration.Config.Tempdir, 0777); err != nil {
		t.Fatal("Unable to create temporary directory", err)
	}
	if err := os.Mkdir(configuration.Config.Dir, 0777); err != nil {
		t.Fatal("Unable to create temporary directory", err)
	}

	// Initialize in memory database
	Init(":memory:")
	defer db.Close()

	reader := strings.NewReader("This is the content")
	hash := "1ae597b48b2e9befd48d39e7c98b14c0bd5b4320f5214a7773f5a3ef321571d318b6b2450d55372269a12549a38f6683d421212fd59adc8873711696cefe09f4"

	before := time.Now().Add(time.Duration(10) * time.Second)
	if err := AddFile(hash, 10, reader); err != nil {
		t.Error("Could not add", hash, err)
	}

	after := time.Now().Add(time.Duration(10) * time.Second)

	var file File
	if info := db.First(&file, &File{Hash: "hash"}); info.Error != nil {
		t.Error("Unable to get", hash, info.Error)
	}

	if !(before.Before(file.Death) && after.After(file.Death)) {
		t.Error("Wrong time")
	}
}

func TestIsHash(t *testing.T) {
	hash := "e9b4a070e1be7dbec8b340ef80744d32f8d3cb9a9d1f89fed225037b9eaf0a271876adadc7485a6090aa2e8c0e30984c26710a501ce889cccb363d1cb28f087b"
	if !isHash(hash, strings.NewReader("Das ist ein Test")) {
		t.Error("Correct hash, but isHash returns false")
	}
	if isHash(hash, strings.NewReader("Das ist noch ein Test")) {
		t.Error("Incorrect hash, but isHash returns true")
	}
}

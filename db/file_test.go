package db

import (
	"testing"
	"time"
)

func TestAddFile(t *testing.T) {
	Init(":memory:")
	defer db.Close()

	before := time.Now().Add(time.Duration(10) * time.Second)
	if !AddFile("Hash1", 10) {
		t.Error("Could not add Hash1")
	}
	after := time.Now().Add(time.Duration(10) * time.Second)

	var file File
	if info := db.First(&file, &File{Hash: "Hash1"}); info.Error != nil {
		t.Error(info.Error)
	}

	if !(before.Before(file.Death) && after.After(file.Death)) {
		t.Error("Wrong time")
	}
	if AddFile("Hash1", 10) {
		t.Error("Could add Hash1 twice")
	}
}

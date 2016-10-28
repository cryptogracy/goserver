package db

import (
	"testing"
)

func TestInit(t *testing.T) {
	err := Init(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	if !db.HasTable(&File{}) {
		t.Error("Table for Files not present")
	}
	if !db.HasTable(&Metadata{}) {
		t.Error("Table for Metadata not present")
	}
}

func TestClose(t *testing.T) {
	Init(":memory:")
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestCleanup(t *testing.T) {
	Init(":memory:")
	defer db.Close()

	AddFile("Hash1", 10)
	AddFile("Hash2", 0)
	AddFile("Hash3", -5)

	affected, err := Cleanup()

	if err != nil {
		t.Error(err)
	}

	if affected != 2 {
		t.Error("Wrong number if deleted Rows")
	}

	var file File
	if info := db.First(&file, "hash = ?", "Hash1"); info.Error != nil {
		t.Error(info.Error)
	}
	if info := db.First(&file, "hash = ?", "Hash2"); info.Error == nil {
		t.Error("Entry Hash2 not deleted")
	}
	if info := db.First(&file, "hash = ?", "Hash3"); info.Error == nil {
		t.Error("Entry Hash2 not deleted")
	}
}

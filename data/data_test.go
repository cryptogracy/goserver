package data

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	db, err := Init(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer db.db.Close()

	if !db.db.HasTable(&File{}) {
		t.Error("Table for Files not present")
	}
	if !db.db.HasTable(&Metadata{}) {
		t.Error("Table for Metadata not present")
	}
}

func TestClose(t *testing.T) {
	db, _ := Init(":memory:")
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestAddFile(t *testing.T) {
	db, _ := Init(":memory:")
	defer db.Close()

	before := time.Now().Add(time.Duration(10) * time.Second)
	if !db.AddFile("Hash1", 10) {
		t.Error("Could not add Hash1")
	}
	after := time.Now().Add(time.Duration(10) * time.Second)

	var file File
	if info := db.db.First(&file, &File{Hash: "Hash1"}); info.Error != nil {
		t.Error(info.Error)
	}

	if !(before.Before(file.Death) && after.After(file.Death)) {
		t.Error("Wrong time")
	}
	if db.AddFile("Hash1", 10) {
		t.Error("Could add Hash1 twice")
	}
}

func TestAddMeta(t *testing.T) {
	db, _ := Init(":memory:")
	defer db.Close()

	db.AddFile("Hash1", 10)
	db.AddFile("Hash2", 10)
	db.AddFile("Hash3", 10)

	if _, success := db.AddMeta(10, []string{"Hash1", "Hash2", "Hash3"}); !success {
		t.Error("Unable to create metadata record")
	}

	if _, success := db.AddMeta(10, []string{"Hash1", "Hash2", "Hash3"}); !success {
		t.Error("Unable to create metadata record twice")
	}
}

func TestCleanup(t *testing.T) {
	db, _ := Init(":memory:")
	defer db.Close()

	db.AddFile("Hash1", 10)
	db.AddFile("Hash2", 0)
	db.AddFile("Hash3", -5)

	affected, err := db.Cleanup()

	if err != nil {
		t.Error(err)
	}

	if affected != 2 {
		t.Error("Wrong number if deleted Rows")
	}

	var file File
	if info := db.db.First(&file, "hash = ?", "Hash1"); info.Error != nil {
		t.Error(info.Error)
	}
	if info := db.db.First(&file, "hash = ?", "Hash2"); info.Error == nil {
		t.Error("Entry Hash2 not deleted")
	}
	if info := db.db.First(&file, "hash = ?", "Hash3"); info.Error == nil {
		t.Error("Entry Hash2 not deleted")
	}
}

func TestRandSeq(t *testing.T) {
	for i := 1; i < 1000; i++ {
		if len(randSeq(i)) != i {
			t.Error("String has wrong length")
		}

		if randSeq(10) == randSeq(10) {
			t.Error("Generate the same sequence twice")
		}
	}
}

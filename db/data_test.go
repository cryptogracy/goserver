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
	Cleanup()
}

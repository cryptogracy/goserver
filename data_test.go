package main

import (
	"testing"
)

func TestDBInit(t *testing.T) {
	db, err := DBInit(":memory:")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var res string
	err = db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='Files'").Scan(&res)
	if err != nil {
		t.Error(err)
	}
}

/* I do not like this test, cause adding and checking are tested in the same
 * function...maybe both are broken */
func TestAddFileandCheck(t *testing.T) {
	db, _ := DBInit(":memory:")
	defer db.Close()

	// Test adding
	if err := db.AddFile("Hash1", 10); err != nil {
		t.Error(err)
	}
	if err := db.AddFile("Hash1", 2); err == nil {
		t.Error(err)
	}
	if err := db.AddFile("Hash2", -5); err != nil {
		t.Error(err)
	}

	// Test checking and proof if adding has worked.
	if !db.Check("Hash1") {
		t.Error("Hash1 not present")
	}
	if !db.Check("Hash2") {
		t.Error("Hash2 not present")
	}
	if db.Check("Hash3") {
		t.Error("Hash3 present")
	}
}

func TestCleanup(t *testing.T) {
	db, _ := DBInit(":memory:")
	defer db.Close()

	// AddFile some Data
	db.AddFile("Hash1", 10)
	db.AddFile("Hash2", 0)
	db.AddFile("Hash3", -5)

	affected, err := db.Cleanup()

	if err != nil {
		t.Error(err)
	}
	if affected != 2 {
		t.Error("Wrong number of deleted Rows")
	}
	if !db.Check("Hash1") {
		t.Error("Hash1 not present")
	}
	if db.Check("Hash2") {
		t.Error("Hash2 present")
	}
	if db.Check("Hash3") {
		t.Error("Hash3 present")
	}
}

package db

import (
	"testing"
)

func TestAddMeta(t *testing.T) {
	Init(":memory:")
	defer db.Close()

	AddFile("Hash1", 10)
	AddFile("Hash2", 10)
	AddFile("Hash3", 10)

	if _, success := AddMeta(10, []string{"Hash1", "Hash2", "Hash3"}); !success {
		t.Error("Unable to create metadata record")
	}

	if _, success := AddMeta(10, []string{"Hash1", "Hash2", "Hash3"}); !success {
		t.Error("Unable to create metadata record twice")
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

package configuration

import (
	"path/filepath"
	"testing"
)

func TestNotExist(t *testing.T) {
	config := filepath.Join("testdata", "notexist")
	if err := Init(config); err != ErrIO {
		t.Error(err)
	}
}

func TestWrongYaml(t *testing.T) {
	config := filepath.Join("testdata", "wrongyaml")
	if err := Init(config); err != ErrYAML {
		t.Error(err)
	}
}

func TestInit(t *testing.T) {
	config := filepath.Join("testdata", "config")
	if err := Init(config); err != nil {
		t.Error(err)
	}

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

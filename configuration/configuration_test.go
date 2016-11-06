package configuration

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNotExist(t *testing.T) {
	config_file = "file-should-not-exist"
	if err := Init(); err != ErrIO {
		t.Error(err)
	}
}

func TestWrongYaml(t *testing.T) {

	content := "Banana"
	config_file = "test-config-file"
	if err := ioutil.WriteFile(config_file, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test-config-file")

	if err := Init(); err != ErrYAML {
		t.Error(err)
	}
}

func TestInit(t *testing.T) {

	content := `
address: 1.2.3.4:1234
static: static_test
dir: dir_test
tempdir: tempdir_test
database: database_test`
	config_file = "test-config-file"
	if err := ioutil.WriteFile(config_file, []byte(content), 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test-config-file")

	if err := Init(); err != nil {
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

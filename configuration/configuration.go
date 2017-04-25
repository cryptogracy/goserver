package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Config configuration = configuration{address, static, dir, tempdir, database}

const (
	address  = "127.0.0.1:8000"
	static   = "simple-frontend"
	dir      = "cache"
	tempdir  = "tmp"
	database = "goserver.db"
)

type configuration struct {
	Address  string `yaml:"address"`
	Static   string `yaml:"static"`
	Dir      string `yaml:"dir"`
	Tempdir  string `yam:"tempdir"`
	Database string `yam:"database"`
}

func Init(config_file string) error {
	content, err := ioutil.ReadFile(config_file)
	if err != nil {
		return fmt.Errorf("Unable to read configuration file, %v", err)
	}
	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		return fmt.Errorf("Unable to parse configuration file, %v", err)
	}
	return nil
}

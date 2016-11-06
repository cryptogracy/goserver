package configuration

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var Config configuration = configuration{address, static, dir, tempdir, database}

var (
	ErrIO   = errors.New("Unable to read from configuration file")
	ErrYAML = errors.New("Unable to parse configuration file")
)

const (
	address  = "127.0.0.1:8000"
	static   = "ui"
	dir      = "cache"
	tempdir  = "tmp"
	database = "goserver.db"
)

var config_file = "config"

type configuration struct {
	Address  string `yaml:"address"`
	Static   string `yaml:"static"`
	Dir      string `yaml:"dir"`
	Tempdir  string `yam:"tempdir"`
	Database string `yam:"database"`
}

func Init() error {
	content, err := ioutil.ReadFile(config_file)
	if err != nil {
		log.Println(ErrIO)
		return ErrIO
	}
	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		log.Println(ErrYAML)
		return ErrYAML
	}
	return nil
}

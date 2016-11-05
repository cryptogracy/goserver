package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"github.com/cryptogracy/goserver/filesystem"
)

var fs filesystem.Filesystem = filesystem.OsFS{}
var Config configuration = configuration{address, static, dir, tempdir, database}

const (
	config_file = "config"
	address     = "127.0.0.1:8000"
	static      = "ui"
	dir         = "cache"
	tempdir     = "tmp"
	database    = "goserver.db"
)

type configuration struct {
	Address  string `yaml:"address"`
	Static   string `yaml:"static"`
	Dir      string `yaml:"dir"`
	Tempdir  string `yam:"tempdir"`
	Database string `yam:"database"`
}

func Init() {
	file, err := fs.Open(config_file)
	if err != nil {
		log.Println("Cannot open configuration file:", err)
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Cannot open configuration file:", err)
		return
	}
	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		panic(err)
	}
}

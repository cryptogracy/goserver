package main

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	"log"
)

const (
	CONFIG_FILE = "config"
	ADDRESS     = "127.0.0.1:8000"
	STATIC      = "ui"
	DIR         = "cache"
	TEMPDIR     = "tmp"
	DATABASE    = "goserver.db"
)

type Configuration struct {
	Address  string `yaml:"address"`
	Static   string `yaml:"static"`
	Dir      string `yaml:"dir"`
	Tempdir  string `yam:"tempdir"`
	Database string `yam:"database"`
}

func readConfiguration() Configuration {
	config := Configuration{ADDRESS, STATIC, DIR, TEMPDIR, DATABASE}

	file, err := fs.Open(CONFIG_FILE)
	if err != nil {
		log.Println("Cannot open configuration file:", err)
		return config
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Cannot open configuration file:", err)
		return config
	}
	err = goyaml.Unmarshal(content, &config)
	errorPanic(err)
	return config
}

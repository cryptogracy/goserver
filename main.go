package main

import (
	"log"
	"net/http"
	"time"
)

type HttpReturn struct {
	Error string
}

var configuration Configuration
var db DB

func main() {

	configuration = readConfiguration()
	log.Println("Using database", configuration.Database)
	log.Println("Static from", configuration.Static)
	log.Println("Serving files from", configuration.Dir)
	log.Println("Listen on", configuration.Address)

	db.Init()

	go db.RemoveFilePeriod(1 * time.Minute)

	log.Fatal(http.ListenAndServe(configuration.Address, routing()))
}

package main

import (
	"log"
	"net/http"
	"time"
)

var configuration Configuration
var data dataControl

func main() {

	configuration = readConfiguration()
	log.Println("Using database", configuration.Database)
	log.Println("Static from", configuration.Static)
	log.Println("Serving files from", configuration.Dir)
	log.Println("Listen on", configuration.Address)

	var err error
	data, err = DBInit()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	go func() {
		err := data.Cleanup()
		if err != nil {
			log.Println(err)
		}
		time.Sleep(1 * time.Minute)
	}()

	log.Fatal(http.ListenAndServe(configuration.Address, routing()))
}

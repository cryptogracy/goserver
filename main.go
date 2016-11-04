package main

import (
	"github.com/cryptogracy/goserver/db"
	"log"
	"net/http"
	"time"
)

var configuration Configuration
var fs fileSystem = osFS{}

func main() {

	configuration = readConfiguration()
	log.Println("Using database", configuration.Database)
	log.Println("Static from", configuration.Static)
	log.Println("Serving files from", configuration.Dir)
	log.Println("Listening on http://" + configuration.Address)

	var err error
	if err = db.Init(configuration.Database); err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.Close()

	go func() {
		for true {
			affected, err := db.Cleanup()
			if err == nil {
				log.Printf("Deleted %v old Files\n", affected)
			} else {
				log.Println(err)
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Fatal(http.ListenAndServe(configuration.Address, newRouter()))

}

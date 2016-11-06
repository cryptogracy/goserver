package main

import (
	"log"
	"net/http"
	"time"

	"github.com/cryptogracy/goserver/configuration"
	"github.com/cryptogracy/goserver/db"
	"github.com/cryptogracy/goserver/routing"
)

func main() {

	configuration.Init()
	log.Println("Using database", configuration.Config.Database)
	log.Println("Static from", configuration.Config.Static)
	log.Println("Serving files from", configuration.Config.Dir)
	log.Println("Listening on http://" + configuration.Config.Address)

	var err error
	if err = db.Init(configuration.Config.Database); err != nil {
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

	log.Fatal(http.ListenAndServe(configuration.Config.Address, routing.Router()))

}

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cryptogracy/goserver/configuration"
	"github.com/cryptogracy/goserver/db"
	"github.com/cryptogracy/goserver/routing"
)

func main() {
	configuration.Init(os.Args[1])
	log.Println("Using database", configuration.Config.Database)
	log.Println("Static from", configuration.Config.Static)
	log.Println("Serving files from", configuration.Config.Dir)
	log.Println("Listening on http://" + configuration.Config.Address)

	if err := db.Init(configuration.Config.Database); err != nil {
		log.Println(err)
		panic(err)
	}
	defer db.Close()

	go func() {
		for {
			db.Cleanup()
			time.Sleep(1 * time.Minute)
		}
	}()

	log.Fatal(http.ListenAndServe(configuration.Config.Address, routing.Router()))

}

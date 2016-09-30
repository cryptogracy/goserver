package main

import (
	"log"
)

func errorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func errorLog(text string, err error) {
	if err != nil {
		log.Println(text, err)
	}
}

package main

import (
	"log"
)

func errorLog(text string, err error) {
	if err != nil {
		log.Println(text, err)
	}
}

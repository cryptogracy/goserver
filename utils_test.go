package main

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"
)

func TestErrorLog(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	errorLog("This is a test", errors.New("testerror"))
}

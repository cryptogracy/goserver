package main

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"
)

func TestErrorPanicNoError(t *testing.T) {
	errorPanic(nil)
}

func TestErrorPanicWithError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")

		}
	}()
	errorPanic(errors.New("testerror"))
}

func TestErrorLog(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	errorLog("This is a test", errors.New("testerror"))
}

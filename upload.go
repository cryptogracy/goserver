package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type HttpReturn struct {
	Error string
}

func jsonAnswer(answer map[string]string) []byte {
	json, err := json.Marshal(answer)
	errorPanic(err)
	return json
}

func isHash(hash string, file io.Reader) bool {
	hasher := sha512.New()
	_, err := io.Copy(hasher, file)
	errorPanic(err)
	if fmt.Sprintf("%x", hasher.Sum(nil)) == hash {
		return true
	}
	return false
}

func Upload(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	w.Header().Set("Content-Type", "application/json") // We always return json

	// File is already uploaded
	if db.CheckPresence(hash) {
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonAnswer(map[string]string{"Error": "File already uploaded"}))
		return
	}

	// Open file
	tempfile := path.Join(configuration.Tempdir, hash)
	out, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	// Upload in Progress
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonAnswer(map[string]string{"Error": "Upload in Progress"}))
		return
	}
	defer out.Close()
	defer os.Remove(tempfile) // We never want to keep the tmpfile

	// Copy to tmp
	_, err = io.Copy(out, r.Body)
	errorPanic(err)

	// Check Hash
	_, err = out.Seek(0, os.SEEK_SET)
	if !isHash(hash, out) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonAnswer(map[string]string{"Error": "Wrong hash"}))
		return
	}

	// Put into configuration Dir
	err = os.Rename(tempfile, path.Join(configuration.Dir, hash))
	errorPanic(err)

	// Get lifespan
	lifespan, err := strconv.Atoi(r.Header.Get("x-http-lifespan"))
	errorPanic(err)

	db.AddFile(hash, lifespan)

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonAnswer(map[string]string{}))
}

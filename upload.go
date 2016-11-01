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
  "github.com/cryptogracy/goserver/db"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	w.Header().Set("Content-Type", "application/json") // We always return json

	// Get lifespan
	lifespan, err := strconv.Atoi(r.Header.Get("x-file-lifespan"))
	if err != nil {
		panic(err)
	}

	// Open file
	tempfile := path.Join(configuration.Tempdir, hash)
	out, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	// File exists, so upload in progress
	if os.IsExist(err) {
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonAnswer(map[string]string{"Error": "Upload in Progress"}))
		return
	} else if err != nil {
		panic(err)
	}
	defer out.Close()
	defer os.Remove(tempfile) // We never want to keep the tmpfile

	// Copy content to tmp
	if _, err := io.Copy(out, r.Body); err != nil {
		panic(err)
	}

	// Check if the hash is correct
	if _, err := out.Seek(0, os.SEEK_SET); err != nil {
    panic(err)
  }
	if !isHash(hash, out) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonAnswer(map[string]string{"Error": "Wrong hash"}))
		return
	}

	// File is already uploaded
	if !db.AddFile(hash, lifespan) {
		w.WriteHeader(http.StatusConflict)
		w.Write(jsonAnswer(map[string]string{"Error": "File already uploaded"}))
		return
	}

	// Put into configuration Dir
	if err := os.Rename(tempfile, path.Join(configuration.Dir, "files", hash)); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonAnswer(map[string]string{}))
}

func jsonAnswer(answer map[string]string) []byte {
	json, err := json.Marshal(answer)
	if err != nil {
		panic(err)
	}
	return json
}

func isHash(hash string, file io.Reader) bool {
	hasher := sha512.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		panic(err)
	}
	if fmt.Sprintf("%x", hasher.Sum(nil)) == hash {
		return true
	}
	return false
}

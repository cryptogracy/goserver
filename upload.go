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
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	probableHash := mux.Vars(r)["hash"]

	if db.CheckPresence(probableHash) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		ret, err := json.Marshal(HttpReturn{"File already uploaded"})
		errorPanic(err)
		w.Write(ret)
		return
	}

	lifespan, err := strconv.Atoi(r.Header.Get("x-http-lifespan"))
	errorPanic(err)

	tempfile := path.Join(configuration.Tempdir, probableHash)

	out, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		ret, err := json.Marshal(HttpReturn{"Upload in Progress"})
		errorPanic(err)
		w.Write(ret)
		return
	}
	defer out.Close()
	defer os.Remove(tempfile)

	_, err = io.Copy(out, r.Body)
	errorPanic(err)

	// Check Hash
	hasher := sha512.New()
	_, err = out.Seek(0, os.SEEK_SET)
	errorPanic(err)
	_, err = io.Copy(hasher, out)
	errorPanic(err)

	generatedHash := fmt.Sprintf("%x", hasher.Sum(nil))
	if probableHash != generatedHash {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		ret, err := json.Marshal(HttpReturn{"Wrong hash"})
		errorPanic(err)
		w.Write(ret)
		return
	}

	err = os.Rename(tempfile, path.Join(configuration.Dir, generatedHash))

	now := time.Now()
	dur := time.Duration(lifespan) * time.Second
	db.AddFile(FileEntry{generatedHash, now, now.Add(dur), dur, ""})

	errorPanic(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	ret, err := json.Marshal(HttpReturn{})
	errorPanic(err)
	w.Write(ret)
}

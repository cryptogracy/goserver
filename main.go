package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

import (
	"github.com/gorilla/mux"
)

type HttpReturn struct {
	Error string
}

var configuration Configuration

func upload(w http.ResponseWriter, r *http.Request) {
	probableHash := mux.Vars(r)["hash"]
	tempfile := path.Join(configuration.Tempdir, probableHash)

	out, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		ret, err := json.Marshal(HttpReturn{"Upload in Progress"})
		if err != nil {
			panic(err)
		}
		w.Write(ret)
		return

	}
	defer out.Close()
	defer os.Remove(tempfile)

	_, err = io.Copy(out, r.Body)
	if err != nil {
		panic(err)
	}

	// Check Hash
	hasher := sha512.New()
	_, err = out.Seek(0, os.SEEK_SET)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(hasher, out)
	if err != nil {
		panic(err)
	}

	generatedHash := fmt.Sprintf("%x", hasher.Sum(nil))
	if probableHash != generatedHash {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		ret, err := json.Marshal(HttpReturn{"Wrong hash"})
		if err != nil {
			panic(err)
		}
		w.Write(ret)
		return
	}

	err = os.Rename(tempfile, path.Join(configuration.Dir, generatedHash))
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	ret, err := json.Marshal(HttpReturn{})
	if err != nil {
		panic(err)
	}
	w.Write(ret)
}

func routing() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/").Subrouter()
	api.PathPrefix("/").Handler(http.FileServer(http.Dir(configuration.Dir))).Methods("GET")
	api.HandleFunc("/{hash:[0-9|a-f]{128}}/", upload).Methods("PUT")

	router.PathPrefix("/ui/").
		Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("ui")))).
		Methods("GET")

	router.PathPrefix("/{hash:[0-9|a-f]{128}}/").
		Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("ui")))).
		Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("ui"))).Methods("GET")

	return router
}

func main() {

	configuration = readConfiguration()
	log.Println("Static from", configuration.Address)
	log.Println("Serving files from", configuration.Address)
	log.Println("Listen on", configuration.Address)

	log.Fatal(http.ListenAndServe(configuration.Address, routing()))
}

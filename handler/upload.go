package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/cryptogracy/goserver/db"
)

func Upload() http.Handler { return http.HandlerFunc(uploadFunc) }

func uploadFunc(w http.ResponseWriter, r *http.Request) {
	hash := mux.Vars(r)["hash"]
	w.Header().Set("Content-Type", "application/json") // We always return json

	// Get lifespan
	lifespan, err := strconv.Atoi(r.Header.Get("x-file-lifespan"))
	if err != nil {
		panic(err)
	}

	// File is already uploaded
	if err := db.AddFile(hash, lifespan, r.Body); err != nil {
		log.Println(err)
		if err == db.ErrFileExist {
			w.WriteHeader(http.StatusConflict)
			w.Write(jsonAnswer(map[string]string{"Error": "File already uploaded"}))
			return
		} else if err == db.ErrHash {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonAnswer(map[string]string{"Error": "Wrong Hash"}))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonAnswer(map[string]string{"Error": "Internal Error"}))
			return
		}
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

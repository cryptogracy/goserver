package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"path"
)

func routing() *mux.Router {
	router := mux.NewRouter()

	// Get upload ui from /
	router.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(configuration.Static, "upload.html"))
	}).Methods("GET")

	// Get download ui to download file with hash <hash> from /<hash>
	router.Path("/{hash:[0-9|a-f]{128}}").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join(configuration.Static, "download.html"))
		}).Methods("GET")

	// Get all files from configuration.Static from /<configuration.Static>
	router.PathPrefix(path.Join("/", configuration.Static)).
		Handler(http.FileServer(http.Dir("."))).Methods("GET")

	// Subrouter for /api/
	api := router.PathPrefix("/api/").Subrouter()

	api.PathPrefix("/").Handler(http.StripPrefix("/api/", http.FileServer(
		http.Dir(configuration.Dir)))).Methods("GET")

	// Put the file (with the lifespan in the header) to /api/<hash>
	api.HandleFunc("/{hash:[0-9|a-f]{128}}", Upload).Methods("PUT").
		HeadersRegexp("x-file-lifespan", "[0-9]+")

	// Push the metadata (with the lifespan in the header) to /api/metadata
	api.HandleFunc("/metadata", PushMeta).
		Methods("PUSH")

	// Get the metadata with the <id> from /api/metadata/<id>
	api.HandleFunc("/metadata/[[:alnum:]]",
		func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

	return router
}

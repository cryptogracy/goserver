package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"path"
)

func routing() *mux.Router {
	router := mux.NewRouter()

	// Get upload ui from /
	router.Path("/").Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join(configuration.Static, "upload.html"))
		})

	// Get download ui to download file with hash <hash> from /<hash>
	router.Path("/{hash:[0-9|a-f]{128}}").Methods("GET").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join(configuration.Static, "download.html"))
		})

	// Get all files from configuration.Static from /<configuration.Static>
	router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/",
		http.FileServer(http.Dir(configuration.Static))))

	router.PathPrefix("/api/").Methods("GET").Handler(http.StripPrefix("/api/",
		http.FileServer(http.Dir(configuration.Dir))))

	// Put the file (with the lifespan in the header) to /api/<hash>
	router.HandleFunc("/api/files/{hash:[0-9|a-f]{128}}", Upload).Methods("PUT").
		HeadersRegexp("x-file-lifespan", "[0-9]+")

	return router
}

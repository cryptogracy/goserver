package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"path"
)

func routing() *mux.Router {
	router := mux.NewRouter()

	ui(router)

  static(router)

	// Put the file (with the lifespan in the header) to /api/<hash>
	router.HandleFunc("/api/files/{hash:[0-9|a-f]{128}}", Upload).Methods("PUT").
		HeadersRegexp("x-file-lifespan", "[0-9]+")

	return router
}

func ui(router *mux.Router) {
	// Get upload ui from /
	upload := router.Path("/").Subrouter()
	upload.Methods("PUT", "PUSH", "PATCH", "DELETE").HandlerFunc(notAllowed)
	upload.Methods("GET").HandlerFunc(serveFile("upload.html"))

	// Get download ui to download ui from /<hash>
	download := router.Path("/{hash:[0-9|a-f]{128}}").Subrouter()
	download.Methods("PUT", "PUSH", "PATCH", "DELETE").HandlerFunc(notAllowed)
	download.Methods("GET").HandlerFunc(serveFile("download.html"))
}

func static(router *mux.Router) {
	// Get all files from configuration.Static from /ui/
	router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/",
		http.FileServer(http.Dir(configuration.Static))))

	// Get all files from configuration.Static from /ui/
  router.PathPrefix("/api/").Methods("GET").Handler(http.StripPrefix("/api/",
		http.FileServer(http.Dir(configuration.Dir))))
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusMethodNotAllowed)
}

func serveFile(file string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(configuration.Static, file))
	}
}

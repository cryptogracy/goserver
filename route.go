package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func routing() *mux.Router {
	router := mux.NewRouter()

	// Routes for API
	api := router.PathPrefix("/api/").Subrouter()
	api.PathPrefix("/").Handler(http.FileServer(http.Dir(configuration.Dir))).Methods("GET")
	api.HandleFunc("/{hash:[0-9|a-f]{128}}/", upload).Methods("PUT").HeadersRegexp("x-http-lifespan", "[1-6]0")

	// Routes for ui upload
	router.PathPrefix("/ui/").
		Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("ui")))).
		Methods("GET")

	// Routes for ui download
	router.PathPrefix("/{hash:[0-9|a-f]{128}}/").
		Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("ui")))).
		Methods("GET")

	// Deliver upload ui also under /
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("ui"))).Methods("GET")

	return router
}

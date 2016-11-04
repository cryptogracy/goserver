package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range Routes() {

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(Logger(route.Handler, route.Name))
	}

	router.Methods("GET").
		PathPrefix("/api/files/").
		Name("ServeFiles").
		Handler(Logger(http.StripPrefix("/api/files/",
			http.FileServer(http.Dir(configuration.Dir))), "ServeFiles"))

	router.Methods("GET").
		PathPrefix("/ui/").
		Name("ServeStatic").
		Handler(Logger(http.StripPrefix("/ui/",
			http.FileServer(http.Dir(configuration.Static))), "ServeStatic"))

	return router
}

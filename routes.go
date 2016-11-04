package main

import (
	"net/http"
	"path"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Dir struct {
	Name    string
	Pattern string
	Dir     string
}

func Routes() []Route {
	return []Route{
		Route{
			"UploadUi",
			"GET",
			"/",
			serveFile("upload.html"),
		},
		Route{
			"DownloadUi",
			"GET",
			"/{hash:[0-9|a-f]{128}",
			serveFile("Download.html"),
		},
		Route{
			"UploadFiles",
			"PUT",
			"/api/files/{hash:[0-9|a-f]{128}}",
			http.HandlerFunc(Upload),
		},
	}
}

func serveFile(file string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path.Join(configuration.Static, file))
		})
}

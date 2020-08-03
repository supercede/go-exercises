package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/supercede/go-exercises/url-shortner/urlshortner"
)

func main() {
	mux := defaultMux()

	format := flag.String("format", "json", "Choose File Type: json or yaml")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/handler-godoc": "https://godoc.org/github.com/gophercises/handler",
		"/yaml-godoc":    "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)

	handler := urlshortner.JSONHandler(mapHandler)

	if *format == "yaml" {
		handler = urlshortner.YAMLHandler(mapHandler)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

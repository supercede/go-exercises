package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/supercede/go-exercises/url-shortner/urlshortner"
)

func main() {
	mux := defaultMux()

	// defaults to links.json
	path := flag.String("format", "links.json", "Choose File Name: must end in json/yaml")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/handler-godoc": "https://godoc.org/github.com/gophercises/handler",
		"/yaml-godoc":    "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortner.MapHandler(pathsToUrls, mux)

	if !strings.HasSuffix(*path, ".json") && !strings.HasSuffix(*path, ".yaml") {
		log.Printf("File error: '%s' is not a valid json or yaml filename", *path)
		return
	}

	var handler http.HandlerFunc
	if strings.HasSuffix(*path, ".yaml") {
		handler = urlshortner.YAMLHandler(mapHandler, *path)
	}

	if strings.HasSuffix(*path, ".json") {
		handler = urlshortner.JSONHandler(mapHandler, *path)
	}

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Address not found")
}

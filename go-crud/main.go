package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/supercede/go-exercises/go-crud/books"
	"github.com/supercede/go-exercises/go-crud/data"
)

// var data []books.Book

func main() {
	path := flag.String("filename", "books.json", "Choose a storage file ending with .json")

	flag.Parse()

	if !strings.HasSuffix(*path, ".json") {
		log.Printf("File error: '%s' is not a valid json filename", *path)
		return
	}

	file := data.NewStore(*path)

	handler := books.Router(file)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

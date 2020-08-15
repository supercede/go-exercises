package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/supercede/go-exercises/go-crud/books"
	"github.com/supercede/go-exercises/go-crud/data"
)

// var data []books.Book

func main() {
	handler := books.Router()

	path := flag.String("filename", "books.json", "Choose a storage file ending with .json")

	flag.Parse()

	fmt.Println(*path)

	if !strings.HasSuffix(*path, ".json") {
		log.Printf("File error: '%s' is not a valid json filename", *path)
		return
	}

	file := data.NewStore(*path)

	fmt.Println(file)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

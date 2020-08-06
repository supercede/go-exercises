package main

import (
	"log"
	"net/http"

	"github.com/supercede/go-exercises/go-crud/books"
)

// var data []books.Book

func main() {
	handler := books.Router()

	log.Fatal(http.ListenAndServe(":8080", handler))
}

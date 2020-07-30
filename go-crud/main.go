package main

import (
	"log"
	"net/http"

	"github.com/supercede/go-exercises/go-crud/controllers"
	"github.com/supercede/go-exercises/go-crud/schemas"
)

var data []schemas.Book

func main() {
	http.HandleFunc("/books", controllers.GetBooks)
	http.HandleFunc("/add-book", controllers.CreateBook)
	http.HandleFunc("/get-book/", controllers.GetBook)
	http.HandleFunc("/delete-book/", controllers.DeleteBook)
	http.HandleFunc("/update-book/", controllers.UpdateBook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

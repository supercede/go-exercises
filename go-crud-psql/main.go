package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/supercede/go-exercises/go-crud-psql/books"
	"github.com/supercede/go-exercises/go-crud-psql/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}

	store := db.Init()
	store.DB.AutoMigrate(&db.Book{})

	handler := books.Router(store)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

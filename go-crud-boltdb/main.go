package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/supercede/go-exercises/go-crud-boltdb/books"
	"github.com/supercede/go-exercises/go-crud-boltdb/db"
)

func main() {
	path := flag.String("dbname", "books.db", "Choose a storage file ending with .db")

	home, err := homedir.Dir()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	flag.Parse()

	if !strings.HasSuffix(*path, ".db") {
		log.Printf("File error: '%s' is not a valid json filename", *path)
		return
	}

	dbPath := filepath.Join(home, *path)

	store, err := db.New(dbPath)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	handler := books.Router(store)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

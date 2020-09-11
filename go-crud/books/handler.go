package books

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/supercede/go-exercises/go-crud/models"

	"github.com/supercede/go-exercises/go-crud/data"
)

type Handler struct {
	store *data.Store
}

func NewHandler(s *data.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	store := h.store

	err := store.ReadFromFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var b models.Book

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Empty Request Body", http.StatusBadRequest)
			return
		}
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = validateBook(b)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := store.AddBook(b)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entry, err := toJSON(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, entry)
}

func (h *Handler) getBooks(w http.ResponseWriter, r *http.Request) {
	store := h.store

	err := store.ReadFromFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	books, err := toJSON(store.Books)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Fprint(w, books)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	store := h.store

	err := store.ReadFromFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	book, err := store.GetBook(id)

	if err != nil {
		http.Error(w, "Book Not found", 404)
		return
	}

	entry, err := toJSON(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, entry)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	store := h.store

	err := store.ReadFromFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	var b models.Book

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Empty Request Body", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := store.UpdateBook(id, b)

	if err != nil {
		http.Error(w, "Book Not found", 404)
		return
	}

	str, err := toJSON(book)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Fprint(w, str)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	store := h.store

	err := store.ReadFromFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	err = store.RemoveBook(id)

	if err != nil {
		http.Error(w, "Book Not found", 404)
		return
	}

	fmt.Fprintf(w, "Book with id %s deleted", id)
}

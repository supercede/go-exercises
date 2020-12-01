package books

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/supercede/go-exercises/go-crud-with-store/models"
	"github.com/supercede/go-exercises/go-crud-with-store/store"
)

type Handler struct {
	store *store.Store
}

func newHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	store := h.store
	var b models.Book

	err := json.NewDecoder(r.Body).Decode(&b)
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
	books, err := store.GetBooks()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	jsonBooks, err := toJSON(books)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Fprint(w, jsonBooks)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	store := h.store
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	book, err := store.GetBook(intID)
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
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
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

	book, err := store.UpdateBook(intID, b)
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
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = store.RemoveBook(intID)
	if err != nil {
		http.Error(w, "Book Not found", 404)
		return
	}

	fmt.Fprintf(w, "Book with id %s deleted", id)
}

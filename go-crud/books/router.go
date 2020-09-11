package books

import (
	"net/http"

	"github.com/supercede/go-exercises/go-crud/data"
)

func (h *Handler) singleBookRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBook(w, r)
		return
	case http.MethodPatch:
		h.updateBook(w, r)
		return
	case http.MethodDelete:
		h.deleteBook(w, r)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) allBooksRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBooks(w, r)
		return
	case http.MethodPost:
		h.createBook(w, r)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Router(s *data.Store) *http.ServeMux {
	handler := NewHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/books", handler.allBooksRouteHandler)
	mux.HandleFunc("/books/", handler.singleBookRouteHandler)
	return mux
}

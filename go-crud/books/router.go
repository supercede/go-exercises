package books

import (
	"net/http"
)

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBook(w, r)
		return
	case http.MethodPatch:
		updateBook(w, r)
		return
	case http.MethodDelete:
		deleteBook(w, r)
		return
	default:
		http.Error(w, "Method not alowed", http.StatusMethodNotAllowed)
	}
}

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/books/add", createBook)
	mux.HandleFunc("/books", getBooks)
	mux.HandleFunc("/books/", resourceHandler)
	return mux
}

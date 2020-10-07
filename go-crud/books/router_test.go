package books

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/supercede/go-exercises/go-crud/data"
	"github.com/supercede/go-exercises/go-crud/models"
)

func makeMapToSlice(m map[string]models.Book) []models.Book {
	v := make([]models.Book, 0, len(m))

	for _, value := range m {
		v = append(v, value)
	}

	return v
}

var s *data.Store = data.NewStore("test.json")
var h *Handler = newHandler(s)

func TestCreateBook(t *testing.T) {
	tt := []struct {
		name       string
		data       []byte
		statusCode int
	}{
		{"incomplete data", []byte(`{"name":"Adam's apple"}`), http.StatusBadRequest},
		{"incompatible type", []byte(`{"name":"Adam's apple", "author": "Bala Samuel", "pubData": {"month":"April","year": "2020"}}`), http.StatusBadRequest},
		{"complete data", []byte(`{"name":"Adam's apple", "author": "Bala Samuel", "pubData": {"month":"April","year": 2020}}`), http.StatusOK},
		{"Empty Request Body", []byte(``), http.StatusBadRequest},
	}
	// Edge cases

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/books", bytes.NewBuffer(tc.data))

			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			h.createBook(rec, req)
			res := rec.Result()
			t.Cleanup(func() { res.Body.Close() })

			if res.StatusCode != tc.statusCode {
				t.Errorf("Expected status %v, got %v", tc.statusCode, res.StatusCode)
			}
		})
	}
}

func TestGetBooks(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/books", nil)

	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()
	h.getBooks(rec, req)
	res := rec.Result()
	t.Cleanup(func() { res.Body.Close() })

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, res.StatusCode)
	}

	books, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}

	m := make(map[string]models.Book)
	err = json.Unmarshal(books, &m)

	iEquals := reflect.DeepEqual(m, s.Books)

	if !iEquals {
		t.Errorf("Expected %v to equal %v", m, s.Books)
	}

}

func TestGetBook(t *testing.T) {
	bookList := makeMapToSlice(s.Books)

	tt := []struct {
		name       string
		id         string
		statusCode int
	}{
		{"Existing ID", bookList[0].Id, http.StatusOK},
		{"Non-existing ID", "162723nfhfhrh47u", http.StatusNotFound},
		{"Empty ID", "", http.StatusNotFound},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			url := "http://localhost:8080/books/" + tc.id
			req, err := http.NewRequest(http.MethodGet, url, nil)

			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			h.getBook(rec, req)
			res := rec.Result()
			t.Cleanup(func() { res.Body.Close() })

			if res.StatusCode != tc.statusCode {
				t.Errorf("Expected status %v, got %v", tc.statusCode, res.StatusCode)
			}

			if res.StatusCode == 200 {
				book, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("Could not read response: %v", err)
				}

				var b models.Book
				err = json.Unmarshal(book, &b)
				iEquals := reflect.DeepEqual(b, s.Books[bookList[0].Id])

				if !iEquals {
					t.Errorf("Expected %v to equal %v", b, s.Books)
				}
			}

			if err != nil {
				t.Fatalf("Could not read response: %v", err)
			}
		})
	}
}

func TestPatchBook(t *testing.T) {
	bookList := makeMapToSlice(s.Books)

	tt := []struct {
		name       string
		id         string
		data       []byte
		statusCode int
	}{
		{"valid data", bookList[0].Id, []byte(`{"name":"Stole something"}`), http.StatusOK},
		{"Non existing Id", "1273bdbvfy44ui3", []byte(`{"name":"Stole something"}`), http.StatusNotFound},
		{"Invalid Data", bookList[0].Id, []byte(`{"pubData": {"month":"April","year": "Somalia"}}`), http.StatusBadRequest},
		{"Empty Request Body", bookList[0].Id, []byte(""), http.StatusBadRequest},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			url := "http://localhost:8080/books/" + tc.id
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(tc.data))

			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			h.updateBook(rec, req)
			res := rec.Result()
			t.Cleanup(func() { res.Body.Close() })

			if res.StatusCode != tc.statusCode {
				t.Errorf("Expected status %v, got %v", tc.statusCode, res.StatusCode)
			}

			if res.StatusCode == 200 {
				book, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("Could not read response: %v", err)
				}

				var b models.Book
				err = json.Unmarshal(book, &b)
				iEquals := reflect.DeepEqual(b, s.Books[bookList[0].Id])

				if !iEquals {
					t.Errorf("Expected %v to equal %v", b, s.Books)
				}
			}

			if err != nil {
				t.Fatalf("Could not read response: %v", err)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	bookList := makeMapToSlice(s.Books)

	tt := []struct {
		name       string
		id         string
		statusCode int
	}{
		{"valid Id", bookList[0].Id, http.StatusOK},
		{"Non existing Id", "1273bdbvfy44ui3", http.StatusNotFound},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			url := "http://localhost:8080/books/" + tc.id
			req, err := http.NewRequest(http.MethodDelete, url, nil)

			if err != nil {
				t.Fatalf("Could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			h.deleteBook(rec, req)
			res := rec.Result()
			t.Cleanup(func() { res.Body.Close() })

			if res.StatusCode != tc.statusCode {
				t.Errorf("Expected status %v, got %v", tc.statusCode, res.StatusCode)
			}

			if res.StatusCode == 200 {
				book, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("Could not read response: %v", err)
				}

				var b models.Book
				err = json.Unmarshal(book, &b)
				iEquals := reflect.DeepEqual(b, s.Books[bookList[0].Id])

				if !iEquals {
					t.Errorf("Expected %v to equal %v", b, s.Books)
				}
			}

			if err != nil {
				t.Fatalf("Could not read response: %v", err)
			}
		})
	}
}

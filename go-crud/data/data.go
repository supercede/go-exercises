package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/segmentio/ksuid"
	"github.com/supercede/go-exercises/go-crud/books"
	// "github.com/supercede/go-exercises/go-crud/books"
)

type Store struct {
	path string

	books map[string]books.Book
	mux   sync.RWMutex
}

func generateId() string {
	id := ksuid.New()
	return id.String()
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) WriteToFile() error {
	file, err := json.MarshalIndent(s.books, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.path, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ReadFromFile() error {
	rawData, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawData, &s.books)
	if string(rawData) == "null" {
		s.books = make(map[string]books.Book, 0)
	}
	if err != nil {
		return err
	}
	fmt.Println(s.books)
	return nil
}

func (s Store) AddBook(b books.Book) error {
	id := generateId()
	s.mux.Lock()
	s.books[id] = b
	s.mux.Unlock()

	err := s.WriteToFile()
	if err != nil {
		return err
	}
	return nil
}

func (s Store) RemoveBook(id string) error {
	_, ok := s.books[id]

	if ok {
		delete(s.books, id)
		return nil
	}
	return fmt.Errorf("Book with Id %s not found", id)
}

func (s Store) GetBook(id string) (books.Book, error) {
	_, ok := s.books[id]

	if ok {
		return s.books[id], nil
	}
	return books.Book{}, fmt.Errorf("Book with Id %s not found", id)
}

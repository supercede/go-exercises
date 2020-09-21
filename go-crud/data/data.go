package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/segmentio/ksuid"
	"github.com/supercede/go-exercises/go-crud/models"
	// "github.com/supercede/go-exercises/go-crud/books"
)

type Store struct {
	path string

	Books map[string]models.Book
	mu    *sync.RWMutex
}

func generateId() string {
	id := ksuid.New()
	return id.String()
}

// Creates a store that stores data in the given path
func NewStore(path string) *Store {
	return &Store{path: path, mu: &sync.RWMutex{}}
}

// Saves data to the store path
func (s *Store) WriteToFile() error {
	file, err := json.MarshalIndent(s.Books, "", "  ")
	if err != nil {
		return err
	}
	s.mu.RLock()
	err = ioutil.WriteFile(s.path, file, 0644)
	if err != nil {
		return err
	}
	s.mu.RUnlock()
	return nil
}

// Reads data from the store path
func (s *Store) ReadFromFile() error {
	_, err := os.Stat(s.path)

	if os.IsNotExist(err) {
		str, err := json.Marshal(s.Books)
		if err != nil {
			return err
		}

		s.mu.RLock()
		err = ioutil.WriteFile(s.path, str, 0644)
		if err != nil {
			return err
		}
		s.mu.RUnlock()
	}

	s.mu.RLock()
	rawData, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}
	s.mu.RUnlock()

	err = json.Unmarshal(rawData, &s.Books)
	if string(rawData) == "null" {
		s.Books = make(map[string]models.Book)
	}
	if err != nil {
		return err
	}

	return nil
}

// Adds a book to the store
func (s *Store) AddBook(b models.Book) (models.Book, error) {
	id := generateId()
	b.Id = id
	s.mu.Lock()
	s.Books[id] = b
	s.mu.Unlock()

	err := s.WriteToFile()
	if err != nil {
		return models.Book{}, err
	}
	return b, nil
}

// Remove a book from the store
func (s *Store) RemoveBook(id string) error {
	_, ok := s.Books[id]

	if ok {
		delete(s.Books, id)
		err := s.WriteToFile()
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Book with Id %s not found", id)
}

// Get a book
func (s *Store) GetBook(id string) (models.Book, error) {
	_, ok := s.Books[id]

	if ok {
		return s.Books[id], nil
	}
	return models.Book{}, fmt.Errorf("Book with Id %s not found", id)
}

// Edit a book
func (s *Store) UpdateBook(id string, b models.Book) (models.Book, error) {
	_, ok := s.Books[id]

	if ok {
		existingBook := s.Books[id]

		if b.Author != "" {
			existingBook.Author = b.Author
		}
		if b.Name != "" {
			existingBook.Name = b.Name
		}
		if b.PubData.Month != "" {
			existingBook.PubData.Month = b.PubData.Month
		}
		if b.PubData.Year != 0 {
			existingBook.PubData.Year = b.PubData.Year
		}

		s.Books[id] = existingBook

		err := s.WriteToFile()

		if err != nil {
			return models.Book{}, err
		}

		return s.Books[id], nil
	}
	return models.Book{}, fmt.Errorf("Book with Id %s not found", id)
}

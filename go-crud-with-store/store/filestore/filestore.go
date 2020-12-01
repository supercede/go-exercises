package filestore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/supercede/go-exercises/go-crud-with-store/models"
)

type FileStore struct {
	path string

	Books map[int]models.Book
	mu    *sync.RWMutex
}

func generateID() int {
	id := time.Now().UnixNano() / (1 << 22)
	return int(id)
}

// NewStore creates a store that stores data in the given path
func NewStore(path string) *FileStore {
	return &FileStore{path: path, mu: &sync.RWMutex{}}
}

// WriteToFile saves data to the store path
func (s *FileStore) WriteToFile() error {
	file, err := json.MarshalIndent(s.Books, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to parse into JSON format")
	}
	s.mu.RLock()
	err = ioutil.WriteFile(s.path, file, 0644)
	if err != nil {
		return errors.Wrap(err, "Error saving to file")
	}
	s.mu.RUnlock()
	return nil
}

// ReadFromFile reads data from the store path
func (s *FileStore) ReadFromFile() error {
	_, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		str, err := json.Marshal(s.Books)
		if err != nil {
			return errors.Wrap(err, "Error reading from file")
		}

		s.mu.RLock()
		err = ioutil.WriteFile(s.path, str, 0644)
		if err != nil {
			return errors.Wrap(err, "Error saving to file")
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
		s.Books = make(map[int]models.Book)
	}
	if err != nil {
		return err
	}

	return nil
}

// AddBook adds a book to the store
func (s *FileStore) AddBook(b models.Book) (models.Book, error) {
	id := generateID()
	b.ID = id
	s.mu.Lock()
	s.Books[id] = b
	s.mu.Unlock()
	return b, nil
}

// RemoveBook removes a book with a given ID from the store
func (s *FileStore) RemoveBook(id int) error {
	_, ok := s.Books[id]

	if ok {
		delete(s.Books, id)
		return nil
	}
	return fmt.Errorf("Book with Id %v not found", id)
}

// GetBook gets a book with a given ID from the store
func (s *FileStore) GetBook(id int) (models.Book, error) {
	_, ok := s.Books[id]

	if ok {
		return s.Books[id], nil
	}
	return models.Book{}, fmt.Errorf("Book with Id %v not found", id)
}

// UpdateBook updates a book with a given ID from the store
func (s *FileStore) UpdateBook(id int, b models.Book) (models.Book, error) {
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
		return s.Books[id], nil
	}
	return models.Book{}, fmt.Errorf("Book with Id %v not found", id)
}

func (s *FileStore) GetBooks() ([]models.Book, error) {
	books := make([]models.Book, 0, len(s.Books))

	for _, book := range s.Books {
		books = append(books, book)
	}

	return books, nil
}

package store

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/supercede/go-exercises/go-crud-with-store/models"
	"github.com/supercede/go-exercises/go-crud-with-store/store/filestore"
	"github.com/supercede/go-exercises/go-crud-with-store/store/shared"
	"github.com/supercede/go-exercises/go-crud-with-store/util"
)

type Store struct {
	Storage shared.Storage
}

func New() (*Store, error) {
	var s shared.Storage
	conf, err := util.GetConfig()
	if err != nil {
		return &Store{}, errors.Wrap(err, "Failed to Read config file")
	}
	switch db := conf.DatabaseType; db {
	case "filestore":
		path := conf.FileDBPath

		if !strings.HasSuffix(path, ".json") {
			return nil, errors.New(fmt.Sprintf("File error: '%s' is not a valid json filename", path))
		}

		s = filestore.NewStore(path)
	default:
		return nil, errors.New("Invalid Database type")
	}

	return &Store{
		Storage: s,
	}, nil
}

func (s *Store) AddBook(b models.Book) (models.Book, error) {
	return s.Storage.AddBook(b)
}

func (s *Store) RemoveBook(id int) error {
	return s.Storage.RemoveBook(id)
}

func (s *Store) UpdateBook(id int, b models.Book) (models.Book, error) {
	return s.Storage.UpdateBook(id, b)
}

func (s *Store) GetBook(id int) (models.Book, error) {
	return s.Storage.GetBook(id)
}

func (s *Store) GetBooks() map[int]models.Book {
	return s.Storage.GetBooks()
}

func (s *Store) ReadFromFile() error {
	return s.Storage.ReadFromFile()
}

func (s *Store) WriteToFile() error {
	return s.Storage.WriteToFile()
}

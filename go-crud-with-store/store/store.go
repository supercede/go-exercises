package store

import (
	"fmt"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/supercede/go-exercises/go-crud-with-store/models"
	"github.com/supercede/go-exercises/go-crud-with-store/store/boltdb"
	"github.com/supercede/go-exercises/go-crud-with-store/store/filestore"
	"github.com/supercede/go-exercises/go-crud-with-store/store/shared"
	"github.com/supercede/go-exercises/go-crud-with-store/util"
)

type Store struct {
	Storage shared.Storage
}

func New() (*Store, error) {
	var s shared.Storage
	var err error
	var home string

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
	case "boltdb":
		home, err = homedir.Dir()
		dbName := conf.BoltDBName

		if !strings.HasSuffix(dbName, ".db") {
			return nil, errors.New(fmt.Sprintf("DB error: '%s' should end in .db", dbName))
		}

		s, err = boltdb.New(filepath.Join(home, dbName))
	default:
		return nil, errors.New("Invalid Database type")
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not initialize the data backend")
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

func (s *Store) GetBooks() ([]models.Book, error) {
	return s.Storage.GetBooks()
}

func (s *Store) ReadFromFile() error {
	return s.Storage.ReadFromFile()
}

func (s *Store) WriteToFile() error {
	return s.Storage.WriteToFile()
}

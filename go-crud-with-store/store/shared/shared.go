package shared

import (
	"github.com/supercede/go-exercises/go-crud-with-store/models"
)

type Storage interface {
	AddBook(b models.Book) (models.Book, error)
	RemoveBook(id int) error
	GetBook(id int) (models.Book, error)
	UpdateBook(id int, b models.Book) (models.Book, error)
	GetBooks() ([]models.Book, error)
	ReadFromFile() error
	WriteToFile() error
}

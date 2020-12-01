package psql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"github.com/supercede/go-exercises/go-crud-with-store/models"
)

type PostgresStore struct {
	DB *gorm.DB
}

func New(user, pass, host, name string) (*PostgresStore, error) {
	credentials := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, name, pass)

	db, err := gorm.Open("postgres", credentials)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}
	db.AutoMigrate(&models.Book{})
	return &PostgresStore{
		DB: db,
	}, nil
}

func (s *PostgresStore) AddBook(b models.Book) (models.Book, error) {
	result := s.DB.Create(&b)
	if result.Error != nil {
		return models.Book{}, result.Error
	}

	return b, nil
}

func (s *PostgresStore) GetBooks() ([]models.Book, error) {
	var books []models.Book

	result := s.DB.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

func (s *PostgresStore) GetBook(id int) (models.Book, error) {
	var book models.Book

	result := s.DB.First(&book, id)
	if result.Error != nil {
		return models.Book{}, result.Error
	}

	return book, nil
}

func (s *PostgresStore) RemoveBook(id int) error {
	var book models.Book
	result := s.DB.First(&book, id)
	if result.Error != nil {
		return result.Error
	}
	result = s.DB.Delete(book)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *PostgresStore) UpdateBook(id int, b models.Book) (models.Book, error) {
	var book models.Book
	result := s.DB.First(&book, id)
	if result.Error != nil {
		return models.Book{}, result.Error
	}

	// result.Updates(Book{b})
	update := s.DB.Model(&book).Updates(b)

	if update.Error != nil {
		return models.Book{}, update.Error
	}

	return book, nil
}

func (b *PostgresStore) ReadFromFile() error {
	return nil
}

func (b *PostgresStore) WriteToFile() error {
	return nil
}

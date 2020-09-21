package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Store struct {
	DB *gorm.DB
}

func Init() *Store {
	user := goDotEnvVariable("DB_USER")
	pass := goDotEnvVariable("DB_PASS")
	host := goDotEnvVariable("DB_HOST")
	name := goDotEnvVariable("DB_NAME")

	credentials := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", host, user, name, pass)

	db, err := gorm.Open("postgres", credentials)
	fmt.Println(credentials)
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	fmt.Println("Done")
	return &Store{
		DB: db,
	}
}

func (s *Store) AddBook(b Book) (*Book, error) {
	result := s.DB.Create(&b)
	if result.Error != nil {
		return nil, result.Error
	}

	return &b, nil
}

func (s *Store) AllBooks() ([]Book, error) {
	var books []Book

	result := s.DB.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}

func (s *Store) GetBook(id int) (*Book, error) {
	var book Book

	result := s.DB.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &book, nil
}

func (s *Store) RemoveBook(id int) error {
	var book Book
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

func (s *Store) UpdateBook(id int, b Book) (*Book, error) {
	var book Book
	result := s.DB.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// result.Updates(Book{b})
	update := s.DB.Model(&book).Updates(b)

	if update.Error != nil {
		return nil, update.Error
	}

	return &book, nil
}

func goDotEnvVariable(key string) string {
	// load .env file
	val, exist := os.LookupEnv(key)

	if !exist {
		log.Fatalf("Error loading %s from .env file", key)
	}

	return val
}

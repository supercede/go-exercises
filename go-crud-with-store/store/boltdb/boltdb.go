package boltdb

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/supercede/go-exercises/go-crud-with-store/models"
)

type BoltStore struct {
	bucket []byte
	db     *bolt.DB
}

func New(dbPath string) (*BoltStore, error) {
	bucket := []byte("book")
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, errors.Wrap(err, "could not open bolt DB database")
	}
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return errors.Wrapf(err, "could not create %s bucket", bucket)
		}
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not create buckets")
	}
	return &BoltStore{
		bucket: bucket,
		db:     db,
	}, nil
}

func (b *BoltStore) GetBook(id int) (models.Book, error) {
	var raw []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		raw = tx.Bucket(b.bucket).Get([]byte(itob(id)))
		if raw == nil {
			return errors.New("No Book found with this ID")
		}
		return nil
	})
	if err != nil {
		return models.Book{}, errors.Wrap(err, "could not view db")
	}
	var book *models.Book
	return *book, json.Unmarshal(raw, &book)
}

func (b *BoltStore) AddBook(book models.Book) (models.Book, error) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		id64, _ := bucket.NextSequence()
		book.ID = int(id64)

		buf, err := json.Marshal(book)
		if err != nil {
			return errors.Wrap(err, "could not marshal entry")
		}
		if err := bucket.Put(itob(book.ID), buf); err != nil {
			return errors.Wrap(err, "could not put data into bucket")
		}
		return nil
	})
	if err != nil {
		return models.Book{}, errors.Wrap(err, "could not update db")
	}

	return book, nil
}

func (b *BoltStore) GetBooks() ([]models.Book, error) {
	var books []models.Book
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		c := bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var book *models.Book
			err := json.Unmarshal(v, &book)

			if err != nil {
				return errors.Wrap(err, "invalid result type")
			}
			books = append(books, *book)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b *BoltStore) UpdateBook(ID int, book models.Book) (models.Book, error) {
	var updatedBook *models.Book

	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		raw := bucket.Get([]byte(itob(ID)))
		if raw == nil {
			return errors.New("No Book found with this ID")
		}
		err := json.Unmarshal(raw, &updatedBook)
		if err != nil {
			return errors.New("Invalid Data")
		}

		if book.Author != "" {
			updatedBook.Author = book.Author
		}
		if book.Name != "" {
			updatedBook.Name = book.Name
		}
		if book.PubData.Month != "" {
			updatedBook.PubData.Month = book.PubData.Month
		}
		if book.PubData.Year != 0 {
			updatedBook.PubData.Year = book.PubData.Year
		}

		buf, err := json.Marshal(updatedBook)
		if err != nil {
			return errors.Wrap(err, "could not marshal entry")
		}

		if err := bucket.Put(itob(updatedBook.ID), buf); err != nil {
			return errors.Wrap(err, "could not put data into bucket")
		}
		return nil
	})
	if err != nil {
		return models.Book{}, errors.Wrap(err, "could not view db")
	}
	return *updatedBook, nil
}

func (b *BoltStore) RemoveBook(ID int) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		return bucket.Delete(itob(ID))
	})
}

func (b *BoltStore) ReadFromFile() error {
	return nil
}

func (b *BoltStore) WriteToFile() error {
	return nil
}

// integer to binary
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

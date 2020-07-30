package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/supercede/go-exercises/go-crud/schemas"
)

func ValidateBook(r schemas.Book) (bool, error) {
	switch true {
	case r.Name == "":
		return false, fmt.Errorf("Name is required")
	case r.Author == "":
		return false, fmt.Errorf("Author is required")
	case r.PubData.Month == "":
		return false, fmt.Errorf("Publication month is required")
	case r.PubData.Year == 0:
		return false, fmt.Errorf("Publication Year is required")
	default:
		return true, nil
	}
}

func ToJson(entry interface{}) (string, error) {
	b, err := json.Marshal(entry)
	if err != nil {
		return "Error", err
	}

	return string(b), nil
}

func Remove(s []schemas.Book, i int) []schemas.Book {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func WriteFile(data []schemas.Book, entry schemas.Book) {
	data = append(data, entry)
	file, _ := json.MarshalIndent(data, "", "  ")
	_ = ioutil.WriteFile("books.json", file, 0644)
}

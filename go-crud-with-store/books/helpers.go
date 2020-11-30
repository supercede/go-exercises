package books

import (
	"encoding/json"
	"fmt"

	"github.com/supercede/go-exercises/go-crud-with-store/models"
)

func validateBook(r models.Book) (bool, error) {
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

func toJSON(entry interface{}) (string, error) {
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return "Error", err
	}

	return string(b), nil
}

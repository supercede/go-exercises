package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/supercede/go-exercises/go-crud/helpers"
	"github.com/supercede/go-exercises/go-crud/schemas"
)

var data []schemas.Book

func LoadBookFile() (bool, error) {
	_, err := os.Stat("books.json")

	if os.IsNotExist(err) {
		str, _ := json.Marshal(data)
		_ = ioutil.WriteFile("books.json", str, 0644)
	}

	rawData, err := ioutil.ReadFile("books.json")

	if err != nil {
		fmt.Println(err)
		return err != nil, err
	}

	// fmt.Print("data:  ", string(rawData))
	err = json.Unmarshal(rawData, &data)
	if string(rawData) == "null" {
		data = make([]schemas.Book, 0)
	}
	if err != nil {
		fmt.Println(err)
		return err != nil, err
	}

	return err != nil, err
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	LoadBookFile()

	var b schemas.Book
	// fmt.Println(r)
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = helpers.ValidateBook(b)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data) == 0 {
		b.Id = 1
	} else {
		b.Id = data[len(data)-1].Id + 1
	}

	helpers.WriteFile(data, b)

	entry, err := helpers.ToJson(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, entry)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	// _, err := loadBookFile()
	_, err := LoadBookFile()

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	books, err := helpers.ToJson(data)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Fprintf(w, books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	LoadBookFile()

	id := strings.TrimPrefix(r.URL.Path, "/get-book/")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, value := range data {
		if value.Id == intID {
			book, err := helpers.ToJson(value)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Fprintf(w, book)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	LoadBookFile()

	id := strings.TrimPrefix(r.URL.Path, "/update-book/")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Id provided", http.StatusBadRequest)
		return
	}

	var b schemas.Book

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range data {
		if data[i].Id == intID {
			if b.Author != "" {
				data[i].Author = b.Author
			}
			if b.Name != "" {
				data[i].Name = b.Name
			}
			if b.PubData.Month != "" {
				data[i].PubData.Month = b.PubData.Month
			}
			if b.PubData.Year != 0 {
				data[i].PubData.Year = b.PubData.Year
			}

			dataStr, _ := json.MarshalIndent(data, "", "  ")
			ioutil.WriteFile("books.json", dataStr, 0644)

			book, err := helpers.ToJson(data[i])

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Fprintf(w, book)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	LoadBookFile()

	id := strings.TrimPrefix(r.URL.Path, "/delete-book/")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for i, value := range data {
		if value.Id == intID {
			data = helpers.Remove(data, i)
			dataStr, _ := json.MarshalIndent(data, "", "  ")
			ioutil.WriteFile("books.json", dataStr, 0644)
			fmt.Fprintf(w, "Book with id %s deleted", id)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

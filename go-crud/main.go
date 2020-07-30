package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type pub struct {
	Year  float64
	Month string
}

type book struct {
	Id      int
	Name    string
	Author  string
	PubData pub
}

var data []book

func validateBook(r book) (bool, error) {
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

func createBook(w http.ResponseWriter, r *http.Request) {
	loadBookFile()

	var b book
	// fmt.Println(r)
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = validateBook(b)

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

	writeFile(b)

	entry, err := toJson(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, entry)
}

func loadBookFile() (bool, error) {
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
		data = make([]book, 0)
	}
	if err != nil {
		fmt.Println(err)
		return err != nil, err
	}
	// fmt.Printf("slice: %+v\n", data)

	return err != nil, err
}

func writeFile(entry book) {
	data = append(data, entry)
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("books.json", file, 0644)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	loadBookFile()

	// id := r.URL.Query().Get("id")

	id := strings.TrimPrefix(r.URL.Path, "/get-book/")

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for _, value := range data {
		if value.Id == intId {
			book, err := toJson(value)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Fprintf(w, book)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	loadBookFile()

	// id := r.URL.Query().Get("id")

	id := strings.TrimPrefix(r.URL.Path, "/update-book/")

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Id provided", http.StatusBadRequest)
		return
	}

	var b book
	// fmt.Println(r)
	// json.NewDecoder(r.Body).Decode(&b)

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(b)

	for i, _ := range data {
		if data[i].Id == intId {
			// fmt.Println("initial value:", data)
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
			// fmt.Println("reviewed value:", data)

			dataStr, _ := json.MarshalIndent(data, "", " ")
			ioutil.WriteFile("books.json", dataStr, 0644)

			book, err := toJson(data[i])

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			fmt.Fprintf(w, book)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func remove(s []book, i int) []book {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	loadBookFile()

	id := strings.TrimPrefix(r.URL.Path, "/delete-book/")

	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	for i, value := range data {
		if value.Id == intId {
			data = remove(data, i)
			// writeFile(data)
			dataStr, _ := json.MarshalIndent(data, "", " ")
			ioutil.WriteFile("books.json", dataStr, 0644)
			fmt.Fprintf(w, "Book with id %s deleted", id)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func toJson(entry interface{}) (string, error) {
	b, err := json.Marshal(entry)
	if err != nil {
		return "Error", err
	}

	return string(b), nil
}

func main() {
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		_, err := loadBookFile()

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}

		// books := data
		books, err := toJson(data)

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}

		fmt.Fprintf(w, books)
	})
	http.HandleFunc("/add-book", createBook)
	http.HandleFunc("/get-book/", getBook)
	http.HandleFunc("/delete-book/", deleteBook)
	http.HandleFunc("/update-book/", updateBook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

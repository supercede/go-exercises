package books

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var data []book

func loadBookFile() error {
	_, err := os.Stat("books.json")

	if os.IsNotExist(err) {
		str, err := json.Marshal(data)

		if err != nil {
			log.Println(err)
			return err
		}
		err = ioutil.WriteFile("books.json", str, 0644)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	rawData, err := ioutil.ReadFile("books.json")

	if err != nil {
		log.Println(err)
		return err
	}

	// fmt.Print("data:  ", string(rawData))
	err = json.Unmarshal(rawData, &data)
	if string(rawData) == "null" {
		data = make([]book, 0)
	}
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := loadBookFile()

		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var b book

		err = json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			if err == io.EOF {
				http.Error(w, "Empty Request Body", http.StatusBadRequest)
				return
			}
			log.Println(err.Error())
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

		writeFile(data, b)

		entry, err := toJSON(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, entry)
		return
	}
	http.Error(w, "Method not alowed", http.StatusMethodNotAllowed)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	// _, err := loadBookFile()
	err := loadBookFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	books, err := toJSON(data)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	fmt.Fprintf(w, books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	err := loadBookFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	intID, err := strconv.Atoi(id)
	if err != nil {
		str := fmt.Sprintf("Cannot convert id: '%s' into integer", id)
		http.Error(w, str, http.StatusBadRequest)
		return
	}

	for _, value := range data {
		if value.Id == intID {
			book, err := toJSON(value)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, book)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	err := loadBookFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid Id provided", http.StatusBadRequest)
		return
	}

	var b book

	err = json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		if err == io.EOF {
			http.Error(w, "Empty Request Body", http.StatusBadRequest)
			return
		}
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

			dataStr, err := json.MarshalIndent(data, "", "  ")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ioutil.WriteFile("books.json", dataStr, 0644)

			book, err := toJSON(data[i])

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, book)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	err := loadBookFile()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/books/")

	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, value := range data {
		if value.Id == intID {
			data = remove(data, i)
			dataStr, err := json.MarshalIndent(data, "", "  ")

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ioutil.WriteFile("books.json", dataStr, 0644)
			fmt.Fprintf(w, "Book with id %s deleted", id)
			return
		}
	}
	http.Error(w, "Book Not found", 404)
}

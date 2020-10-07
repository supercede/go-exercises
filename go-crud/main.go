package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/supercede/go-exercises/go-crud/books"
	"github.com/supercede/go-exercises/go-crud/data"
)

func main() {
	path := flag.String("filename", "books.json", "Choose a storage file ending with .json")

	flag.Parse()

	if !strings.HasSuffix(*path, ".json") {
		log.Printf("File error: '%s' is not a valid json filename", *path)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	file := data.NewStore(*path)
	file.ReadFromFile()

	// ticker := time.NewTicker(5 * time.Second)

	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			fmt.Println("here")
			file.WriteToFile()
		}
	}()

	handler := books.Router(file)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		file.WriteToFile()
		// os.Exit(0)
		tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(tc)
	}()

	// handler := books.Router(file)

	// s := &http.Server{
	// 	Addr:         ":8080",
	// 	Handler:      handler,
	// 	WriteTimeout: 2 * time.Second,
	// }
	// go func() {
	log.Fatal(s.ListenAndServe())
	// }()

	// Using this, the program doesn't shut down on first try
	// for {
	// select {
	// case <-ticker.C:
	// 	fmt.Println("here")
	// 	file.WriteToFile()
	// case <-sigs:
	// 	sig := <-sigs
	// 	fmt.Println(sig)
	// 	file.WriteToFile()
	// 	os.Exit(0)
	// 	// case
	// }
	// }
}

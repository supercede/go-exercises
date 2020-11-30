package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/supercede/go-exercises/go-crud-with-store/books"
	"github.com/supercede/go-exercises/go-crud-with-store/store"
	"github.com/supercede/go-exercises/go-crud-with-store/util"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}

	conf, err := util.GetConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	database := conf.DatabaseType

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	newStore, err := store.New()
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}
	if database == "filestore" {
		newStore.ReadFromFile()
	}

	// ticker := time.NewTicker(5 * time.Second)

	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			fmt.Println("here")

			if database == "filestore" {
				newStore.WriteToFile()
			}
		}
	}()

	handler := books.Router(newStore)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		if database == "filestore" {
			newStore.WriteToFile()
		}
		// os.Exit(0)
		tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(tc)
	}()

	log.Fatal(s.ListenAndServe())
}

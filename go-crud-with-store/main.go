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
	"github.com/pkg/errors"
	"github.com/supercede/go-exercises/go-crud-with-store/books"
	"github.com/supercede/go-exercises/go-crud-with-store/store"
	"github.com/supercede/go-exercises/go-crud-with-store/util"
)

func main() {
	conf, err := getConfig()
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

	if database == "filestore" {
		go func() {
			ticker := time.NewTicker(5 * time.Second)

			for range ticker.C {
				fmt.Println("here")

				newStore.WriteToFile()
			}
		}()
	}

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
		tc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(tc)
	}()

	log.Fatal(s.ListenAndServe())
}

func getConfig() (util.EnvVariables, error) {
	if err := godotenv.Load(); err != nil {
		return util.EnvVariables{}, errors.Wrap(err, "failed to load env file")
	}

	conf, err := util.GetConfig()
	if err != nil {
		return util.EnvVariables{}, errors.Wrap(err, "failed to load config vars")
	}

	return conf, nil
}

// func initApp() () {}

# go-crud

## Introduction

A go CRUD application that saves data to and reads from a json file

## Requirements

- [Go](https://golang.org) - v1.11 above

## Installation

```
go get github.com/supercede/go-exercises/go-crud
```

## Build & Run

```
cd $GOPATH/src/github.com/supercede/go-exercises/go-crud

go build .
./go-crud
```

## Routes

- `/books` - Get all books
- `/add-book` - Add a book
- `/get-book/{id}` - Get book with given path Id
- `/delete-book/{id}` - delete book with given path Id
- `/update-book/{id}` - Update book with given path Id

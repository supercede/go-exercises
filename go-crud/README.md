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

- `GET` `/books` - Get all books
- `POST` `/books` - Add a book
- `GET` `/books/{id}` - Get book with given path Id
- `DELETE` `/books/{id}` - delete book with given path Id
- `PATCH` `/books/{id}` - Update book with given path Id

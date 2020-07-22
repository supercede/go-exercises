# csv-quiz : A Go CSV quiz parser

## Overview [![GoDoc](https://godoc.org/github.com/supercede/go-exercises/csv-quiz?status.svg)](https://godoc.org/github.com/supercede/go-exercises/csv-quiz)

A command-line go application that parses quizzes from a CSV file. The quiz file defaults to problems.csv, however, the filename can be changed/edited using the `--csv` flag while running the application.

## Install

```
go get github.com/supercede/go-exercises/csv-quiz
```

## Build & Run

```
cd $GOPATH/src/github.com/supercede/go-exercises/csv-quiz

go build .
./csv-quiz (--csv filename.csv)
```

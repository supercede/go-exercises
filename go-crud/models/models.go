package models

type Pub struct {
	Year  float64
	Month string
}

type Book struct {
	Id      string
	Name    string
	Author  string
	PubData Pub
}

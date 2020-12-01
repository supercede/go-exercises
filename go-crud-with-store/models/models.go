package models

type Pub struct {
	Year  float64
	Month string
}

type Book struct {
	ID      int
	Name    string
	Author  string
	PubData Pub
}

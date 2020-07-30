package schemas

type Pub struct {
	Year  float64
	Month string
}

type Book struct {
	Id      int
	Name    string
	Author  string
	PubData Pub
}

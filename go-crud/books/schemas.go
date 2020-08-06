package books

type pub struct {
	Year  float64
	Month string
}

type book struct {
	Id      int
	Name    string
	Author  string
	PubData pub
}

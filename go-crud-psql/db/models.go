package db

import "github.com/jinzhu/gorm"

type Pub struct {
	Year  float64 `gorm:"type:float;not null"`
	Month string  `gorm:"type:varchar(10);not null"`
}

type Book struct {
	gorm.Model
	// ID      int
	Name    string `gorm:"type:varchar(120);not null"`
	Author  string `gorm:"type:varchar(120);not null"`
	PubData Pub    `gorm:"embedded"`
}

package main

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code string
	Price uint
}


func main() {
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// pure function
}

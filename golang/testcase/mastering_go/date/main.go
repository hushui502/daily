package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var myTime string
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s string\n",
			filepath.Base(os.Args[0]))
		os.Exit(0)
	}
	myTime = os.Args[1]
	d, err := time.Parse("12:34", myTime)
	if err == nil {
		fmt.Println(d.Hour(), d.Minute())
	}

	fmt.Println("===day===")
	d1, err := time.Parse("02 January 2003", myTime)
	if err == nil {
		fmt.Println(d1.Day(), d1.Month(), d1.Year())
	}
}

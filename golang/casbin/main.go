package main

import (
	"fmt"
	"log"

	"github.com/mitchellh/go-homedir"
)

func main() {
	homedir.DisableCache = false

	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	homedir.Expand("~/goland")
	fmt.Println("Home dir:", dir)

}
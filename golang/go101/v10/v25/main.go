package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	upperDirpattern := "./*"

	matches, err := filepath.Glob(upperDirpattern)
	if err != nil {
		panic(err)
	}
	for _, file := range matches {
		fmt.Println(file)
	}
}

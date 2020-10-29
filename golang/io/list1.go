package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	r, _ := http.Get(os.Args[1])
	file, _ := os.Create(os.Args[2])
	defer file.Close()
	defer r.Body.Close()
	dest := io.MultiWriter(os.Stdout, file)
	io.Copy(dest, r.Body)
}

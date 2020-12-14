package main

import (
	"net/http"
	"os"
	"strings"
)

func sharDir(dirName string, port string, ch chan bool) {
	h := http.FileServer(http.Dir(dirName))
	err := http.ListenAndServe(":"+port, h)
	if err != nil {
		ch <- false
	}
}

func main() {
	ch := make(chan bool)
	port := "8000"
	if len(os.Args) > 1 {
		port = strings.Join(os.Args[1:2], "")
	}
	go sharDir(".", port, ch)
	res := <-ch
	if res == false {
		//
	}
}

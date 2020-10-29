package main

import (
	"database/sql"
	"io"
	"net/http"
)

func SayHello(w http.ResponseWriter, r *http.Request) {
	w.Header()["Content-Type"] = []string{"application/json"}
	io.WriteString(w, "hello")
}

func main() {
	http.HandleFunc("/", SayHello)
	http.ListenAndServe(":9090", nil)
}

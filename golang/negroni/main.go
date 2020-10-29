package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "welcome new page")
	})

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3030", n)
}

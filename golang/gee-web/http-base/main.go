package main

import (
	"awesomeProject2/gee-web/http-base/gee"
	"fmt"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.PATH = %q\n", r.URL.Path)
	})
	r.Run(":9999")
}

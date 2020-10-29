package main

import "net/http"

func helloHandler(wr http.ResponseWriter, req *http.Request) {

}

func main() {
	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func WithServerHeader(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithServerHeader")
		w.Header().Set("Server", "HelloServer")
		h(w, r)
	}
}

func WithAuthCookie(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithAuthCookie")
		cookie := &http.Cookie{Name: "Auth", Value: "Pass", Path: "/"}
		http.SetCookie(w, cookie)
		h(w, r)
	}
}

func WithBasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithBasicAuth")
		cookie, err := r.Cookie("Auth")
		if err != nil || cookie.Value != "Pass" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h(w, r)
	}
}

func WithDebugLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("--->WithDebugLog")
		r.ParseForm()
		log.Println(r.Form)
		log.Println("path", r.URL.Path)
		log.Println("schema", r.URL.Scheme)
		log.Println(r.Form["url_long"])
		for k, v := range r.Form {
			log.Println("key: ", k)
			log.Println(v, strings.Join(v, ""))
		}
		h(w, r)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved Request %s from %s\n", r.URL.Path, r.RemoteAddr)
	fmt.Fprintf(w, "Hello, World! "+r.URL.Path)
}

type HttpHandlerDecorator func(http.HandlerFunc) http.HandlerFunc

func Handler(h http.HandlerFunc, decors ...HttpHandlerDecorator) http.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i]
		h = d(h)
	}
	return h
}

func main() {
	// pipleline
	http.HandleFunc("/v1/hello", Handler(hello,
		WithAuthCookie, WithBasicAuth, WithDebugLog, WithServerHeader))
}

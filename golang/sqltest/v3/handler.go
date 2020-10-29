package main

import (
	"context"
	"fmt"
	"net/http"
)

type RequestContextKey string

func requestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestContextKey("id"), "uuid")
		next(w, r.WithContext(ctx))
	}
}

func test1(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(RequestContextKey("id"))
	w.Write([]byte("reuqest_id = " + id.(string)))
}

func test2(w http.ResponseWriter, r *http.Request) {
	id := fmt.Sprintf("%v", r.Context().Value(RequestContextKey("id")))
	w.Write([]byte("reuqest_id = " + id))
}

func main() {
	http.Handle("/test1", requestID(test1))
	http.HandleFunc("/test2", test2)
	http.ListenAndServe(":8080", nil)
}
package main

import (
"context"
"github.com/julienschmidt/httprouter"
"net/http"
)

type RouterContextKey string

func compatible(next http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(r.Context(), RouterContextKey("params"), p)
		next(w, r.WithContext(ctx))
	}
}

func user(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(RouterContextKey("params")).(httprouter.Params)
	w.Write([]byte(p.ByName("name")))
}

func main() {
	router := httprouter.New()
	router.GET("/user/:name", compatible(user))
	http.ListenAndServe(":8080", router)
}
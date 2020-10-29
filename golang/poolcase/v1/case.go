package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

func slowQuery() error {
	// timeout sql
	_, err := db.Exec("")

	return err
}

func main() {
	var err error
	db, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", exampleHandler)

	log.Println("Listening...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	err := slowQuery()
	if err != nil {
		serverError(w, err)
		return
	}
	fmt.Fprintf(w, "OK")
}

func serverError(w http.ResponseWriter, err error) {
	log.Fatal("ERROR: %s", err.Error())
	http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
}

func setTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

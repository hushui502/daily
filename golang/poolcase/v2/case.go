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

func slowQuery(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

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

	if err = db.Ping(); err != nil {
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
	err := slowQuery(r.Context())
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

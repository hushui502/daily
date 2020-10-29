package main

import (
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

func main() {
	limiter := rate.NewLimiter(rate.Every(100*time.Millisecond), 10)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			log.Println("say hello")
		}
	})

	_ = http.ListenAndServe(":8080", nil)
}

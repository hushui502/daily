package main

import (
	"fmt"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	expriation := time.Now()
	expriation = expriation.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "name", Value: "hufan", Expires: expriation}
	http.SetCookie(w, &cookie)
}

func getname(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("name")
	fmt.Fprint(w, cookie)
}

func main() {
	http.ListenAndServe(":8080", nil)
	http.HandleFunc("/login", login)
	http.HandleFunc("/getname", getname)
}

package main

import (
	"encoding/json"
	"net/http"
)

type Author struct {
	name string `json:"name"`
	age  int    `json:"age"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {

	a := Author{name: "hufan", age: 12}
	json.Marshal(a)
	//json.Unmarshal()

}

package main

import (
	"net/http"
	_ "net/http/pprof"
)

var datas []string

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	datas = append(datas, sData)

	return sData
}

func main() {
	go func() {
		for {
			Add("hello")
		}
	}()

	http.ListenAndServe(":6060", nil)
}

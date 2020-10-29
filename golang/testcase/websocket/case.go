package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	//hello := "hello"
	//enbyte := base64encode([]byte(hello))
	//fmt.Println(enbyte)
	//
	//debyte, _ := base64decode(enbyte)
	//fmt.Println(string(debyte))
	http.HandleFunc("/", locale)
	http.ListenAndServe(":8080", nil)

}

func base64encode(src []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(src))
}

func base64decode(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func locale(w http.ResponseWriter, r *http.Request) {
	//r.Header.Get()
	fmt.Println(r.URL.Query().Get("locale"))
	errors.New()
}

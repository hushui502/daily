package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func main() {
	http.Handle("/", websocket.Handler(Echo))
	err := http.ListenAndServe(":8080", nil)

	// TLS
	//err := http.ListenAndServeTLS(":8080", "hu.name.pem", "private.pem", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func Echo(ws *websocket.Conn) {
	for n := 0; n < 10; n++ {
		msg := "Hello " + string(n+48)
		fmt.Println("Sending to client: " + msg)
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can not send")
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can not receive")
			break
		}

		fmt.Println("Received back from client: " + reply)
	}
}

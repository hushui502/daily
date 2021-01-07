package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"os"
)


// Normal websocket ====> curl ws://localhost:8080
// TLS websocket ====> curl wss://localhost:8080

func main() {

	// url
	service := "ws://localhost:8080/"
	conn, err := websocket.Dial(service, "", "http://localhost")
	if err != nil {
		fmt.Println(err)
	}

	// JSON
	//websocket.JSON.Send(conn, "")
	//websocket.JSON.Receive(conn, "")

	var msg string
	for {
		err := websocket.Message.Receive(conn, &msg)
		if err != nil {
			if err == io.EOF {
				// graceful shutdown by server
				break
			}
			fmt.Println("Could not receive msg " + err.Error())
			break
		}
		fmt.Println("Received from server: " + msg)
		// return the msg
		err = websocket.Message.Send(conn, msg)
		if err != nil {
			fmt.Println("Could not return msg")
			break
		}
	}

	os.Exit(1)
}

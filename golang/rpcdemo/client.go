package rpcdemo

import (
	"log"
)

func main() {
	client, err := DialHelloService("tcp", "localhost:9999")
	if err != nil {
		log.Fatal("dialing...", err)
	}
	var reply string
	err = client.Hello("hello", &reply)
	if err != nil {
		log.Fatal(err)
	}
}

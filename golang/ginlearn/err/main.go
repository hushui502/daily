package main

import (
	"awesomeProject2/ginlearn/err/sub1"
	"fmt"
)

func main() {
	err := sub1.Diff(1, 3)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}

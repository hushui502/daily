package main

import (
	"encoding/json"
	"fmt"
)

type AA struct {
	name string `json:"name"`
}
func main() {
	aa := AA{name:"file"}
	res, _ := json.Marshal(&aa)

	fmt.Println(string(res))
}

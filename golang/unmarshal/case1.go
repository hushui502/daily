package main

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Foo string `json:"foo"`
}

func main() {
	content := []byte(`{"foo": "bar"}`)
	var result1, result2 *Result

	err := json.Unmarshal(content, &result1)
	fmt.Println(result1, err)

	err = json.Unmarshal(content, result2)
	fmt.Println(result2, err)
}

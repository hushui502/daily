package main

import "fmt"

const templae = `Warning: you are using %d bytes of storage`

func main() {
	warn := fmt.Sprintf(templae, 12)
	fmt.Println(warn)
}

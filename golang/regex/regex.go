package main

import (
	"fmt"
	"regexp"
)

const text = "My email is cc@gmail.com"

func main() {
	re := regexp.MustCompile(`([a1-zA-Z0-9]+)@([a1-zA-Z0-9]+)(\.[a1-zA-Z0-9.]+)`)
	match := re.FindAllStringSubmatch(text, -1)
	for _, match := range match {
		fmt.Println(match)
	}
}

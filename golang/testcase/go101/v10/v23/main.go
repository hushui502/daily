package main

import (
	"fmt"
	"os"
)

func fileExit(filename string) bool {
	info, err := os.Stat(filename)
	fmt.Println(info.Size())
	return err == nil
}
func main() {
	a := fileExit("D:\\project\\go\\src\\awesomeProject2\\go101\\v10\\v21\\blockchain.go")
	println(a)
}

package main

import (
	"fmt"
	"os"
)

// file mode
func main() {
	info, _ := os.Stat("")
	mode := info.Mode()
	fmt.Println("filename", " mode is ", mode.String()[1:10])
}



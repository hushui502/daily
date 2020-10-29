package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func savePid() {
	pidFilename := "/pid/" + filepath.Base(os.Args[0]) + ".pid"
	pid := os.Getpid()
	ioutil.WriteFile(pidFilename, []byte(strconv.Itoa(pid)), 0755)
}

func main() {
	a := filepath.Base(os.Args[0])
	b := filepath.Dir(os.Args[0])
	fmt.Println(a)
	fmt.Println(b)
}

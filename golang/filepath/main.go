package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.Println(GetAppPath())
}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

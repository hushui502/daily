package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func main() {
	var stderr io.ReadCloser
	var err error
	cmd := exec.Command("ping", "localhost")

	stderr, err = cmd.StderrPipe()
	if err != nil {

	}
	err = cmd.Start()
	if stderr != nil {
		reader := bufio.NewReader(stderr)
		go func() {
			for {
				line, err := reader.ReadString(byte('\n'))
				if err != nil || io.EOF == err {
					stderr.Close()
					break
				}
				fmt.Println(line)
			}
		}()
	}
	err = cmd.Wait()
	if err != nil {
		return
	}
}

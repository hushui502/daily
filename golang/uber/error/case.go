package main

import "fmt"

type errNotFound struct {
	file string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("file %q not found", e.file)
}

func open(file string) error {
	return errNotFound{file: file}
}

func IsNotFoundError(err error) bool {
	_, ok := err.(errNotFound)
	return ok
}

func main() {
	if err := open("hhh"); err != nil {
		if IsNotFoundError(err) {
			println("errNotFound")
		} else {
			panic("ss")
		}
	}
}

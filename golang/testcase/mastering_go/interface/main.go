package main

import "fmt"

type A struct {
	
}

type B struct {
	
}

type C struct {

}

func tellInterface(x interface{}) {
	switch v := x.(type) {
	case A:
		fmt.Println("This is a A")
	case B:
		fmt.Println("This is a B")
	default:
		fmt.Println("This is a ", v)
	}
}
func main() {
	a := A{}
	b := B{}
	c := C{}
	tellInterface(a)
	tellInterface(b)
	tellInterface(c)
}

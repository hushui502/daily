package main

func test() *int {
	x := new(int)
	*x = 3
	return x
}

type A struct {

}
func main() {
	test()
	a := new(A)
	println(a)
}
package main

import "fmt"

func funReturnFun() func() int {
	i := 0
	return func() int {
		i++
		return i * i
	}
}

func func1(i int) int {
	return i * i
}

func func2(i int) int {
	return i + i
}

func funFunc(f func(int) int, v int) int {
	return f(v)
}

func main() {
	i := funReturnFun()
	j := funReturnFun()
	fmt.Println(i())
	fmt.Println(i())
	fmt.Println(j())
	fmt.Println(j())

	fmt.Println("--")
	fmt.Println(funFunc(func1, 3))
	fmt.Println(funFunc(func2, 3))
}

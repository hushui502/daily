package main

import (
	"fmt"
)

func main() {
	ar := make([]int64, 2, 2)
	fmt.Println(ar)

	arDeriv1 := ar
	arDeriv1 = append(arDeriv1, 0)
	fmt.Println(ar, arDeriv1)

	arDeriv2 := ar
	arDeriv2 = append(arDeriv2, 1)
	fmt.Println(ar, arDeriv1, arDeriv2)
}
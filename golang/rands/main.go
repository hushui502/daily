package main

import (
	"fmt"
	"math/rand"
)

func uniqRands(quantity, maxVal int) []int {
	if maxVal < quantity {
		quantity = maxVal
	}

	intSlice := make([]int, maxVal)
	for i := 0; i < maxVal; i++ {
		intSlice[i] = i
	}

	for i := 0; i < quantity; i++ {
		j := rand.Int()%maxVal + i
		intSlice[i], intSlice[j] = intSlice[j], intSlice[i]
		maxVal--
	}

	return intSlice[0:quantity]
}

func main() {
	nums := uniqRands(5, 5)
	fmt.Println(nums)
}

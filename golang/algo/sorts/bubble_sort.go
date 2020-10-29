package main

import "fmt"

func bubbleSort(array []int) []int {
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(array)-1; i++ {
			if array[i+1] < array[i] {
				array[i+1], array[i] = array[i], array[i+1]
				swapped = true
			}
		}
	}
	return array
}

func main() {
	a := []int{2, 1, 5, 3}
	b := bubbleSort(a)
	fmt.Println(b)
}

func b1(arr []int) []int {
	swapped := true
	for swapped {
		swapped = false
		for i := 0 ; i < len(arr)-1; i++ {
			if arr[i+1] < arr[i] {
				swapped = true
				arr[i+1], arr[i] = arr[i], arr[i+1]
			}
		}
	}

	return arr
}
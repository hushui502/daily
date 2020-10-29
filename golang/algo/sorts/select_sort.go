package main

func selectionSort(arr []int) []int {

	for i := 0; i < len(arr); i++ {
		min := i
		for j := i + 1; j < len(arr); j++ {
			if arr[min] > arr[j] {
				min = j
			}
		}

		tmp := arr[i]
		arr[i] = arr[min]
		arr[min] = tmp
	}

	return arr
}

func se(arr []int) []int {
	for i := 0; i < len(arr); i++ {
		min := i
		for j := i+1; j < len(arr); j++ {
			if arr[min] > arr[j] {
				min = j
			}
		}
		arr[i], arr[min] = arr[min], arr[i]
	}

	return arr
}
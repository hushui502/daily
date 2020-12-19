package sort

import "math"

func bubbleSort(a []int, n int) {
	if n <= 1 {
		return
	}

	for i := 0; i < n; i++ {
		flag := false
		for j := 0; j < n-i-1; j++ {
			if a[j] > a[j+1] {
				a[j], a[j+1] = a[j+1], a[j]
				flag = true
			}
		}
		if !flag {
			break
		}
	}
}

func insertSort(a []int, n int) {
	if n <= 1 {
		return
	}

	for i := 1; i < n; i++ {
		value := a[i]
		j := i - 1
		for ; j >= 0; j-- {
			if a[j] > value {
				a[j+1] = a[j]
			} else {
				break
			}
		}
		a[j+1] = value
	}
}

func selectSort(a []int, n int) {
	if n <= 1 {
		return
	}

	for i := 0; i < n; i++ {
		minIndex := i
		for j := i+1; j < n; j++ {
			if a[j] < a[minIndex] {
				minIndex = j
			}
		}
		a[i], a[minIndex] = a[minIndex], a[i]
	}
}

func MergeSort(arr []int) {
	arrLen := len(arr)
	if arrLen <= 1 {
		return
	}

	mergeSort(arr, 0, arrLen-1)
	
}

func mergeSort(arr []int, start, end int) {
	if start >= end {
		return
	}
	mid := start + (end-start) / 2
	mergeSort(arr, start, mid)
	mergeSort(arr, mid+1, end)
	merge(arr, start, mid, end)
}

func merge(arr []int, start, mid, end int) {
	tmpArr := make([]int, end-start+1)

	i := start
	j := mid + 1
	k := 0

	for ; i <= mid && j <= end; k++ {
		if arr[i] <= arr[j] {
			tmpArr[k] = arr[i]
			i++
		} else {
			tmpArr[k] = arr[j]
			j++
		}
	}

	for ; i <= mid; i++ {
		tmpArr[k] = arr[i]
		k++
	}
	for ; j <= end; j++ {
		tmpArr[k] = arr[j]
		k++
	}

	copy(arr[start:end+1], tmpArr)
}

func QuickSort(arr []int) {
	separateSort(arr, 0, len(arr)-1)
}

func separateSort(arr []int, start, end int) {
	if start >= end {
		return
	}

	i := partition(arr, start, end)
	separateSort(arr, start, i-1)
	separateSort(arr, i+1, end)
}

func partition(arr []int, start, end int) int {
	pivot := arr[end]

	var i = start
	for j := start; j < end; j++ {
		if arr[j] < pivot {
			if !(i == j) {
				arr[i], arr[j] = arr[j], arr[i]
			}
			i++
		}
	}

	arr[i], arr[end] = arr[end], arr[i]

	return i
}

func bucketSortSimple(source []int) {
	if len(source) <= 1 {
		return
	}
	array := make([]int, getMax(source)+1)
	for i := 0; i < len(source); i++ {
		array[source[i]]++
	}

	c := make([]int, 0)
	for i := 0; i < len(array); i++ {
		for array[i] != 0 {
			c = append(c, i)
			array[i]--
		}
	}

	copy(source, c)
}

func getMax(a []int) int {
	max := a[0]
	for i := 1; i < len(a); i++ {
		if a[i] > max {
			max = a[i]
		}
	}

	return max
}

func countingSort(a []int, n int) {
	if n <= 1 {
		return
	}

	var max int = math.MinInt32

	for i := range a {
		if a[i] > max {
			max = a[i]
		}
	}

	c := make([]int, max+1)
	for i := range a {
		c[a[i]]++
	}

	r := make([]int, n)
	for i := n-1; i >= 0; i-- {
		index := c[a[i]] - 1
		r[index] = a[i]
		c[a[i]]--
	}

	copy(a, r)
}
package main

func merge(a, b []int) []int {

	r := make([]int, len(a)+len(b))
	i := 0
	j := 0

	for i < len(a) && j < len(b) {
		if a[i] <= b[j] {
			r[i+j] = a[i]
			i++
		} else {
			r[i+j] = b[j]
			j++
		}
	}

	for i < len(a) {
		r[i+j] = a[i]
		i++
	}

	for j < len(b) {
		r[i+j] = b[j]
		j++
	}

	return r
}

func MergeSort(items []int) []int {
	if len(items) < 2 {
		return items
	}

	middle := len(items) / 2
	a := MergeSort(items[:middle])
	b := MergeSort(items[middle:])

	return merge(a, b)
}

package alg

func separateSort(a []int, start, end int) {
	i := partition(a, start, end)
	separateSort(a, start, i-1)
	separateSort(a, i+1, end)
}

func partition(a []int, start, end int) int {
	pivot := a[end]
	i := start

	for j := start; j < end; j++ {
		if a[j] < pivot {
			if !(i == j) {
				a[i], a[j] = a[j], a[i]
			}
			i++
		}
	}
	a[i], a[end] = a[end], a[i]
	return i
}

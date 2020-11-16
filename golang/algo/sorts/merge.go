package alg

func mergeSort(a []int, start, end int) {
	if start >= end {
		return
	}

	mid := start + (end - start) / 2
	mergeSort(a, start, mid)
	mergeSort(a, mid+1, end)
	merge(a, start, mid, end)
}

func merge(a []int, start, mid, end int) {
	tmpArr := make([]int, end-start+1)
	i := start
	j := mid+1
	k := 0

	for ; i <= mid && j <= end; k++ {
		if a[i] < a[j] {
			tmpArr[k] = a[i]
			i++
		} else {
			tmpArr[k] = a[j]
			j++
		}
	}

	for ; i <= mid; i++ {
		tmpArr[k] = a[i]
		k++
	}
	for ; j <= end; j++ {
		tmpArr[k] = a[j]
		k++
	}

	copy(a[start:end+1], tmpArr)
}

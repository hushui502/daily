package alg

func bubble(arr []int) {
	n := len(arr)
	isSwap := false

	for i := n - 1; i > 0; i++ {
		for j := 0; j < i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				isSwap = true
			}
		}
		if !isSwap {
			return
		}
	}
}

func selectSort(arr []int) {
	n := len(arr)
	minIndex := 0
	minVal := arr[minIndex]

	for i := 0; i < n - 1; i++ {
		minVal = arr[i]
		minIndex = i
		for j := i + 1; j < n; j++ {
			if arr[j] < minVal {
				minIndex = j
				minVal = arr[j]
			}
		}
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
}

func InsertSort(arr []int) {
	n := len(arr)

	for i := 1; i < n; i++ {
		deal := arr[i]
		j := i - 1

		if deal < arr[j] {
			for ; deal < arr[j] && j >= 0; j-- {
				arr[j+1] = arr[j]
			}
			arr[j+1] = deal
		}
	}
}

func MergeSort(arr []int, begin int, end int) {
	if end - begin > 1 {
		mid := begin + (end - begin) / 2
		MergeSort(arr, begin, mid)
		MergeSort(arr, mid, end)
		Merge(arr, begin, mid, end)
	}
}

func Merge(arr []int, begin, mid, end int) {
	leftSize := mid - begin
	rightSize := end - mid
	result := make([]int, 0, leftSize+rightSize)

	l, r := 0, 0

	for l < leftSize && r < rightSize {
		lValue := arr[begin+l]
		rValue := arr[mid+r]
		if lValue < rValue {
			result = append(result, lValue)
			l++
		} else {
			result = append(result, rValue)
			r++
		}
	}

	result = append(result, arr[begin+l:mid]...)
	result = append(result, arr[mid+r:end]...)
}
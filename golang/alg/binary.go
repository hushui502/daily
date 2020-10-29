package alg

func bs(a []int, v int, low, high int) int {
	if low >= high {
		return -1
	}
	mid := low + (high - low) / 2
	if a[mid] == v {
		return mid
	} else if a[mid] > v {
		return bs(a, v, low, mid-1)
	} else {
		return bs(a, v, mid+1, high)
	}
}

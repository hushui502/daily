package alg

func InsertionSort(a []int, n int) {
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

func Insert2(a []int) {
	n := len(a)
	if n <= 1 {
		return
	}

	for i := 1; i <= n-1; i++ {
		deal := a[i]
		j := i - 1

		if deal < a[j] {
			for ; j >= 0 && deal < a[j]; j-- {
				a[j+1] = a[j]
			}
			a[j+1] = deal
		}
	}
}

package alg

func BubbleSort(a []int, n int) {
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


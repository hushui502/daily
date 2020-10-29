package main

func countBits(num int) []int {
	res := make([]int, num+1)
	for i := 1; i <= num; i++ {
		res[i] = res[i&(i-1)] + 1
	}

	return res
}

func counts(num int) int {
	res := 0
	for num&(num-1) != 0 {
		res++
	}

	return res
}

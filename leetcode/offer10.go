package leetcode

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}


func fib1(n int) int {
	f0, f1 := 0, 1
	for i := 0; i < n; i++ {
		f0, f1 = f1, (f0+f1) % 1000000007
	}

	return f0
}



func fib2(n int) int {
	if n <= 1 {
		return n
	}

	res := make([]int, n + 1)
	res[0] = 0
	res[1] = 1
	for i := 2; i <= n; i++ {
		res[i] = (res[i-1] + res[i-2]) % 1000000007
	}

	return res[n]
}
package main

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const a, b = 10, 20

func main() {
	if max(a, b) == a {
		println("a")
	} else {
		println("b")
	}
}

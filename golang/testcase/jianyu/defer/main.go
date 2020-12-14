package main

func DoDefer(key, value string) {
	defer func(key, value string) {
		_ = key + value
	}(key, value)
}

func DoNotDefer(key, value string) {
	_ = key + value
}

func main() {
	defer println("a1")
}

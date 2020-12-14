package main

func comp1(s1, s2 []byte) bool {
	return string(s1) == string(s2)
}

func main() {
	println("ss")
}

func comp2(s1, s2 []byte) bool {
	return conv(s1) == conv(s2)
}

func conv(s []byte) string {
	return string(s)
}

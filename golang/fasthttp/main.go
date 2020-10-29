package main


func main() {
	s := []byte("ss")
	str := "hello"
	s = append(s, str...)
}

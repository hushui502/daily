package main

func Caller1(c chan string) {
	c <- "hello caller1"
	close(c)
}

func Caller2(c chan string) {
	c <- "hello caller2"
	close(c)
}

//func main() {
//	a1, b := make(chan string), make(chan string)
//	go Caller1(a1)
//	go Caller2(b)
//
//	var msg string
//	ok1, ok2 := true, true
//	for ok1 || ok2 {
//		select {
//		case msg, ok1 = <-a1:
//			if ok1 {
//				fmt.Printf("%s from A\n", msg)
//			}
//		case msg, ok2 = <-b:
//			if ok2 {
//				fmt.Printf("%s from B\n", msg)
//			}
//		}
//	}
//}
//
//

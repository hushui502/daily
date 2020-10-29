package main

func main() {
	ch := make(chan struct{})
	ch <- struct{}{}

	<-ch
}

func Outer() {
OuterLoop:
	for i := 0; i < 10; i++ {
		switch {
		case i == 2:
			break OuterLoop
		}
	}
}

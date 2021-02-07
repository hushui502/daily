package main

import (
	"fmt"
	"time"
)

type Message struct {
	str string
	block chan struct{}
}

func main() {
	ch := fanIn(generator("HELLO"), generator("WORLD"))

	msg1 := <-ch
	msg2 := <-ch

	<-msg1.block
	<-msg2.block
}

func fanIn(ch1, ch2 <-chan Message) <-chan Message {
	new_ch := make(chan Message)
	go func() {
		for {
			new_ch <- <-ch1
		}
	}()
	go func() {
		for {
			new_ch <- <-ch2
		}
	}()

	return new_ch
}

func generator(msg string) <-chan Message {
	ch := make(chan Message)
	blockStep := make(chan struct{})
	go func() {
		for i := 0; ; i++ {
			ch <- Message{fmt.Sprintf("%s %d", msg, i), blockStep}
			time.Sleep(time.Second)
			blockStep <- struct{}{}
		}
	}()
	
	return ch
}





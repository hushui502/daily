package main

import (
	"fmt"
	"time"
)

func main() {
	// no buffer channel
	ch := make(chan string)
	go channel_print("hello", ch)
	for i := 0; i < 5; i++ {
		// 这里<-ch会一直阻塞到接受到信息
		fmt.Println(<-ch)
	}
	fmt.Println("Done!")
}

// chan<-  单向通道，意思是这个chan是需要接受外部传来的信息的，这种单向的写法往往是为了让代码逻辑更清晰
func channel_print(msg string, ch chan<- string) {
	for i := 0; ; i++ {
		ch <- fmt.Sprintf("%s %d", msg, i)
		time.Sleep(time.Second)
	}
}

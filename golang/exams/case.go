package main

import (
	"fmt"
	"sync"
	"time"
)

/*
写代码实现两个 goroutine，其中⼀个产⽣随机数并写⼊到 go channel 中，另外⼀
个从 channel 中读取数字并打印到标准输出。最终输出五个随机数。
*/
//func main() {
//	wg := sync.WaitGroup{}
//	out := make(chan int)
//
//	wg.Add(2)
//
//	go func() {
//		defer wg.Done()
//		for i := 0; i < 5; i++ {
//			out <- rand.Intn(5)
//		}
//		close(out)
//	}()
//
//	go func() {
//		defer wg.Done()
//		for i := range out {
//			fmt.Println(i)
//		}
//	}()
//
//	wg.Wait()
//}

//func main() {
//	go func() {
//		t := time.NewTicker(time.Duration(1) * time.Second)
//		for {
//			select {
//			case <-t.C:
//				go func() {
//					defer func() {
//						if err := recover(); err != nil {
//							fmt.Println(err)
//						}
//					}()
//				}()
//				proc()
//			}
//		}
//	}()
//
//	select {}
//}

func proc() {
	println("hello")
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	wg.Add(1)

	ch := make(chan bool)

	go time.AfterFunc(timeout, func() {
		ch <- true
	})

	go func() {
		wg.Wait()
		ch <- false
	}()

	return <-ch
}

const (
	a = iota
	b = iota
)

const (
	name = "hello"
	c = iota
	d = iota
)

type Student struct {
	Name string
}

func main() {
	str1 := []string{"a", "b", "c"}
	str2 := str1[1:]

	str2[1] = "new"

	fmt.Println(str1)
	str2 = append(str2, "z", "x")
	fmt.Println(str1)

	fmt.Println(Student{Name: "menglu"} == Student{Name: "menglu"})
}

// dddcdccdab
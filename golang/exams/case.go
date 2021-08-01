package main

import (
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
	c    = iota
	d    = iota
)

type Student struct {
	Name string
}

//func main() {
//	str1 := []string{"a", "b", "c"}
//	str2 := str1[1:]
//
//	str2[1] = "new"
//
//	fmt.Println(str1)
//	str2 = append(str2, "z", "x")
//	fmt.Println(str1)
//
//	fmt.Println(Student{Name: "menglu"} == Student{Name: "menglu"})
//}
//

//func main() {
//	timer := time.NewTimer(time.Duration(5) * time.Second)
//	data := []int{1, 2, 3, 10, 999, 8, 345, 7, 98, 33, 66, 77, 88, 68, 96}
//
//	dataLen := len(data)
//	size := 3
//	target := 345
//	ctx, cancel := context.WithCancel(context.Background())
//	resultChan := make(chan bool)
//
//	for i := 0; i < dataLen; i += size {
//		end := i + size
//		if end >= dataLen {
//			end = dataLen-1
//		}
//		go searchTarget(ctx, data[i:end], target, resultChan)
//	}
//
//	select {
//	case <-timer.C:
//		fmt.Fprintln(os.Stderr, "Timeout! Not Found")
//		cancel()
//	case <- resultChan:
//		fmt.Fprintf(os.Stdout, "Found it!\n")
//		cancel()
//	}
//
//	time.Sleep(time.Second * 2)
//}
//
//func searchTarget(ctx context.Context, data []int, target int, resultChan chan bool) {
//	for _, v := range data {
//		select {
//		case <-ctx.Done():
//			fmt.Fprintf(os.Stdout, "Task cancelded! \n")
//			return
//		default:
//		}
//		fmt.Fprintf(os.Stdout, "v: %d \n", v)
//		time.Sleep(time.Duration(1500) * time.Millisecond)
//		if target == v {
//			resultChan <- true
//			return
//		}
//	}
//}

func main() {

}

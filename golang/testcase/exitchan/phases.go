package exitchan

import (
	"fmt"
	"time"
)

// 尽量使用非阻塞 I/O（非阻塞 I/O 常用来实现高性能的网络库），阻塞 I/O 很可能导致 goroutine 在某个调用一直等待，而无法正确结束。
// 业务逻辑总是考虑退出机制，避免死循环。
// 任务分段执行，超时后即时退出，避免 goroutine 无用的执行过多，浪费资源。

func do2phases(phase1, done chan bool) {
	time.Sleep(time.Second)	//	the first phase
	select {
	case phase1 <- true:
	default:
		return
	}
	time.Sleep(time.Second)	// the second phase
	done <- true
}

func timeoutFirstPhase() error {
	phase1 := make(chan bool)
	done := make(chan bool)

	go do2phases(phase1, done)
	select {
	case <-phase1:
		<-done
		fmt.Println("done")
		return nil
	case <-time.After(time.Microsecond):
		return fmt.Errorf("timeout")
	}
}





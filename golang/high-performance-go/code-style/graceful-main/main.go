package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// parent-child model
//func main() {
//	// create a sig channel, capture syscall signal to this channel
//	sig := make(chan os.Signal, 1)
//	stopCh := make(chan struct{})
//	finishCh := make(chan struct{})
//	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
//
//	go func(stopCh, finishCh chan struct{}) {
//		for {
//			select {
//			case <-stopCh:
//				fmt.Println("stopped")
//				finishCh <- struct{}{}
//				return
//			default:
//				time.Sleep(time.Second)
//			}
//		}
//	}(stopCh, finishCh)
//
//	<-sig
//	stopCh <- struct{}{}
//	<-finishCh
//	fmt.Println("finished")
//}

//func main() {
//	sig := make(chan os.Signal)
//	stopCh := make(chan chan struct{})
//	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
//
//	go func(stopCh chan chan struct{}) {
//		for {
//			select {
//			// channel是进行值拷贝的
//			case ch := <-stopCh:
//				fmt.Println("stopped")
//				ch <- struct{}{}
//				return
//			default:
//				time.Sleep(time.Second)
//			}
//		}
//	}(stopCh)
//
//	<-sig
//	ch := make(chan struct{})
//	stopCh <- ch
//	<-ch
//	fmt.Println("finished")
//}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())
	finishedCh := make(chan struct{})

	go func(ctx context.Context, finishedCh chan struct{}) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopped")
				finishedCh <- struct{}{}
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}(ctx, finishedCh)

	<-sig
	cancel()
	<-finishedCh
	fmt.Println("finished")
}
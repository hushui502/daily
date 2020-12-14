package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleSignal(signal os.Signal) {
	fmt.Println("handleSignal() Caught", signal)
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case os.Interrupt:
				fmt.Println("caught: ", sig)
			case syscall.SIGINT:
				handleSignal(sig)
				return
			}
		}
	}()

	for  {
		fmt.Println(",")
		time.Sleep(20 * time.Second)
	}
}

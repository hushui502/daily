/*
	Five silent philosophers sit at a round table with bowls of spaghetti.
	Forks are placed between each pair of adjacent philosophers.
	Each philosopher must alternately think and eat. However,
	a philosopher can only eat spaghetti when they have both left and right forks.
	Each fork can be held by only one philosopher and so a philosopher can use the fork only if it is not being used by
	another philosopher.
	After an individual philosopher finishes eating,
	they need to put down both forks so that the forks become available to others.
	A philosopher can take the fork on their right or the one on their left as they become available,
	but cannot start eating before getting both forks.
	Eating is not limited by the remaining amounts of spaghetti or stomach space;
	an infinite supply and an infinite demand are assumed.
	Design a discipline of behaviour (a concurrent algorithm) such that no philosopher will starve;
	i.e., each can forever continue to alternate between eating and thinking,
	assuming that no philosopher can know when others may want to eat or think.
*/

package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
	"time"
)

const N = 5

var forks [N]chan int

func requestForks(ID int) {
	
	leftFork, rightFork := (ID+N-1)%N, (ID+N+1)%N
	
	for {
		select {
		case forks[leftFork] <- ID:
			select {
			case forks[rightFork] <- ID:
				
				sleepTime := rand.Intn(10) + 10
				fmt.Printf("Philosopher %d picks both forks and eats for %d milliseconds.\n", ID, sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Millisecond)
				
				<-forks[leftFork]
				<-forks[rightFork]
				
				return
			default:
				fmt.Printf("Philosopher %d can't pick his right fork then gave up.\n", ID)
				<-forks[leftFork]
			}
		default:
			fmt.Printf("Philosopher %d can't pick both forks then gave up.\n", ID)
		}
		
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}
}

func main() {
	_ = trace.Start(os.Stderr)
	defer trace.Stop()

	for i := 0; i < N; i++ {
		forks[i] = make(chan int, 1)
	}

	for i := 0; i < 100; i++ {
		go requestForks(rand.Intn(N))
	}

	time.Sleep(time.Duration(1000) * time.Millisecond)
}
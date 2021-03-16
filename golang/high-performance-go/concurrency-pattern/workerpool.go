package main

import (
	"fmt"
	"sync"
	"time"
)

func workerEfficient(id int, jobs <-chan int, results chan<- int) {
	var wg sync.WaitGroup
	for j := range jobs {
		wg.Add(1)
		go func(job int) {
			fmt.Println("worker", id, "started job", job)
			time.Sleep(time.Second)
			fmt.Println("worker", id, "finished job", job)
			results <- job * 2
			wg.Done()
		}(j)
	}

	wg.Wait()
}

func main() {
	const numJobs = 8
	jobs := make(chan int, numJobs)
	resulls := make(chan int, numJobs)

	// start the worker
	for w := 1; w <= 3; w++ {
		go workerEfficient(w, jobs, resulls)
	}

	// send the work
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}

	close(jobs)
	fmt.Println("closed job")

	// finish job
	for a := 1; a <= numJobs; a++ {
		// block channel
		<-resulls
	}
	close(resulls)
}

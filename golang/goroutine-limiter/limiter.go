package main

import "sync/atomic"

const (
	// DefaultLimit is the default concurrency limit
	DefaultLimit = 100
)

type ConcurrencyLimiter struct {
	limit int
	// just use chan for limiting the num of goroutine
	tickets chan int
	// currently running goroutines
	numInProgress int32
}

// NewConcurrencyLimiter allocates a new ConcurrencyLimiter
func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	if limit <= 0 {
		limit = DefaultLimit
	}

	// allocate a limiter instance
	c := &ConcurrencyLimiter{
		limit:   limit,
		tickets: make(chan int, limit),
	}

	// allocate the tickets:
	for i := 0; i < c.limit; i++ {
		c.tickets <- 1
	}

	return c
}

// ExecuteWithTicket adds a job into an execution queue and returns a ticket id.
func (c *ConcurrencyLimiter) Execute(job func()) int {
	ticket := <-c.tickets
	atomic.AddInt32(&c.numInProgress, 1)
	go func() {
		defer func() {
			c.tickets <- ticket
			atomic.AddInt32(&c.numInProgress, -1)
		}()

		job()
	}()

	return ticket
}

// ExecuteWithTicket adds a job into an execution queue and returns a ticket id.
func (c *ConcurrencyLimiter) ExecuteWithTicket(job func(ticket int)) int {
	ticket := <-c.tickets
	atomic.AddInt32(&c.numInProgress, 1)
	go func() {
		defer func() {
			c.tickets <- ticket
			atomic.AddInt32(&c.numInProgress, -1)
		}()

		// run the job
		job(ticket)
	}()
	return ticket
}

// Wait will block all the previously Executed jobs completed running.
func (c *ConcurrencyLimiter) Wait() {
	for i := 0; i < c.limit; i++ {
		_ = <-c.tickets
	}
}

// GetNumInProgress returns a (racy) counter of how many go routines are active right now
func (c *ConcurrencyLimiter) GetNumInProgress() int32 {
	return atomic.LoadInt32(&c.numInProgress)
}

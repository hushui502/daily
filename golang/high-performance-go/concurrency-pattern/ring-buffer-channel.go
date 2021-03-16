package main

import "log"

// ringBuffer throttle buffer for implement async channel
type ringBuffer struct {
	inCh  chan int
	outCh chan int
}

func NewRingBuffer(inCh, outCh chan int) *ringBuffer {
	return &ringBuffer{
		inCh:  inCh,
		outCh: outCh,
	}
}

func (r *ringBuffer) Run() {
	for v := range r.inCh {
		select {
		case r.outCh <- v:
		default:
			<-r.outCh
			r.outCh <- v
		}
	}

	// because the inCh is already closed, so this operation is safe
	close(r.outCh)
}

func main() {
	inCh := make(chan int)
	// try to change outCh buffer to understand the result
	outCh := make(chan int, 1)

	rb := NewRingBuffer(inCh, outCh)
	go rb.Run()

	for i := 0; i < 10; i++ {
		inCh <- i
	}

	// close inCh => close outCh safely
	close(inCh)

	for res := range outCh {
		log.Println(res)
	}

}

package queue

import "fmt"

type ArrayQueue struct {
	q []interface{}
	capacity int
	head int
	tail int
}

func NewArrayQueue(n int) *ArrayQueue {
	return &ArrayQueue{
		q:        make([]interface{}, n),
		capacity: n,
		head:     0,
		tail:     0,
	}
}

func (queue *ArrayQueue) EnQueue(v interface{}) bool {
	// full
	if queue.tail == queue.capacity {
		return false
	}

	queue.q[queue.tail] = v
	queue.tail++

	return true
}

func (queue *ArrayQueue) DeQueue() interface{} {
	if queue.tail == queue.head {
		return nil
	}

	v := queue.q[queue.head]
	queue.head++

	return v
}

func (this *ArrayQueue) String() string {
	if this.head == this.tail {
		return "empty queue"
	}
	result := "head"
	for i := this.head; i <= this.tail-1; i++ {
		result += fmt.Sprintf("<-%+v", this.q[i])
	}
	result += "<-tail"
	return result
}
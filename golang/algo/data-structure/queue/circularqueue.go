package queue

import "fmt"

type CircularQueue struct {
	q []interface{}
	capacity int
	head int
	tail int
}

func NewCircularQueue(n int) *CircularQueue {
	if n == 0 {
		return nil
	}

	return &CircularQueue{make([]interface{}, n), n, 0, 0}
}

func (c *CircularQueue) IsEmpty() bool {
	if c.head == c.tail {
		return true
	}

	return false
}

func (this *CircularQueue) IsFull() bool {
	if this.head == (this.tail+1)%this.capacity {
		return true
	}
	return false
}

func (c *CircularQueue) EnQueue(v interface{}) bool {
	if c.IsFull() {
		return false
	}
	c.q[c.tail] = v
	c.tail = (c.tail + 1) % c.capacity

	return true
}

func (c *CircularQueue) DeQueue() interface{} {
	if c.IsEmpty() {
		return false
	}

	v := c.q[c.head]
	c.head = (c.head + 1) % c.capacity

	return v
}

func (this *CircularQueue) String() string {
	if this.IsEmpty() {
		return "empty queue"
	}
	result := "head"
	var i = this.head
	for true {
		result += fmt.Sprintf("<-%+v", this.q[i])
		i = (i + 1) % this.capacity
		if i == this.tail {
			break
		}
	}
	result += "<-tail"
	return result
}



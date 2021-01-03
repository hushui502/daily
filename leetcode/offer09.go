package leetcode

import "container/list"

type CQueue struct {
	stack1, stack2 *list.List
}

func Constructor() CQueue {
	return CQueue{
		stack1: list.New(),
		stack2: list.New(),
	}
}

func (cq *CQueue) AppendTail(value int) {
	cq.stack1.PushBack(value)
}

func (cq *CQueue) DeleteHead() int {
	if cq.stack2.Len() == 0 {
		for cq.stack1.Len() > 0 {
			cq.stack2.PushBack(cq.stack1.Remove(cq.stack1.Back()))
		}
	}

	if cq.stack2.Len() != 0 {
		e := cq.stack2.Back()
		cq.stack2.Remove(e)

		return e.Value.(int)
	}

	return -1
}

type CQueue1 struct {
	stackA []int
	stackB []int
}

func NewQueue() CQueue1 {
	return CQueue1{}
}

func (this *CQueue1) AppendTail(value int) {
	this.stackA = append(this.stackA, value)
}

func (this *CQueue1) DeleteHead() int {
	if len(this.stackB) == 0 {
		if len(this.stackA) == 0 {
			return -1
		}
		for len(this.stackA) > 0 {
			index := len(this.stackA) - 1
			value := this.stackA[index]
			this.stackB = append(this.stackB, value)
			this.stackA = this.stackA[:index]
		}
	}

	index := len(this.stackB) - 1
	value := this.stackB[index]
	this.stackB = this.stackB[:index]

	return value
}
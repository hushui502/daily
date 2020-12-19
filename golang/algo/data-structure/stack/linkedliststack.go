package stack

import "fmt"

type node struct {
	next *node
	val interface{}
}

type LinkedListStack struct {
	topNode *node
}

func NewLinedListStack() *LinkedListStack {
	return &LinkedListStack{nil}
}

func (stack *LinkedListStack) IsEmpty() bool {
	return stack.topNode == nil
}

func (stack *LinkedListStack) Push(v interface{}) {
	stack.topNode = &node{next: stack.topNode, val: v}
}

func (stack *LinkedListStack) Pop() interface{} {
	if stack.IsEmpty() {
		return nil
	}
	v := stack.topNode.val
	stack.topNode = stack.topNode.next

	return v
}

func (stack *LinkedListStack) Top() interface{} {
	if stack.IsEmpty() {
		return nil
	}

	return stack.topNode.val
}

func (stack *LinkedListStack) Flush() {
	stack.topNode = nil
}

func (stack *LinkedListStack) Print() {
	if stack.IsEmpty() {
		fmt.Println("empty stack")
	} else {
		cur := stack.topNode
		for cur != nil {
			fmt.Println(cur.val)
			cur = cur.next
		}
	}
}


package linked_list

import "fmt"

type node struct {
	val  int
	next *node
}

type singlelinkedlist struct {
	head *node
}

func newNode(val int) *node {
	return &node{val, nil}
}

func (l1 *singlelinkedlist) addAtBeg(val int) {
	n := newNode(val)
	n.next = l1.head
	l1.head = n
}

func (l1 *singlelinkedlist) addAtEnd(val int) {
	n := newNode(val)

	if l1.head == nil {
		l1.head = n
		return
	}

	cur := l1.head
	for ; cur.next != nil; cur = cur.next {
	}

	cur.next = n
}

func (l1 *singlelinkedlist) delAtBeg() int {
	if l1.head == nil {
		return -1
	}

	cur := l1.head
	l1.head = cur.next
	return cur.val
}

func (l1 *singlelinkedlist) delAtEnd() int {
	if l1.head == nil {
		return -1
	}

	if l1.head.next == nil {
		return l1.delAtBeg()
	}

	cur := l1.head
	for ; cur.next.next != nil; cur = cur.next {
	}
	retval := cur.next.val
	cur.next = nil
	return retval
}

func (l1 *singlelinkedlist) count() int {
	var ctr int
	for cur := l1.head; cur != nil; cur = cur.next {
		ctr += 1
	}

	return ctr
}

func (l1 *singlelinkedlist) reverse() {
	var prev, next *node
	cur := l1.head
	for cur != nil {
		next = cur.next
		cur.next = prev
		prev = cur
		cur = next
	}

	l1.head = prev
}

func (l1 *singlelinkedlist) display() {
	for cur := l1.head; cur != nil; cur = cur.next {
		fmt.Print(cur.val, " ")
	}
	fmt.Print("\n")
}

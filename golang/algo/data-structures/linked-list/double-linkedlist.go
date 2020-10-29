package linked_list

type node struct {
	val  int
	next *node
	prev *node
}

type doublelinkedlist struct {
	head *node
}

func newNode(val int) *node {
	n := &node{}
	n.val = val
	n.next = nil
	n.prev = nil
	return n
}

func (l1 *doublelinkedlist) addAtBeg(val int) {
	n := newNode(val)
	n.next = l1.head
	l1.head = n
}

func (l1 *doublelinkedlist) addAtEnd(val int) {
	n := newNode(val)
	if l1.head == nil {
		l1.head = n
		return
	}

	cur := l1.head
	for ; cur.next != nil; cur = cur.next {
	}

	cur.next = n
	n.prev = cur
}

func (l1 *doublelinkedlist) delAtBeg() int {
	if l1.head == nil {
		return -1
	}

	cur := l1.head
	l1.head = cur.next

	if l1.head != nil {
		l1.head.prev = nil
	}

	return cur.val
}

func (l1 *doublelinkedlist) delAtEnd() int {
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

func (l1 *doublelinkedlist) count() int {
	var ctr int = 0

	for cur := l1.head; cur != nil; cur = cur.next {
		ctr += 1
	}

	return ctr
}

func (l1 *doublelinkedlist) reverse() {
	var prev, next *node
	cur := l1.head

	for cur != nil {
		next = cur.next
		cur.next = prev
		cur.prev = next
		prev = cur
		cur = next
	}
	l1.head = prev
}

package linked_list

type node struct {
	v    int
	next *node
}

var head *node = nil

func (l *node) pushFront(val int) *node {
	if head == nil {
		l.v = val
		l.next = nil
		head = l
		return l
	} else {
		nnode := &node{}
		nnode = head

		nnode2 := &node{
			v:    val,
			next: nnode,
		}

		head = nnode2
		return head
	}
}

func (l *node) pushBack(val int) *node {
	if head == nil {
		l.v = val
		l.next = nil
		head = l
		return l
	} else {
		for l.next != nil {
			l = l.next
		}
		l.next = &node{}
		l.next.v = val
		l.next.next = nil
		return l
	}
}

func (l *node) popFront() *node {
	if head == nil {
		return head
	}

	cpnode := &node{}
	cpnode = head.next

	head = cpnode
	return head
}

func (l *node) popBack() *node {
	if head == nil {
		return head
	}

	cpnode := &node{}
	cpnode = head

	for cpnode.next.next != nil {
		cpnode = cpnode.next
	}
	cpnode.next = nil
	return head
}

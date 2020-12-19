package linkedlist

import "fmt"

type ListNode struct {
	next *ListNode
	value interface{}
}

type LinkedList struct {
	head *ListNode
	length uint
}

func NewListNode(v interface{}) *ListNode {
	return &ListNode{nil, v}
}

func (node *ListNode) GetNext() *ListNode {
	return node.next
}

func (node *ListNode) GetValue() interface{} {
	return node.value
}

func NewLinedList() *LinkedList {
	return &LinkedList{NewListNode(0), 0}
}

func (ls *LinkedList) InsertAfter(p *ListNode, v interface{}) bool {
	if p == nil {
		return false
	}
	newNode := NewListNode(v)
	oldNext := p.next
	p.next = newNode
	newNode.next = oldNext
	ls.length++

	return true
}

func (ls *LinkedList) InsertBefore(p *ListNode, v interface{}) bool {
	if p == nil || p == ls.head {
		return false
	}

	cur := ls.head.next
	pre := ls.head

	for cur != nil {
		if cur == p {
			break
		}
		pre = cur
		cur = cur.next
	}

	if cur == nil {
		return false
	}

	newNode := NewListNode(v)
	pre.next = newNode
	newNode.next = cur
	ls.length++

	return true
}

func (ls *LinkedList) InsertToTail(v interface{}) bool {
	cur := ls.head
	for cur.next != nil {
		cur = cur.next
	}

	return ls.InsertAfter(cur, v)
}

func (ls *LinkedList) FindByIndex(index uint) *ListNode {
	if index >= ls.length {
		return nil
	}
	cur := ls.head.next
	for i := 0; uint(i) < index; i++ {
		cur = cur.next
	}

	return cur
}

func (ls *LinkedList) DeleteNode(p *ListNode) bool {
	if p == nil {
		return false
	}

	cur := ls.head.next
	pre := ls.head

	for cur != nil {
		if cur == p {
			break
		}
		cur = cur.next
	}

	if cur == nil {
		return false
	}
	pre.next = p.next
	p = nil
	ls.length--

	return true
}

func (ls *LinkedList) Reverse() {
	if ls.head == nil || ls.head.next == nil || ls.head.next.next == nil {
		return
	}
	var pre *ListNode = nil
	cur := ls.head.next
	for cur != nil {
		tmp := cur.next
		cur.next = pre
		pre = cur
		cur = tmp
	}

	ls.head.next = pre
}

func (ls *LinkedList) HasCycle() bool {
	if ls.head != nil {
		slow := ls.head
		fast := ls.head
		for slow != nil && fast != nil {
			slow = slow.next
			fast = fast.next.next
			if slow == fast {
				return true
			}
		}
	}

	return false
}

func MergeSortedList(l1, l2 *LinkedList) *LinkedList {
	if l1 == nil || l1.head == nil || l1.head.next == nil {
		return l2
	}
	if l2 == nil || l2.head == nil || l2.head.next == nil {
		return l1
	}

	l := &LinkedList{head:&ListNode{}}
	cur := l.head
	l1cur := l1.head.next
	l2cur := l2.head.next

	for l1cur != nil && l2cur != nil {
		if l1cur.value.(int) < l2cur.value.(int) {
			cur.next = l1cur
			l1cur = l1cur.next
		} else {
			cur.next = l2cur
			l2cur = l2cur.next
		}
		cur = cur.next
	}

	if l1cur != nil {
		cur.next = l1cur
	}
	if l2cur != nil {
		cur.next = l2cur
	}

	return l
}

func (ls *LinkedList) DeleteBottomN(n int) {
	if n <= 0 || ls.head == nil || ls.head.next == nil {
		return
	}

	fast := ls.head
	for i := 1; i < n && fast != nil; i++ {
		fast = fast.next
	}

	if fast == nil {
		return
	}

	slow := ls.head
	for fast.next != nil {
		slow = slow.next
		fast = fast.next
	}

	// skip node which need to be deleted
	slow.next = slow.next.next
}

func (ls *LinkedList) FindMiddleNode() *ListNode {
	if ls.head == nil || ls.head.next == nil {
		return nil
	}
	if ls.head.next.next == nil {
		return ls.head.next
	}

	slow, fast := ls.head, ls.head
	for slow != nil && fast != nil {
		fast = fast.next.next
		slow = slow.next
	}

	return slow
}

func (ls *LinkedList) Print() {
	cur := ls.head.next
	format := ""
	for cur != nil {
		format += fmt.Sprintf("%+v", cur.GetValue())
		cur = cur.next
		if cur != nil {
			format += "->"
		}
	}

	fmt.Println(format)
}

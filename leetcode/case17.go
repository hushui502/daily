package main

func mergeLists(lists []*ListNode) *ListNode {
	length := len(lists)
	if length < 1 {
		return nil
	}
	if length == 1 {
		return lists[0]
	}

	num := length / 2
	left := mergeLists(lists[:num])
	right := mergeLists(lists[num:])
	return mergeTwoList1(left, right)
}

func mergeTwoList1(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

	if l1.Val < l2.Val {
		l1.Next = mergeTwoList1(l1.Next, l2)
		return l1
	}
	l2.Next = mergeTwoList1(l1, l2.Next)
	return l2
}
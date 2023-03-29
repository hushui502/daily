package main

func addTwoNumbers(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

	x, y := 0, 0
	carry := 0

	head := &ListNode{}
	dummy := head

	for l1 != nil || l2 != nil || carry != 0 {
		if l1 != nil {
			x = l1.Val
			l1 = l1.Next
		} else {
			x = 0
		}
		if l2 != nil {
			y = l2.Val
			l2 = l2.Next
		} else {
			y = 0
		}

		tmp := x + y + carry
		carry = tmp / 10
		dummy.Next = &ListNode{Val: tmp % 10}
		dummy = dummy.Next
	}

	return head.Next
}

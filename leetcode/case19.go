package main

func reverseGroup(head *ListNode, k int) *ListNode {
	node := head
	for i := 0; i < k; i++ {
		if node == nil {
			return nil
		}
		node = node.Next
	}
	newHead := reverse1(head, node)
	head.Next = reverseGroup(node, k)
	return newHead
}

func reverse1(first *ListNode, last *ListNode) *ListNode {
	prev := last
	for first != last {
		tmp := first.Next
		first.Next = prev
		prev = first
		first = tmp
	}

	return prev
}
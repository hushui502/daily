package leetcode

func deleteDuplicates3(head *ListNode) *ListNode {
	if head == nil {
		return head
	}

	nilNode := &ListNode{Val:0, Next:head}
	head = nilNode

	lastVal := 0
	for head.Next != nil && head.Next.Next != nil {
		if head.Next.Val == head.Next.Next.Val {
			lastVal = head.Next.Val
			for head.Next != nil && lastVal == head.Next.Next.Val {
				head.Next = head.Next.Next
			}
		} else {
			head = head.Next
		}
	}

	return nilNode.Next
}

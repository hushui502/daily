package leetcode

func deleteDuplicates(head *ListNode) *ListNode {
	cur := head
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return nil
	}

	for cur.Next != nil {
		if cur.Next.Val == cur.Val {
			cur.Next = cur.Next.Next
		} else {
			cur = cur.Next
		}
	}

	return head
}

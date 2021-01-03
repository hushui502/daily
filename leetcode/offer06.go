package leetcode

func reversePrint(head *ListNode) []int {
	if head == nil {
		return nil
	}
	
	return appendData(head)
}

func appendData(head *ListNode) []int {
	if head.Next != nil {
		list := appendData(head.Next)
		list = append(list, head.Val)
		return list
	}

	return []int{head.Val}
}


func reversePrint2(head *ListNode) []int {
	if head == nil {
		return nil
	}

	var newHead *ListNode
	res := []int{}

	for head != nil {
		node := head.Next
		head.Next = newHead
		newHead = head
		head = node
	}

	for newHead != nil {
		res = append(res, newHead.Val)
		newHead = newHead.Next
	}

	return res
}

func reversePrint3(head *ListNode) []int {
	if head == nil {
		return nil
	}

	res := []int{}
	for head != nil {
		res = append(res, head.Val)
		head = head.Next
	}

	for i, j := 0, len(res) - 1; i < j; {
		res[i], res[j] = res[j], res[i]
		i++
		j--
	}

	return res
}
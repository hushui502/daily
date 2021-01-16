package leetcode

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	// 1. 新建一个哨兵节点
	result := &ListNode{}
	result.Next = head
	// 2. 给cur和pre赋值
	cur := result
	var pre *ListNode
	i := 1
	// 3. 开始遍历
	for head != nil {
		// 因为是倒数第几个删除，所以用这种方式，有点像快指针
		if i >= n {
			pre = cur
			cur = cur.Next
		}
		head = head.Next
		i++
	}

	// 4. 删除元素，因为都是哨兵节点开始的，所以可以直接返回哨兵节点
	pre.Next = pre.Next.Next
	return result.Next
}

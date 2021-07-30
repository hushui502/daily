package main

func twoSum(nums []int, target int) []int {
	m := make(map[int]int)

	for i, v := range nums {
		if _, ok := m[target-v]; ok {
			return []int{i, m[target-v]}
		}
		m[v] = i
	}

	return []int{}
}

type ListNode struct {
	Val int
	Next *ListNode
}

func addTwoSum(l1, l2 *ListNode) *ListNode {
	res := &ListNode{}
	dummy := res

	n1, n2, carry := 0, 0, 0
	for l1 != nil || l2 != nil || carry != 0 {
		if l1 == nil {
			n1 = 0
		} else {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 == nil {
			n2 = 0
		} else {
			n2 = l2.Val
			l2 = l2.Next
		}

		num := n1 + n2 + carry
		dummy.Next = &ListNode{Val: num % 10}
		carry = num / 10
		dummy = dummy.Next
	}

	if carry != 0 {
		dummy.Next = &ListNode{Val: carry}
	}

	return res.Next
}


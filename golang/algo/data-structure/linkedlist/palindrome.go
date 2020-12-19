package linkedlist

func isPalindrome(l *LinkedList) bool {
	len := l.length
	if len == 0 {
		return false
	}
	if len == 1 {
		return true
	}

	s := make([]string, 0, len/2)
	cur := l.head
	for i := uint(1); i <= len; i++ {
		cur = cur.next
		if len%2 != 0 && i == (len/2+1) {
			continue
		}
		if i < len/2 {
			s = append(s, cur.GetValue().(string))
		} else {
			if s[len-i] != cur.GetValue().(string) {
				return false
			}
		}
	}

	return true
}

func isPalindrome2(l *LinkedList) bool {
	len := l.length
	if len == 0 {
		return false
	}
	if len == 1 {
		return true
	}

	var isPalindrome = true
	step := len / 2
	var pre *ListNode = nil
	cur := l.head.next
	next := l.head.next.next
	for i := uint(1); i <= step; i++ {
		tmp := cur.GetNext()
		cur.next = pre
		pre = cur
		cur = tmp
		next = cur.GetNext()
	}
	mid := cur

	var left, right *ListNode = pre, nil
	if len%2 == 0 {
		right = mid
	} else {
		right = mid.GetNext()
	}

	for left != nil && right != nil {
		if left.GetValue().(string) != right.GetValue().(string) {
			isPalindrome = false
			break
		}
		left = left.GetNext()
		right = right.GetNext()
	}

	cur = pre
	pre = mid

	for cur != nil {
		next = cur.GetNext()
		cur.next = pre
		pre = cur
		cur = next
	}

	l.head.next = pre

	return isPalindrome
}
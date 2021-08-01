package main

import "sort"

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
	Val  int
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

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func zTraversalTree(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	res := [][]int{}
	q := []*TreeNode{root}
	curNum, nexNum := 1, 0
	tmp := []int{}
	index := 0

	for len(q) != 0 {
		if curNum != 0 {
			node := q[0]
			if node.Left != nil {
				q = append(q, node.Left)
				nexNum++
			}
			if node.Right != nil {
				q = append(q, node.Right)
				nexNum++
			}
			q = q[1:]
			curNum--
			tmp = append(tmp, node.Val)
		}

		if curNum == 0 {
			if index%2 != 0 {
				for i, j := 0, len(tmp)-1; i < j; {
					tmp[i], tmp[j] = tmp[j], tmp[i]
					i++
					j--
				}
			}
			res = append(res, tmp)
			curNum = nexNum
			nexNum = 0
			tmp = []int{}
		}
	}

	return res
}

func searchRotateSortedArray(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	low, high := 0, len(nums)-1

	for low <= high {
		mid := low + (high-low)>>1
		if nums[mid] == target {
			return mid
		} else if nums[mid] > nums[low] {
			if nums[low] <= target && nums[mid] > target {
				high = mid - 1
			} else {
				low = mid + 1
			}
		} else if nums[mid] < nums[high] {
			if nums[high] >= target && nums[mid] < target {
				low = mid + 1
			} else {
				high = mid - 1
			}
		} else {
			if nums[mid] == nums[low] {
				low++
			}
			if nums[mid] == nums[high] {
				high--
			}
		}
	}

	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func lengthOfLongestSubstring(s string) int {
	var freq [256]int

	res, left, right := 0, 0, -1

	for left < len(s) {
		if right+1 < len(s) && freq[s[right+1]-'a'] == 0 {
			freq[s[right+1]-'a'] = 1
			right++
		} else {
			freq[s[left]-'a'] = 0
			left++
		}
		res = max(res, right-left+1)
	}

	return res
}

func threeSum(nums []int) [][]int {
	var res [][]int

	m := make(map[int]int)

	for _, v := range nums {
		m[v]++
	}

	uniqueNum := []int{}
	for k := range m {
		uniqueNum = append(uniqueNum, k)
	}
	sort.Ints(uniqueNum)

	for i := 0; i < len(uniqueNum); i++ {
		if uniqueNum[i]*3 == 0 && m[uniqueNum[i]] >= 3 {
			res = append(res, []int{uniqueNum[i], uniqueNum[i], uniqueNum[i]})
		}
		for j := i + 1; j < len(uniqueNum); j++ {
			if uniqueNum[i]*2+uniqueNum[j] == 0 && m[uniqueNum[i]] >= 2 {
				res = append(res, []int{uniqueNum[i], uniqueNum[i], uniqueNum[j]})
			}
			if uniqueNum[j]*2+uniqueNum[i] == 0 && m[uniqueNum[j]] >= 2 {
				res = append(res, []int{uniqueNum[j], uniqueNum[j], uniqueNum[i]})
			}
			c := 0 - uniqueNum[j] - uniqueNum[i]
			if c > uniqueNum[j] && m[c] >= 1 {
				res = append(res, []int{uniqueNum[i], uniqueNum[j], c})
			}
		}
	}

	return res
}

func expand(s string, left, right int) (int, int) {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		left--
		right++
	}

	return left + 1, right - 1
}

func longestPalindrome(s string) string {
	start, end := 0, 0

	for i := 0; i < len(s); i++ {
		l1, r1 := expand(s, i, i)
		l2, r2 := expand(s, i, i+1)

		if r1-l1 > end-start {
			start, end = l1, r1
		}
		if r2-l2 > end-start {
			start, end = l2, r2
		}
	}

	return s[start : end+1]
}

func maxArea(nums []int) int {
	var res int
	var hight int

	left, right := 0, len(nums)-1

	for left < right {
		wid := right - left
		if nums[left] < nums[right] {
			hight = nums[left]
			left++
		} else {
			hight = nums[right]
			right--
		}

		res = max(res, wid*hight)
	}

	return res
}

var letterMap = []string{
	" ",
	"",
	"abc",
	"def",
	"ghi",
	"jkl",
	"mno",
	"pqrs",
	"tuv",
	"wxyz",
}

var result []string

func letterCombinations(digits string) []string {
	if digits == "" {
		return []string{}
	}

	result = []string{}
	findCombinations(&digits, 0, "")

	return result
}

func findCombinations(digits *string, index int, s string) {
	if index == len(*digits) {
		result = append(result, s)
		return
	}
	num := (*digits)[index]
	letters := letterMap[num-'0']
	for i := 0; i < len(letters); i++ {
		findCombinations(digits, index+1, s+string(letters[i]))
	}
}

func removeNthElementFromEnd(head *ListNode, n int) *ListNode {
	fast, slow := head, head

	for i := 0; i < n; i++ {
		fast = fast.Next
	}

	if fast == nil {
		head = head.Next
		return head
	}

	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}

	slow.Next = slow.Next.Next

	return head
}

func isValid(s string) bool {
	if len(s) == 0 {
		return true
	}

	stack := []rune{}

	for _, v := range s {
		if v == '[' || v == '{' || v == '(' {
			stack = append(stack, v)
		} else if v == ']' && len(stack) > 0 && stack[len(stack)-1] == ']' || v == '}' && len(stack) > 0 && stack[len(stack)-1] == '{' || v == ')' && len(stack) > 0 && stack[len(stack)-1] == '(' {
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
	}

	return len(stack) == 0
}

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			l1.Next = mergeTwoLists(l1.Next, l2)
			return l1
		}
		l2.Next = mergeTwoLists(l2.Next, l1)
		return l2
	}

	return l1
}

func mergeTwoLists2(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	head := dummy

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			head.Next = l1
			l1 = l1.Next
		} else {
			head.Next = l2
			l2 = l2.Next
		}

		head = head.Next
	}

	if l1 != nil {
		head.Next = l1
	}
	if l2 != nil {
		head.Next = l2
	}

	return dummy.Next
}

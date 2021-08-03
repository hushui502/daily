package main

import (
	"math"
	"sort"
)

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

func mergeKListNodes(list []*ListNode) *ListNode {
	len := len(list)
	if len == 0 {
		return nil
	}
	if len == 1 {
		return list[0]
	}

	left := mergeKListNodes(list[:len/2])
	right := mergeKListNodes(list[len/2:])

	root := mergeTwoLists(left, right)

	return root
}

func generateParenthesis(n int) []string {
	var res []string
	if n == 0 {
		return []string{}
	}

	findParenthesis(n, n, "", &res)

	return res
}

func findParenthesis(lindex, rindex int, str string, res *[]string) {
	if lindex == 0 && rindex == 0 {
		*res = append(*res, str)
		return
	}

	if lindex > 0 {
		findParenthesis(lindex-1, rindex, str+"(", res)
	}
	if lindex < rindex && rindex > 0 {
		findParenthesis(lindex, rindex-1, str+")", res)
	}
}

func trap(height []int) int {
	left, right := 0, len(height)-1

	maxLeft, maxRight := 0, 0

	var res int

	for left <= right {
		if height[left] < height[right] {
			if height[left] < maxLeft {
				res += maxLeft - height[left]
			} else {
				maxLeft = height[left]
			}
			left++
		} else {
			if height[right] < maxRight {
				res += maxRight - height[right]
			} else {
				maxRight = height[right]
			}
			right--
		}
	}

	return res
}

func combinationSum(nums []int, target int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}

	var res [][]int
	var c []int

	sort.Ints(nums)

	findCombinationSum(nums, 0, target, c, &res)

	return res
}

func findCombinationSum(nums []int, index int, target int, c []int, res *[][]int) {
	if target <= 0 {
		if target == 0 {
			tmp := make([]int, len(c))
			copy(tmp, c)
			*res = append(*res, tmp)
		}
		return
	}

	for i := index; i < len(nums); i++ {
		if nums[i] > target {
			break
		}
		c = append(c, nums[i])
		findCombinationSum(nums, i, target-nums[i], c, res)
		c = c[:len(c)-1]
	}
}

func searchFirst(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)>>1

		if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			if mid == 0 || nums[mid-1] != target {
				return mid
			}
			right = mid - 1
		}
	}

	return -1
}

func searchLast(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := left + (right-left)>>1

		if nums[mid] > target {
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			if (mid == len(nums)-1) || (nums[mid+1] != target) {
				return mid
			}
			left = mid + 1
		}
	}

	return -1
}

func searchRange(nums []int, target int) []int {
	return []int{searchFirst(nums, target), searchLast(nums, target)}
}

func findMedianSortedArray(nums1, nums2 []float64) float64 {
	res := append(nums1, nums2...)

	sort.Float64s(res)

	if len(res)%2 == 1 {
		return res[len(res)/2]
	}

	return (res[len(res)/2-1] + res[len(res)/2]) / 2
}

func climbStairs(n int) int {
	if n < 2 {
		return 0
	}

	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

func nextPermutation(nums []int) []int {
	n := len(nums)
	i := n - 2

	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	}

	if i >= 0 {
		j := n - 1
		for j >= 0 && nums[j] <= nums[i] {
			j--
		}
		nums[i], nums[j] = nums[j], nums[i]
	}

	reverse(nums[i+1:])

	return nums
}

func reverse(nums []int) {
	for i, n := 0, len(nums); i < n/2; i++ {
		nums[i], nums[n-i-1] = nums[n-i-1], nums[i]
	}
}

func search(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)>>1
		if nums[mid] == target {
			return mid
		} else if nums[mid] > nums[left] {
			if nums[left] <= target && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else if nums[mid] < nums[right] {
			if nums[right] >= target && nums[mid] < target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		} else {
			if nums[mid] == nums[left] {
				left++
			}
			if nums[mid] == nums[right] {
				right++
			}
		}
	}

	return -1
}

func findPermute(nums []int, index int, c []int, used *[]bool, res *[][]int) {
	if index == len(nums) {
		tmp := make([]int, len(c))
		copy(tmp, c)
		*res = append(*res, c)
		return
	}

	for i := 0; i < len(nums); i++ {
		if !(*used)[i] {
			(*used)[i] = true
			c = append(c, nums[i])
			findPermute(nums, index+1, c, used, res)
			c = c[:len(c)-1]
			(*used)[i] = false
		}
	}
}

func permute(nums []int) [][]int {
	var res [][]int

	var c []int
	used := make([]bool, len(nums))

	findPermute(nums, 0, c, &used, &res)

	return res
}

func rotate(matrix [][]int) [][]int {
	row := len(matrix)
	if row == 0 {
		return [][]int{}
	}

	col := len(matrix[0])

	for i := 0; i < row; i++ {
		for j := i; j < col; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	halfCol := col / 2
	for i := 0; i < row; i++ {
		for j := 0; j < halfCol; j++ {
			matrix[i][j], matrix[i][col-j-1] = matrix[i][col-j-1], matrix[i][j]
		}
	}

	return matrix
}

func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	fast, slow := head, head

	for slow != nil && fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}

	return false
}

func singleNumber(nums []int) int {
	var res int

	for _, v := range nums {
		res ^= v
	}

	return res
}

func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}

	min, maxProfit := prices[0], 0

	for i := 0; i < len(prices); i++ {
		if prices[i]-min > maxProfit {
			maxProfit = prices[i] - min
		}
		if min > prices[i] {
			min = prices[i]
		}
	}

	return maxProfit
}

func preOrder(root *TreeNode, res *[]int) {
	if root == nil {
		return
	}

	*res = append(*res, root.Val)
	preOrder(root.Left, res)
	preOrder(root.Right, res)
}

func flatten(root *TreeNode) {
	res := []int{}

	preOrder(root, &res)

	cur := root
	for i := 1; i < len(res); i++ {
		cur.Left = nil
		cur.Right = &TreeNode{Val: res[i], Left: nil, Right: nil}
		cur = cur.Right
	}

	return
}

func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	return max(maxDepth(root.Left), maxDepth(root.Right)) + 1
}

func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	q := []*TreeNode{root}
	curNum, nextNum := 1, 0
	tmp := []int{}
	res := [][]int{}

	for len(q) > 0 {
		if curNum != 0 {
			node := q[0]
			if node.Left != nil {
				q = append(q, node.Left)
				nextNum++
			}
			if node.Right != nil {
				q = append(q, node.Right)
				nextNum++
			}
			curNum--
			tmp = append(tmp, node.Val)
			q = q[1:]
		}
		if curNum == 0 {
			res = append(res, tmp)
			tmp = []int{}
			curNum = nextNum
			nextNum = 0
		}
	}

	return res
}

func isSameTree(t1, t2 *TreeNode) bool {
	if t1 == nil && t2 == nil {
		return true
	} else if t1 != nil && t2 != nil {
		if t1.Val != t2.Val {
			return false
		}
		return isSameTree(t1.Left, t2.Left) && isSameTree(t1.Right, t2.Right)
	} else {
		return false
	}
}

func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return isSameTree(root.Left, invertTree(root.Right))
}

func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	invertTree(root.Left)
	invertTree(root.Right)

	root.Left, root.Right = root.Right, root.Left

	return root
}

func isValidBST(root *TreeNode) bool {
	return dfs(root, math.MinInt64, math.MaxFloat64)
}

func dfs(root *TreeNode, min, max float64) bool {
	if root == nil {
		return true
	}
	v := float64(root.Val)

	return v > min && v < max && dfs(root.Left, min, v) && dfs(root.Right, v, max)
}

// 139 128 105 95 79 78 75 76

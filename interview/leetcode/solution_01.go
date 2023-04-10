package leetcode

import (
	"container/list"
	"sort"
)

func twoSum(nums []int, target int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	m := make(map[int]int)

	for i, v := range nums {
		if _, ok := m[target-v]; ok {
			return []int{m[target-v], i}
		}
		m[v] = i
	}

	return []int{}
}

func longestPalindrome(s string) string {
	if len(s) == 0 {
		return ""
	}

	start, end := 0, 0

	for i := 0; i < len(s); i++ {
		left1, right1 := expandAroundCenter(s, i, i)
		left2, right2 := expandAroundCenter(s, i, i+1)

		if right1-left1 > end-start {
			start, end = left1, right1
		}
		if right2-left2 > end-start {
			start, end = left2, right2
		}
	}

	return s[start : end+1]
}

func expandAroundCenter(s string, left, right int) (int, int) {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		left--
		right++
	}

	return left + 1, right - 1
}

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	nums3 := append(nums1, nums2...)
	quickSortHelper(nums3, 0, len(nums3)-1)

	if len(nums3)%2 == 0 {
		return float64(nums3[len(nums3)/2]+nums3[len(nums3)/2-1]) / 2
	} else {
		return float64(nums3[len(nums3)/2])
	}
}

func runningSum(nums []int) []int {
	for i := 1; i < len(nums); i++ {
		nums[i] += nums[i-1]
	}

	return nums
}

func runningSum2(nums []int) []int {
	dp := make([]int, len(nums)+1)
	dp[0] = 0

	for i := 1; i < len(nums); i++ {
		dp[i] = dp[i-1] + nums[i-1]
	}

	return dp[1:]
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode

	for head != nil {
		next := head.Next
		head.Next = prev
		prev = head
		head = next
	}

	return prev
}

func zigzagLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	var result [][]int
	queue := []*TreeNode{root}
	level := 0

	for len(queue) > 0 {
		var tmp []int
		size := len(queue)

		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]

			if level%2 == 0 {
				tmp = append(tmp, node.Val)
			} else {
				tmp = append([]int{node.Val}, tmp...)
			}

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		result = append(result, tmp)
		level++
	}

	return result
}

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil || p == root || q == root {
		return root
	}

	left := lowestCommonAncestor(root.Left, p, q)
	right := lowestCommonAncestor(root.Right, p, q)

	if left != nil && right != nil {
		return root
	}

	if left == nil {
		return right
	}

	return left
}

func numIslands(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}

	count := 0

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == '1' {
				count++
				dfs(grid, i, j)
			}
		}
	}

	return count
}

func dfs(grid [][]byte, i, j int) {
	if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[0]) || grid[i][j] == '0' {
		return
	}

	grid[i][j] = '0'

	dfs(grid, i-1, j)
	dfs(grid, i+1, j)
	dfs(grid, i, j-1)
	dfs(grid, i, j+1)
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	i, j, k := m-1, n-1, m+n-1

	for i >= 0 && j >= 0 {
		if nums1[i] > nums2[j] {
			nums1[k] = nums1[i]
			i--
		} else {
			nums1[k] = nums2[j]
			j--
		}
		k--
	}

	for j >= 0 {
		nums1[k] = nums2[j]
		j--
		k--
	}
}

func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}

	minPrice := prices[0]
	maxProfit := 0

	for i := 1; i < len(prices); i++ {
		if prices[i] < minPrice {
			minPrice = prices[i]
		} else if prices[i]-minPrice > maxProfit {
			maxProfit = prices[i] - minPrice
		}
	}

	return maxProfit
}

type entry struct {
	key, value int
}

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	ll       *list.List
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		ll:       list.New(),
	}
}

func (this *LRUCache) Get(key int) int {
	if e, ok := this.cache[key]; ok {
		this.ll.MoveToFront(e)
		return e.Value.(*entry).value
	}

	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if e, ok := this.cache[key]; ok {
		e.Value.(*entry).value = value
		this.ll.MoveToFront(e)
		return
	}

	e := this.ll.PushFront(&entry{key, value})
	this.cache[key] = e

	if this.ll.Len() > this.capacity {
		e := this.ll.Back()
		this.ll.Remove(e)
		delete(this.cache, e.Value.(*entry).key)
	}
}

func findKthLargest(nums []int, k int) int {
	quickSort(nums)

	return nums[len(nums)-k]
}

func hasCycle(head *ListNode) bool {
	if head == nil {
		return false
	}

	slow, fast := head, head.Next

	for fast != nil && fast.Next != nil {
		if slow == fast {
			return true
		}

		slow = slow.Next
		fast = fast.Next.Next
	}

	return false
}

func isValid(s string) bool {
	var stack []byte

	for i := 0; i < len(s); i++ {
		if s[i] == '(' || s[i] == '[' || s[i] == '{' {
			stack = append(stack, s[i])
		} else {
			if len(stack) == 0 {
				return false
			}

			if s[i] == ')' && stack[len(stack)-1] != '(' {
				return false
			}

			if s[i] == ']' && stack[len(stack)-1] != '[' {
				return false
			}

			if s[i] == '}' && stack[len(stack)-1] != '{' {
				return false
			}

			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func search(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)>>1

		if nums[mid] == target {
			return mid
		}

		if nums[mid] >= nums[left] {
			if target >= nums[left] && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			if target > nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}

	var res [][]int
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		var level []int
		size := len(queue)

		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)

			if node.Left != nil {
				queue = append(queue, node.Left)
			}

			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		res = append(res, level)
	}

	return res
}

func twoSumPlus(nums []int, target int) []int {
	m := make(map[int]int)

	for i, num := range nums {
		if j, ok := m[target-num]; ok {
			return []int{j, i}
		}

		m[num] = i
	}

	return nil
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	cur := dummy

	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			cur.Next = list1
			list1 = list1.Next
		} else {
			cur.Next = list2
			list2 = list2.Next
		}

		cur = cur.Next
	}

	if list1 != nil {
		cur.Next = list1
	}

	if list2 != nil {
		cur.Next = list2
	}

	return dummy.Next
}

func sortArray(nums []int) []int {
	quickSort(nums)

	return nums
}

func maxSubArray(nums []int) int {
	maxSum, curSum := nums[0], nums[0]

	for i := 1; i < len(nums); i++ {
		if curSum < 0 {
			curSum = nums[i]
		} else {
			curSum += nums[i]
		}

		if curSum > maxSum {
			maxSum = curSum
		}
	}

	return maxSum
}

func threeSum(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}

	m := make(map[int]int)
	for _, v := range nums {
		m[v]++
	}

	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var res [][]int

	for i := 0; i < len(keys); i++ {
		if keys[i]*3 == 0 && m[keys[i]] >= 3 {
			res = append(res, []int{keys[i], keys[i], keys[i]})
		}
		for j := i + 1; j < len(keys); j++ {
			if keys[i]*2+keys[j] == 0 && m[keys[i]] >= 2 {
				res = append(res, []int{keys[i], keys[i], keys[j]})
			}
			if keys[j]*2+keys[i] == 0 && m[keys[j]] >= 2 {
				res = append(res, []int{keys[j], keys[j], keys[i]})
			}
			c := 0 - keys[i] - keys[j]
			if c < keys[i] && m[c] >= 1 {
				res = append(res, []int{keys[i], keys[j], c})
			}
		}
	}

	return res
}

func permute(nums []int) [][]int {
	var res [][]int
	permuteHelper(nums, 0, &res)

	return res
}

func permuteHelper(nums []int, start int, res *[][]int) {
	if start == len(nums) {
		*res = append(*res, append([]int{}, nums...))
		return
	}

	for i := start; i < len(nums); i++ {
		nums[i], nums[start] = nums[start], nums[i]
		permuteHelper(nums, start+1, res)
		nums[i], nums[start] = nums[start], nums[i]
	}
}

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}

	a, b := headA, headB

	for a != b {
		if a == nil {
			a = headB
		} else {
			a = a.Next
		}

		if b == nil {
			b = headA
		} else {
			b = b.Next
		}
	}

	return a
}

func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return nil
	}

	var res []int
	left, right, top, bottom := 0, len(matrix[0])-1, 0, len(matrix)-1

	for left <= right && top <= bottom {
		for i := left; i <= right; i++ {
			res = append(res, matrix[top][i])
		}

		for i := top + 1; i <= bottom; i++ {
			res = append(res, matrix[i][right])
		}

		for i := right - 1; i >= left && top != bottom; i-- {
			res = append(res, matrix[bottom][i])
		}

		for i := bottom - 1; i > top && left != right; i-- {
			res = append(res, matrix[i][left])
		}

		left++
		right--
		top++
		bottom--
	}

	return res
}

func reverseBetween(head *ListNode, left int, right int) *ListNode {
	dummy := &ListNode{Next: head}
	pre := dummy

	for i := 0; i < left-1; i++ {
		pre = pre.Next
	}

	cur := pre.Next
	for i := 0; i < right-left; i++ {
		next := cur.Next
		cur.Next = next.Next
		next.Next = pre.Next
		pre.Next = next
	}

	return dummy.Next
}

func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	if len(lists) == 1 {
		return lists[0]
	}

	return mergeTwoLists(mergeKLists(lists[:len(lists)/2]), mergeKLists(lists[len(lists)/2:]))
}

func addStrings(num1 string, num2 string) string {
	var res []byte
	i, j := len(num1)-1, len(num2)-1
	carry := 0

	for i >= 0 || j >= 0 {
		var x, y int

		if i >= 0 {
			x = int(num1[i] - '0')
			i--
		}

		if j >= 0 {
			y = int(num2[j] - '0')
			j--
		}

		sum := x + y + carry
		res = append(res, byte(sum%10+'0'))
		carry = sum / 10
	}

	if carry > 0 {
		res = append(res, byte(carry+'0'))
	}

	reverse(res)

	return string(res)
}

func reverse(nums []byte) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func detectCycle(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	slow, fast := head, head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			break
		}
	}

	if fast == nil || fast.Next == nil {
		return nil
	}

	slow = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}

	return slow
}

// todo 未完成
func mergeArrays(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return nil
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var res [][]int
	res = append(res, intervals[0])

	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] <= res[len(res)-1][1] {
			res[len(res)-1][1] = max(intervals[i][1], res[len(res)-1][1])
		} else {
			res = append(res, intervals[i])
		}
	}

	return res
}

func rightSideView(root *TreeNode) []int {
	if root == nil {
		return nil
	}

	var res []int
	q := []*TreeNode{root}

	for len(q) > 0 {
		res = append(res, q[len(q)-1].Val)

		var next []*TreeNode
		for _, node := range q {
			if node.Left != nil {
				next = append(next, node.Left)
			}
			if node.Right != nil {
				next = append(next, node.Right)
			}
		}

		q = next
	}

	return res
}

func lengthOfLIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	dp := make([]int, len(nums))
	dp[0] = 1

	for i := 1; i < len(nums); i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
	}

	return max(dp...)
}

func trap(height []int) int {
	if len(height) == 0 {
		return 0
	}

	var res int
	left, right := 0, len(height)-1
	leftMax, rightMax := height[left], height[right]

	for left < right {
		if height[left] < height[right] {
			left++
			leftMax = max(leftMax, height[left])
			res += leftMax - height[left]
		} else {
			right--
			rightMax = max(rightMax, height[right])
			res += rightMax - height[right]
		}
	}

	return res
}

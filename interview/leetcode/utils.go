package leetcode

type ListNode struct {
	Val  int
	Next *ListNode
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func max(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}

	return res
}

func min(a, b int) int {
	if a > b {
		return b
	}

	return a
}

func abs(a int) int {
	if a < 0 {
		return -a
	}

	return a
}

func quickSort(nums []int) {
	if len(nums) == 0 {
		return
	}

	quickSortHelper(nums, 0, len(nums)-1)
}

func quickSortHelper(nums []int, left, right int) {
	if left >= right {
		return
	}

	pivot := partition(nums, left, right)
	quickSortHelper(nums, left, pivot-1)
	quickSortHelper(nums, pivot+1, right)
}

func partition(nums []int, left, right int) int {
	pivot := nums[right]
	i := left
	for j := left; j < right; j++ {
		if nums[j] < pivot {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[right] = nums[right], nums[i]

	return i
}

// maxPathSum returns the max path sum of the tree rooted at root.
func maxPathSum(root *TreeNode) int {
	res := root.Val
	maxPathSumHelper(root, &res)

	return res
}

// maxPathSumHelper returns the max path sum of the subtree rooted at root.
func maxPathSumHelper(root *TreeNode, res *int) int {
	if root == nil {
		return 0
	}

	left := max(0, maxPathSumHelper(root.Left, res))
	right := max(0, maxPathSumHelper(root.Right, res))

	*res = max(*res, left+right+root.Val)

	return max(left, right) + root.Val
}

// no recursion
func inorderTraversal(root *TreeNode) []int {
	var res []int
	var stack []*TreeNode

	for root != nil || len(stack) != 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}

		root = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = append(res, root.Val)
		root = root.Right
	}

	return res
}

// recursion
func inorderTraversal2(root *TreeNode) []int {
	var res []int
	inorderTraversalHelper(root, &res)

	return res
}

func inorderTraversalHelper(root *TreeNode, res *[]int) {
	if root == nil {
		return
	}

	inorderTraversalHelper(root.Left, res)
	*res = append(*res, root.Val)
	inorderTraversalHelper(root.Right, res)
}

type MyQueue struct {
	stack []int
	queue []int
}

func MyQueueConstructor() MyQueue {
	return MyQueue{
		stack: make([]int, 0),
		queue: make([]int, 0),
	}
}

func (this *MyQueue) Push(x int) {
	this.stack = append(this.stack, x)
}

func (this *MyQueue) Pop() int {
	if len(this.queue) == 0 {
		this.fromStackToQueue(this.stack)
	}

	res := this.queue[len(this.queue)-1]
	this.queue = this.queue[:len(this.queue)-1]

	return res
}

func (this *MyQueue) Peek() int {
	if len(this.queue) == 0 {
		this.fromStackToQueue(this.stack)
	}

	return this.queue[len(this.queue)-1]
}

func (this *MyQueue) Empty() bool {
	return len(this.stack) == 0 && len(this.queue) == 0
}

func (this *MyQueue) fromStackToQueue(stack []int) []int {
	var queue []int
	for len(stack) != 0 {
		queue = append(queue, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return queue
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	first, second := head, dummy
	for i := 0; i < n; i++ {
		first = first.Next
	}

	for first != nil {
		first = first.Next
		second = second.Next
	}

	second.Next = second.Next.Next

	return dummy.Next
}

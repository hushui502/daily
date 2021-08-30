func permutation(s string) []string {
	var res []string
	bytes := []byte(s)

	var dfs func(int)

	dfs = func(x int) {
		if x == len(bytes)-1 {
			res = append(res, string(bytes))
		}
		dict := make(map[byte]bool)

		for i := x; i < len(bytes); i++ {
			!dict[bytes[i]] {
				dict[bytes[i]] = true
				bytes[i], bytes[x] = bytes[x], bytes[i]
				dfs(x+1)
				bytes[i], bytes[x] = bytes[x], bytes[i]
			}
		}
	}


	dfs(0)

	return res
}


func getKthFromEnd(head *ListNode, n int) *ListNode {
	if head == nil {
		return nil
	}

	slow, fast := head, head

	for i := 0; i < n; i++ {
		fast = fast.Next
	}

	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}

	return slow
}

type Cqueue struct {
	stack1, stack2 *list.List
}

func NewCqueue() Cqueue {
	return Cqueue{
		stack1: list.New(),
		stack2: list.New(),
	}
}

func (this *Cqueue) AppendTail(value int) {
	this.stack1.PushBack(value)
}

func (this *Cqueue) DeleteHead() int {
	if this.stack2.Len() == 0 {
		for this.stack1.Len() > 0 {
			this.stack2.PushBack(this.stack1.Remove(this.stack1.Back()))
		}
	}

	if this.stack2.Len() != 0 {
		e := this.stack2.Back()
		this.stack2.Remove(e)

		return e.Value.(int)
	}

	return -1
}

// 12345
// 34512
func minArray(nums []int) int {
	low := 0
	high := len(nums) - 1

	for low < high {
		mid := low + (high - low)>>1

		if nums[mid] > nums[high] {
			low = mid + 1
		} else if nums[mid] < nums[high] {
			high = mid
		} else {
			high--
		}
	}

	return nums[low]
}

func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	root := &TreeNode{
		Val: preorder[0], 
		Left: nil, 
		Right: nil,
	}

	k := 0
	for i := 0; i < len(inorder); i++ {
		if preorder[0] == inorder[I] {
			break
		}
		k++ 
	}

	root.Left = buildTree(preorder[1:len(inorder[:k])+1], inorder[:k])
	root.Right = buildTree(preorder[len(inorder[:k])+1:], inorder[k+1:])

	return root
}

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

func maxProfit(prices []int) int {
	if len(prices) <= 1 {
		return 0
	}
	min, maxProfit := prices[0], 0

	for i := 0; i < len(prices); i++ {
		if prices[i]-min > maxProfit {
			maxProfit = prices[i]-min
		}

		if min > prices[i] {
			min = prices[i]
		}
	}

	return maxProfit
}

func isStraight(nums []int) bool {
	sort.Ints(nums)

	sub := 0

	for i := 0; i < 4; i++ {
		if nums[i] == 0 {
			continue
		}

		if nums[i] == nums[i+1] {
			return false
		}

		sub += nums[i+1] - nums[i]
	}

	return sub < 5
}

func pathSum(root *TreeNode, target int) [][]int {
	var res [][]int
	var path []int

	var dfs func(*TreeNode, int)

	dfs = func(node *TreeNode, left int) {
		if node == nil {
			return
		}
		left -= node.Val
		path = append(path, node.Val)
		defer func() { path = path[:len(path)-1] }()

		if node.Left == nil && node.Right == nil && left == 0 {
			res = append(res, append([]int(nil), path...))
			return
		}

		dfs(node.Left, left)
		dfs(node.Right, left)
	}


	dfs(root, target)
	return res
}

func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return []int{}
	}

	q := []*TreeNode{root}

	tmp := []int{}
	res := [][]int{}

	curNum, nextNum := 1, 0

	for len(q) != 0 {
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
			tmp = append(tmp, node.Val)
			curNum--
			q = q[1:]
		}
		if curNum == 0 {
			curNum = nextNum
			res = append(res, tmp)
			tmp = []int{}
			nextNum = 0
		}
	}

	return res
}

func maxSubArray(nums []int) int {
	m := make([]int, len(nums))

	for i, v := range nums {
		m[i] = v
	}

	for i := 1; i < len(nums); i++ {
		if m[i] + m[i-1] > m[i] {
			m[i] = m[i] + m[i-1]
		}
	}

	max := m[0]

	for _, v := range m {
		if v > max {
			max = v
		}
	}


	return max
}

func translateNum(num int) int {
	src := strconv.Itoa(num)

	p, q, r := 0, 0, 1

	for i := 0; i < len(src); i++ {
		p , q, r = q, r, 0
		r += q
		if i == 0 {
			continue
		}
		pre := src[i-1:i+1]
		if pre <= "25" && pre >= "10" {
			r += q
		}
	}

	return r
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}

	res := make([]int, n+1)

	res[0] = 1
	res[1] = 1
	for i := 2; i <= n; i++ {
		res[i] = (res[i-1] + res[i-2]) % 1000000007
	}

	return res[n]
}

func findNumberIn2DArray(matrix [][]int, target int) bool {
	if len(matrix) == 0 {
		retrun false
	}

	n := len(matrix)
	m := len(matrix[0])

	i, j := 0, m-1

	for i < n && j >= 0 {
		if matrix[i][j] == target {
			return true
		} else if matrix[i][j] > target {
			j--
		} else {
			i++
		}
	}

	return false
}

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

	p := &ListNode{}
	dummyHead = p

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			dummyHead.Next = l1
			l1 = l1.Next
		} else {
			dummyHead.Next = l2
			l2 = l2.Next
		}
		dummyHead = dummyHead.Next
	}

	if l1 != nil {
		dummyHead.Next = l1
	}
	if l2 != nil {
		dummyHead.Next = l2
	}

	return p.Next
}

func maxSlide(nums []int, k int) []int {
	res := make([]int, 0, k)

	n := len(nums)
	if n == 0 {
		return []int{}
	}

	for i := 0; i < n-k; i++ {
		max := nums[i]
		for j := 1; j < k; j++ {
			if max < nums[i+j] {
				max = nums[i+j]
			}
		}
		res = append(res, max)
	}

	return res
}

func lengthOfLongestSubstring(s string) int {
	var freq [256]byte


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


func lastRemaining(n, m int) int {
	idx := 0

	for i := 2; i < n; i++ {
		idx = (idx+m) % i
	}

	return idx
}

func isSubStructure(A, B *TreeNode) bool {
	if A == nil && B == nil {
		return true
	}
	if A == nil || B == nil {
		return false
	}

	return helper(A, B) || isSubStructure(A.Right, B) || isSubStructure(A.Left, B)
}

func helper(a, b *TreeNode) bool {
	if a == nil {
		return false
	}
	if b == nil {
		return true
	}

	if a.Val != b.Val {
		return false
	}

	return helper(a.Left, b.Left) && helper(a.Right, b.Right)
}

func replaceSpace(s string) string {
	res := make([]rune, len(s)*3)

	i := 0

	for _, v := range s {
		if v != ' ' {
			res[i] = v
			i++
		} else {
			res[i] = '%'
			res[i+1] = '2'
			res[i+2] = '0'
			i += 3
		}
	}

	return string(res)[:i]
}

func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	var prev *ListNode

	curr := head

	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	return prev
}

type entry struct {
	key int
	value int
}

type LRUCache struct {
	cap int
	cache map[int]*list.Element
	lst *list.List
}

func Constructor(cap int) LRUCache {
	return LRUCache {
		cap,
		make(map[int]*list.Element),
		list.New(),
	}
}

func (c *LRUCache) Get(key int) int {
	if e := c.cache[key]; e != nil {
		c.lst.MoveToFront(e)
		return e.Value.(entry).value
	}

	return -1
}

func (c *LRUCache) Put(key, value int) {
	if e := c.cache[key]; e != nil {
		e.Value = entry{key, value}
		c.lst.MoveToFront(e)
		return
	}

	c.cahce[key] = c.lst.PushFront(entry{key, value})
	if c.cap < len(c.cache) {
		delete(c.cache, c.lst.Remove(c.lst.Back()).(entry).key)
	}
}

func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	var prev *ListNode
	cur := head

	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}

	return prev
}

func findRepeatNumber(nums []int) int {
	m := make(map[int]int)

	for _, v := range nums {
		if _, ok := m[v]; ok {
			return v
		}
		m[v] = 1
	}

	return -1
}

func permutaion(s string) []string {
	var res []string
	bytes := []byte(s)

	var dfs func(int)

	dfs = func(x int) {
		if x == len(bytes)-1 {
			res = append(res, string(bytes))
		}

		dict := make(map[byte]bool)

		for i := x; i < len(bytes); i++ {
			if !dict[bytes[i]] {
				dict[bytes[i]] = true
				bytes[i], bytes[x] = bytes[x], bytes[i]
				dfs(x+1)
				bytes[i], bytes[x] = bytes[x], bytes[i]
			}
		}
	}

	dfs(0)

	return res
}

type CQueue struct {
	stack1, stack2 *list.List
}

func Constructor() CQueue {
	return CQueue{
		stack1: list.New(),
		stack2: list.New(),
	}
}

func (c *CQueue) AppendTail(value int) {
	c.stack1.PushBack(value)
}

func (c *CQueue) DeleteHead() int {
	if c.stack2.Len() == 0 {
		for c.stack1.Len() > 0 {
			c.stack2.PushBack(c.stack1.Remove(c.stack1.Back()))
		}
	}

	if c.stack2.Len() != 0 {
		e := c.stack2.Back()
		c.stack2.Remove(e)

		return e.Value.(int)
	}

	return -1
}

func getKthFromEnd(head *ListNode, n int) *ListNode {
	if head == nil {
		return nil
	}

	slow, fast := head, head

	for i := 0; i < n; i++ {
		fast = fast.Next
	}

	for fast != nil {
		fast = fast.Next
		slow = slow.Next
	}

	return slow
}

12345
34512
func minArray(nums []int) int {
	low, high := 0, len(nums)-1

	for low < high {
		mid := low + (high-low)>>1

		if nums[mid] > nums[high] {
			low = mid + 1
		} else if nums[mid] < nums[high] {
			high = mid
		} else {
			high--
		}
	}

	return nums[low]
}

func buildTree(preorder, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}

	root := &TreeNode{
		Val: preorder[0],
		Left: nil,
		Right: nil,
	}

	k := 0

	for i := 0; i < len(inorder); i++ {
		if preorder[0] == inorder[i] {
			break
		}
		k++
	}

	root.Left = buildTree(preorder[1:len(inorder[:k])+1], inorder[:k])
	root.Right = buildTree(preorder[len(inorder[:k])+1:], inorder[k+1:])

	return root
}

func reversePrint(head *ListNode) []int {
	var res []int

	for head != nil {
		res = append(res, head.Val)
		head = head.Next
	}

	for i := 0; i < len(res)/2; i++ {
		res[i], res[len(res)-i-1] = res[len(res)-i-1], res[i]
	}

	return res
}

func fib(n int) int {
	if n == 0 || n == 1 {
		return n
	}

	dp := make([]int, n+1)

	dp[0] = 0
	dp[1] = 1

	for i := 2; i <= n; i++ {
		dp[i] = (dp[i-1] + dp[i-2]) % 1000000007
	}

	return dp[n]
}

func findNumberIn2DArray(matrix [][]int, target int) bool {
	if len(matrix) == 0 {
		return false
	}

	n, m := len(matrix), len(matrix[0])

	i, j := 0, m-1

	for i < n && j >= 0 {
		if matrix[i][j] == target {
			return true
		} else if matrix[i][j] > target {
			j--
		} else {
			i++
		}
	}

	return false
}

func numWays(n int) int {
	if n < 2 {
		return 1
	}
	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1

	for i := 2; i <= n; i++ {
		dp[i] = (dp[i-1] + dp[i-2]) % 1000000007
	}

	return dp[n]
}

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

	dummyHead := &ListNode{}
	head := dummyHead
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

	return dummyHead.Next
}

func maxSlideWindow(nums []int, k int) []int {
	var res []int

	n := len(nums)

	for i := 0; i <= n-k; i++ {
		max := nums[i]
		for j := 1; j < k; j++ {
			if max < nums[i+j] {
				max = nums[i+j]
			}
		}
		res = append(res, max)
	}

	return res
}

func lengthOfLongestSubstring(s string) int {
	var freq [256]byte

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


type entry struct {
	key int
	value int
}

type cache struct {
	cap int
	cache map[int]*list.Elment
	lst *list.List
}

func Constructor(cap int) cache {
	return cache{
		cap,
		make(map[int]*list.Element),
		list.New(),
	}
}

func (c *cache) Get(key int) int {
	if e := c.cache[key]; e != nil {
		c.lst.MoveToFront(e)
		return e.Value.(entry).value
	}

	return -1
}

func (c *cache) Put(key, value int) {
	if e := c.cache[key]; e != nil {
		e = entry{key, value}
		c.lst.MoveToFront(e)
		return
	}

	c.cache[key] = c.lst.PushFront(entry{key, value})

	if c.cap < len(c.cache) {
		delete(c.cache, c.lst.Remove(c.lst.Back()).(entry).key)
	}
}

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

func reverseListK(head *ListNode, k int) {
	tail := head

	for i := 0; i < k; i++ {
		if tail == nil {
			return head
		}
		tail = tail.Next
	}

	prev := reverseListK(tail, k)

	for i := 0; i < k; i++ {
		next := head.Next
		head.Next = prev
		prev = head
		head = next
	}

	return prev
}

func maxSubArray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}

	res := nums[0]

	dp := make([]int, len(nums))
	dp[0] = nums[0]

	for i := 1; i < len(nums); i++ {
		if dp[i-1] < 0 {
			dp[i] = nums[i]
		} else {
			dp[i] = dp[i-1] + nums[i]
		}

		res = max(res, dp[i])
	}

	return res
}

func threeSum(nums []int) [][]int {
	var res [][]int
	if len(nums) == 0 {
		return res
	}

	m := make(map[int]int, len(nums))

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
		for j := i+1; j < len(uniqueNum); j++ {
			if uniqueNum[i]*2+uniqueNum[j] == 0 && m[uniqueNum[i]] >=2 {
				res = append(res, []int{uniqueNum[i], uniqueNum[i], uniqueNum[j]})
			}
			if uniqueNum[j]*2+uniqueNum[i] == 0 && m[uniqueNum[j]] >=2 {
				res = append(res, []int{uniqueNum[i], uniqueNum[j], uniqueNum[j]})
			}
			c := 0 - uniqueNum[i] - uniqueNum[j]
			if c > uniqueNum[j] && m[c] >= 1 {
				res = append(res, []int{uniqueNum[i], uniqueNum[j], c})
			}
		}
	}

	return res
} 

func mergeTwoLists(l1, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}

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
		l2.Next = mergeTwoLists(l1, l2.Next)
		return l2
	}

	return l1
}

func hasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}
	slow, fast := head, head

	for slow != nil && fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}

	return false
}

func getIntersectionNode(l1, l2 *ListNode) *ListNode {
	if l1 == nil || l2 == nil {
		return nil
	}

	a, b := l1, l2

	for a != b {
		if a == nil {
			a = l2
		} else {
			a = a.Next
		}
		if b == nil {
			b = l1
		} else {
			b = b.Next
		}
	}

	return a
}

func LevelOrder(root *TreeNode) [][]res {
	if root == nil {
		return [][]int{}
	}

	var res [][]int
	var flag bool

	q := []*TreeNode{root}
	curNum, nextNum := 1, 0
	tmp := []int{}

	for len(q) != 0 {
		if curNum != 0 {
			node := q[0]
			q = q[1:]

			if node.Left != nil {
				q = append(q, node.Left)
				nextNum++
			}
			if node.Right != nil {
				q = append(q, node.Right)
				nextNum++
			}

			tmp = append(tmp, node.Val)
			curNum--

		}
		if curNum == 0 {
			if !flag {
				for i := 0; i < len(tmp)/2; i++ {
					tmp[i], tmp[len(tmp)-i-1] = tmp[len(tmp)-i-1], tmp[i]
				}
			}

			res = append(res, tmp)
			tmp = []int{}
			curNum = nextNum
			nextNum = 0
			flag = !flag
		}
	}

	return res
}

func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}

	min, maxProfit := prices[0], 0

	for i := 1; i < len(prices); i++ {
		if prices[i]-min > maxProfit {
			maxProfit = prices[i]-min
		}
		if min > prices[i] {
			min = prices[i]
		}
	}

	return maxProfit
}

func addStrings(s1, s2 string) string {
	var res string

	n1, n2 := len(s1)-1, len(s2)-1
	carry := 0

	for ; n1 >= 0 || n2 >= 0 || carry != 0; {
		x, y := 0, 0 
		if n1 >= 0 {
			x = int(s1[n1]-'0')
		}
		if n2 >= 0 {
			y = int(s2[n2]-'0')
		}

		tmp := x + y + carry
		num := tmp % 10
		carry = tmp / 10
		res = strconv.Itoa(num) + res
		n1--
		n2--
	}


	return res
}

func merge(nums1 []int, m int, nums2 []int, n int) {
	i := m-1
	j := n-1
	k := m+n-1

	for ; i >= 0 && j >= 0; k-- {
		if nums1[i] < nums2[j] {
			nums1[k] = nums2[j]
			j--
		} else {
			nums1[k] = nums1[i]
			i--
		}
	}

	for ; j >= 0; k-- {
		nums[k] = nums2[j]
		j--
	}
}

func isValid(s string) bool {
	if s == "" {
		return true
	}

	stack := []rune{}

	for _, v := range s {
		if v == '(' || v == '{' || v == '[' {
			stack = append(stack, v)
		} else if v == ')' && len(stack) > 0 && stack[len(stack)-1] == '(' || v == '}' && len(stack) > 0 && stack[len(stack)-1] == '{' || v == ']' && len(stack) > 0 && stack[len(stack)-1] == '[' {
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
	}

	return len(stack) == 0
}

func sortArray(nums []int) []int {

	quickSort(nums, 0, len(nums)-1)

	return nums
}

func quickSort(nums []int, lo, hi int) {
	if lo >= hi {
		return
	}

	p := patition(nums, lo, hi)
	quickSort(nums, lo, p-1)
	quickSort(nums, p+1, hi)
}

func patition(nums []int, lo, hi int) int {
	pivot := nums[hi]
	i := lo - 1

	for j := lo; j < hi; j++ {
		if nums[j] < pivot {
			i++
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	nums[i+1], nums[hi] = nums[hi], nums[i+1]

	return i+1
}

func findKthLargest(nums []int, k int) int {
	sortArray(nums)

	return nums[len(nums)-k-1]
}

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == q || root == p {
		return root
	}
	left, right := lowestCommonAncestor(root.Left, p, q), lowestCommonAncestor(root.Right, p, q)
	if left != nil && right != nil {
		return root
	}
	if left == nil {
		return right
	}
	return left
}

func hasCycle(head *ListNode) (bool, *ListNode) {
	if head == nil {
		return false, nil
	}

	slow, fast := head, head

	for slow != nil && fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true, slow
		}
	}

	return false, nil
}

func detectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}

	isCycle, slow := hasCycle(head)

	if !isCycle {
		return nil
	}

	fast := head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}

	return fast
}

func findPermute(nums []int, index int, c []int, used *[]bool, ans *[][]int) {
	if index == len(nums) {
		p := make([]int, len(nums))
		copy(p, c)
		*ans = append(*ans, p)
		return
	}

	for i := 0; i < len(nums); i++ {
		if !(*used)[i] {
			(*used)[i] = true
			c = append(c, nums[i])
			findPermute(nums, c, index+1, used, ans)
			c = c[:len(c)-1]
			(*used)[i] = false
		}
	}
}

func permute(nums []int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}

	var res [][]int
	var c []int

	used := make([]bool, len(nums))

	findPermute(nums, 0, c, &used, &res)

	return res
}

func binarySearch(nums []int) int {
	low, high := 0, len(nums)-1

	for low <= high {
		mid := low + (high - low)>>1
		if target == nums[mid] {
			return mid
		} else if target > nums[mid] {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	var res []int

	helper(root, &res)

	return res
}

func helper(root *TreeNode, ans *[]int) {
	if root == nil {
		return nil
	}

	helper(root.Left, ans)
	*ans = append(*ans, root.Val)
	helper(root.Right, ans)
}

func inorderTraversal(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	var res []int

	stack := []*TreeNode{}

	for root != nil || len(stack) != 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}

		root = stack[len(stack)-1]
		res = append(res, root.Val)
		stack = stack[:len(stack)-1]

		root = root.Right
	}

	return res
}

func expand(s string, left, right int) (int, int) {
	for ; left >= 0 && right < len(s) && s[left] == s[right]; {
		left--
		right++
	}

	return left+1, right-1
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

	return s[start:end+1]
}


func spiralOrder(matrix [][]int) []int {
	m := len(matrix)
	if m == 0 {
		return nil
	}

	n := len(matrix[0])
	if n == 0 {
		return nil
	}

	left, top, bottom, right := 0, 0, m-1, n-1

	sum := m*n
	count := 0

	var res []int

	for count < sum {
		i, j := top, left
		for ; j <= right && count < sum; {
			res = append(res, matrix[i][j])
			j++
			count++
		}
		i, j = top+1, right
		for ; i <= bottom && count < sum; {
			res = append(res, matrix[i][j])
			i++
			count++
		}
		i, j = bottom, right-1
		for ; j >= left && count < sum; {
			res = append(res, matrix[i][j])
			j--
			count++
		}
		i, j = bottom-1, left
		for ; i > top && count < sum; {
			res = append(res, matrix[i][j])
			i--
			count++
		}

		left, top, bottom, right = left+1, top+1, bottom-1, right-1
	}

	return res
}

func search(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right - left)>>1
		if nums[mid] == target {
			return mid
		} else if nums[mid] > nums[left] {
			if target >= nums[left] && target < nums[mid] {
				right = mid-1
			} else {
				left = mid+1
			}
		} else if nums[mid] < nums[right] {
			if target <= nums[right] && target > nums[mid] {
				left = mid+1
			} else {
				right = mid-1
			}
		} else {
			if nums[mid] == nums[left] {
				left++
			}
			if nums[mid] == nums[right] {
				right--
			}
		}
	}

	return -1
}


var dir [][]int{
	{1, 0},
	{-1, 0},
	{0, 1},
	{0, -1},
}

func searchIsIlands(grid [][]byte, visisted *[][]bool, x, y int) {
	(*visisted)[x][y] = true
	for i := 0; i < 4; i++ {
		nx := x + dir[i][0]
		ny := y + dir[i][1]
		if isInBoard(grid, nx, ny) && !(*visisted)[nx][ny] && grid[nx][ny] == '1' {
			searchIsIlands(grid, visisted, nx, ny)
		}
	}
}

func numIslands(grid [][]byte) int {
	m := len(grid)
	if m == 0 {
		return 0
	}
	n := len(grid[0])
	if n == 0 {
		return 0
	}

	res, visisted := 0, make([][]bool, m)
	for i := 0; i < m; i++ {
		visisted[i] = make([]bool, n)
	}

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' && !visisted[i][j] {
				searchIsIlands(grid, &visited, i, j)
				res++
			}
		}
	}

	return res
}

func isInBoard(board [][]byte, x, y int) bool {
	return x >= 0 && x < len(board) && y >= 0 && y < len(board[0])
}

func lengthOfLIS(nums []int) int {
	dp := []int{}

	for _, num := range nums {
		i := sort.SearchInts(dp, num)
		if i == len(dp) {
			dp = append(dp, num)
		} else {
			dp[i] = num
		}
	}

	return len(dp)
}

type MyQueue struct {
	Stack *[]int
	Queue *[]int
}

func Constructor() MyQueue {
	stack := []int{}
	queue := []int{}

	return MyQueue{&stack, &queue}
}

func (q *MyQueue) Push(x int) {
	*q.Stack = append(*q.Stack, x)
}

func (q *MyQueue) Pop() int {
	if len(*q.Queue) == 0 {
		q.fromStackToQueue(q.Stack, q.Queue)
	}

	poped := (*q.Queue)[len(*q.Queue)-1]
	*q.Queue = (*q.Queue)[:len(*q.Queue)-1]

	return poped
}

func (q *MyQueue) Peek() int {
	if len(*q.Queue) == 0 {
		q.fromStackToQueue(q.Stack, q.Queue)
	}

	return (*q.Queue)[len(*q.Queue)-1]
}

func (q *MyQueue) Empty() bool {
	return len(*q.Stack) == 0 && len(*q.Queue) == 0
}

func (q *MyQueue) fromStackToQueue(stack, queue *[]int) {
	for len(*q.Stack) > 0 {
		poped := (*q.Stack)[len(*q.Stack)-1]
		*q.Stack = (*q.Stack)[:len(*q.Stack)-1]
		*q.Queue = append(*q.Queue, poped)
	}
}


func trap(nums []int) int {
	var res int

	left, right := 0, len(nums)-1

	maxLeft, maxRight := 0, 0

	for left <= right {
		if nums[left] < nums[right] {
			if nums[left] > maxLeft {
				maxLeft = nums[left]
			} else {
				res += maxLeft - nums[left]
			}
			left++
		} else {
			if nums[right] > maxRight {
				maxRight = nums[right]
			} else {
				res += maxRight - nums[right]
			}
			right--
		}
	}

	return res
}

func quickSort(nums []int, left, right int) int {
	pivot := nums[right]
	i := left - 1

	for j := lo; j < hi; j++ {
		i++
		if nums[j] < pivot {
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	nums[i+1], nums[right] = nums[right], nums[i+1]

	return i+1
}

func clean(s string) (int, string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return 1, ""
	}

	sign, numStr := 1, ""
	switch s[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		sign, numStr = 1, s
	case '+':
		sign, numStr = 1, s[1:]
	case '-':
		sign, numStr = -1, s[1:]
	default:
		sign, numStr = 1, ""
	}

	for i, v := range numStr {
		if v < '0' || v > '9' {
			numStr = numStr[:i]
			break
		}
	}

	return sign, numStr
}

func convert(sign int, numStr string) int {
	var absNum int

	for _, v := range numStr {
		absNum = absNum * 10 + int(v-'0')
		if absNum > math.MaxInt32 {
			if sign == 1 {
				return math.MaxInt32
			} else {
				return math.MinInt32
			}
		}
	}

	return absNum * sign
}

func myAtoi(s string) int {
	return convert(clean(s))
}

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	if x < 10 {
		return true
	}

	s := strconv.Itoa(x)
	length := len(s)

	for i := 0; i < length/2; i++ {
		if s[i] != s[length-i-1] {
			return false
		}
	}

	return true
}

func longestCommonPrefix(strs []string) string {

	prefix := strs[0]

	for i := 0; i < len(strs); i++ {
		index := 0
		length := min(len(prefix), len(strs[i]))

		for index < length && prefix[index] == strs[i][index] {
			index++
		}

		prefix = prefix[:index]
		if len(prefix) == 0 {
			return ""
		}
	}

	return prefix
}

func threeClosest(nums []int, target int) int {
	var res int

	diff := math.MaxInt32
	for i := 0; i < len(nums); i++ {
		for j := i+1; j < len(nums); j++ {
			for k := j+1; k < len(nums); k++ {
				if abs(nums[i]+nums[j]+nums[k]-target) < diff {
					diff = abs(nums[i], nums[j], nums[k]-target)
					res = nums[i] + nums[j] + nums[k]
				}
			}
		}
	}

	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
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


func letterCombinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}
	var res []string

	helper(&digits, 0, "", &res)

	return res
}

func helper(digits *string, index int, c string, res *[]string) {
	if len(*digits) == index {
		*res = append(*res, c)
		return
	}

	num := (*digits)[index]
	letters := letterMap[num-'0']
	for i := 0; i < len(letters); i++ {
		helper(digits, index+1, c+string(letters[i]), res)
	}
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	if head == nil {
		return nil
	}

	slow, fast := head, head

	for i := 0; i < n; i++ {
		fast = fast.Next
	}

	if fast == nil {
		return head.Next
	}

	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}

	slow.Next = slow.Next.Next

	return head
}

func findGeneration(n int) []string {
	if n == 0 {
		return []string{}
	}

	res := []string{}
	helper(n, n, "", &res)
}

func helper(left, right int, c string, res *[]string) {
	if left == 0 && right == 0 {
		*res = append(*res, c)
		return
	}

	if left > 0 {
		helper(left-1, right, c+"(", res)
	}
	if right > 0 && left < right {
		helper(left, right-1, c+")", res)
	}
}

func reverseGroup(head *ListNode, k int) *ListNode {
	tail := head

	for i := 0; i < k; i++ {
		if tail == nil {
			return head
		}
		tail = tail.Next
	}

	prev := reverseGroup(tail, k)

	for i := 0; i < k; i++ {
		next := head.Next
		head.Next = prev
		prev = head
		head = next
	}

	return prev
}

func removeDuplicates(nums []int) int {
	for i := len(nums)-1; i > 0; i-- {
		if nums[i] == nums[i-1] {
			nums = append(nums[:i], nums[i+1:]...)
		}
	}
	return len(nums)
}

func strStr(haystack string, needle string) int {
	for i := 0; ; i++ {
		for j := 0; ; j++ {
			if j == len(needle) {
				return i
			}
			if i+j == len(haystack) {
				return -1
			}
			if needle[j] != haystack[i+j] {
				break
			}
		}
	}

	return -1
}

func search(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)<<1

		if nums[mid] == target {
			return mid
		} else if nums[mid] > nums[left] {
			if target >= nums[left] && nums[mid] > target {
				right = mid-1
			} else {
				left = mid + 1
			}
		} else if nums[mid] < nums[right] {
			if target <= nums[right] && nums[mid] < target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		} else {
			if nums[mid] == nums[left] {
				left++
			}
			if nums[mid] == nums[right] {
				right--
			}
		}
	}

	return -1
}

func searchFirst(nums []int, target int) int {
	if len(nums) == 0 {
		return -1
	}

	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)>>1
		if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		} else {
			if (mid == 0) || nums[mid-1] != target {
				return mid
			}
			right = mid-1
		}
	}
}

func combinationSum(nums []int, target int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}

	var res [][]Int
	var c []int

	sort.Ints(nums)

	helper(nums, 0, target, c, &res)

	return res
}

func helper(nums []int, index int, target int, c []int, res *[][]int) {
	if target <= 0 {
		if target == 0 {
			tmp := make([]int, len(c))
			copy(tmp, c)
			*res = append(*res, p)
		}
		return
	}

	for i := index; i < len(nums); i++ {
		if nums[i] > target {
			break
		}
		c = append(c, nums[i])
		helper(nums, i, target-nums[i], c, res)
		c = c[:len(c)-1]
	}
}

func findPermute(nums []int, index int, c []int, used *[]bool, res *[][]int) {
	if len(nums) == index {
		p := make([]int, len(nums))
		copy(p, c)
		*res = append(*res, p)
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
	if len(nums) == 0 {
		return [][]int{}
	}

	var res [][]int
	var c []int

	used := make([]bool, len(nums))

	findPermute(nums, 0, c, &used, &res)

	return res
}

func maxArea(nums []int) int {

	width := 0
	height := 0
	res := 0

	left, right := 0, len(nums)-1

	for left <= right {
		width = right - left
		if nums[left] < nums[right] {
			height = nums[left]
			left++
		} else {
			height = nums[right]
			right--
		}
		res = max(res, width*height)
	}

	return res
}

func compressString(s string) string {
	slen := len(s)

	res := []byte{}

	for i := 0; i < slen; i++ {
		count := 1

		for i+1 < slen; S[i] == S[i+1] {
			i++
			count++
		}

		num := strconv.Itoa(count)
		res = append(res, S[i])
		res = append(res, []byte(num)...)
	}

	return string(res)
}

func isFlipedString(s1, s2 string) bool {
	if len(s1) == 0 && len(s2) == 0 {
		return true
	}

	s := s1+s2

	return len(s2) > 0 && strings.Index(s, s2) > -1
}

func removeDuplicateNodes(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	slow := head

	for slow != nil {
		fast := slow 
		for fast.Next != nil {
			if slow.Val == fast.Next.Val {
				fast.Next = fast.Next.Next
			} else {
				fast = fast.Next
			}
		}
		slow = slow.Next
	}

}

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}

	dummy := &ListNode{}

	for headA != nil {
		dummy = headB
		for dummy != nil {
			if dummy == headA {
				return dummy
			}
			dummy = dummy.Next
		}
		headA = headA.Next
	}

	return nil
}

type MinSatck struct {
	stack []int
	min []int
}

func Constructor() MinStack {
	return MinStack{
		stack: []int,
		min: []int,
	}
}

func (m *MinSatck) Push(val int) {
	if len(m.min) == 0 {
		m.min = append(m.min, val)
	} else {
		minVal := m.min[len(m.min)-1]
		if minVal < val {
			m.min = append(m.min, minVal)
		} else {
			m.min = append(m.min, val)
		}
	}

	m.stack = append(m.stack, val)
}

func (this *MinStack) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
	this.min = this.min[:len(this.min)-1]
}


func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1]
}


func (this *MinStack) GetMin() int {
	return this.min[len(this.min)-1]
}

func sortedArrayToBST(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}
	if len(nums) == 1 {
		return &TreeNode{Val: nums[0]}
	}

	mid := len(nums) / 2
	head := &TreeNode{Val: nums[mid]}
	head.Left = sortedArrayToBST(nums[:mid])
	head.Right = sortedArrayToBST(nums[mid+1:])

	return head
}

func isBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return abs(height(root.Left) - height(root.Right)) <= 1 && isBalanced(root.Left) && isBalanced(root.Right)
}

func height(root *TreeNode) int {
	if root == nil {
		return 0
	}

	return max(height(root.Left), height(root.Right)) + 1
}

func inorderSuccessor(root *TreeNode, p *TreeNode) *TreeNode {
	var (
		ans *TreeNode
		flag bool
		f func(*TreeNode)
	)

	f = func(r *TreeNode) {
		if r != nil && ans == nil {
			f(r.Left)
			if p == r {
				flag = true
			} else if flag {
				ans, flag = r, false
			}
			f(r.Right)
		}
	}

	f(root)

	return ans
}

func exchangeBits(num int) int {
	even := num & 0xaaaaaaa
	odd := num & 0x55555555

	return even >> 1 || odd << 1
}

func canPermutePalindrome(s string) bool {
	m := make(map[rune]int)

	for _, v := range s {
		m[v]++
	}

	var n int

	for _, v := range m {
		if v != 2 {
			n++
			if n == 2 {
				return false
			}
		}
	}

	return true
}

func climbStair(n int) int {
	if n < 2 {
		return 1
	}

	dp := make([]int, n+1)

	for i := 2; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n]
}

func mySqrt(x int) int {
	if x == 0 {
		return 0
	}

	left, right, res := 1, x, -1

	for left <= right {
		mid := left + (right-left)>>1
		if mid == x/mid {
			return mid
		} else if mid < x/mid {
			left = mid+1
			res = mid
		} else {
			right = mid-1
		}
	}

	return res
}

func main() {
	letter, number := make(chan struct{}), make(chan struct{})

	wg := sync.WaitGroup{}

	go func() {
		i := 1
		for {
			select {
			case <- number:
				println(i)
				i++
				println(i)
				i++
				letter <- struct {}{}
			}
		}
	}()
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		i := 'A'
		for {
			select {
			case <- letter:
				if i >= 'Z' {
					wg.Done()
					return
				}
				println(string(i))
				i++
				println(string(i))
				i++
				number <- struct{}{}
			}
		}
	}(&wg)

	number <- struct{}{}
	wg.Wait()

}

func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	return max(maxDepth(root.Left), maxDepth(root.Right)) + 1
}

func nextPermute(nums []int) {
	n := len(nums)
	i := n-2

	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	} 

	if i >= 0 {
		j := n-1
		for j >= 0 && nums[i] >= nums[j] {
			j--
		}
		nums[i], nums[j] = nums[j], nums[i]
	}

	reverse(nums[i+1:])
}


func wallsAndGates(rooms [][]int)  {
	if len(rooms) == 0 || len(rooms[0]) == 0 {
		return
	}

	m, n := len(rooms), len(rooms[0])

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if rooms[i][j] == 0 {
				helper(rooms, i, j, 0)
			}
		}
	}
}

func helper(rooms [][]int, i, j, step int) {
	m, n := len(rooms), len(rooms[0])
	if i < 0 || i >= m || j < 0 || j >= n || rooms[i][j] < step {
		return
	}

	rooms[i][j] = step

	helper(rooms, i+1, j, step+1)
	helper(rooms, i-1, j, step+1)
	helper(rooms, i, j+1, step+1)
	helper(rooms, i, j-1, step+1)
}

func mergeTwoTrees(t1, t2 *TreeNode) *TreeNode {
	if t1 == nil {
		return t2
	}
	if t2 == nil {
		return t1
	}

	t1.Val += t2.Val
	t1.Left = mergeTwoTrees(t1.Left, t2.Left)
	t1.Right = mergeTwoTrees(t1.Right, t2.Right)

	return t1
}

func twoSum(nums []int, target int) [][]int {
	if len(nums) == 0 {
		return [][]int{}
	}
	sort.Ints(nums)

	var res [][]int

	left, right := 0, len(nums)-1

	for left < right {
		sum := nums[left]+nums[right]
		if sum == target {
			res = append(res, []int{nums[left], nums[right]})
			left++
			right--
			for ; left < right && nums[left] == nums[left+1]; {
				left++
			}
			for ; left < right && nums[right] == nums[right-1]; {
				right--
			}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}

	return res
}

func convert(s string, numRow int) string {
	if s == "" {
		return ""
	}

	if numRow < 2 {
		return s
	}

	i, flag := 0, -1
	res = make([]string, len(s))

	for _, v := range s {
		res[i] += string(v)
		if i == 0 || i == numRow-1 {
			flag = -flag
		}
		i += flag
	}

	return string.Join(res, "")
}

func intToRoman(num int) string {
	roman := []string{"I", "IV", "V", "IX", "X", "XL", "L", "XC", "C", "CD", "D", "CM", "M"}
	intSlice := []int{1, 4, 5, 9, 10, 40, 50, 90, 100, 400, 500, 900, 1000}


	var res string

	for i := len(intSlice)-1; i >= 0; i-- {
		if num < intSlice[i] {
			continue
		}
		repeats := num/intSlice[i]
		res += strings.Repeat(roman[i], repeats)
		num -= intSlice[i] * repeats

		if num == 0 {
			break
		}
	}

	return res
}


var roman = map[string]int{
	"I": 1,
	"V": 5,
	"X": 10,
	"L": 50,
	"C": 100,
	"D": 500,
	"M": 1000,
}

func romanToInt(s string) int {
	if s == "" {
		return 0
	}	

	lastInt, num, total := 0, 0, 0

	for i := 0; i < len(s); i++ {
		c := s[len(s)-(i+1):len(s)-i]
		num = roman[c]
		if num < lastInt {
			total = total - num
		} else {
			total = total + num
		}

		lastInt = num
	}

	return total	
}

func longestCommonPrefix(strs []string) string {
	prefix := strs[0]

	for i := 0; i < len(strs); i++ {
		index := 0
		length := min(len(strs[i]), len(prefix))

		for index < length && strs[i][index] == prefix[index] {
			index++
		}

		prefix = prefix[:index]
		if prefix == "" {
			return ""
		} 
	}

	return prefix
}

func helper(digits string, index int, c string, res *[]string) {
	if len(digits) == index {
		*res = append(*res, c)
	}

	num := digits[index]
	letters := letterMap[num-'0']
	for i := 0; i < len(letters); i++ {
		helper(digits, index+1, c+string(letters[i]), res)
	}
}


func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := head.Next
	head.Next = swapPairs(newHead.Next)
	newHead.Next = head

	return newHead
}

func reverseGroup(head *ListNode, k int) *ListNode {
	if head == nil {
		return nil
	}

	tail := head

	for i := 0; i < k; i++ {
		if tail == nil {
			return head
		}
		tail = tail.Next
	}

	prev := reverseGroup(tail, k)

	for i := 0; i < k; i++ {
		next := head.Next
		head.Next = prev
		prev = head
		head = next
	}

	return prev
}

func removeDuplicates(nums []int) int {

	if len(nums) == 0 {
		return 0
	}

	for i := len(nums)-1; i > 0; i-- {
		if nums[i] == nums[i-1] {
			nums = append(nums[:i], nums[i+1:]...)
		}
	}

	return len(nums)
}

func strStr(a, b string) int {
	for i := 0; ; i++ {
		for j := 0; ; j++ {
			if j == len(b) {
				return i
			}
			if i+j == len(a) {
				return -1
			}
			if a[i+j] != b[j] {
				break
			}
		}
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
		} else if nums[left] < nums[mid] {
			if nums[left] <= target && target < nums[mid] {
				right = mid-1
			} else {
				left = mid+1
			}
		} else if nums[right] > nums[mid] {
			if nums[right] >= target && target > nums[mid] {
				left = mid+1
			} else {
				right = mid-1
			}
		} else {
			if nums[mid] == nums[left] {
				left++
			}
			if nums[mid] == nums[right] {
				right--
			}
		}
	}

	return -1
}


func firstMissingPositive(nums []int) int {
	sort.Ints(nums)

	v := 1
	for i := 0; i < len(nums); i++ {
		if v == nums[i] {
			v++
		}
	}

	return v
}

func findPermute(nums []int, index int, c []int, used *[]bool, res *[][]int) {
	if len(nums) == index {
		q := make([]int, len(nums))
		copy(q, c)
		*res = append(*res, q)
	}

	for i := 0; i < len(nums); i++ {
		if !(*used)[i] {
			(*used)[i] = true
			c = append(c, nums[i])
			findPermute(nums, index+1, c, used, res)
			(*used)[i] = false
			c = c[:len(c)-1]
		}
	}
}

func rotate(matrix [][]int) {
	row := len(matrix)
	if row == 0 {
		return 
	}

	col := len(matrix[0])

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	halfCol := col/2

	for i := 0; i < row; i++ {
		for j := 0; j < halfCol; j++ {
			matrix[i][j], matrix[i][col-j-1] = matrix[i][col-j-1], matrix[i][j]
		}
	}
}

func myPow(x, n int) float64 {
	if n >= 0 {
		return quickPow(x, n)
	}
	return 1.0 / quickPow(x, n)
}

func quickPow(x, n int) float64 {
	if n == 0 {
		return 1
	}
	y := quickPow(x, n/2)
	if n%2 == 0 {
		return y*y
	}
	return y*y*x
}

func maxSubarray(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	dp := make([]int, len(nums))

	dp[0] = nums[0]
	res := nums[0]

	for i := 1; i < len(nums); i++ {
		if dp[i-1] < 0 {
			dp[i] = nums[i]
		} else {
			dp[i] = dp[i-1] + nums[i]
		} 
		res = max(res, dp[i])
	}

	return res
}

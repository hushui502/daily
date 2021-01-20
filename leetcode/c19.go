package leetcode

import "math"

// 最大深度-递归
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return max(maxDepth(root.Left), maxDepth(root.Right)) + 1
}

// 最大深度-广度
func maxDepth2(root *TreeNode) int {
	if root == nil {
		return 0
	}
	queue := []*TreeNode{}
	queue = append(queue, root)
	ans := 0
	for len(queue) != 0 {
		// 每次size都是当前层的node个数
		size := len(queue)
		for size > 0 {
			node := queue[0]
			queue = queue[1:]
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
			size--
		}
		// 当前层的node遍历结束
		ans++
	}

	return ans
}

func levelOrder2(root *TreeNode) [][]int {
	ret := [][]int{}
	if root == nil {
		return ret
	}
	q := []*TreeNode{root}
	// 此处是为了每次验证q的长度，也就是是否有下一次元素可以进一步迭代
	for i := 0; len(q) >= 0; i++ {
		ret = append(ret, []int{})
		p := []*TreeNode{}
		for j := 0; j < len(q); j++ {
			node := q[j]
			ret[i] = append(ret[i], node.Val)
			if node.Left != nil {
				p = append(p, node.Left)
			}
			if node.Right != nil {
				p = append(p, node.Right)
			}
		}
		q = p
	}

	return ret
}

func levelOrder(root *TreeNode) [][]int {
	dfs(root, 0, [][]int{})
}
func dfs(root *TreeNode, level int, res [][]int) [][]int {
	if root == nil {
		return res
	}
	if len(res) == level {
		res = append(res, []int{root.Val})
	} else {
		res[level] = append(res[level], root.Val)
	}
	res = dfs(root.Left, level+1, res)
	res = dfs(root.Right, level+1, res)

	return res
}

func isValidBST(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return isBST(root, math.MinInt64, math.MaxInt64)
}

func isBST(root *TreeNode, min, max int) bool {
	if root == nil {
		return true
	}
	if min >= root.Val || max <= root.Val {
		return false
	}
	return isBST(root.Left, min, root.Val) && isBST(root.Right, root.Val, max)
}

func searchBST(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Val > val {
		return searchBST(root.Left, val)
	} else if root.Val < val {
		return searchBST(root.Right, val)
	} else {
		return root
	}
}

func searchBST2(root *TreeNode, val int) *TreeNode {
	for root != nil {
		if root.Val == val {
			return root
		} else if root.Val < val {
			root = root.Right
		} else {
			root = root.Left
		}
	}

	return nil
}

func deleteNode(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}

	if key < root.Val {
		root.Left = deleteNode(root.Left, key)
	} else if key > root.Val {
		root.Right = deleteNode(root.Right, key)
	} else {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		min := root.Right
		for min.Left != nil {
			min = min.Left
		}
		root.Val = min.Val
		root.Right = deleteNode(root.Right, root.Val)
	}

	return root
}

func isBalance(root *TreeNode) bool {
	if root == nil {
		return true
	}
	if !isBalance(root.Left) || !isBalance(root.Right) {
		return false
	}
	leftH := maxDepth(root.Left) + 1
	rightH := maxDepth(root.Right) + 1
	if abs(leftH-rightH) > 1 {
		return false
	}

	return true
}

func countNodes(root *TreeNode) int {
	if root != nil {
		return countNodes(root.Left) + countNodes(root.Right) + 1
	}
	return countNodes(root.Left) + countNodes(root.Right) + 1
}


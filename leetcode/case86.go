package leetcode

func minDepath(root *TreeNode) int {
	if root == nil {
		return 0
	}

	if root.Left == nil {
		return minDepath(root.Right) + 1
	}

	if root.Right == nil {
		return minDepath(root.Left) + 1
	}

	return min(minDepath(root.Left), minDepath(root.Right)) + 1
}

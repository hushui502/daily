package leetcode

func isBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	leftHight := depth(root.Left)
	rightHight := depth(root.Right)

	return abs(leftHight-rightHight) <= 1 && isBalanced(root.Left) && isBalanced(root.Right)
}

func abs(x int) int {
	// todo 溢出
	if x < 0 {
		return -x
	}

	return x
}


func depth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	return max(depth(root.Left), depth(root.Right)) + 1
}

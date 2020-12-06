package leetcode

func isSymmetic(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return isSameTree(invertTree(root.Left), root.Right)
}

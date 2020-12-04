package leetcode

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderTraverseal(root *TreeNode) []int {
	var result []int
	inorder(root, &result)
	return result
}

func inorder(root *TreeNode, output *[]int) {
	if root != nil {
		inorder(root.Left, output)
		*output = append(*output, root.Val)
		inorder(root.Right, output)
	}
}



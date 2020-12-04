package leetcode

import "math"

func isValidbst(root *TreeNode, min, max float64) bool {
	if root == nil {
		return true
	}

	v := float64(root.Val)

	return v < max && v > min && isValidbst(root.Left, min, v) && isValidbst(root.Right, v, max)
}

func isValidBST(root *TreeNode) bool {
	return isValidbst(root, math.Inf(-1), math.Inf(1))
}

func isValidBST1(root *TreeNode) bool {
	arr := []int{}
	inOrder(root, &arr)
	for i := 1; i < len(arr); i++ {
		if arr[i-1] > arr[i] {
			return false
		}
	}

	return true
}

func inOrder(root *TreeNode, arr *[]int) {
	if root == nil {
		return
	}

	inOrder(root.Left, arr)
	*arr = append(*arr, root.Val)
	inorder(root.Right, arr)
}

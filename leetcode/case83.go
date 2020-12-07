package leetcode

func levelOrderBottom(root *TreeNode) [][]int {
	tmp := levelOrder(root)
	res := [][]int{}

	for i := len(tmp)-1; i >= 0; i-- {
		res = append(res, tmp[i])
	}

	return res
}

package leetcode

func flatten(root *TreeNode) {
	list, cur := []int{}, &TreeNode{}
	preorder(root, &list)
	cur = root
	for i := 1; i < len(list); i++ {
		cur.Left = nil
		cur.Right = &TreeNode{Val:list[i], Left:nil, Right:nil}
		cur = cur.Right
	}

	return
}

func flatten1(root *TreeNode) {
	if root == nil || (root.Left == nil && root.Right == nil) {
		return
	}

	flatten(root.Left)
	flatten(root.Right)
	currRight := root.Right
	root.Right = root.Left
	root.Left = nil

	for root.Right != nil {
		root = root.Right
	}

	root.Right = currRight
}
package leetcode

func buildTree1(preorder []int, inorder []int) *TreeNode {
	for k := range inorder {
		if inorder[k] == preorder[0] {
			return &TreeNode{
				Val:preorder[0],
				Left:buildTree1(preorder[1:k+1], inorder[0:k]),
				Right:buildTree1(preorder[k+1:], inorder[k+1:]),
			}
		}
	}

	return nil
}
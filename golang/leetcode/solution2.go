package main

import (
	"math"
)

//Definition for a1 binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

var ans int

func diameterOfBinaryTree(root *TreeNode) int {
	ans = 1
	depth(root)
	return ans - 1
}

func depth(node *TreeNode) int {
	if node == nil {
		return 0
	}
	L := depth(node.Left)
	R := depth(node.Right)
	ans = int(math.Max(float64(ans), float64(R+L+1)))
	return int(math.Max(float64(L), float64(R)) + 1)
}

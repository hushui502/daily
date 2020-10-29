package binary_tree

import "fmt"

type node struct {
	val   int
	left  *node
	right *node
}

type btree struct {
	root *node
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func newNode(val int) *node {
	n := &node{val: val, left: nil, right: nil}
	return n
}

func inorder(n *node) {
	if n == nil {
		return
	}
	fmt.Print(n.val, " ")
	inorder(n.left)
	inorder(n.right)
}

func postorder(n *node) {
	if n == nil {
		return
	}
	inorder(n.left)
	inorder(n.right)
	fmt.Print(n.val, " ")
}

func levelorder(root *node) {
	var q []*node
	var n *node

	q = append(q, root)
	for len(q) != 0 {
		n, q = q[0], q[1:]
		fmt.Print(n.val, " ")
		if n.left != nil {
			q = append(q, n.left)
		}
		if n.right != nil {
			q = append(q, n.right)
		}
	}
}

func _calculate_depth(n *node, depth int) int {
	if n == nil {
		return depth
	}
	return max(_calculate_depth(n.left, depth+1), _calculate_depth(n.right, depth+1))
}

func (t *btree) depth() int {
	return _calculate_depth(t.root, 0)
}

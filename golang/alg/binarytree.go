package alg

type BinarySearchTreeNode struct {
	value int
	times int
	left *BinarySearchTreeNode
	right *BinarySearchTreeNode
}

type BinarySearchTree struct {
	Root *BinarySearchTreeNode
}

func NewBinarySearchTree() *BinarySearchTree {
	return new(BinarySearchTree)
}

func (node *BinarySearchTreeNode) Add(value int) {
	if value < node.value {

	}
}

func (node *BinarySearchTreeNode) Search(value int)  {

}
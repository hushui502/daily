package trie

type Node struct {
	children map[rune]*Node
	isLeaf   bool
}

func NewNode() *Node {
	n := &Node{}
	n.children = make(map[rune]*Node)
	n.isLeaf = false
	return n
}

func (n *Node) Insert(s string) {
	curr := n
	for _, c := range s {
		next, ok := curr.children[c]
		if !ok {
			next = NewNode()
			curr.children[c] = next
		}
		curr = next
	}
	curr.isLeaf = true
}

func (n *Node) Find(s string) bool {
	curr := n
	for _, c := range s {
		next, ok := curr.children[c]
		if !ok {
			return false
		}
		curr = next
	}
	return true
}

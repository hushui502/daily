package skip

import "data-structure/common"

type widths []uint64

type nodes []*node

type node struct {
	forward nodes
	widths widths
	entry common.Comparator
}

func (n *node) Compare(e common.Comparator) int {
	return n.entry.Compare(e)
}

func newNode(cmp common.Comparator, maxLevels uint8) *node {
	return &node{
		forward: make(nodes, maxLevels),
		widths:  make(widths, maxLevels),
		entry:   cmp,
	}
}



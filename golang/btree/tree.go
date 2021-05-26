package main

type BTree map[int]node

func NewBTree() *BTree {
	bt := BTree{}
	leaf := newLeafNode(nil)
	r := newInteriorNode(nil, leaf)
	leaf.p = r
	bt[-1] = r
	bt[0] = leaf
	return &bt
}

func (bt *BTree) Count() int {
	count := 0
	leaf := (*bt)[0].(*leafNode)
	for leaf != nil {
		count += leaf.countNum()
		leaf = leaf.next
	}
	return count
}

// 返回根结点
func (bt *BTree) Root() node {
	return (*bt)[-1]
}

// 返回 第一个叶子结点
func (bt *BTree) First() node {
	return (*bt)[0]
}

func (bt *BTree) Values() []*leafNode {
	nodes := make([]*leafNode, 0)
	leaf := (*bt)[0].(*leafNode)
	for leaf != nil {
		nodes = append(nodes, leaf)
		leaf = leaf.next
	}

	return nodes
}

// 在 B+ 树中，插入 key-value
func (bt *BTree) Insert(key int, value string) {
	_, oldIndex, leaf := search((*bt)[-1], key)
	p := leaf.parent()
	mid, nextLeaf, bump := leaf.insert(key, value)
	// 未分裂，则直接返回
	if !bump {
		return
	}

	(*bt)[mid] = nextLeaf
	// 分裂则继续插入中间节点
	var midNode node
	midNode = leaf
	// 设置父亲节点指向分裂出的子（叶子）节点
	p.kcs[oldIndex].child = leaf.next
	// 新分裂出的节点设置该中间节点为父亲节点
	leaf.next.setParent(p)
	// 赋值，获取该中间节点和其父节点
	interior, interiorP := p, p.parent()
	// 迭代向上判断父亲节点是否需要分裂
	for {
		var oldIndex int
		var newNode *interiorNode
		// 判断是否到达根节点
		isRoot := interiorP == nil
		// 未到达根节点，在父亲节点中查询元素的位置
		if !isRoot {
			oldIndex, _ = interiorP.find(key)
		}
		// 将叶子节点分裂后产生的中间元素同时传给父亲中间节点，并传入分裂后的原始叶子节点
		// 同时返回分裂后产生的中间节点和中间元素的 key
		mid, newNode, bump = interior.insert(mid, midNode)
		// 未分裂，直接返回
		if !bump {
			return
		}
		// 填充 map
		(*bt)[newNode.kcs[0].key] = newNode

		if !isRoot {
			// 未到达根节点，将元素插入父亲节点，基本过程同上
			interiorP.kcs[oldIndex].child = newNode
			newNode.setParent(interiorP)

			midNode = interior
		} else {
			// 更新 map 中的 root 节点
			(*bt)[interior.kcs[0].key] = (*bt)[-1]
			// 到达根节点，根节点上移，并插入原始中间节点
			(*bt)[-1] = newInteriorNode(nil, newNode)
			node := (*bt)[-1].(*interiorNode)
			node.insert(mid, interior)
			(*bt)[-1] = node
			newNode.setParent(node)

			return
		}
		// 赋值，获取该中间节点的父亲节点和其父亲的父节点
		interior, interiorP = interiorP, interior.parent()
	}
}

func (bt *BTree) Search(key int) (string, bool) {
	kv, _, _ := search((*bt)[-1], key)
	if kv == nil {
		return "", false
	}
	return kv.value, true
}

func search(n node, key int) (*kv, int, *leafNode) {
	curr := n
	oldIndex := -1

	for {
		switch t := curr.(type) {
		case *leafNode:
			i, ok := t.find(key)
			if !ok {
				return nil, oldIndex, t
			}
			return &t.kvs[i], oldIndex, t
		case *interiorNode:
			i, _ := t.find(key)
			curr = t.kcs[i].child
			oldIndex = i
		default:
			panic("")
		}
	}
}

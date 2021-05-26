package main

import "sort"

type kc struct {
	key   int
	child node
}

// 预留一个空槽
type kcs [MaxKC + 1]kc

func (a *kcs) Len() int {
	return len(a)
}

func (a *kcs) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a *kcs) Less(i, j int) bool {
	if a[i].key == 0 {
		return false
	}
	if a[j].key == 0 {
		return true
	}

	return a[i].key < a[j].key
}

// 中间节点
type interiorNode struct {
	// 存储元素
	kcs kcs
	// 实际存储元素数目
	count int
	// 父亲节点
	p *interiorNode
}

func newInteriorNode(p *interiorNode, largestChild node) *interiorNode {
	i := &interiorNode{
		p:     p,
		count: 1,
	}
	if largestChild != nil {
		i.kcs[0].child = largestChild
	}

	return i
}

func (in *interiorNode) find(key int) (int, bool) {
	c := func(i int) bool {
		return in.kcs[i].key > key
	}

	i := sort.Search(in.count-1, c)

	if i == -1 {
		return -1, false
	}

	return i, true
}

func (in *interiorNode) insert(key int, child node) (int, *interiorNode, bool) {
	// 确定key应该在中间节点中的位置
	i, _ := in.find(key)

	if !in.full() {
		copy(in.kcs[i+1:], in.kcs[i:in.count])
		in.kcs[i].key = key
		in.kcs[i].child = child
		child.setParent(in)
		in.count++
		return 0, nil, false
	}

	// 达到数量限制，在最后侧的空槽追加
	in.kcs[MaxKC].key = key
	in.kcs[MaxKC].child = child
	child.setParent(in)
	// 中间节点分裂
	next, midKey := in.split()

	return midKey, next, true
}

func (in *interiorNode) split() (*interiorNode, int) {
	sort.Sort(&in.kcs)

	// 获取中间元素的位置
	midIndex := MaxKC / 2
	midChild := in.kcs[midIndex].child
	midKey := in.kcs[midIndex].key

	// 创建一个新没有父亲节点(第一个 nil)的中间节点
	next := newInteriorNode(nil, nil)
	// 将中间元素的右侧数组拷贝到新的分裂节点
	copy(next.kcs[0:], in.kcs[midIndex+1:])
	// 初始化原始节点的右半部分
	in.initArray(midIndex + 1)
	next.count = MaxKC - midIndex
	// 更新分裂节点中所有元素子节点的父亲节点
	for i := 0; i < next.count; i++ {
		next.kcs[i].child.setParent(next)
	}

	// 更新原始节点的参数，将中间元素放进原始节点
	in.count = midIndex + 1
	in.kcs[in.count-1].key = 0
	in.kcs[in.count-1].child = midChild
	midChild.setParent(in)

	// 返回分裂后产生的中间节点和中间元素的 key，供父亲节点插入
	return next, midKey
}

// 判断是否达到中间节点的最大元素数目限制 MaxKC
func (in *interiorNode) full() bool {
	return in.count == MaxKC
}

// 返回中间节点的父亲节点
func (in *interiorNode) parent() *interiorNode {
	return in.p
}

func (in *interiorNode) setParent(p *interiorNode) {
	in.p = p
}

func (in *interiorNode) countNum() int {
	return in.count
}

// 初始化数组从 num 起的元素为空结构
func (in *interiorNode) initArray(num int) {
	for i := num; i < len(in.kcs); i++ {
		in.kcs[i] = kc{}
	}
}

package main

const (
	// 叶子节点最大的存储数目
	MaxKV = 255
	// 中间节点最大的存储数目
	MaxKC = 511
)

type node interface {
	// 确定元素在节点中的位置
	find(key int) (int, bool)
	// 获取父亲节点
	parent() *interiorNode
	// 设置父亲节点
	setParent(*interiorNode)
	// 是否达到最大数目
	full() bool
	// 元素数目统计
	countNum() int
}

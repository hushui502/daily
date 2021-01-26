package main

import (
	"sync/atomic"
	"unsafe"
)

type LFStack struct {
	head unsafe.Pointer	// 栈顶
}

type Node struct {
	val int32
	next unsafe.Pointer
}

func NewLFStack() *LFStack {
	n := unsafe.Pointer(&Node{})
	return &LFStack{head: n}
}

func (s *LFStack) push(v int32) {
	n := &Node{val: v}
	for {
		// 取出栈顶
		old := atomic.LoadPointer(&s.head)
		// 替换 推入栈顶
		if atomic.CompareAndSwapPointer(&s.head, old, unsafe.Pointer(n)) {
			return
		}
	}
}

func (s *LFStack) Pop() int32 {
	for {
		old := atomic.LoadPointer(&s.head)
		if old == nil {
			return 0
		}
		oldNode := (*Node)(old)
		// 取下一个节点
		next := atomic.LoadPointer(&oldNode.next)
		// 重置栈顶
		if atomic.CompareAndSwapPointer(&s.head, old, next) {
			return oldNode.val
		}
	}
}
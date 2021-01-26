# Atomic

## 案例
```
type Config struct {
	v atomic.Value // 假设 data 就是整个 config 了
}

// Get get config data
func (c *Config) Get() []int {
	return (*c.v.Load().(*[]int))
}

// Set set config data
func (c *Config) Set(n []int) {
	c.v.Store(&n)
}
```
atomic的性能是高于加锁的，但这种情况适合写情况很少的场景，因为他是copy on write


- 写入时复制（英语：Copy-on-write，简称 COW）是一种计算机程序设计领域的优化策略。其核心思想是，如果有多个调用者（callers）同时请求相同资源（如内存或磁盘上的数据存储），他们会共同获取相同的指针指向相同的资源，直到某个调用者试图修改资源的内容时，系统才会真正复制一份专用副本（private copy）给该调用者，而其他调用者所见到的最初的资源仍然保持不变。这过程对其他的调用者都是透明的。此作法主要的优点是如果调用者没有修改该资源，就不会有副本（private copy）被创建，因此多个调用者只是读取操作时可以共享同一份资源。（wiki）

这种的问题是如果有部分goroutine使用老对象，那么老对象不会被回收，如果这时候大量写入，会产生大量的副本，性能反而会很差。

## 源码分析
- 第一类AddXXX
```
func AddInt32(addr *int32, delta int32) (new int32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddUint64(addr *uint64, delta uint64) (new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
```

- 第二类CompareAndSwapXXX
```
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
```
- 第三类LoadXXX
```
func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)
func LoadUintptr(addr *uintptr) (val uintptr)
```
- 第四类StoreXXX
```
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)
func StoreUintptr(addr *uintptr, val uintptr)
```
- 第五类SwapXXX
```
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)
```
- 第六类Value，可用于任意类型的store load，实例就是这种
```
type Value
func (v *Value) Load() (x interface{})
func (v *Value) Store(x interface{})
```

Store
```
func (v *Value) Store(x interface{}) {
	if x == nil {
		panic("sync/atomic: store of nil value into Value")
	}
    // ifaceWords 其实就是定义了一下 interface 的结构，包含 data 和 type 两部分
    // 这里 vp 是原有值
    // xp 是传入的值
	vp := (*ifaceWords)(unsafe.Pointer(v))
	xp := (*ifaceWords)(unsafe.Pointer(&x))
    // for 循环不断尝试
	for {
        // 这里先用原子方法取一下老的类型值
		typ := LoadPointer(&vp.typ)
		if typ == nil {
            // 等于 nil 就说明这是第一次 store
            // 调用 runtime 的方法禁止抢占，避免操作完成一半就被抢占了
            // 同时可以避免 GC 的时候看到 unsafe.Pointer(^uintptr(0)) 这个中间状态的值
			runtime_procPin()
			if !CompareAndSwapPointer(&vp.typ, nil, unsafe.Pointer(^uintptr(0))) {
				runtime_procUnpin()
				continue
			}

			// 分别把值和类型保存下来
			StorePointer(&vp.data, xp.data)
			StorePointer(&vp.typ, xp.typ)
			runtime_procUnpin()
			return
		}

		if uintptr(typ) == ^uintptr(0) {
            // 如果判断发现这个类型是这个固定值，说明当前第一次赋值还没有完成，所以进入自旋等待
			continue
		}
		// 第一次赋值已经完成，判断新的赋值的类型和之前是否一致，如果不一致就直接 panic
		if typ != xp.typ {
			panic("sync/atomic: store of inconsistently typed value into Value")
		}
        // 保存值
		StorePointer(&vp.data, xp.data)
		return
	}
}
```

Load
```
func (v *Value) Load() (x interface{}) {
	vp := (*ifaceWords)(unsafe.Pointer(v))
    // 先拿到类型值
	typ := LoadPointer(&vp.typ)
    // 这个说明还没有第一次 store 或者是第一次 store 还没有完成
	if typ == nil || uintptr(typ) == ^uintptr(0) {
		// First store not yet completed.
		return nil
	}
    // 获取值
	data := LoadPointer(&vp.data)
    // 构造 x 类型
	xp := (*ifaceWords)(unsafe.Pointer(&x))
	xp.typ = typ
	xp.data = data
	return
}
```

## 实现一个无锁栈
```
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
```
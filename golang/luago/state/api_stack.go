package state

func (l *luaState) GetTop() int {
	return l.stack.top
}

func (l *luaState) AbsIndex(idx int) int {
	return l.stack.absIndex(idx)
}

func (l *luaState) CheckStack(n int) bool {
	l.stack.check(n)

	return true
}

func (l *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		l.stack.pop()
	}

	// other impl
	//l.SetTop(-n-1)
}

func (l *luaState) Copy(fromIdx, toIdx int) {
	val := l.stack.get(fromIdx)
	l.stack.set(toIdx, val)
}

func (l *luaState) PushValue(idx int) {
	val := l.stack.get(idx)
	l.stack.push(val)
}

// 将idx处元素替换为pop元素
func (l *luaState) Replace(idx int) {
	val := l.stack.pop()
	l.stack.set(idx, val)
}

func (l *luaState) Insert(idx int) {
	l.Rotate(idx, 1)
}

func (l *luaState) Remove(idx int) {
	l.Rotate(idx, -1)
	l.Pop(1)
}

func (l *luaState) Rotate(idx, n int) {
	t := l.stack.top - 1
	p := l.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	l.stack.reverse(p, m)
	l.stack.reverse(m+1, t)
	l.stack.reverse(p, t)
}

// 设置新的栈顶，不满足的被弹出
func (l *luaState) SetTop(idx int) {
	newTop := l.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	n := l.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			l.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			l.stack.push(nil)
		}
	}
}

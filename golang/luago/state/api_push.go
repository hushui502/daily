package state

import . "luago/api"

func (l *luaState) PushNil() {
	l.stack.push(nil)
}

func (l *luaState) PushBoolean(b bool) {
	l.stack.push(b)
}

func (l *luaState) PushInteger(n int64) {
	l.stack.push(n)
}

func (l *luaState) PushNumber(n float64) {
	l.stack.push(n)
}

func (l *luaState) PushString(s string) {
	l.stack.push(s)
}

func (l *luaState) PushGoFunction(f GoFunction) {
	l.stack.push(newGoClosure(f, 0))
}

func (l *luaState) PushGoClosure(f GoFunction, n int) {
	closure := newGoClosure(f, n)
	for i := n; i > 0; i-- {
		val := l.stack.pop()
		closure.upvals[n-1] = &upvalue{&val}
	}
	l.stack.push(closure)
}
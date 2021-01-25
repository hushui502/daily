package state

import (
	. "luago/api"
	"luago/number"
)

// 获取idx索引处的len，并压入栈顶
// TODO 考虑字符串之外的情况
func (self *luaState) Len(idx int) {
	val := self.stack.get(idx)
	if s, ok := val.(string); ok {
		self.stack.push(len(s))
	} else if result, ok := callMetamethod(val, val, "__len", self); ok {
		self.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		self.stack.push(int64(t.len()))
	} else {
		panic("length error!")
	}
}

func (self *luaState) Concat(n int) {
	if n == 0 {
		self.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if self.IsString(-1) && self.IsString(-2) {
				s2 := self.ToString(-1)
				s1 := self.ToString(-2)
				self.stack.pop()
				self.stack.pop()
				self.stack.push(s1 + s2)
				continue
			}
			b := self.stack.pop()
			a := self.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", self); ok {
				self.stack.push(result)
				continue
			}
			panic("concatenation error!")
		}
	}
	// n == 1 do nothing
}

func (self *luaState) Next(idx int) bool {
	val := self.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := self.stack.pop()
		if nextKey := t.nextKey(key); nextKey != nil {
			self.stack.push(nextKey)
			self.stack.push(t.get(nextKey))
			return true
		}
		return false
	}
	panic("table expected!")
}

// 返现栈顶错误，直接panic模拟抛
func (self *luaState) Error() int {
	err := self.stack.pop()
	panic(err)
}

func (self *luaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := self.stack
	status = LUA_ERRRUN

	// catch error
	defer func() {
		if err := recover(); err != nil {
			for self.stack != caller {
				self.popLuaStack()
			}
			self.stack.push(err)
		}
	}()

	self.Call(nArgs, nResults)
	// if no catch error
	status = LUA_OK

	return
}

func (self *luaState) StringToNumber(s string) bool {
	if n, ok := number.ParseInteger(s); ok {
		self.PushInteger(n)
		return true
	}
	if n, ok := number.ParseFloat(s); ok {
		self.PushNumber(n)
		return true
	}

	return false
}



























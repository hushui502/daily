package state

// 获取idx索引处的len，并压入栈顶
// TODO 考虑字符串之外的情况
func (self *luaState) Len(idx int) {
	val := self.stack.get(idx)
	if s, ok := val.(string); ok {
		self.stack.push(len(s))
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
			panic("concatenation error!")
		}
	}
	// n == 1 do nothing
}

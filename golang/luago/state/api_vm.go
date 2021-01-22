package state

func (self *luaState) PC() int {
	return self.stack.pc
}

func (self *luaState) AddPC(n int) {
	self.stack.pc += n
}

// 根据PC索引从函数原型的指令表中取出一条指令，并且将PC+1
func (self *luaState) Fetch() uint32 {
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc++

	return i
}

// 根据索引从常量表取出一个常量值，并推入栈顶
func (self *luaState) GetConst(idx int) {
	c := self.stack.closure.proto.Constants[idx]
	self.stack.push(c)
}

// Lua虚拟机指令操作数里携带的寄存器是从0开始，但是API索引是从1开始，因此要+1
func (self *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		self.GetConst(rk & 0xFF)
	} else { // register
		self.PushValue(rk + 1)
	}
}

package state

func (self *luaState) PC() int {
	return self.pc
}

func (selft *luaState) AddPC(n int) {
	selft.pc += n
}

// 根据PC索引从函数原型的指令表中取出一条指令，并且将PC+1
func (self *luaState) Fetch() uint32 {
	i := self.proto.Code[self.pc]
	self.pc++

	return i
}

// 根据索引从常量表取出一个常量值，并推入栈顶
func (self *luaState) GetConst(idx int) {
	c := self.proto.Constants[idx]
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

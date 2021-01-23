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

func (self *luaState) LoadProto(idx int) {
	stack := self.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)

	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {			// 当前函数的局部变量
			if stack.openuvs == nil {
				stack.openuvs = map[int]*upvalue{}
			}
			if openuv, found := stack.openuvs[uvIdx]; found { 			// 当前的局部变量在栈上则直接访问
				closure.upvals[i] = openuv
			} else {	// 需要把变量保存在某地openvs，局部变量退出作用域
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else {	// 外围函数的局部变量
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

func (self *luaState) RegisterCount() int {
	return int(self.stack.closure.proto.MaxStackSize)
}

func (self *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(self.stack.varargs)
	}

	self.stack.check(n)
	self.stack.pushN(self.stack.varargs, n)
}

// 处于开启状态的Upvalue引用了还在寄存器中的Lua值，将这些值复制出来。更新upvalue
func (self *luaState) CloseUpvalues(a int) {
	for i, openuv := range self.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(self.stack.openuvs, i)
		}
	}
}
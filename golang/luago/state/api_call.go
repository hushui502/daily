package state

import (
	. "luago/api"
	"luago/binchunk"
	"luago/vm"
)

func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	self.stack.push(c)

	if len(proto.Upvalues) > 0 {
		env := self.registry.get(LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}

	return 0
}

func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		if c.proto != nil {
			self.callLuaClosure(nArgs, nResults, c)
		} else {
			self.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not function!")
	}
}

func (self *luaState) callLuaClosure(nArags, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// 新的调用栈
	newStack := newLuaStack(nArags+LUA_MINSTACK, self)
	newStack.closure = c

	// 把函数和参数值一次性弹出
	funcAndArgs := self.stack.popN(nArags + 1)
	// 调用新帧的pushN方法将固定的参数值推入栈
	newStack.pushN(funcAndArgs[1:], nParams)
	// 参数推入完毕需要修改栈顶指针，让其指向最后一个寄存器
	newStack.top = nRegs
	if nArags > nParams && isVararg {
		// 暂存多出的参数值
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// 将被调用栈推入调用栈顶执行
	self.pushLuaStack(newStack)
	self.runLuaClosure()
	self.popLuaStack()

	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// 新建一个调用栈
	newStack := newLuaStack(nArgs+LUA_MINSTACK, self)
	newStack.closure = c

	// 将参数值从主调栈弹出
	args := self.stack.popN(nArgs)
	// 将参数值推入调用栈
	newStack.pushN(args, nArgs)
	// 至于主调用栈中的闭包直接扔掉即可
	self.stack.pop()

	// 参数传递完毕，再次将被调栈推入调用栈，让它成为当前帧
	self.pushLuaStack(newStack)
	// 执行Go函数
	r := c.goFunc(self)
	// 执行完毕将调用栈从主调栈中弹出，主调栈又成为当前栈
	self.popLuaStack()

	// 将返回值从被调栈推入主调栈
	if nResults != 0 {
		results := newStack.popN(r)
		self.stack.check(len(results))
		self.stack.pushN(results, nResults)
	}
}

func (self *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(self.Fetch())
		inst.Execute(self)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}

	}
}

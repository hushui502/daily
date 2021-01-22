package state

import (
	"fmt"
	"luago/binchunk"
	"luago/vm"
)

func (self *luaState) Load(chunk []byte, chunkName, mode string) int {
	proto := binchunk.Undump(chunk)
	c := newLuaClosure(proto)
	self.stack.push(c)

	return 0
}

func (self *luaState) Call(nArgs, nResults int) {
	val := self.stack.get(-(nArgs + 1))
	if c, ok := val.(*closure); ok {
		fmt.Printf("call %s<%d, %d>\n", c.proto.Source,
			c.proto.LineDefined, c.proto.LastLineDefined)
		self.callLuaClosure(nArgs, nResults, c)
	} else {
		panic("not function!")
	}
}

func (self *luaState) callLuaClosure(nArags, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// 新的调用栈
	newStack := newLuaStack(nRegs + 20)
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
		results := newStack.popN(newStack.top-nRegs)
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
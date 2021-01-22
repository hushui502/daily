package state

import . "luago/api"

// 创建一个table并且推入栈顶
func (self *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	self.stack.push(t)
}

func (self *luaState) NewTable() {
	self.CreateTable(0, 0)
}

func (self *luaState) GetTable(idx int) LuaType {
	// 获取一个table
	t := self.stack.get(idx)
	// 弹出栈顶=键元素
	k := self.stack.pop()

	return self.getTable(t, k)
}

func (self *luaState) getTable(t, k luaValue) LuaType {
	// t必须为table
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		// 将值推入栈顶
		self.stack.push(v)
		return typeOf(v)
	}
	panic("not a table!")
}

func (self *luaState) GetField(idx int, k string) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, k)
}

// 键是数字，主要用于table中的数组arr
func (self *luaState) GetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i)
}

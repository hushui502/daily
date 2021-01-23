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

// 取得一个tbl，然后弹出栈顶元素作为键来去tbl查找val并推入栈顶，返回该val的type
func (self *luaState) GetTable(idx int) LuaType {
	// 获取一个table
	t := self.stack.get(idx)
	// 弹出栈顶=键元素
	k := self.stack.pop()

	return self.getTable(t, k, false)
}

// raw == true 忽略元方法
func (self *luaState) getTable(t, k luaValue, raw bool) LuaType {
	// t必须为table
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if raw || v != nil || !tbl.hasMetafield("__index") {
			// 将值推入栈顶
			self.stack.push(v)
			return typeOf(v)
		}
		if !raw {
			if mf := getMetafield(t, "__index", self); mf != nil {
				switch x := mf.(type) {
				case *luaTable:
					return self.getTable(x, k, false)
				case *closure:
					self.stack.push(mf)
					self.stack.push(t)
					self.stack.push(k)
					self.Call(2, 1)
					v := self.stack.get(-1)
					return typeOf(v)
				}
			}
		}

	}
	panic("not a table!")
}

func (self *luaState) GetField(idx int, k string) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, k, false)
}

// 键是数字，主要用于table中的数组arr
func (self *luaState) GetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i, false)
}

func (self *luaState) RawGet(idx int) LuaType {
	t := self.stack.get(idx)
	k := self.stack.pop()

	return self.getTable(t, k, true)
}

func (self *luaState) RawGetI(idx int, i int64) LuaType {
	t := self.stack.get(idx)
	return self.getTable(t, i, true)
}

func (self *luaState) GetMetatable(idx int) bool {
	val := self.stack.get(idx)

	if mt := getMetatable(val, self); mt != nil {
		self.stack.push(mt)
		return true
	} else {
		return false
	}
}

// 获取全局环境中的某个字段并推入栈顶
func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)

	return self.getTable(t, name, false)
}



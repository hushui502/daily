package state

import . "luago/api"

func (self *luaState) SetTable(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()

	self.setTable(t, k, v, false)
}

func (self *luaState) setTable(t, k, v interface{}, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
		if !raw {
			if mf := getMetafield(t, "__newindex", self); mf != nil {
				switch x := mf.(type) {
				case *luaTable:
					self.setTable(x, k, v, false)
				case *closure:
					self.stack.push(mf)
					self.stack.push(t)
					self.stack.push(k)
					self.stack.push(v)
					self.Call(3, 0)
					return
				}
			}
		}
	}
	panic("not a table")
}

func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, k, v, false)
}

func (self *luaState) SetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, false)
}

func (self *luaState) RawSet(idx int) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	k := self.stack.pop()
	self.setTable(t, k, v, true)
}

func (self *luaState) RawSetI(idx int, i int64) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v, true)
}

func (self *luaState) SetMetatable(idx int) {
	val := self.stack.get(idx)
	mtVal := self.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, self)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, self)
	} else {
		panic("table expected!")
	}
}

func (self *luaState) SetGlobal(name string) {
	v := self.stack.pop()
	t := self.registry.get(LUA_REGISTRYINDEX)

	self.setTable(t, name, v, false)
}

func (self *luaState) Register(name string, f GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

package state

import . "luago/api"

type luaState struct {
	registry *luaTable
	stack    *luaStack
}

func New() *luaState {
	// 创建注册表
	registry := newLuaTable(0, 0)
	// 放入一个全局环境，其实全局环境就是一个普通的Lua表
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))

	return ls
}

// 从栈顶推入一个调用栈
func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}

// 从栈顶弹出一个调用栈
func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}

// 将全局环境推入栈顶
func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

// 获取全局环境中的某个字段并推入栈顶
func (self *luaState) GetGlobal(name string) LuaType {
	t := self.registry.get(LUA_RIDX_GLOBALS)

	return self.getTable(t, name)
}

func (self *luaState) SetGlobal(name string) {
	v := self.stack.pop()
	t := self.registry.get(LUA_RIDX_GLOBALS)

	self.setTable(t, name, v)
}

func (self *luaState) Register(name string, f GoFunction) {
	self.PushGoFunction(f)
	self.SetGlobal(name)
}

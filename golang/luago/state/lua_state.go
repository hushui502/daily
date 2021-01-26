package state

import . "luago/api"

type luaState struct {
	registry *luaTable
	stack    *luaStack
	/* coroutine */
	coStatus int
	coCaller *luaState
	coChan chan int
}

func New() *luaState {
	ls := &luaState{}

	// 创建注册表
	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, ls)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(0, 20))
	// 放入一个全局环境，其实全局环境就是一个普通的Lua表
	ls.registry = registry
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

// 是否是主线程 也就是main
func (self *luaState) isMainThread() bool {
	return self.registry.get(LUA_RIDX_MAINTHREAD) == self
}

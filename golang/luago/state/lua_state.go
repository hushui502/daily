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

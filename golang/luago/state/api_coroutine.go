package state

import . "luago/api"

// 创建一个线程并推入栈顶
func (self *luaState) NewThread() LuaState {
	t := &luaState{registry: self.registry}
	t.pushLuaStack(newLuaStack(LUA_MINSTACK, t))
	self.stack.push(t)

	return t
}

func (self *luaState) Resume(from LuaState, nArgs int) int {
	lsFrom := from.(*luaState)
	if lsFrom.coChan == nil {
		lsFrom.coChan = make(chan int)
	}

	// 两个线程通过coChan来相互协作，如果首次用到该字段则需要先初始化
	if self.coChan == nil {		// start coroutine
		self.coChan = make(chan int)
		self.coCaller = lsFrom
		go func() {
			self.coStatus = self.PCall(nArgs, -1, 0)
			lsFrom.coChan <- 1
		}()
	} else {		// resume coroutine
		self.coStatus = LUA_OK
		self.coChan <- 1
	}

	<-lsFrom.coChan			// wait coroutine to finish or yield

	return self.coStatus
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yield
func (self *luaState) Yield(nResults int) int {
	if self.coCaller == nil { // todo
		panic("attempt to yield from outside a coroutine")
	}
	self.coStatus = LUA_YIELD
	self.coCaller.coChan <- 1
	<-self.coChan
	return self.GetTop()
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isyieldable
func (self *luaState) IsYieldable() bool {
	if self.isMainThread() {
		return false
	}
	return self.coStatus != LUA_YIELD // todo
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_status
// lua-5.3.4/src/lapi.c#lua_status()
func (self *luaState) Status() int {
	return self.coStatus
}

// debug
func (self *luaState) GetStack() bool {
	return self.stack.prev != nil
}

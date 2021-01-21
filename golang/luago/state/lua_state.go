package state

import . "luago/binchunk"

type luaState struct {
	stack *luaStack
	proto *Prototype
	pc    int
}

func New(stackSize int, proto *Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(20),
		proto: proto,
		pc:    0,
	}
}

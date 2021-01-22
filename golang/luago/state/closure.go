package state

import "luago/binchunk"

type closure struct {
	// 函数原型
	proto *binchunk.Prototype
}

func newLuaClosure(proto *binchunk.Prototype) *closure {
	return &closure{proto: proto}
}

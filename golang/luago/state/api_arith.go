package main

import (
	. "luago/api"
	. "luago/number"
	"math"
)

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var (
	iadd  = func(a, b int64) int64 { return a + b }
	fadd  = func(a, b float64) float64 { return a + b }
	isub  = func(a, b int64) int64 { return a - b }
	fsub  = func(a, b float64) float64 { return a - b }
	imul  = func(a, b int64) int64 { return a * b }
	fmul  = func(a, b float64) float64 { return a * b }
	imod  = IMod
	fmod  = FMod
	pow   = math.Pow
	div   = func(a, b float64) float64 { return a / b }
	iidiv = IFloorDiv
	fidiv = FFloorDiv
	band  = func(a, b int64) int64 { return a & b }
	bor   = func(a, b int64) int64 { return a | b }
	bxor  = func(a, b int64) int64 { return a ^ b }
	shl   = ShiftLeft
	shr   = ShiftRight
	iunm  = func(a, _ int64) int64 { return -a }
	funm  = func(a, _ float64) float64 { return -a }
	bnot  = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	operator{iadd, fadd},
	operator{isub, fsub},
	operator{imul, fmul},
	operator{imod, fmod},
	operator{nil, pow},
	operator{nil, div},
	operator{iidiv, fidiv},
	operator{band, nil},
	operator{bor, nil},
	operator{bxor, nil},
	operator{shl, nil},
	operator{shr, nil},
	operator{iunm, funm},
	operator{bnot, nil},
}

func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = self.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = self.stack.pop()
	} else {
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		self.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

// 加减乘除取反期望是转换为整数进行运算，对于其他操作就转换成浮点数运算
func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil {
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		// add sub mul mod idiv unm
		if op.integerFunc != nil {
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}

	return nil
}

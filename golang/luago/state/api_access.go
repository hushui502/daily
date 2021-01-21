package main

import (
	"fmt"
	. "luago/api"
)

// 将lua类型转换为对应的字符串
func (l *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTABLE:
		return "table"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// 根据索引返回其类型
func (l *luaState) Type(idx int) LuaType {
	if l.stack.isValid(idx) {
		val := l.stack.get(idx)
		return typeOf(val)
	}

	return LUA_TNONE
}

func (l *luaState) IsNone(idx int) bool {
	return l.Type(idx) == LUA_TNONE
}

func (l *luaState) IsNil(idx int) bool {
	return l.Type(idx) == LUA_TNIL
}

func (l *luaState) IsNoneOrNil(idx int) bool {
	return l.Type(idx) <= LUA_TNIL
}

func (l *luaState) IsBoolean(idx int) bool {
	return l.Type(idx) == LUA_TBOOLEAN
}

func (l *luaState) IsInteger(idx int) bool {
	val := l.stack.get(idx)
	_, ok := val.(int64)

	return ok
}

func (l *luaState) IsNumber(idx int) bool {
	_, ok := l.ToNumberX(idx)

	return ok
}

func (l *luaState) IsString(idx int) bool {
	t := l.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (l *luaState) IsTable(idx int) bool {
	return l.Type(idx) == LUA_TTABLE
}

func (l *luaState) IsThread(idx int) bool {
	return l.Type(idx) == LUA_TTHREAD
}

func (l *luaState) IsFunction(idx int) bool {
	return l.Type(idx) == LUA_TFUNCTION
}

func (l *luaState) ToBoolean(idx int) bool {
	val := l.stack.get(idx)
	return convertToBoolean(val)
}

func (l *luaState) ToInteger(idx int) int64 {
	i, _ := l.ToIntegerX(idx)

	return i
}

func (l *luaState) ToIntegerX(idx int) (int64, bool) {
	val := l.stack.get(idx)

	return convertToInteger(val)
}

func (l *luaState) ToNumber(idx int) float64 {
	n, _ := l.ToNumberX(idx)

	return n
}

func (l *luaState) ToNumberX(idx int) (float64, bool) {
	val := l.stack.get(idx)

	return convertToFloat(val)
}

func (l *luaState) ToString(idx int) string {
	s, _ := l.ToStringX(idx)

	return s
}

func (l *luaState) ToStringX(idx int) (string, bool) {
	val := l.stack.get(idx)
	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		l.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

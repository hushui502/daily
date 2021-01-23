package state

import (
	"fmt"
	"luago/number"
	"math"
)

type luaTable struct {
	metatable *luaTable // 元表
	arr  []luaValue
	_map map[luaValue]luaValue
	keys map[luaValue]luaValue
	changed bool
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}

	return t
}

func setMetatable(val luaValue, mt *luaTable, ls *luaState) {
	 if t, ok := val.(*luaTable); ok {
	 	t.metatable = mt
	 	return
	 }
	 key := fmt.Sprintf("_MT%d", typeOf(val))
	 ls.registry.put(key, mt)
}

func getMetatable(val luaValue, ls *luaState) *luaTable {
	if t, ok := val.(*luaTable); ok {
		return t.metatable
	}
	key := fmt.Sprintf("_MT%d", typeOf(val))
	if mt := ls.registry.get(key); mt != nil {
		return mt.(*luaTable)
	}

	return nil
}

// 首先查询arr，然后再map
func (self *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(self.arr)) {
			return self.arr[idx-1]
		}
	}

	return self._map[key]
}

func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}

	return key
}

func (self *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(self.arr))
		// 这里的说明此时正在做数组元素的更新，因此会直接return，不需要更新map
		if idx <= arrLen {
			self.arr[idx-1] = val
			if idx == arrLen && val == nil {
				self._shrinkArray()
			}
			return
		}
		// 这里是刚好要扩容数组，这里会触发一次对map元素的遍历添加到arr
		if idx == arrLen+1 {
			delete(self._map, key)
			if val != nil {
				self.arr = append(self.arr, val)
				self._expandArray()
			}
		}
	}

	// 一般添加元素
	if val != nil {
		if self._map == nil {
			self._map = make(map[luaValue]luaValue)
		}
		self._map[key] = val
	} else { // 如果是nil就把该键从map删除节约空间
		delete(self._map, key)
	}
}

// 删除val = nil的数组slot
func (self *luaTable) _shrinkArray() {
	for i := len(self.arr) - 1; i >= 0; i-- {
		if self.arr[i] == nil {
			self.arr = self.arr[0:i]
		}
	}
}

// 数组扩容后，把原来存在于哈希表的某些元素移动到数组内
func (self *luaTable) _expandArray() {
	for idx := int64(len(self.arr)) + 1; true; idx++ {
		if val, found := self._map[idx]; found {
			delete(self._map, idx)
			self.arr = append(self.arr, val)
		} else {
			break
		}
	}
}

func (self *luaTable) len() int {
	return len(self.arr)
}

func (self *luaTable) hasMetafield(fieldName string) bool {
	return self.metatable != nil && self.metatable.get(fieldName) != nil
}

func (self *luaTable) nextKey(key luaValue) luaValue {
	if self.keys == nil || key == nil {
		self.initKeys()
		self.changed = false
	}

	return self.keys[key]
}

func (self *luaTable) initKeys() {
	self.keys = make(map[luaValue]luaValue)
	var key luaValue
	for i, v := range self.arr {
		if v != nil {
			self.keys[key]= int64(i + 1)
			key = int64(i + 1)
		}
	}
	for k, v := range self._map {
		if v != nil {
			self.keys[key] = k
			key = k
		}
	}
}

















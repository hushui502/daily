package stdlib

import "sort"
import "strings"
import . "luago/api"

const MAX_LEN = 1000000 // TODO

const (
	TAB_R  = 1               /* read */
	TAB_W  = 2               /* write */
	TAB_L  = 4               /* length */
	TAB_RW = TAB_R | TAB_W /* read/write */
)

var tabFuncs = map[string]GoFunction{
	"move":   tabMove,
	"insert": tabInsert,
	"remove": tabRemove,
	"sort":   tabSort,
	"concat": tabConcat,
	"pack":   tabPack,
	"unpack": tabUnpack,
}

func OpenTableLib(ls LuaState) int {
	ls.NewLib(tabFuncs)
	return 1
}


// table.move (a1, f, e, t [,a2])
func tabMove(ls LuaState) int {
	f := ls.CheckInteger(2)	// 需要移动表的起始位置
	e := ls.CheckInteger(3)	// 需要移动表的结束位置
	t := ls.CheckInteger(4)	// 需要插入表的起始位置

	tt := 1	// a1 表
	if !ls.IsNoneOrNil(5) {
		tt = 5		// 需要插入表a2
	}
	_checkTab(ls, 1, TAB_R)	// 保证src表可读
	_checkTab(ls, tt, TAB_W)		// 保证dst表可写
	if e >= f {		// 正常的逻辑，重写或者不需要移动
		var n, i int64
		ls.ArgCheck(f > 0 || e < LUAI_MAXSTACK+f, 3,
			"too many elements to move")
		n = e - f + 1	// number of elements to move
		ls.ArgCheck(t <= LUAI_MAXSTACK-n+1, 4,
			"destination wrap around")
		if t > e || t <= f || (tt != 1 && !ls.Compare(1, tt, LUA_OPEQ)) {
			for i = 0; i < n; i++ {
				ls.GetI(1, f+i)
				ls.SetI(tt, t+i)
			}
		} else {
			for i := n-1; i >= 0; i-- {
				ls.GetI(1, f+i)
				ls.SetI(tt, t+i)
			}
		}
	}
	
	ls.PushValue(tt)
	return 1
}

func tabInsert(ls LuaState) int {
	e := _auxGetN(ls, 1, TAB_RW) + 1 /* first empty element */
	var pos int64                    /* where to insert new element */
	switch ls.GetTop() {
	case 2: /* called with only 2 arguments */
		pos = e /* insert new element at the end */
	case 3:
		pos = ls.CheckInteger(2) /* 2nd argument is the position */
		ls.ArgCheck(1 <= pos && pos <= e, 2, "position out of bounds")
		for i := e; i > pos; i-- { /* move up elements */
			ls.GetI(1, i-1)
			ls.SetI(1, i) /* t[i] = t[i - 1] */
		}
	default:
		return ls.Error2("wrong number of arguments to 'insert'")
	}
	ls.SetI(1, pos) /* t[pos] = v */
	return 0
}

func tabRemove(ls LuaState) int {
	size := _auxGetN(ls, 1, TAB_RW)
	pos := ls.OptInteger(2, size)
	if pos != size { /* validate 'pos' if given */
		ls.ArgCheck(1 <= pos && pos <= size+1, 1, "position out of bounds")
	}
	ls.GetI(1, pos) /* result = t[pos] */
	for ; pos < size; pos++ {
		ls.GetI(1, pos+1)
		ls.SetI(1, pos) /* t[pos] = t[pos + 1] */
	}
	ls.PushNil()
	ls.SetI(1, pos) /* t[pos] = nil */
	return 1
}

func tabConcat(ls LuaState) int {
	tabLen := _auxGetN(ls, 1, TAB_R)
	sep := ls.OptString(2, "")
	i := ls.OptInteger(3, 1)
	j := ls.OptInteger(4, tabLen)

	if i > j {
		ls.PushString("")
		return 1
	}

	buf := make([]string, j-i+1)
	for k := i; k > 0 && k <= j; k++ {
		ls.GetI(1, k)
		if !ls.IsString(-1) {
			ls.Error2("invalid value (%s) at index %d in table for 'concat'",
				ls.TypeName2(-1), i)
		}
		buf[k-i] = ls.ToString(-1)
		ls.Pop(1)
	}
	ls.PushString(strings.Join(buf, sep))

	return 1
}

func _auxGetN(ls LuaState, n, w int) int64 {
	_checkTab(ls, n, w|TAB_L)
	return ls.Len2(n)
}

//
func _checkTab(ls LuaState, arg, what int) {
	if ls.Type(arg) != LUA_TTABLE { /* is it not a table? */
		n := 1                     /* number of elements to pop */
		if ls.GetMetatable(arg) && /* must have metatable */
			(what&TAB_R != 0 || _checkField(ls, "__index", &n)) &&
			(what&TAB_W != 0 || _checkField(ls, "__newindex", &n)) &&
			(what&TAB_L != 0 || _checkField(ls, "__len", &n)) {
			ls.Pop(n) /* pop metatable and tested metamethods */
		} else {
			ls.CheckType(arg, LUA_TTABLE) /* force an error */
		}
	}
}

func _checkField(ls LuaState, key string, n *int) bool {
	ls.PushString(key)
	*n++
	return ls.RawGet(-*n) != LUA_TNIL
}

func tabPack(ls LuaState) int {
	n := int64(ls.GetTop())   /* number of elements to pack */
	ls.CreateTable(int(n), 1) /* create result table */
	ls.Insert(1)              /* put it at index 1 */
	for i := n; i >= 1; i-- { /* assign elements */
		ls.SetI(1, i)
	}
	ls.PushInteger(n)
	ls.SetField(1, "n") /* t.n = number of elements */
	return 1            /* return table */
}

func tabUnpack(ls LuaState) int {
	i := ls.OptInteger(2, 1)
	e := ls.OptInteger(3, ls.Len2(1))
	if i > e { /* empty range */
		return 0
	}

	n := int(e - i + 1)
	if n <= 0 || n >= MAX_LEN || !ls.CheckStack(n) {
		return ls.Error2("too many results to unpack")
	}

	for ; i < e; i++ { /* push arg[i..e - 1] (to avoid overflows) */
		ls.GetI(1, i)
	}
	ls.GetI(1, e) /* push last element */
	return n
}

func tabSort(ls LuaState) int {
	w := wrapper{ls}
	ls.ArgCheck(w.Len() < MAX_LEN, 1, "arrray too big")
	sort.Sort(w)

	return 0
}

type wrapper struct {
	ls LuaState
}

// 这里是为了实现sort的接口，否则自定义类型无法使用sort.Sort()
func (self wrapper) Len() int {
	return int(self.ls.Len2(1))
}

func (self wrapper) Less(i, j int) bool {
	ls := self.ls
	if ls.IsFunction(2) {
		ls.PushValue(2)
		ls.GetI(1, int64(i+1))
		ls.GetI(1, int64(j+1))
		ls.Call(2, 1)
		b := ls.ToBoolean(-1)
		ls.Pop(1)
		return b
	} else {
		ls.GetI(1, int64(i+1))
		ls.GetI(1, int64(j+1))
		b := ls.Compare(-2, -1, LUA_OPLT)
		ls.Pop(2)
		return b
	}
}

func (self wrapper) Swap(i, j int) {
	ls := self.ls

	ls.GetI(1, int64(i+1))
	ls.GetI(1, int64(j+1))
	ls.SetI(1, int64(i+1))
	ls.SetI(1, int64(j+1))
}
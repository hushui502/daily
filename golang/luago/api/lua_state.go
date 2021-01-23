package api

type LuaType = int                 // lua的类型
type ArithOp = int                 // 运算类型
type CompareOp = int               // 比较类型
type GoFunction func(LuaState) int // Go函数类型

// 注册表伪索引减去该索引对应的就是Upvalue伪索引
func LuaUpvalueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}

type LuaState interface {
	/* basic stack manipulation */
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)
	/* access functions (stack -> Go) */
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	IsTable(idx int) bool
	IsThread(idx int) bool
	IsFunction(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	/* push functions (Go -> stack) */
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	/* arith */
	Arith(op ArithOp)                          // 算数运算和位运算
	Compare(idx1, idx2 int, op CompareOp) bool // 比较运算
	Len(idx int)                               // 取长度运算
	Concat(n int)                              // 字符串拼接运算
	/* get functios (Lua-stack) */
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	/* set functions (stack->lua) */
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, n int64)
	/* function call */
	// 加载二进制chunk，把主函数原型转化为闭包并推入栈顶
	// mode 加载模式 可选b t bt 代表chunk需要为字节还是文本还是都可以
	// 返回状态码 0 成功 非0加载失败
	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)
	/* Go function call */
	// Go函数进入Lua栈，变成Go闭包才能为Lua使用，此方法就是Go函数转换撑Go闭包推入栈顶
	PushGoFunction(f GoFunction)
	// 判断索引处的值能否转成Go函数，不改变栈的状态
	IsGoFunction(idx int) bool
	// 把栈索引处的值转换成Go函数并返回，如果不可，返沪nil
	ToGoFunction(idx int) GoFunction
	/* 操作全局变量 */
	// 将全局环境表推入栈顶
	PushGlobalTable()
	GetGlobal(name string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)
	/* go闭包支持 */
	// 将Go函数变成闭包推入栈顶，需要先从栈顶弹出n个LuaValue，这些值会成为Go闭包的Upvalue
	PushGoClosure(f GoFunction, n int)
}

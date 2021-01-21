package api

type LuaVM interface {
	LuaState
	PC() int          // 返回当前的pc
	AddPC(n int)      // 修改PC，用于实现跳转指令
	Fetch() uint32    // 取出当前指令，将PC指向下一条指令，循环常用
	GetConst(idx int) // 从常量表取出常量并推入栈顶
	GetRK(rk int)     // 从常量表或者栈里取出常量并推入栈顶
}

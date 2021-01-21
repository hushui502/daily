package vm

import . "luago/api"

// 加载指令，将nil 布尔值 常量表里的常量加载到寄存器
func loadNil(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	vm.PushNil()
	for i := a; a <= a+b; i++ {
		vm.Copy(-1, i)
	}
	vm.Pop(1)
}

func loadBool(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	vm.PushBoolean(b != 0)
	vm.Replace(a)
	if c != 0 {
		vm.AddPC(1)
	}
}

// local a, b, c = nil, 1, "foo"
// lua的编译器会将字面量（主要是数字和字符串）收集起来，放进常量表
func loadK(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	// 根据索引从常量表中取出常量推入栈顶
	vm.GetConst(bx)
	// 将栈顶的元素移动到指定索引处
	vm.Replace(a)
}

func loadKx(i Instruction, vm LuaVM) {
	a, _ := i.ABx()
	a += 1
	ax := Instruction(vm.Fetch()).Ax()

	vm.GetConst(ax)
	vm.Replace(a)
}

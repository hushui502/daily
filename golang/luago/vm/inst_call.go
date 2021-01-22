package vm

import . "luago/api"

func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1

	vm.LoadProto(bx)
	vm.Replace(a)
}

func call(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1

	// 将被调函数和参数推入栈顶
	nArgs := _pushFuncAndArgs(a, b, vm)
	vm.Call(nArgs, c-1)
	// 将返回值移动到适合的vm
	_popResults(a, c, vm)
}

func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	if b >= 1 {
		vm.CheckStack(b)
		for i := a; i < b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		// 这种情况类似 f(1, 2, g())
		// 这里的g()虽然没有val和res但是仍然要让f()接收
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _popResults(a, c int, vm LuaVM) {
	if c == 1 { // no results
	} else if c > 1 { // c-1 results
		for i := a + c - 2; i >= a; i-- {
			// 将返回值依次推入
			vm.Replace(i)
		}
	} else {
		// 暂时将返回值留在栈顶
		vm.CheckStack(1)
		// 标记这些返回值原本是要移动到哪个寄存器中
		vm.PushInteger(int64(a))
	}
}

func _fixStack(a int, vm LuaVM) {
	x := int(vm.ToInteger(-1))
	vm.Pop(1)

	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	// 后半部分参数值已经在栈顶了，所以只需要把函数和前半部分参数值推入栈顶然后旋转
	vm.Rotate(vm.RegisterCount()+1, x-a)
}

// 将返回值推入栈顶
// b == 1 不需要返回任何值
// b > 1 需要返回b-1个值
// b == 0 一部分返回值已经在栈顶，将另一部分也推入栈顶
func _return(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 { // no return values
	} else if b > 1 {
		vm.CheckStack(b)
		for i := a; i <= a+b-2; i++ {
			vm.PushValue(i)
		}
	} else {
		_fixStack(a, vm)
	}
}

func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b != 1 {
		// 将全部的vararg参数推入栈顶
		vm.LoadVararg(b - 1)
		// 将这些参数移动到适合的寄存器
		_popResults(a, b, vm)
	}
}

// 尾递归优化
func tailCall(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	c := 0

	// 压入参数值
	nArgs := _pushFuncAndArgs(a, b, vm)
	// 执行
	vm.Call(nArgs, c-1)
	// 弹出，避免继续递归调用新建调用栈
	_popResults(a, c, vm)
}

func self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

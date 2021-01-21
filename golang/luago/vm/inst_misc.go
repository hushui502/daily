package vm

import . "luago/api"

// local a, b, c, d;  d = b
// MOVE 3 1
// 局部变量存到寄存器中，但是MOVE等操作的操作数，也就是ABC()中的A是8位，所以局部变量不能超过255
func move(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	// 因为lua api的栈索引从1开始
	a += 1
	b += 1
	vm.Copy(b, a)
}

// 跳转指令，就是改变PC
func jmp(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	vm.AddPC(sBx)
	if a != 0 {
		panic("--todo!")
	}
}

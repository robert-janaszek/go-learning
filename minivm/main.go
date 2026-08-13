package main

import (
	"fmt"
	"minivm/vm"
)

func main() {
	mem := vm.NewMemory(1024)
	machine := vm.NewVM(mem)

	err := machine.Call(100, 1)

	if err != nil {
		fmt.Println(err)
		return
	}
	machine.Push(5)

	err = machine.Call(121, 2)

	if err != nil {
		fmt.Println(err)
		return
	}

	machine.Push(321)
	machine.Push(4)

	fmt.Println(mem.Dump(992, 1024))

	machine.Ret()
	machine.Ret()

	fmt.Println(machine.StackDepth())
}

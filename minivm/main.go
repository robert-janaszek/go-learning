package main

import (
	"fmt"
	"minivm/vm"
)

func main() {
	mem := vm.NewMemory(1024)

	err := mem.Store(4, 1023)

	if err != nil {
		fmt.Println(err)
		return
	}

	val, err := mem.Load(4)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%x\n", val)
	fmt.Printf("%d\n", val)

	fmt.Println(mem.Dump(0, 14))

	machine := vm.NewVM(mem)

	err = machine.Push(2032)

	if err != nil {
		fmt.Println(err)
		return
	}

	val, err = machine.Peek()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(machine.StackDepth())
	fmt.Println(val)

	val, err = machine.Pop()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(val)
	fmt.Println(machine.StackDepth())
}

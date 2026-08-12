package main

import (
	"fmt"
	"minivm/vm"
)

func main() {
	mem := vm.NewMemory(1024)

	err := mem.Store(8, 123)
	if err != nil {
		fmt.Println(err)
		return
	}
	machine := vm.NewVM(mem)

	err = machine.Push(8)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = machine.LoadIndirect()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(mem.Dump(vm.Addr(mem.Size()-8), vm.Addr(mem.Size())))

	err = machine.Push(64)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = machine.Push(34)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = machine.StoreIndirect()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(mem.Dump(64, 68))
	fmt.Println(mem.Dump(vm.Addr(mem.Size()-8), vm.Addr(mem.Size())))
}

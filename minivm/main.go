package main

import (
	"fmt"
	"minivm/vm"
)

func main() {
	mem := vm.NewMemory(1024)
	machine := vm.NewVM(mem)

	addr, err := machine.Alloc(5)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("first addr: %d\n", addr)

	addr, err = machine.Alloc(8)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("second addr: %d\n", addr)
	fmt.Println(mem.Dump(64, 92))

	err = machine.Free(addr)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("after free")
	fmt.Println(mem.Dump(64, 92))
}

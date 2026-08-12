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
}

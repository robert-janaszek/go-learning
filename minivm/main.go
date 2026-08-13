package main

import (
	"fmt"
	"minivm/vm"
)

func main() {
	runProgramA()
	runProgramB()
	runProgramC()
}

func runProgramA() {
	fmt.Println("=== Program A: stack arithmetic (40 + 2) ===")
	mem := vm.NewMemory(1024)
	machine := vm.NewVM(mem)

	code := []vm.Instr{
		{Op: vm.OpPush, Arg: 40},
		{Op: vm.OpPush, Arg: 2},
		{Op: vm.OpAdd},
		{Op: vm.OpPrint},
		{Op: vm.OpHalt},
	}

	if err := machine.Execute(code); err != nil {
		fmt.Println("Execute error:", err)
		return
	}
	dumpState(machine, mem)
}

func runProgramB() {
	fmt.Println("=== Program B: heap Alloc / Store / Load / Free ===")
	mem := vm.NewMemory(1024)
	machine := vm.NewVM(mem)

	// Alloc 4, store 7, load it back, print, free.
	code := []vm.Instr{
		{Op: vm.OpPush, Arg: 4},
		{Op: vm.OpAlloc},
		{Op: vm.OpDup},
		{Op: vm.OpPush, Arg: 7},
		{Op: vm.OpStore},
		{Op: vm.OpDup},
		{Op: vm.OpLoad},
		{Op: vm.OpPrint},
		{Op: vm.OpFree},
		{Op: vm.OpHalt},
	}

	if err := machine.Execute(code); err != nil {
		fmt.Println("Execute error:", err)
		return
	}
	dumpState(machine, mem)
}

func runProgramC() {
	fmt.Println("=== Program C: Call / Ret (+1 via scratch at 0x10) ===")
	mem := vm.NewMemory(1024)
	machine := vm.NewVM(mem)

	// Ret discards locals (sp = fp), so the callee writes through reserved addr 0x10.
	// Main: store 41 at 0x10, Call incr, load 0x10, Print.
	// incr: leave addr under the new value, then Store (pop val, pop addr).
	code := []vm.Instr{
		{Op: vm.OpPush, Arg: 0x10}, // 0
		{Op: vm.OpPush, Arg: 41},   // 1
		{Op: vm.OpStore},           // 2
		{Op: vm.OpCall, Arg: 8},    // 3 → incr
		{Op: vm.OpPush, Arg: 0x10}, // 4
		{Op: vm.OpLoad},            // 5
		{Op: vm.OpPrint},           // 6 → 42
		{Op: vm.OpHalt},            // 7
		// incr:
		{Op: vm.OpPush, Arg: 0x10}, // 8  addr (stays under result)
		{Op: vm.OpPush, Arg: 0x10}, // 9
		{Op: vm.OpLoad},            // 10
		{Op: vm.OpPush, Arg: 1},    // 11
		{Op: vm.OpAdd},             // 12  stack: [0x10, 42]
		{Op: vm.OpStore},           // 13
		{Op: vm.OpRet},             // 14
	}

	if err := machine.Execute(code); err != nil {
		fmt.Println("Execute error:", err)
		return
	}
	dumpState(machine, mem)
}

func dumpState(machine *vm.VM, mem *vm.Memory) {
	fmt.Printf("SP=%d FP=%d heapBrk=%d depth=%d\n",
		machine.SP(), machine.FP(), machine.HeapBrk(), machine.StackDepth())
	fmt.Println(mem.Dump(0x00, 0x50))
	fmt.Println(mem.Dump(AddrHigh(mem), vm.Addr(mem.Size())))
	fmt.Println()
}

func AddrHigh(mem *vm.Memory) vm.Addr {
	size := mem.Size()
	if size < 64 {
		return 0
	}
	return vm.Addr(size - 64)
}

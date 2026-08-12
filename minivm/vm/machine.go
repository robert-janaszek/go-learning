package vm

import (
	"encoding/binary"
	"errors"
)

type VM struct {
	mem *Memory
	sp  Addr
	Heap
}

func NewVM(mem *Memory) *VM {
	return &VM{
		mem: mem,
		sp:  Addr(mem.Size()),
		Heap: Heap{
			heapStart: 0x40,
			heapBrk:   0x40,
		},
	}
}

func (v *VM) Push(value uint32) error {
	addr := v.sp - WordSize

	if addr <= 0 {
		return OutOfMemoryErr
	}

	err := v.mem.validateAddress(addr)

	if err != nil {
		return err
	}

	v.sp = addr

	binary.LittleEndian.PutUint32(v.mem.data[addr:addr+WordSize], value)

	return nil
}

func (v *VM) Pop() (uint32, error) {
	addr := v.sp

	if addr+WordSize > Addr(v.mem.Size()) {
		return 0, errors.New("cannot pop from empty stack")
	}

	err := v.mem.validateAddress(addr)

	if err != nil {
		return 0, err
	}

	v.sp += WordSize

	return binary.LittleEndian.Uint32(v.mem.data[addr : addr+WordSize]), nil
}

func (v *VM) Peek() (uint32, error) {
	addr := v.sp
	if addr+WordSize > Addr(v.mem.Size()) {
		return 0, errors.New("cannot peek from empty stack")
	}

	err := v.mem.validateAddress(addr)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(v.mem.data[addr : addr+WordSize]), nil
}

func (v *VM) StackDepth() int {
	return (v.mem.Size() - int(v.sp)) / WordSize
}

func (v *VM) LoadIndirect() error {
	addr, err := v.Pop()

	if err != nil {
		return err
	}

	val, err := v.mem.Load(Addr(addr))

	if err != nil {
		return err
	}

	err = v.Push(val)

	if err != nil {
		return err
	}

	return nil
}

func (v *VM) StoreIndirect() error {
	val, err := v.Pop()

	if err != nil {
		return err
	}

	addr, err := v.Pop()

	if err != nil {
		return err
	}

	err = v.mem.Store(Addr(addr), val)

	if err != nil {
		return err
	}

	return nil
}

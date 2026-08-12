package vm

import (
	"encoding/binary"
	"errors"
)

type VM struct {
	mem *Memory
	sp  Addr
}

func NewVM(mem *Memory) *VM {
	return &VM{
		mem: mem,
		sp:  Addr(mem.Size()),
	}
}

func (v *VM) Push(value uint32) error {
	addr := v.sp - WordSize

	if addr <= 0 {
		return errors.New("out of memory")
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

package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Addr uint32

const WordSize = 4

type Memory struct {
	data []byte
}

func NewMemory(size int) *Memory {
	if size%WordSize != 0 {
		panic(fmt.Sprintf("memory size should be divisible by word size: %d", WordSize))
	}

	memoryData := make([]byte, size)

	return &Memory{
		data: memoryData,
	}
}

func (m *Memory) Size() int {
	return len(m.data)
}

func (m *Memory) validateAddress(addr Addr) error {
	if addr == 0 {
		return errors.New("cannot access address 0")
	}

	if uint32(addr)+WordSize > uint32(m.Size()) {
		return errors.New("address out of bounds")
	}

	if uint32(addr)%WordSize != 0 {
		return errors.New("incorrect address, should be divisible by word size")
	}

	return nil
}

func (m *Memory) Load(addr Addr) (uint32, error) {
	err := m.validateAddress(addr)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(m.data[addr : addr+WordSize]), nil
}

func (m *Memory) Store(addr Addr, value uint32) error {
	err := m.validateAddress(addr)

	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(m.data[addr:addr+WordSize], value)

	return nil
}

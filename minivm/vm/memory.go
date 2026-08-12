package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
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

func (m *Memory) Dump(start, end Addr) string {
	refinedEnd := min(uint32(end), uint32(m.Size()))
	endRem := refinedEnd % WordSize
	endClamped := refinedEnd - endRem

	startRem := uint32(start) % WordSize
	refinedStart := uint32(start) - startRem

	if refinedStart >= endClamped {
		return ""
	}

	builder := strings.Builder{}

	for i := refinedStart; i < endClamped; i = i + WordSize {
		_, err := fmt.Fprintf(&builder, "%04x: ", i)
		if err != nil {
			return ""
		}

		_, err = fmt.Fprintf(&builder, "%02x %02x %02x %02x ", m.data[i], m.data[i+1], m.data[i+2], m.data[i+3])
		if err != nil {
			return ""
		}

		_, err = fmt.Fprintf(&builder, "%d\n", binary.LittleEndian.Uint32(m.data[i:i+WordSize]))
		if err != nil {
			return ""
		}
	}

	return builder.String()
}

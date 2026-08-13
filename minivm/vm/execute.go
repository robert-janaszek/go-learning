package vm

import (
	"errors"
	"fmt"
)

func withIP(ip int, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("ip %d: %w", ip, err)
}

func (v *VM) Execute(code []Instr) error {
	v.ip = 0

	for range 100_000 {
		if v.ip < 0 {
			return withIP(v.ip, errors.New("negative ip"))
		}

		if v.ip >= len(code) {
			return withIP(v.ip, fmt.Errorf("ip exceeds code len=%d", len(code)))
		}

		instr := code[v.ip]
		cur := v.ip
		v.ip++

		switch instr.Op {
		case OpHalt:
			return nil
		case OpPush:
			if err := v.Push(instr.Arg); err != nil {
				return withIP(cur, err)
			}
		case OpPop:
			if _, err := v.Pop(); err != nil {
				return withIP(cur, err)
			}
		case OpAdd:
			a, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			b, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			if err = v.Push(b + a); err != nil {
				return withIP(cur, err)
			}
		case OpLoad:
			addr, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			val, err := v.mem.Load(Addr(addr))
			if err != nil {
				return withIP(cur, err)
			}
			if err := v.Push(val); err != nil {
				return withIP(cur, err)
			}
		case OpStore:
			val, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			addr, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			if err := v.mem.Store(Addr(addr), val); err != nil {
				return withIP(cur, err)
			}
		case OpAlloc:
			nbytes, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			ptr, err := v.Alloc(nbytes)
			if err != nil {
				return withIP(cur, err)
			}
			if err := v.Push(uint32(ptr)); err != nil {
				return withIP(cur, err)
			}
		case OpFree:
			ptr, err := v.Pop()
			if err != nil {
				return withIP(cur, err)
			}
			if err = v.Free(Addr(ptr)); err != nil {
				return withIP(cur, err)
			}
		case OpCall:
			if err := v.Call(uint32(v.ip), 0); err != nil {
				return withIP(cur, err)
			}
			v.ip = int(instr.Arg)
		case OpRet:
			addr, err := v.Ret()
			if err != nil {
				return withIP(cur, err)
			}
			v.ip = int(addr)
		case OpDup:
			val, err := v.Peek()
			if err != nil {
				return withIP(cur, err)
			}
			if err := v.Push(val); err != nil {
				return withIP(cur, err)
			}
		case OpPrint:
			val, err := v.Pop()

			if err != nil {
				return withIP(cur, err)
			}

			fmt.Println(val)
		default:
			return withIP(cur, errors.New("unknown opcode"))
		}
	}

	return withIP(v.ip, errors.New("infinite loop detected"))
}

package vm

import (
	"errors"
	"strings"
	"testing"
)

func mustExecute(t *testing.T, v *VM, code []Instr) {
	t.Helper()
	if err := v.Execute(code); err != nil {
		t.Fatalf("Execute failed: %v (ip=%d sp=%d)", err, v.ip, v.sp)
	}
}

func TestExecute_Add40Plus2(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	mustExecute(t, v, []Instr{
		{Op: OpPush, Arg: 40},
		{Op: OpPush, Arg: 2},
		{Op: OpAdd},
		{Op: OpHalt},
	})

	got, err := v.Peek()
	if err != nil {
		t.Fatalf("Peek failed: %v", err)
	}
	if got != 42 {
		t.Fatalf("wanted top of stack 42, got %d (sp=%d)", got, v.sp)
	}
}

func TestExecute_AllocStoreLoadFree(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	const scratch Addr = 0x10

	const result Addr = 0x14

	mustExecute(t, v, []Instr{
		{Op: OpPush, Arg: uint32(scratch)},
		{Op: OpPush, Arg: 4},
		{Op: OpAlloc},
		{Op: OpStore}, // mem[0x10] = ptr

		{Op: OpPush, Arg: uint32(scratch)},
		{Op: OpLoad},
		{Op: OpDup},
		{Op: OpPush, Arg: 7},
		{Op: OpStore}, // mem[ptr] = 7

		{Op: OpPush, Arg: uint32(result)},
		{Op: OpPush, Arg: uint32(scratch)},
		{Op: OpLoad},
		{Op: OpLoad},
		{Op: OpStore}, // mem[0x14] = 7 (payload is clobbered by Free's next ptr)

		{Op: OpPush, Arg: uint32(scratch)},
		{Op: OpLoad},
		{Op: OpFree},
		{Op: OpHalt},
	})

	ptr, err := mem.Load(scratch)
	if err != nil {
		t.Fatalf("Load scratch failed: %v", err)
	}

	got, err := mem.Load(result)
	if err != nil {
		t.Fatalf("Load result failed: %v", err)
	}
	if got != 7 {
		t.Fatalf("wanted stored value 7, got %d", got)
	}

	brk := v.heapBrk
	reused, err := v.Alloc(4)
	if err != nil {
		t.Fatalf("Alloc after bytecode Free failed: %v", err)
	}
	if reused != Addr(ptr) {
		t.Fatalf("wanted Free in bytecode to reuse ptr %d, got %d", ptr, reused)
	}
	if v.heapBrk != brk {
		t.Fatalf("wanted heapBrk unchanged at %d, got %d", brk, v.heapBrk)
	}
}

func TestExecute_CallRetNested(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	// 0: Call 2
	// 1: Halt
	// 2: Call 4
	// 3: Ret
	// 4: Ret
	code := []Instr{
		{Op: OpCall, Arg: 2},
		{Op: OpHalt},
		{Op: OpCall, Arg: 4},
		{Op: OpRet},
		{Op: OpRet},
	}

	mustExecute(t, v, code)

	if got := v.StackDepth(); got != 0 {
		t.Fatalf("wanted empty stack after nested Call/Ret, got depth %d", got)
	}
	if v.fp != Addr(mem.Size()) {
		t.Fatalf("wanted fp restored to %d, got %d", mem.Size(), v.fp)
	}
}

func TestExecute_StackOverflowOnPush(t *testing.T) {
	// size=8 → one valid stack word (addr 4); second Push must fail.
	mem := NewMemory(8)
	v := NewVM(mem)

	code := []Instr{
		{Op: OpPush, Arg: 1},
		{Op: OpPush, Arg: 2},
		{Op: OpHalt},
	}

	err := v.Execute(code)
	if err == nil {
		t.Fatal("wanted stack overflow error, got nil")
	}
	if !errors.Is(err, ErrOutOfMemory) && !strings.Contains(err.Error(), "ip ") {
		t.Fatalf("wanted wrapped OOM/overflow with ip, got %v", err)
	}
	if !strings.Contains(err.Error(), "ip 1") {
		t.Fatalf("wanted error at ip 1 (second Push), got %v", err)
	}
}

func TestExecute_AllocOOM(t *testing.T) {
	mem := NewMemory(128)
	v := NewVM(mem)

	// Keep allocating 32-byte payloads until heap collides with stack.
	code := make([]Instr, 0, 64)
	for range 40 {
		code = append(code,
			Instr{Op: OpPush, Arg: 32},
			Instr{Op: OpAlloc},
			Instr{Op: OpPop}, // discard ptr so stack does not grow
		)
	}
	code = append(code, Instr{Op: OpHalt})

	err := v.Execute(code)
	if err == nil {
		t.Fatal("wanted OOM from Alloc, got nil")
	}
	if !errors.Is(err, ErrOutOfMemory) {
		t.Fatalf("wanted ErrOutOfMemory, got %v", err)
	}
	if !strings.Contains(err.Error(), "ip ") {
		t.Fatalf("wanted ip in error, got %v", err)
	}
}

func TestExecute_NullLoad_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	code := []Instr{
		{Op: OpPush, Arg: 0},
		{Op: OpLoad},
		{Op: OpHalt},
	}

	err := v.Execute(code)
	if err == nil {
		t.Fatal("wanted error for Load of null address, got nil")
	}
	if !strings.Contains(err.Error(), "ip 1") {
		t.Fatalf("wanted error at ip 1 (OpLoad), got %v", err)
	}
}

func TestExecute_UnknownOpcode_ReturnsError(t *testing.T) {
	mem := NewMemory(64)
	v := NewVM(mem)

	code := []Instr{
		{Op: Op(255)},
	}

	err := v.Execute(code)
	if err == nil {
		t.Fatal("wanted error for unknown opcode, got nil")
	}
	if !strings.Contains(err.Error(), "unknown opcode") {
		t.Fatalf("wanted unknown opcode error, got %v", err)
	}
}

func TestExecute_Dup(t *testing.T) {
	mem := NewMemory(64)
	v := NewVM(mem)

	mustExecute(t, v, []Instr{
		{Op: OpPush, Arg: 9},
		{Op: OpDup},
		{Op: OpHalt},
	})

	if got := v.StackDepth(); got != 2 {
		t.Fatalf("wanted StackDepth=2 after Dup, got %d", got)
	}
	top, err := v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	under, err := v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if top != 9 || under != 9 {
		t.Fatalf("wanted Dup of 9, got %d and %d", top, under)
	}
}

func TestExecute_PushLoop_StackOverflow(t *testing.T) {
	mem := NewMemory(64)
	v := NewVM(mem)

	code := make([]Instr, 0, 33)
	for range 32 {
		code = append(code, Instr{Op: OpPush, Arg: 1})
	}
	code = append(code, Instr{Op: OpHalt})

	err := v.Execute(code)
	if err == nil {
		t.Fatal("wanted stack overflow from Push loop, got nil")
	}
	if !errors.Is(err, ErrOutOfMemory) {
		t.Fatalf("wanted ErrOutOfMemory, got %v (ip=%d sp=%d)", err, v.ip, v.sp)
	}
	if !strings.Contains(err.Error(), "ip ") {
		t.Fatalf("wanted ip in error, got %v", err)
	}
}

func TestExecute_CallIncrViaScratch(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	// Ret discards callee stack; result goes through reserved 0x10 (same as Program C).
	mustExecute(t, v, []Instr{
		{Op: OpPush, Arg: 0x10},
		{Op: OpPush, Arg: 41},
		{Op: OpStore},
		{Op: OpCall, Arg: 7},
		{Op: OpPush, Arg: 0x10},
		{Op: OpLoad},
		{Op: OpHalt},
		{Op: OpPush, Arg: 0x10},
		{Op: OpPush, Arg: 0x10},
		{Op: OpLoad},
		{Op: OpPush, Arg: 1},
		{Op: OpAdd},
		{Op: OpStore},
		{Op: OpRet},
	})

	got, err := v.Peek()
	if err != nil {
		t.Fatalf("Peek failed: %v", err)
	}
	if got != 42 {
		t.Fatalf("wanted 41+1=42, got %d (sp=%d)", got, v.sp)
	}
}

func TestExecute_MissingHalt_ExceedsCode(t *testing.T) {
	mem := NewMemory(64)
	v := NewVM(mem)

	err := v.Execute([]Instr{
		{Op: OpPush, Arg: 1},
	})
	if err == nil {
		t.Fatal("wanted error when program falls off the end, got nil")
	}
	if !strings.Contains(err.Error(), "ip 1") {
		t.Fatalf("wanted ip 1 (past last instr), got %v", err)
	}
}

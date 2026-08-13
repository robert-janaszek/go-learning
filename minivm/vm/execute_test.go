package vm

import (
	"errors"
	"strings"
	"testing"
)

func TestExecute_Add40Plus2(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	code := []Instr{
		{Op: OpPush, Arg: 40},
		{Op: OpPush, Arg: 2},
		{Op: OpAdd},
		{Op: OpHalt},
	}

	if err := v.Execute(code); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

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

	// Alloc 4 → Dup → Push 7 → Store → Dup → Load → leave 7 on stack; Free ptr under it via swap pattern:
	// After Load stack is [ptr, 7]. Pop 7, Free ptr, Push 7 back for assert.
	code := []Instr{
		{Op: OpPush, Arg: 4},
		{Op: OpAlloc},
		{Op: OpDup},
		{Op: OpPush, Arg: 7},
		{Op: OpStore},
		{Op: OpDup},
		{Op: OpLoad},
		{Op: OpHalt},
	}

	if err := v.Execute(code); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	got, err := v.Pop()
	if err != nil {
		t.Fatalf("Pop value failed: %v", err)
	}
	if got != 7 {
		t.Fatalf("wanted loaded value 7, got %d", got)
	}

	ptr, err := v.Pop()
	if err != nil {
		t.Fatalf("Pop ptr failed: %v", err)
	}
	if err := v.Free(Addr(ptr)); err != nil {
		t.Fatalf("Free failed: %v", err)
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

	if err := v.Execute(code); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

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

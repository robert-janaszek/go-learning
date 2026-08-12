package vm

import (
	"strings"
	"testing"
)

func TestVM_PushPop_LIFO(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	if err := v.Push(1); err != nil {
		t.Fatalf("Push(1) failed: %v", err)
	}
	if err := v.Push(2); err != nil {
		t.Fatalf("Push(2) failed: %v", err)
	}
	if err := v.Push(3); err != nil {
		t.Fatalf("Push(3) failed: %v", err)
	}

	got, err := v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if got != 3 {
		t.Fatalf("wanted %d, got %d", 3, got)
	}

	got, err = v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if got != 2 {
		t.Fatalf("wanted %d, got %d", 2, got)
	}

	got, err = v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if got != 1 {
		t.Fatalf("wanted %d, got %d", 1, got)
	}
}

func TestVM_PopOnEmptyStack_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	_, err := v.Pop()
	if err == nil {
		t.Fatal("wanted error for Pop on empty stack, found nil")
	}
}

func TestVM_PushOverflow_ReturnsError(t *testing.T) {
	// For size=8 (and WordSize=4) there is space for exactly 1 valid word:
	// valid aligned addresses are: 4 (address 0 is forbidden by validateAddress).
	mem := NewMemory(8)
	v := NewVM(mem)

	if err := v.Push(123); err != nil {
		t.Fatalf("Push(123) failed: %v", err)
	}

	if err := v.Push(456); err == nil {
		t.Fatal("wanted overflow error for Push on full stack, found nil")
	}
}

func TestVM_DumpShowsValuesAtHighAddresses(t *testing.T) {
	// With size=16 and WordSize=4:
	// initial sp = 16
	// Push(1) stores at addr=12 (0x000c), which is the "high end" of memory.
	mem := NewMemory(16)
	v := NewVM(mem)

	if err := v.Push(1); err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	dump := mem.Dump(0, Addr(mem.Size()))

	// Memory is little-endian, so uint32(1) => bytes: 01 00 00 00.
	// Dump prints: "%04x: %02x %02x %02x %02x %d\n"
	if !strings.Contains(dump, "000c: 01 00 00 00 1\n") {
		t.Fatalf("expected dump line for addr=000c not found.\nDump:\n%s", dump)
	}
}

func TestVM_PeekAndStackDepth(t *testing.T) {
	mem := NewMemory(64)
	v := NewVM(mem)

	if got := v.StackDepth(); got != 0 {
		t.Fatalf("wanted initial StackDepth=0, got %d", got)
	}

	// Top of stack should be the last pushed value.
	if err := v.Push(10); err != nil {
		t.Fatalf("Push failed: %v", err)
	}
	if err := v.Push(20); err != nil {
		t.Fatalf("Push failed: %v", err)
	}

	if got := v.StackDepth(); got != 2 {
		t.Fatalf("wanted StackDepth=2, got %d", got)
	}

	peeked, err := v.Peek()
	if err != nil {
		t.Fatalf("Peek failed: %v", err)
	}
	if peeked != 20 {
		t.Fatalf("wanted Peek=20, got %d", peeked)
	}

	// Peek must not modify the stack.
	if got := v.StackDepth(); got != 2 {
		t.Fatalf("wanted StackDepth still 2 after Peek, got %d", got)
	}

	_, err = v.Pop()
	if err != nil {
		t.Fatalf("Pop failed: %v", err)
	}
	if got := v.StackDepth(); got != 1 {
		t.Fatalf("wanted StackDepth=1 after Pop, got %d", got)
	}
}

func TestVM_PeekOnEmptyStack_ReturnsError(t *testing.T) {
	mem := NewMemory(64)
	v := NewVM(mem)

	_, err := v.Peek()
	if err == nil {
		t.Fatal("wanted error for Peek on empty stack, found nil")
	}
}


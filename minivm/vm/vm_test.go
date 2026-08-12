package vm

import "testing"

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


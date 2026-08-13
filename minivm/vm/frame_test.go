package vm

import "testing"

func TestVM_CallRet_ReturnsSavedAddr(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	if err := v.Call(100, 0); err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	got, err := v.Ret()
	if err != nil {
		t.Fatalf("Ret failed: %v", err)
	}
	if got != 100 {
		t.Fatalf("wanted return addr 100, got %d", got)
	}

	if got := v.StackDepth(); got != 0 {
		t.Fatalf("wanted StackDepth=0 after Ret, got %d", got)
	}
}

func TestVM_CallRet_NestedThreeLevels_LocalsDoNotClobber(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	if err := v.Call(100, 0); err != nil {
		t.Fatalf("Call(100) failed: %v", err)
	}
	if err := v.Push(11); err != nil {
		t.Fatalf("Push outer local failed: %v", err)
	}

	if err := v.Call(200, 0); err != nil {
		t.Fatalf("Call(200) failed: %v", err)
	}
	if err := v.Push(22); err != nil {
		t.Fatalf("Push middle local failed: %v", err)
	}

	if err := v.Call(300, 0); err != nil {
		t.Fatalf("Call(300) failed: %v", err)
	}
	if err := v.Push(33); err != nil {
		t.Fatalf("Push inner local failed: %v", err)
	}

	got, err := v.Peek()
	if err != nil {
		t.Fatalf("Peek inner local failed: %v", err)
	}
	if got != 33 {
		t.Fatalf("wanted inner local 33, got %d", got)
	}

	got, err = v.Ret()
	if err != nil {
		t.Fatalf("inner Ret failed: %v", err)
	}
	if got != 300 {
		t.Fatalf("wanted return addr 300, got %d", got)
	}

	got, err = v.Peek()
	if err != nil {
		t.Fatalf("Peek middle local failed: %v", err)
	}
	if got != 22 {
		t.Fatalf("wanted middle local 22 after inner Ret, got %d", got)
	}

	got, err = v.Ret()
	if err != nil {
		t.Fatalf("middle Ret failed: %v", err)
	}
	if got != 200 {
		t.Fatalf("wanted return addr 200, got %d", got)
	}

	got, err = v.Peek()
	if err != nil {
		t.Fatalf("Peek outer local failed: %v", err)
	}
	if got != 11 {
		t.Fatalf("wanted outer local 11 after middle Ret, got %d", got)
	}

	got, err = v.Ret()
	if err != nil {
		t.Fatalf("outer Ret failed: %v", err)
	}
	if got != 100 {
		t.Fatalf("wanted return addr 100, got %d", got)
	}

	if got := v.StackDepth(); got != 0 {
		t.Fatalf("wanted StackDepth=0 after all Ret, got %d", got)
	}
}

func TestVM_Ret_DiscardsLocalsBelowFrame(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	if err := v.Call(50, 0); err != nil {
		t.Fatalf("Call failed: %v", err)
	}
	if err := v.Push(1); err != nil {
		t.Fatalf("Push local 1 failed: %v", err)
	}
	if err := v.Push(2); err != nil {
		t.Fatalf("Push local 2 failed: %v", err)
	}
	if err := v.Push(3); err != nil {
		t.Fatalf("Push local 3 failed: %v", err)
	}

	got, err := v.Ret()
	if err != nil {
		t.Fatalf("Ret failed: %v", err)
	}
	if got != 50 {
		t.Fatalf("wanted return addr 50, got %d", got)
	}

	if got := v.StackDepth(); got != 0 {
		t.Fatalf("wanted StackDepth=0 after Ret discarded locals, got %d", got)
	}
}

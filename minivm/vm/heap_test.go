package vm

import (
	"errors"
	"testing"
)

func TestHeap_AllocTwoBlocks_AddressSpacing(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr1, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("first Alloc failed: %v", err)
	}

	ptr2, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("second Alloc failed: %v", err)
	}

	wantDelta := Addr(WordSize + 8)
	gotDelta := ptr2 - ptr1
	if gotDelta != wantDelta {
		t.Fatalf("wanted address delta %d, got %d", wantDelta, gotDelta)
	}
}

func TestHeap_Alloc_StoresSizeInHeader(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	headerAddr := ptr - WordSize
	size, err := mem.Load(headerAddr)
	if err != nil {
		t.Fatalf("Load header failed: %v", err)
	}

	if size != 8 {
		t.Fatalf("wanted header size 8, got %d", size)
	}
}

func TestHeap_Alloc_StoreLoadThroughPayload(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	const value uint32 = 4242
	if err := mem.Store(ptr, value); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	got, err := mem.Load(ptr)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if got != value {
		t.Fatalf("wanted %d, got %d", value, got)
	}
}

func TestHeap_AllocZeroNbytes_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	_, err := v.Alloc(0)
	if err == nil {
		t.Fatal("wanted error for Alloc(0), found nil")
	}

	if !errors.Is(err, ErrZeroNbytes) {
		t.Fatalf("wanted ErrZeroNbytes, got %v", err)
	}
}

func TestHeap_AllocUnalignedNbytes_RoundsUpHeader(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(5)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	headerAddr := ptr - WordSize
	size, err := mem.Load(headerAddr)
	if err != nil {
		t.Fatalf("Load header failed: %v", err)
	}

	if size != 8 {
		t.Fatalf("wanted rounded header size 8, got %d", size)
	}
}

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

func TestHeap_FreeThenAlloc_ReusesAddress(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr1, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}
	brkAfterFirst := v.heapBrk

	if err := v.Free(ptr1); err != nil {
		t.Fatalf("Free failed: %v", err)
	}

	ptr2, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc after Free failed: %v", err)
	}

	if ptr2 != ptr1 {
		t.Fatalf("wanted reuse of %d, got %d", ptr1, ptr2)
	}
	if v.heapBrk != brkAfterFirst {
		t.Fatalf("wanted heapBrk unchanged at %d, got %d", brkAfterFirst, v.heapBrk)
	}
}

func TestHeap_FreeThenAlloc_ClearsFreeBit(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}
	if err := v.Free(ptr); err != nil {
		t.Fatalf("Free failed: %v", err)
	}

	reused, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc after Free failed: %v", err)
	}
	if reused != ptr {
		t.Fatalf("wanted reuse of %d, got %d", ptr, reused)
	}

	header, err := mem.Load(ptr - WordSize)
	if err != nil {
		t.Fatalf("Load header failed: %v", err)
	}
	if header&0x80000000 != 0 {
		t.Fatalf("wanted free bit cleared after reuse, got header %#x", header)
	}
}

func TestHeap_FreeZero_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	err := v.Free(0)
	if !errors.Is(err, ErrInvalidFree) {
		t.Fatalf("wanted ErrInvalidFree, got %v", err)
	}
}

func TestHeap_DoubleFree_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}
	if err := v.Free(ptr); err != nil {
		t.Fatalf("first Free failed: %v", err)
	}

	err = v.Free(ptr)
	if !errors.Is(err, ErrDoubleFree) {
		t.Fatalf("wanted ErrDoubleFree, got %v", err)
	}
}

func TestHeap_FreeInsideBlock_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	err = v.Free(ptr + WordSize) // middle of payload, not block start
	if !errors.Is(err, ErrInvalidFree) {
		t.Fatalf("wanted ErrInvalidFree, got %v", err)
	}
}

func TestHeap_AllocUntilStackCollision_ReturnsOOM(t *testing.T) {
	// size=128: heap starts at 0x40, stack sentinel at 128.
	// After filling heap toward SP, next Alloc must fail.
	mem := NewMemory(128)
	v := NewVM(mem)

	for {
		_, err := v.Alloc(8)
		if err != nil {
			if !errors.Is(err, ErrOutOfMemory) {
				t.Fatalf("wanted ErrOutOfMemory, got %v", err)
			}
			return
		}
	}
}

func TestHeap_AllocRespectsStackPointer(t *testing.T) {
	mem := NewMemory(256)
	v := NewVM(mem)

	// Fill most of the upper memory with stack words.
	for {
		if err := v.Push(1); err != nil {
			break
		}
	}

	_, err := v.Alloc(8)
	if !errors.Is(err, ErrOutOfMemory) {
		t.Fatalf("wanted ErrOutOfMemory when stack is full, got %v", err)
	}
}

func TestHeap_FreeList_FirstFitSkipsTooSmall(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	small, err := v.Alloc(4)
	if err != nil {
		t.Fatalf("Alloc(4) failed: %v", err)
	}
	big, err := v.Alloc(16)
	if err != nil {
		t.Fatalf("Alloc(16) failed: %v", err)
	}

	if err := v.Free(small); err != nil {
		t.Fatalf("Free(small) failed: %v", err)
	}
	if err := v.Free(big); err != nil {
		t.Fatalf("Free(big) failed: %v", err)
	}

	// LIFO free-list: head=big, then small.
	// Alloc(16) should take big (first fit), not fail or bump past.
	got, err := v.Alloc(16)
	if err != nil {
		t.Fatalf("Alloc(16) from free-list failed: %v", err)
	}
	if got != big {
		t.Fatalf("wanted first-fit reuse of big=%d, got %d", big, got)
	}
}

func TestHeap_Alloc_UnlinksNonHeadFromFreeList(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	small, err := v.Alloc(4)
	if err != nil {
		t.Fatalf("Alloc(4) failed: %v", err)
	}
	big, err := v.Alloc(16)
	if err != nil {
		t.Fatalf("Alloc(16) failed: %v", err)
	}

	// free-list after these Frees: small -> big (LIFO push on head)
	if err := v.Free(big); err != nil {
		t.Fatalf("Free(big) failed: %v", err)
	}
	if err := v.Free(small); err != nil {
		t.Fatalf("Free(small) failed: %v", err)
	}

	brkBefore := v.heapBrk

	got, err := v.Alloc(16)
	if err != nil {
		t.Fatalf("Alloc(16) failed: %v", err)
	}
	if got != big {
		t.Fatalf("wanted reuse of big=%d, got %d", big, got)
	}
	if v.heapBrk != brkBefore {
		t.Fatalf("wanted heapBrk unchanged at %d, got %d", brkBefore, v.heapBrk)
	}

	// small remains on free-list; its next must no longer point to big
	if v.freeHead != small {
		t.Fatalf("wanted freeHead=%d, got %d", small, v.freeHead)
	}
	next, err := mem.Load(small)
	if err != nil {
		t.Fatalf("Load small.next failed: %v", err)
	}
	if next != 0 {
		t.Fatalf("wanted small.next=0 after unlinking big, got %d", next)
	}
}

func TestHeap_Alloc_ReusesLargerFreeBlockForSmallerRequest(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(16)
	if err != nil {
		t.Fatalf("Alloc(16) failed: %v", err)
	}
	if err := v.Free(ptr); err != nil {
		t.Fatalf("Free failed: %v", err)
	}

	reused, err := v.Alloc(4)
	if err != nil {
		t.Fatalf("Alloc(4) after Free(16) failed: %v", err)
	}
	if reused != ptr {
		t.Fatalf("wanted reuse of %d, got %d", ptr, reused)
	}
}

func TestHeap_Alloc_FallsBackToBumpWhenFreeBlockTooSmall(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	small, err := v.Alloc(4)
	if err != nil {
		t.Fatalf("Alloc(4) failed: %v", err)
	}
	brkAfterSmall := v.heapBrk

	if err := v.Free(small); err != nil {
		t.Fatalf("Free(small) failed: %v", err)
	}

	big, err := v.Alloc(16)
	if err != nil {
		t.Fatalf("Alloc(16) failed: %v", err)
	}
	if big == small {
		t.Fatalf("wanted bump allocation, reused too-small free block at %d", small)
	}
	if v.heapBrk <= brkAfterSmall {
		t.Fatalf("wanted heapBrk to grow past %d, got %d", brkAfterSmall, v.heapBrk)
	}
	if v.freeHead != small {
		t.Fatalf("wanted small block to remain on free-list, freeHead=%d", v.freeHead)
	}
}

func TestHeap_Alloc_NbytesTooLarge_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	// alignUp keeps MSB set — cannot store size without colliding with free-bit
	_, err := v.Alloc(0x80000001)
	if !errors.Is(err, ErrNbytesTooLarge) {
		t.Fatalf("wanted ErrNbytesTooLarge, got %v", err)
	}
}

func TestHeap_Alloc_BumpBlockHasNoFreeBit(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	header, err := mem.Load(headerOf(ptr))
	if err != nil {
		t.Fatalf("Load header failed: %v", err)
	}
	if isFree(header) {
		t.Fatalf("wanted allocated bump block, got free header %#x", header)
	}
}

func TestHeap_Free_SetsFreeBitInHeader(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	if err := v.Free(ptr); err != nil {
		t.Fatalf("Free failed: %v", err)
	}

	header, err := mem.Load(headerOf(ptr))
	if err != nil {
		t.Fatalf("Load header failed: %v", err)
	}
	if !isFree(header) {
		t.Fatalf("wanted free bit set, got header %#x", header)
	}
	if getSizeFromHeader(header) != 8 {
		t.Fatalf("wanted size preserved as 8, got %d", getSizeFromHeader(header))
	}
}

func TestHeap_Free_LinksBlockOnFreeList(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	if err := v.Free(ptr); err != nil {
		t.Fatalf("Free failed: %v", err)
	}

	if v.freeHead != ptr {
		t.Fatalf("wanted freeHead=%d, got %d", ptr, v.freeHead)
	}

	next, err := mem.Load(ptr)
	if err != nil {
		t.Fatalf("Load next failed: %v", err)
	}
	if next != 0 {
		t.Fatalf("wanted next=0 on first free, got %d", next)
	}
}

func TestHeap_Free_SecondBlockBecomesListHead(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	first, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc(first) failed: %v", err)
	}
	second, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc(second) failed: %v", err)
	}

	if err := v.Free(first); err != nil {
		t.Fatalf("Free(first) failed: %v", err)
	}
	if err := v.Free(second); err != nil {
		t.Fatalf("Free(second) failed: %v", err)
	}

	if v.freeHead != second {
		t.Fatalf("wanted freeHead=%d, got %d", second, v.freeHead)
	}

	next, err := mem.Load(second)
	if err != nil {
		t.Fatalf("Load second.next failed: %v", err)
	}
	if Addr(next) != first {
		t.Fatalf("wanted second.next=%d, got %d", first, next)
	}
}

func TestHeap_Free_HeaderPointer_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	ptr, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	err = v.Free(headerOf(ptr))
	if !errors.Is(err, ErrInvalidFree) {
		t.Fatalf("wanted ErrInvalidFree for header pointer, got %v", err)
	}
}

func TestHeap_Free_PointerBeyondHeapBrk_ReturnsError(t *testing.T) {
	mem := NewMemory(1024)
	v := NewVM(mem)

	_, err := v.Alloc(8)
	if err != nil {
		t.Fatalf("Alloc failed: %v", err)
	}

	err = v.Free(v.heapBrk)
	if !errors.Is(err, ErrInvalidFree) {
		t.Fatalf("wanted ErrInvalidFree for ptr beyond heapBrk, got %v", err)
	}
}

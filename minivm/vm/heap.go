package vm

type Heap struct {
	heapStart Addr
	heapBrk   Addr
}

// Block layout (bump allocator):
//
//	[ size:u32 ][ payload... ]
//	            ^
//	            returned Addr
func (v *VM) Alloc(nbytes uint32) (Addr, error) {
	if nbytes == 0 {
		return 0, ErrZeroNbytes
	}

	rem := nbytes % WordSize
	nbytesRefined := nbytes

	if rem != 0 {
		nbytesRefined = nbytes + (WordSize - rem)
	}

	addr := v.heapBrk + WordSize

	if addr+Addr(nbytesRefined) >= v.sp {
		return 0, OutOfMemoryErr
	}

	err := v.mem.Store(v.heapBrk, nbytesRefined)

	if err != nil {
		return 0, err
	}

	v.heapBrk = addr + Addr(nbytesRefined)

	return addr, nil
}

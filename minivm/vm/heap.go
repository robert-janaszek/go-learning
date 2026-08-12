package vm

type Heap struct {
	heapStart Addr
	heapBrk   Addr
}

func (v *VM) Alloc(nbytes uint32) (Addr, error) {
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

package vm

type Heap struct {
	heapStart Addr
	heapBrk   Addr
	freeHead  Addr
}

func isFree(value uint32) bool {
	return value&0x80000000 != 0
}

func getSizeFromHeader(header uint32) uint32 {
	return header & 0x7fffffff
}

func setAllocated(header uint32) uint32 {
	return header & 0x7fffffff
}

func setFree(header uint32) uint32 {
	return header | 0x80000000
}

func headerOf(payload Addr) Addr {
	return payload - WordSize
}

func payloadOf(header Addr) Addr {
	return header + WordSize
}

func alignUp(nbytes uint32) uint32 {
	rem := nbytes % WordSize

	if rem != 0 {
		return nbytes + (WordSize - rem)
	}

	return nbytes
}

func blockSize(size uint32) Addr {
	return Addr(WordSize + size)
}

// free word pointer:
// [ free-bit+size][ next | ... ]
// ^ header        ^ payload
func (v *VM) Alloc(nbytes uint32) (Addr, error) {
	if nbytes == 0 {
		return 0, ErrZeroNbytes
	}

	nbytesRefined := alignUp(nbytes)

	prevPayloadPointer := v.freeHead
	freePayloadPointer := v.freeHead

	for freePayloadPointer != 0 {
		freeHeaderPointer := headerOf(freePayloadPointer)
		freeHeader, err := v.mem.Load(freeHeaderPointer)

		if err != nil {
			return 0, err
		}

		free := isFree(freeHeader)
		size := getSizeFromHeader(freeHeader)

		if !free {
			return 0, ErrCorruptHeap
		}

		if nbytesRefined <= size {
			nextPointer, err := v.mem.Load(freePayloadPointer)

			if err != nil {
				return 0, err
			}

			if freePayloadPointer == v.freeHead {
				v.freeHead = Addr(nextPointer)
			} else {
				err := v.mem.Store(prevPayloadPointer, nextPointer)

				if err != nil {
					return 0, err
				}
			}

			err = v.mem.Store(freeHeaderPointer, setAllocated(freeHeader))

			if err != nil {
				return 0, err
			}

			return freePayloadPointer, nil
		}

		nextPointer, err := v.mem.Load(freePayloadPointer)

		if err != nil {
			return 0, err
		}

		prevPayloadPointer = freePayloadPointer
		freePayloadPointer = Addr(nextPointer)
	}

	msb := nbytesRefined & 0x80000000
	if msb > 0 {
		return 0, ErrNbytesTooLarge
	}

	addr := v.heapBrk + WordSize

	if addr+Addr(nbytesRefined) >= v.sp {
		return 0, ErrOutOfMemory
	}

	err := v.mem.Store(v.heapBrk, nbytesRefined)

	if err != nil {
		return 0, err
	}

	v.heapBrk = addr + Addr(nbytesRefined)

	return addr, nil
}

func (v *VM) Free(ptr Addr) error {
	if ptr < v.heapStart {
		return ErrInvalidFree
	}

	headerAddr := v.heapStart

	for headerAddr < v.heapBrk {
		header, err := v.mem.Load(headerAddr)

		if err != nil {
			return err
		}

		free := isFree(header)
		size := getSizeFromHeader(header)
		payload := payloadOf(headerAddr)

		if payload == ptr {
			if free {
				return ErrDoubleFree
			}

			err := v.mem.Store(headerAddr, setFree(header))

			if err != nil {
				return err
			}

			err = v.mem.Store(payload, uint32(v.freeHead))

			if err != nil {
				return err
			}

			v.freeHead = ptr

			return nil
		}
		if payload > ptr {
			return ErrInvalidFree
		}

		headerAddr += blockSize(size)
	}

	return ErrInvalidFree
}

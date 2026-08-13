package vm

import "errors"

var (
	ErrOutOfMemory    = errors.New("out of memory")
	ErrZeroNbytes     = errors.New("nbytes must be greater than 0")
	ErrNbytesTooLarge = errors.New("nbytes too large")
	ErrDoubleFree     = errors.New("cannot double-free a block")
	ErrInvalidFree    = errors.New("invalid pointer")
	ErrCorruptHeap    = errors.New("corrupted heap")
)

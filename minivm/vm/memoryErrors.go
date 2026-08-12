package vm

import "errors"

var (
	OutOfMemoryErr = errors.New("out of memory")
	ErrZeroNbytes  = errors.New("nbytes must be greater than 0")
)

package vm

func (v *VM) Call(returnAddr uint32, nargs int) error {
	err := v.Push(returnAddr)

	if err != nil {
		return err
	}

	err = v.Push(uint32(v.fp))

	if err != nil {
		return err
	}

	v.fp = v.sp

	return nil
}

func (v *VM) Ret() (uint32, error) {
	v.sp = v.fp

	callerFP, err := v.Pop()

	if err != nil {
		return 0, err
	}

	returnAddr, err := v.Pop()

	if err != nil {
		return 0, err
	}

	v.fp = Addr(callerFP)

	return returnAddr, nil
}

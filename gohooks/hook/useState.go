package hook

func UseState[T any](initial T) (T, func(T)) {
	r := runtime
	index := r.hookIndex
	defer func() {
		r.hookIndex++
	}()

	setter := func(state T) {
		r.hookState[index].state = state
		select {
		case r.updates <- struct{}{}:
		default:
		}
	}

	if len(r.hookState) <= index {
		r.hookState = append(r.hookState, hookState{state: initial})
		return initial, setter
	}

	stateStruct := r.hookState[index]
	return stateStruct.state.(T), setter
}

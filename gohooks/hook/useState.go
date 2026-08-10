package hook

import "reflect"

func UseState[T any](initial T) (T, func(T)) {
	r := runtime
	index := r.hookIndex
	defer func() {
		r.hookIndex++
	}()

	setter := func(state T) {
		if reflect.DeepEqual(r.hookState[index].state, state) {
			return
		}

		r.hookState[index].consecutiveSchedules++
		r.hookState[index].lastScheduled = true
		r.hookState[index].state = state
		select {
		case r.updates <- struct{}{}:
		default:
		}
	}

	if len(r.hookState) <= index {
		r.hookState = append(r.hookState, hookState{
			state:                initial,
			lastScheduled:        false,
			consecutiveSchedules: 0,
		})
		return initial, setter
	}

	if !r.hookState[index].lastScheduled {
		r.hookState[index].consecutiveSchedules = 0
	} else {
		r.hookState[index].lastScheduled = false
	}

	stateStruct := r.hookState[index]
	return stateStruct.state.(T), setter
}

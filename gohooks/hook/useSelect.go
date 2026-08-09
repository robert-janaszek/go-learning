package hook

import "reflect"

func UseSelect[S, T any](get func() S, selectFn func(s S) T) T {
	r := runtime
	index := r.hookIndex
	defer func() {
		r.hookIndex++
	}()

	store := get()
	val := selectFn(store)

	if len(r.hookState) <= index {
		r.hookState = append(r.hookState, hookState{state: val})
		return val
	}

	if reflect.DeepEqual(r.hookState[index].state, val) {
		return val
	}

	r.hookState[index].state = val
	// select { // TODO: to be added when there is tree of components
	// case r.updates <- struct{}{}:
	// default:
	// }

	return val
}

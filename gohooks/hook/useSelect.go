package hook

import "reflect"

// UseSelect caches the selected slice T for this hook slot.
// Today it only compares during Render — it does not skip component renders by itself.
//
// TODO: add a subscribable store (alongside the component tree). On store notify,
// recompute get+selectFn outside c(); schedule Render only when T changed.
// Persist get/selectFn (and instance) on the slot for subscribe/unsubscribe on mount/unmount.
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
	return val
}

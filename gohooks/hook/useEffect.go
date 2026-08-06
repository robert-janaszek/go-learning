package hook

import "reflect"

func UseEffect(effect func() func(), deps []any) {
	r := runtime
	index := r.effectIndex
	defer func() {
		r.effectIndex++
	}()

	if len(r.effectState) <= index {
		r.effectState = append(r.effectState, effectState{
			effect:         effect,
			deps:           deps,
			runAfterRender: true,
		})

		return
	}

	r.effectState[index].runAfterRender = false

	if !reflect.DeepEqual(r.effectState[index].deps, deps) {
		r.effectState[index].runAfterRender = true
		r.effectState[index].effect = effect
		r.effectState[index].deps = deps
	}
}

package components

import (
	"fmt"
	"gohooks/hook"
)

func innerCounter(id string) hook.Result {
	v, s := hook.UseState(0)

	fmt.Printf("%s: %d\n", id, v)
	if v == 0 {
		s(v + 1)
	}

	return hook.Result{}
}

func NestedCounter() hook.Result {
	v, s := hook.UseState(0)
	cancel := hook.UseCancel()

	s(v + 1)
	if v == 1 {
		cancel()
	}

	child1 := hook.Element{
		Key:       "c1",
		Component: func() hook.Result { return innerCounter("c1") },
	}
	child2 := hook.Element{
		Key:       "c2",
		Component: func() hook.Result { return innerCounter("c2") },
	}
	child3 := hook.Element{
		Key:       "c3",
		Component: func() hook.Result { return innerCounter("c3") },
	}

	if v == 0 {
		return hook.Result{
			Children: []hook.Element{
				child1, child2, child3,
			},
		}
	}

	return hook.Result{
		Children: []hook.Element{
			child1, child3,
		},
	}
}

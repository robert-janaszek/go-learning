package components

import (
	"fmt"
	"gohooks/hook"
)

func Counter() {
	value, set := hook.UseState(0)
	cancel := hook.UseCancel()

	if value < 20 {
		set(value + 1)
	} else {
		cancel()
	}

	lower10 := min(value, 10)

	deps := []any{lower10}
	hook.UseEffect(func() func() {
		fmt.Printf("value changed: %d\n", value)

		return func() {
			fmt.Printf("cleanup %d\n", value)
		}
	}, deps)
}

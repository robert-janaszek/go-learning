package components

import (
	"fmt"
	"gohooks/hook"
)

func Counter() {
	value, set := hook.UseState(0)

	fmt.Println(value)
	if value < 20 {
		set(value + 1)
	}
}

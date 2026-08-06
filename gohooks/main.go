package main

import (
	"fmt"
	"gohooks/hook"
)

func main() {
	runtime := hook.Runtime{}
	runtime.Render(func() {
		value, set := hook.UseState(0)
		fmt.Println(value)
		set(1)
		fmt.Println(runtime)
	})
}

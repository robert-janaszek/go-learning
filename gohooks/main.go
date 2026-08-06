package main

import (
	"fmt"
	"gohooks/hook"
)

func main() {
	runtime := hook.Runtime{}
	runtime.Render(func() {
		fmt.Println("hello")
	})
}

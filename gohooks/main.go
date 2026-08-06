package main

import (
	"gohooks/components"
	"gohooks/hook"
)

func main() {
	runtime := hook.Runtime{}
	runtime.Run(components.Counter)
}

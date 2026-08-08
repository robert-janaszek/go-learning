package main

import (
	"gohooks/components"
	"gohooks/hook"
)

func main() {
	runtime := hook.CreateRuntime()
	runtime.Run(components.Counter)
	runtime.Unmount()
}

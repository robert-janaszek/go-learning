package main

import (
	"context"
	"gohooks/components"
	"gohooks/hook"
)

func main() {
	ctx := context.Background()
	runtime := hook.CreateRuntime()
	runtime.Run(ctx, components.NestedCounter)
}

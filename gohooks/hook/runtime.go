package hook

import "context"

type Runtime struct {
	hookState   []hookState
	hookIndex   int
	effectState []effectState
	effectIndex int
	updates     chan struct{}
	cancel      context.CancelFunc
}

func CreateRuntime() *Runtime {
	rt := &Runtime{
		updates: make(chan struct{}, 1),
	}

	rt.updates <- struct{}{}

	return rt
}

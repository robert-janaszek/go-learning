package hook

type Runtime struct {
	hookState       []hookState
	hookIndex       int
	effectState     []effectState
	effectIndex     int
	component       Component
	numberOfHooks   int
	numberOfEffects int
	updates         chan struct{}
}

func CreateRuntime() *Runtime {
	rt := &Runtime{
		updates: make(chan struct{}, 1),
	}

	rt.updates <- struct{}{}

	return rt
}

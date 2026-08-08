package hook

var runtime *Runtime

func (r *Runtime) Render(c Component) {
	runtime = r
	defer func() {
		runtime = nil
	}()

	r.hookIndex = 0
	r.effectIndex = 0
	r.component = c

	c()

	for i := range r.effectState {
		effect := &r.effectState[i]
		if effect.runAfterRender {
			if effect.cleanup != nil {
				effect.cleanup()
			}
			cleanup := effect.effect()
			effect.cleanup = cleanup
		}
	}
}

func (r *Runtime) Run(c Component) {
	i := 0
	for ; i < 50; i++ {
		select {
		case <-r.updates:
		default: // TODO: change for asynchronous handling
			return // as long as there is synchronous set state keep looping
		}
		r.Render(c)

		if i == 0 {
			r.numberOfHooks = len(r.hookState)
			r.numberOfEffects = len(r.effectState)
		}

		if r.hookIndex != r.numberOfHooks {
			panic("hooks order mismatch")
		}

		if r.effectIndex != r.numberOfEffects {
			panic("effects order mismatch")
		}
	}

	if i == 50 {
		panic("infinite loop found")
	}
}

func (r *Runtime) Unmount() {
	for i := range r.effectState {
		effect := &r.effectState[i]

		if effect.cleanup != nil {
			effect.cleanup()
		}
	}

	r.hookState = nil
	r.effectState = nil
}

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
		r.dirty = false
		r.Render(c)

		if !r.dirty {
			break
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

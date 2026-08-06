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

	for _, effect := range r.effectState {
		if effect.runAfterRender {
			effect.effect()
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

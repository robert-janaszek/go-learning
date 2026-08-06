package hook

func (r *Runtime) Render(c Component) {
	r.hookIndex = 0
	r.component = c

	c()
}

package hook

type Runtime struct {
	hookState []hookState
	hookIndex int
	component Component
	dirty     bool
}

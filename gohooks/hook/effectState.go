package hook

type effectState struct {
	effect         func() func()
	deps           []any
	runAfterRender bool
	cleanup        func()
}

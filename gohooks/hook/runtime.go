package hook

type Runtime struct {
	hookState       []hookState
	hookIndex       int
	effectState     []effectState
	effectIndex     int
	component       Component
	dirty           bool
	numberOfHooks   int
	numberOfEffects int
}

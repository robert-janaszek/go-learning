package hook

type Instance struct {
	hookState       []hookState
	numberOfHooks   int
	effectState     []effectState
	numberOfEffects int
	children        map[string]*Instance
	initialized     bool
}

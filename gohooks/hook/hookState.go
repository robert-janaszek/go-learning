package hook

type hookState struct {
	state                any
	lastScheduled        bool
	consecutiveSchedules int
}

package hook

type Element struct {
	Key       string
	Component Component
}

type Result struct {
	Children []Element
	Out      string
}

type Component func() Result

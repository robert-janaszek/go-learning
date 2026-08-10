package hook

import (
	"context"
	"fmt"
)

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

	for i := range r.hookState {
		if r.hookState[i].consecutiveSchedules >= 10 {
			panic(fmt.Sprintf("infinite state update found, hook index %d", i))
		}
	}
}

func (r *Runtime) Run(ctx context.Context, c Component) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r.cancel = cancel
	for i := 0; ; i++ {
		select {
		case <-r.updates:
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
		case <-ctx.Done():
			return
		}
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

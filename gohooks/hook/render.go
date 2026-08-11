package hook

import (
	"context"
	"fmt"
)

var runtime *Runtime

func (r *Runtime) Render(instance *Instance, root Component) {
	r.effectState = instance.effectState
	r.hookState = instance.hookState

	result := r.RenderInstance(root)

	instance.hookState = r.hookState
	instance.effectState = r.effectState

	if !instance.initialized {
		instance.initialized = true
		instance.numberOfHooks = len(r.hookState)
		instance.numberOfEffects = len(r.effectState)
	}

	if r.hookIndex != instance.numberOfHooks {
		panic("hooks order mismatch")
	}

	if r.effectIndex != instance.numberOfEffects {
		panic("effects order mismatch")
	}

	for i := range result.Children {
		child := result.Children[i]

		inst := instance.children[child.Key]

		if inst == nil {
			inst = &Instance{
				children: make(map[string]*Instance),
			}
			instance.children[child.Key] = inst
		}

		r.Render(inst, result.Children[i].Component)
	}
}

func (r *Runtime) RenderInstance(c Component) Result {
	r.hookIndex = 0
	r.effectIndex = 0

	result := c()

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

	return result
}

func (r *Runtime) Run(ctx context.Context, root Component) {
	runtime = r
	defer func() {
		runtime = nil
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	r.cancel = cancel

	instance := Instance{
		children: make(map[string]*Instance),
	}

	for i := 0; ; i++ {
		select {
		case <-r.updates:
			r.Render(&instance, root)
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

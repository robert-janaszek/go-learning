package hook

import (
	"context"
	"testing"
)

func TestRuntimeRun(t *testing.T) {
	type got struct {
		renders         int
		hookCount       int
		state           int
		effectRunCount  int
		cleanupRunCount int
	}

	tests := []struct {
		name string
		run  func() got
		want got
	}{
		{
			name: "no-change",
			run: func() got {
				rt := CreateRuntime()
				renders := 0
				effectRunCount := 0
				cleanupRunCount := 0

				rt.Run(context.Background(), func() Result {
					value, _ := UseState(10)
					cancel := UseCancel()
					UseEffect(func() func() {
						effectRunCount++
						return func() {
							cleanupRunCount++
						}
					}, []any{value})
					renders++
					cancel()

					return Result{}
				})

				return got{
					renders:         renders,
					hookCount:       len(rt.hookState),
					state:           rt.hookState[0].state.(int),
					effectRunCount:  effectRunCount,
					cleanupRunCount: cleanupRunCount,
				}
			},
			// +1 cleanup from Unmount on ctx.Done
			want: got{renders: 1, hookCount: 1, state: 10, effectRunCount: 1, cleanupRunCount: 1},
		},
		{
			name: "one-state-update-with-effect",
			run: func() got {
				rt := CreateRuntime()
				renders := 0
				effectRunCount := 0
				cleanupRunCount := 0

				rt.Run(context.Background(), func() Result {
					value, set := UseState(10)
					cancel := UseCancel()
					if value == 10 {
						set(value + 1)
					} else {
						cancel()
					}
					UseEffect(func() func() {
						effectRunCount++
						return func() {
							cleanupRunCount++
						}
					}, []any{value})
					renders++

					return Result{}
				})

				return got{
					renders:         renders,
					hookCount:       len(rt.hookState),
					state:           rt.hookState[0].state.(int),
					effectRunCount:  effectRunCount,
					cleanupRunCount: cleanupRunCount,
				}
			},
			// deps cleanup + Unmount cleanup
			want: got{renders: 2, hookCount: 1, state: 11, effectRunCount: 2, cleanupRunCount: 2},
		},
		{
			name: "one-state-update",
			run: func() got {
				rt := CreateRuntime()
				renders := 0
				effectRunCount := 0
				cleanupRunCount := 0

				rt.Run(context.Background(), func() Result {
					value, set := UseState(10)
					cancel := UseCancel()
					if value == 10 {
						set(value + 1)
					} else {
						cancel()
					}
					UseEffect(func() func() {
						effectRunCount++
						return func() {
							cleanupRunCount++
						}
					}, []any{})
					renders++

					return Result{}
				})

				return got{
					renders:         renders,
					hookCount:       len(rt.hookState),
					state:           rt.hookState[0].state.(int),
					effectRunCount:  effectRunCount,
					cleanupRunCount: cleanupRunCount,
				}
			},
			// empty deps: only Unmount cleanup
			want: got{renders: 2, hookCount: 1, state: 11, effectRunCount: 1, cleanupRunCount: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.run()
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestUseStateRerenderOnChange(t *testing.T) {
	tests := []struct {
		name        string
		initial     int
		sets        []int // values set during first render
		wantRenders int
		wantState   int
	}{
		{
			name:        "same-value-skips-rerender",
			initial:     10,
			sets:        []int{10, 10},
			wantRenders: 1,
			wantState:   10,
		},
		{
			name:        "different-value-rerenders-same-after-does-not",
			initial:     10,
			sets:        []int{11},
			wantRenders: 2,
			wantState:   11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := CreateRuntime()
			renders := 0

			rt.Run(context.Background(), func() Result {
				value, set := UseState(tt.initial)
				cancel := UseCancel()
				renders++

				if renders == 1 {
					for _, v := range tt.sets {
						set(v)
					}
					// same-value case: no signal → must cancel or Run blocks
					if tt.wantRenders == 1 {
						cancel()
					}
					return Result{}
				}

				// after a real update, setting the same value again must not loop
				set(value)
				set(value)
				cancel()

				return Result{}
			})

			if renders != tt.wantRenders {
				t.Errorf("renders: got %d, want %d", renders, tt.wantRenders)
			}
			if got := rt.hookState[0].state.(int); got != tt.wantState {
				t.Errorf("state: got %d, want %d", got, tt.wantState)
			}
		})
	}
}

func TestRuntimeCleanup(t *testing.T) {
	runLog := []string{}

	rt := CreateRuntime()
	rt.Run(context.Background(), func() Result {
		value, set := UseState(10)
		cancel := UseCancel()
		if value == 10 {
			set(value + 1)
		} else {
			cancel()
		}
		UseEffect(func() func() {
			runLog = append(runLog, "effect")
			return func() {
				runLog = append(runLog, "cleanup")
			}
		}, []any{value})

		return Result{}
	})

	// mount effect → deps cleanup+effect → Unmount cleanup
	want := []string{"effect", "cleanup", "effect", "cleanup"}
	if len(runLog) != len(want) {
		t.Fatalf("wanted len = %d; got %d (%v)", len(want), len(runLog), runLog)
	}
	for i := range want {
		if runLog[i] != want[i] {
			t.Errorf("log[%d]: want %q; got %q", i, want[i], runLog[i])
		}
	}
}

func TestTreeRemoveChildRunsCleanup(t *testing.T) {
	cleanups := map[string]int{}
	lastState := map[string]int{}

	makeChild := func(id string) Component {
		return func() Result {
			v, set := UseState(0)
			UseEffect(func() func() {
				return func() {
					cleanups[id]++
				}
			}, []any{})
			if v == 0 {
				set(1)
			}
			lastState[id] = v
			return Result{}
		}
	}

	rt := CreateRuntime()
	rt.Run(context.Background(), func() Result {
		phase, setPhase := UseState(0)
		cancel := UseCancel()

		a := Element{Key: "a", Component: makeChild("a")}
		b := Element{Key: "b", Component: makeChild("b")}
		c := Element{Key: "c", Component: makeChild("c")}

		switch phase {
		case 0:
			setPhase(1)
			return Result{Children: []Element{a, b, c}}
		case 1:
			setPhase(2)
			// drop middle child — should Unmount b before rendering a,c
			return Result{Children: []Element{a, c}}
		default:
			cancel()
			return Result{Children: []Element{a, c}}
		}
	})

	if cleanups["b"] != 1 {
		t.Fatalf("child b cleanup: want 1 (removed from tree), got %d (%v)", cleanups["b"], cleanups)
	}
	// a and c cleaned on final root Unmount
	if cleanups["a"] != 1 || cleanups["c"] != 1 {
		t.Fatalf("want a=1 c=1 final unmount cleanups, got %v", cleanups)
	}
	// surviving children kept state through the removal frame
	if lastState["a"] != 1 || lastState["c"] != 1 {
		t.Fatalf("want a,c state 1 after surviving remove, got %v", lastState)
	}
	if _, ok := lastState["b"]; ok && lastState["b"] != 0 && lastState["b"] != 1 {
		t.Fatalf("unexpected last state for b: %v", lastState["b"])
	}
}

func TestTreeRootUnmountCleansAllChildren(t *testing.T) {
	cleanups := 0

	child := func() Result {
		UseEffect(func() func() {
			return func() {
				cleanups++
			}
		}, []any{})
		return Result{}
	}

	rt := CreateRuntime()
	rt.Run(context.Background(), func() Result {
		step, setStep := UseState(0)
		cancel := UseCancel()

		if step == 0 {
			setStep(1)
			return Result{Children: []Element{
				{Key: "x", Component: child},
				{Key: "y", Component: child},
			}}
		}

		cancel()
		return Result{Children: []Element{
			{Key: "x", Component: child},
			{Key: "y", Component: child},
		}}
	})

	if cleanups != 2 {
		t.Fatalf("want 2 child cleanups on root Unmount, got %d", cleanups)
	}
}

func TestTreeSiblingStateIsolation(t *testing.T) {
	left, right := -1, -1

	rt := CreateRuntime()
	rt.Run(context.Background(), func() Result {
		step, setStep := UseState(0)
		cancel := UseCancel()

		leftChild := Element{
			Key: "left",
			Component: func() Result {
				v, set := UseState(0)
				if v == 0 {
					set(7)
				}
				left = v
				return Result{}
			},
		}
		rightChild := Element{
			Key: "right",
			Component: func() Result {
				v, set := UseState(0)
				if v == 0 {
					set(9)
				}
				right = v
				return Result{}
			},
		}

		if step == 0 {
			setStep(1)
			return Result{Children: []Element{leftChild, rightChild}}
		}
		cancel()
		return Result{Children: []Element{leftChild, rightChild}}
	})

	if left != 7 || right != 9 {
		t.Fatalf("want independent child state left=7 right=9, got left=%d right=%d", left, right)
	}
}

func TestTreeReorderWithAndWithoutKeys(t *testing.T) {
	t.Run("with keys state follows identity", func(t *testing.T) {
		var order []int

		child := func(initial int) Component {
			return func() Result {
				v, _ := UseState(initial)
				order = append(order, v)
				return Result{}
			}
		}

		rt := CreateRuntime()
		rt.Run(context.Background(), func() Result {
			step, setStep := UseState(0)
			cancel := UseCancel()
			a := Element{Key: "a", Component: child(10)}
			b := Element{Key: "b", Component: child(20)}

			order = nil
			if step == 0 {
				setStep(1)
				return Result{Children: []Element{a, b}}
			}
			cancel()
			// swap order; keys should keep 10 on a and 20 on b
			return Result{Children: []Element{b, a}}
		})

		if len(order) != 2 || order[0] != 20 || order[1] != 10 {
			t.Fatalf("want render order values [20 10] after keyed swap, got %v", order)
		}
	})

	t.Run("without keys state follows index", func(t *testing.T) {
		var order []int

		child := func(initial int) Component {
			return func() Result {
				v, _ := UseState(initial)
				order = append(order, v)
				return Result{}
			}
		}

		rt := CreateRuntime()
		rt.Run(context.Background(), func() Result {
			step, setStep := UseState(0)
			cancel := UseCancel()
			// empty Key → refineKey falls back to index
			first := Element{Component: child(10)}
			second := Element{Component: child(20)}

			order = nil
			if step == 0 {
				setStep(1)
				return Result{Children: []Element{first, second}}
			}
			cancel()
			// swap components; index slots keep old state (10 stays at i:0)
			return Result{Children: []Element{second, first}}
		})

		if len(order) != 2 || order[0] != 10 || order[1] != 20 {
			t.Fatalf("want index-stable values [10 20] after keyless swap, got %v", order)
		}
	})
}

func TestRuntimeCancelFromParent(t *testing.T) {
	rt := CreateRuntime()
	ctx, cancel := context.WithCancel(context.Background())

	renders := 0
	mounted := make(chan struct{})
	done := make(chan struct{})
	go func() {
		rt.Run(ctx, func() Result {
			renders++
			_, _ = UseState(0)
			if renders == 1 {
				close(mounted)
			}

			return Result{}
		})
		close(done)
	}()

	<-mounted
	cancel()
	<-done

	if renders != 1 {
		t.Fatalf("want 1 render before idle+parent cancel, got %d", renders)
	}
}

func TestRuntimeStateUpdatePanic(t *testing.T) {
	rt := CreateRuntime()
	ctx := context.Background()
	renders := 0

	defer func() {
		if r := recover(); r != nil {
			if renders != 10 {
				t.Fatalf("expected panic after 10 renders, found %d", renders)
			}
		} else {
			t.Fatalf("missing panic for infinite updates")
		}
	}()

	rt.Run(ctx, func() Result {
		renders++
		val, set := UseState(0)

		set(val + 1)

		return Result{}
	})
}

func TestRuntimeStateStorm(t *testing.T) {
	rt := CreateRuntime()
	ctx, cancel := context.WithCancel(context.Background())
	renders := 0

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic")
		}
	}()

	cancelled := false

	rt.Run(ctx, func() Result {
		renders++
		val, set := UseState(0)
		val1, set1 := UseState(1)

		if (val+val1)%2 == 0 {
			set(val + 1)
		} else {
			set1(val1 + 1)
		}

		if val+val1 > 50 {
			cancelled = true
			cancel()
		}

		return Result{}
	})

	if !cancelled {
		t.Fatal("expected cancel, but not found")
	}
}

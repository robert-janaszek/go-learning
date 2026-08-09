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

				rt.Run(context.Background(), func() {
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
				})

				return got{
					renders:         renders,
					hookCount:       len(rt.hookState),
					state:           rt.hookState[0].state.(int),
					effectRunCount:  effectRunCount,
					cleanupRunCount: cleanupRunCount,
				}
			},
			want: got{renders: 1, hookCount: 1, state: 10, effectRunCount: 1, cleanupRunCount: 0},
		},
		{
			name: "one-state-update-with-effect",
			run: func() got {
				rt := CreateRuntime()
				renders := 0
				effectRunCount := 0
				cleanupRunCount := 0

				rt.Run(context.Background(), func() {
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
				})

				return got{
					renders:         renders,
					hookCount:       len(rt.hookState),
					state:           rt.hookState[0].state.(int),
					effectRunCount:  effectRunCount,
					cleanupRunCount: cleanupRunCount,
				}
			},
			want: got{renders: 2, hookCount: 1, state: 11, effectRunCount: 2, cleanupRunCount: 1},
		},
		{
			name: "one-state-update",
			run: func() got {
				rt := CreateRuntime()
				renders := 0
				effectRunCount := 0
				cleanupRunCount := 0

				rt.Run(context.Background(), func() {
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
				})

				return got{
					renders:         renders,
					hookCount:       len(rt.hookState),
					state:           rt.hookState[0].state.(int),
					effectRunCount:  effectRunCount,
					cleanupRunCount: cleanupRunCount,
				}
			},
			want: got{renders: 2, hookCount: 1, state: 11, effectRunCount: 1, cleanupRunCount: 0},
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

			rt.Run(context.Background(), func() {
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
					return
				}

				// after a real update, setting the same value again must not loop
				set(value)
				set(value)
				cancel()
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
	rt.Run(context.Background(), func() {
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
	})
	rt.Unmount()

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

func TestRuntimeCancelFromParent(t *testing.T) {
	rt := CreateRuntime()
	ctx, cancel := context.WithCancel(context.Background())

	renders := 0
	mounted := make(chan struct{})
	done := make(chan struct{})
	go func() {
		rt.Run(ctx, func() {
			renders++
			_, _ = UseState(0)
			if renders == 1 {
				close(mounted)
			}
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

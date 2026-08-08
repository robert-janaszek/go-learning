package hook

import (
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

				rt.Run(func() {
					value, _ := UseState(10)
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
			want: got{renders: 1, hookCount: 1, state: 10, effectRunCount: 1, cleanupRunCount: 0},
		},
		{
			name: "one-state-update-with-effect",
			run: func() got {
				rt := CreateRuntime()
				renders := 0
				effectRunCount := 0
				cleanupRunCount := 0

				rt.Run(func() {
					value, set := UseState(10)
					if value == 10 {
						set(value + 1)
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

				rt.Run(func() {
					value, set := UseState(10)
					if value == 10 {
						set(value + 1)
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

func TestRuntimeCleanup(t *testing.T) {
	runLog := []string{}

	rt := CreateRuntime()
	rt.Run(func() {
		value, set := UseState(10)
		if value == 10 {
			set(value + 1)
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

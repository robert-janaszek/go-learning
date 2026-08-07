Here are **18 project exercises: A React Hooks–style mini-runtime (in Go)**.

In React you have `useState`, `useEffect`, and function components. Here you will build a **simplified counterpart of that model in Go** — not to replace React, but to understand:

* how hook call order works,
* where state lives between “renders”,
* how an effect depends on a dependency list,
* how closures interact with mutable runtime state.

This is an architecture exercise. You implement a small engine plus one sample component (e.g. a counter).

---

```text
gohooks/
├── go.mod
├── main.go
├── hook/
│   ├── runtime.go     # render context, hook queue, Render()
│   ├── state.go       # UseState
│   ├── effect.go      # UseEffect
│   └── runtime_test.go
└── components/
    └── counter.go     # sample component that uses hooks
```

You may tweak names slightly, but keep the split: **runtime / state / effect / component**.

---

## Phase 1: Render model and runtime (Steps 1–5)

### Step 1: Module and layout

Create a `gohooks` directory (next to `taskmanager`) and initialize the module:

```bash
go mod init gohooks
```

Add the `hook/` and `components/` folders.

### Step 2: Component type

In the `hook` package, define a component function type, e.g.:

```go
type Component func()
```

For now the component returns nothing — a “render” is a side effect (`fmt.Println` as console UI). Later you can change it to `func() string` if you prefer returning text to print.

### Step 3: Single hook slot state

Define an internal hook slot struct, e.g. `hookState`, that can hold:

* a state value (`any` / `interface{}`),
* optionally effect data (deps, cleanup) — you may split this into separate types in later steps.

To start, a single value field (`any`) is enough.

### Step 4: Render runtime

Create a `Runtime` struct (name is up to you) that holds:

* a list / slice of hook slots (`[]hookState` or similar),
* the current hook index during one render (`hookIndex int`),
* a reference to the currently rendered `Component`,
* a “needs re-render” flag (e.g. `dirty bool`) or an update queue.

Same rule as React: **during a single component call, hooks are read/created in a fixed order**, and `hookIndex` advances on every `UseState` / `UseEffect`.

### Step 5: `Render` function

Write `func (r *Runtime) Render(c Component)` (or `Mount`) that:

1. Resets `hookIndex` to `0` (order from the start).
2. Sets the current component.
3. Calls `c()`.
4. After the render, runs effects (for now you can leave this empty — you will fill it in Phase 3).

From `main`, run an empty component `func() { fmt.Println("hello") }` and confirm it prints.

---

## Phase 2: `UseState` (Steps 6–10)

### Step 6: Access to the current runtime

Hooks must know *which* runtime they run in. Pick one approach (simpler → harder):

1. **A package-level pointer** to the active `*Runtime` (set inside `Render`) — simplest to start, like early React.
2. **Context passed some other way** — later; don’t complicate it yet.

Implement option 1: at the start of `Render` set `activeRuntime = r`, clear it on exit (e.g. `defer`).

### Step 7: `UseState` — first call vs later calls

Implement:

```go
func UseState[T any](initial T) (T, func(T))
```

(or without generics: `UseState(initial any) (any, func(any))` if you want to avoid type parameters for now).

Behavior:

* On the **first** render, when there is no slot at `hookIndex` yet: create the slot, store `initial`.
* On **later** renders: read the value from the slot at `hookIndex` (do not reset to `initial`).
* Always increment `hookIndex` by 1 before returning.

### Step 8: Setter and scheduling a re-render

The setter returned from `UseState` should:

1. Write the new value into the correct slot (note: the setter outlives a single render — it must remember the **slot index**, not rely on the current `hookIndex`).
2. Mark the runtime as “dirty” / schedule another render.

At this stage the setter **does not** need to call `Render` synchronously — a flag is enough. You will add the loop in the next step.

### Step 9: Re-render loop

Add a method like `func (r *Runtime) Run(c Component)` that:

1. Calls `Render(c)`.
2. While the “dirty” flag is set: clears the flag and calls `Render(c)` again.
3. Guard against infinite loops (e.g. a setter called unconditionally in the component body) — you may add a render limit (e.g. max 50) and `panic` / return an error if exceeded.

### Step 10: `Counter` component

In `components/counter.go`, write a component that:

* calls `UseState(0)` for the counter,
* prints the current value,
* somehow increments state (e.g. only once on the first render — or via a simple “event simulation” from `main`).

A `main` test scenario is enough for now:

* mount `Counter`,
* call the setter from outside (if you expose it) **or** inside the component after printing, schedule `+1` only when the value is `< N`.

Goal: see console values `0, 1, 2, …` coming from re-renders, not from a manual `for` loop in the component that bypasses hooks.

---

## Phase 3: `UseEffect` (Steps 11–15)

### Step 11: Effect slot

Extend the slot model (or add a separate type) with effect data:

* effect function: `func() func()` (the effect returns an optional `cleanup`; cleanup may be `nil`),
* previous dependencies: e.g. `[]any`,
* a flag “run after this render”.

Target signature:

```go
func UseEffect(effect func() func(), deps ...any)
```

If you prefer no variadic: `UseEffect(effect func() func(), deps []any)`.

### Step 12: Comparing the dependency list

During `UseEffect`:

1. Get / create the slot at the current `hookIndex`, then advance the index.
2. Compare the new `deps` with the stored previous ones.
3. If this is the first time **or** deps changed → mark the effect to run after the render.
4. If deps did **not** change → do not run the effect (like React).

Comparison: simple `==` after casting, or `reflect.DeepEqual`, is fine for learning.

**React note:** `deps == nil` vs `deps == []` — you may simplify: always require a deps list (even empty). An empty slice = “run once after mount”.

### Step 13: Running effects after render

After calling the component in `Render`:

1. Walk slots marked to run.
2. If the slot had a previous `cleanup` — call it **before** running the new effect.
3. Call the effect, store the returned cleanup.
4. Store the new deps in the slot.

### Step 14: Effect in `Counter`

Extend the component:

* `count, setCount := UseState(0)`
* `UseEffect(func() func() { fmt.Println("effect: count =", count); return nil }, count)`

Or with cleanup:

```text
effect start for count=X
… re-render …
cleanup for count=X
effect start for count=X+1
```

Confirm that when `count` does not change, the effect does **not** run again (do two renders with the same value — e.g. setter with the same number — and observe no new log).

### Step 15: Unmount / final cleanup

Add `func (r *Runtime) Unmount()`:

* call every stored effect `cleanup`,
* clear slots / mark the runtime as unmounted.

From `main`: `Run` → a few updates → `Unmount` → confirm the last cleanups ran.

---

## Phase 4: Rules of Hooks, tests, polish (Steps 16–18)

### Step 16: Call order rule (Rules of Hooks)

Deliberately break a component: call `UseState` inside `if count > 0 { ... }`. Run several renders and observe what happens (wrong values, panic, index desync).

Then add a simple runtime guard, e.g.:

* remember the hook count from the first render,
* on later renders, if the final `hookIndex` differs — `panic("hooks order mismatch")`.

This is one of the most important takeaways of the whole project.

### Step 17: Unit tests

In `hook/runtime_test.go`, write tests (table-driven where it makes sense):

1. `UseState` keeps its value across renders.
2. The setter causes another render with the new value.
3. `UseEffect` runs after the first render.
4. `UseEffect` does not run when deps did not change.
5. Cleanup runs before the next effect and on `Unmount`.

Use counters (`int` in closures / test fields), not only `fmt.Println`.

```bash
go test ./...
go test -cover ./...
```

### Step 18: `main` as a mini demo + build

In `main.go`, assemble a readable demo:

1. Create a `Runtime`.
2. Run the `Counter` component (state + effect logs).
3. Simulate a few “clicks” (setters / scheduled updates).
4. `Unmount`.

Build the binary:

```bash
go build -o gohooks main.go
./gohooks
```

---

## “Done” checklist

The project works when:

* [ ] A function component renders through the runtime
* [ ] `UseState` keeps state across renders and the setter schedules an update
* [ ] There is a re-render loop (with a safe limit)
* [ ] `UseEffect` honors deps + cleanup
* [ ] `Unmount` cleans up effects
* [ ] Hook order violations are detected
* [ ] Tests cover state and effect
* [ ] The CLI demo behaves close to React intuition

---

## Hints (no full solutions)

* A hook must close over the **slot index** in the setter/cleanup — do not use the global `hookIndex` at the time the setter runs later.
* Generics (`UseState[T]`) are convenient, but you can start with `any` and refactor later.
* Do not build a Virtual DOM — “render” = printing state is enough.
* Do not aim for full React 18 compatibility (batching, concurrent) — this is a deliberate teaching model.
* If you hit an import cycle: `components` → `hook`, never the reverse; `main` wires both.

---

## Optional extensions (after finishing 1–18)

Short list below — the **full extension track** (channels, `UseSelect`, per-hook limit, component tree) lives in [`gohooks-advanced-en.md`](./gohooks-advanced-en.md).

* `UseRef` — mutable cell without a re-render
* `UseMemo` — cache a value while deps are unchanged
* setter batching (many `setState` = one re-render)
* `UseState` with an updater `func(prev T) T`

Do not do these until you finish the critical path 1–18.

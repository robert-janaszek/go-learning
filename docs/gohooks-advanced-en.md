# GoHooks — advanced extensions (after base steps 1–18)

This document continues [`gohooks-en.md`](./gohooks-en.md). Assume you already have a working runtime: `UseState`, `UseEffect` (deps + cleanup), `Run` / `Unmount`, hook-order protection, and basic tests.

Goal: move from a “dirty loop with a max of 50 renders” to an **event-driven** model (channels), smarter re-render scheduling, and a **component tree**.

Do not chase a full React Fiber clone — this remains a deliberate teaching engine.

---

## What changes vs the base course

| Base (1–18) | Extension goal |
|-------------|----------------|
| `Run` spins in a loop, ~50 render cap | Loop can live longer; it wakes when there is work |
| Global render budget | **Per-hook** limit: the same slot cannot schedule an update N times in a row |
| `dirty` flag + busy loop inside one `Run` call | Channel signal → re-render only after a real change |
| `Component func()` with no children | A component may return a list of children to render |

---

## Phase A: Channel-based re-render (no busy-loop)

### Step A1: “Needs re-render” signal

Instead of (or in addition to) a bare `dirty` flag, the runtime owns a channel, e.g. `updates chan struct{}` (or `chan update` carrying a slot index — see A3/A4).

The `UseState` setter, after writing a new value, must **not** assume a `for` loop on the same call stack will flush renders. It sends a signal (non-blocking or via a buffer of 1 — your design choice).

Think through:

* buffered `1` vs unbuffered channel,
* setters called from another goroutine outside `Render` (for a start you may require: only from the runtime loop / after `Render`),
* whether multiple `set`s before the next render collapse into **one** signal (early batching).

### Step A2: Event loop `Run`

`Run` is roughly:

1. First `Render` (mount).
2. Then: `for { select { case <-updates: Render; case <-ctx.Done()/quit: return } }`.
3. Shutdown: close channel / `Unmount` / cancel context — pick one clear API strategy and document it.

Important: **do not** busy-spin checking `dirty`. No signal = the goroutine blocks on the channel.

### Step A3: Re-render only when state actually changed

In the setter:

* if the new value equals the previous one (`==` for comparable types, or `DeepEqual` for learning) → **do not** signal,
* otherwise write and signal.

Add a test: `set(sameValue)` twice → one extra render (or zero), not two.

### Step A4: `UseSelect` — selector without useless renders

Add a hook shaped like:

```go
func UseSelect[S, T any](get func() S, selectFn func(S) T) T
```

or a simpler variant on your state model (select from a slot / small store).

Behavior:

* compute `selected`,
* keep the previous `selected` in the slot,
* if `selected` did **not** change → do not schedule a re-render,
* if it did → store and schedule (channel).

This teaches **subscribing to a slice of state**, not the whole world — `useSelector`-style intuition without Redux.

Demo: UI shows `count/2` or `count > 10`; increments that leave the selected value unchanged must not print another UI frame.

---

## Phase B: “10 in a row” limit (per hook), not a global 50

### Step B1: Per-slot consecutive schedule counter

Remove the global “max 50 renders” as the primary guard (optional hard-cap is fine).

Main protection:

* every slot that can schedule a re-render (`UseState`, `UseSelect`, …) has a `consecutiveSchedules` counter,
* when a slot **successfully** schedules another render, increment that slot’s counter,
* when **another** slot schedules, or after a “calm” render — reset counters (define the rule and stick to it),
* after **10** successful schedules from the **same** slot in a row → `panic` / error mentioning the hook index.

Teaching goal: catch unconditional `set(x+1)` in the component body without a crude “50 renders and stop”.

### Step B2: Storm test

Write a test that:

1. deliberately creates an infinite update from one setter,
2. expects panic / error after ~10,
3. separately: alternating updates from two slots should **not** trip one slot’s limit (or should — depending on your reset rule; document the choice).

### Step B3 (optional): schedule vs render

One channel signal may equal one render while many setters were batched. “In a row” often counts **flushes scheduled by that slot**, not every `set` line inside a batch. Decide and document.

---

## Phase C: Component returns a list of children

### Step C1: New component type

Move from `func()` to something like:

```go
type Component func() []Component
```

or:

```go
type Component func() Result
// Result holds children + optional console “output”
```

Empty / `nil` = leaf (like today’s Counter, with an explicit return).

### Step C2: Tree render

After calling a component, `Render` receives children and renders each **recursively** (or via DFS/BFS queue).

Hook rules for trees are the hard part:

* **Simpler (recommended first):** one active node at a time; on enter child, push hook context (own slots / own `hookIndex`), on exit pop. Each component instance owns storage.
* **Harder:** one global hook queue like React (order = DFS call order) — easier to get wrong, closer to the real model.

In the base course, hook order was per single component — here you need **instance identity**.

### Step C3: Child identity (key)

When a parent returns:

```text
[ChildA, ChildB] → [ChildB, ChildA]
```

without keys you corrupt state (as in React). Add something simple:

```go
type Element struct {
    Key       string
    Component Component
}
```

Parent returns `[]Element`. Runtime maps `key → hook storage`. Missing key = list index (document the limitation).

### Step C4: Mount / update / unmount children

On the parent’s next render:

* new key → mount (empty slots, then effects),
* missing key → unmount branch (effect cleanups),
* same key → re-render with preserved state.

Demo: three mini-counters; remove the middle one; assert its cleanup ran and neighbors kept state.

### Step C5: Tree tests

1. Parent with two children — each has its own `UseState`.
2. Reorder without key vs with key — observe the difference.
3. Unmounting a parent unmounts descendants.

---

## Phase D: API polish and demo

### Step D1: `main` as a small event-driven app

* start `Run` in a goroutine or `Run(ctx)`,
* from `main` / a ticker / `stdin`, send “clicks” (safe setter exposure or an event queue),
* on quit → `Unmount` and exit.

### Step D2: Setter batching

Many synchronous `set`s in one event = **one** re-render. Natural with a buffer-1 channel and an “already scheduled” flag.

### Step D3: “Done” criteria (extensions)

* [ ] Runtime loop waits on a channel; no busy-loop
* [ ] Setter with unchanged value does not wake a render
* [ ] `UseSelect` (or equivalent) skips useless renders
* [ ] ~10 consecutive schedules from the **same** hook; no global “50 and done” as the only guard
* [ ] Component may return children; children have isolated state
* [ ] Key (or documented index fallback) + branch unmount
* [ ] Tests cover channel, select, storm limit, tree

---

## What else is worth adding? (suggested order)

Not required to finish A–D, but they fit this model well:

1. **`UseRef`** — mutable cell with no channel signal (contrast to `UseState`).
2. **`UseMemo` / `UseCallback`** — cache while deps unchanged.
3. **`UseState(updater func(prev T) T)`** — safer with batching and channel events.
4. **Simple Context** (`UseProvide` / `UseContext`) — avoids prop drilling in the Phase C tree; implement with a stack on the runtime during DFS.
5. **Single-goroutine runtime rule** — document: all hooks/renders on one goroutine; other goroutines only via `r.Enqueue(fn)` on a channel. Avoids races without mutexes everywhere.
6. **`UseReducer`** — `dispatch(action)` instead of many setters; sits nicely on an event channel.
7. **Dev-only tree log** (DFS indent) — helps debug child mount/unmount.
8. **Skip on purpose:** HTML Virtual DOM, concurrent rendering, Suspense — high cost; A–D is already dense.

If you pick only one thing after A–D: **Context + single-goroutine rule** — makes the component tree practical.

---

## Hints (no full solutions)

* A setter still closes over `*Runtime` + slot index; it also uses that runtime’s update channel.
* `select` + `ctx.Done()` in `Run` teaches cancellation better than only `close(chan)`.
* For children, do “per-instance storage + key” first; only then consider one global hook queue.
* Test the “10 in a row” limit with `defer recover` — idiomatic Go.
* Do not mix in one PR: channel refactor **and** full child trees — split milestones (Phase A+B, then C).

---

## Suggested file layout

```text
gohooks/
├── hook/
│   ├── runtime.go      # Run on select/chan, Unmount
│   ├── state.go        # UseState + signal + equality
│   ├── select.go       # UseSelect
│   ├── effect.go
│   ├── tree.go         # Element, key, mount/unmount children
│   └── *_test.go
└── components/
    ├── counter.go
    └── list_demo.go    # parent returning []Element
```

Names are flexible — keep **signals / select / tree** clearly separated from effects.

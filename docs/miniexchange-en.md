Here are **20 project exercises: A mini exchange / Order Matching Engine (in Go)**.

You already know goroutines, channels, `select`, `WaitGroup`, and mutexes from Day 7, plus projects like the parser / gohooks / minivm. Here you will build a **simplified order matching engine** — an abstraction of what sits behind an exchange:

* traders producing orders in parallel,
* a single goroutine that owns the order book,
* matching bids/asks into `Trade`s,
* portfolios / balances guarded by a mutex,
* fan-out of market data to subscribers over channels,
* a clean session shutdown (`close` + `WaitGroup` + no data races).

This is not meant to be production-ready (not FIX, no persistence, no real money). The point is to **see where a real system uses channels vs mutexes** — and why an order book usually has a single owner.

---

```text
miniexchange/
├── main.go
└── exchange/
    ├── order.go         # Side, Order, Trade
    ├── book.go          # order book + Match
    ├── portfolio.go     # balances under a mutex
    ├── engine.go        # engine loop on channels
    ├── trader.go        # parallel traders
    ├── marketdata.go    # market-data fan-out
    └── exchange_test.go
```

You may tweak file names slightly, but keep the split: **model / book / portfolio / engine / traders / market data**.

---

## Domain model (before you code)

Adopt one convention and stick to it for the whole project:

```text
                    ┌─────────────┐
   trader A ──┐     │   orders    │      ┌──────────────┐
   trader B ──┼────►│   chan      │─────►│    Engine    │
   trader C ──┘     └─────────────┘      │  (1 goroutine)│
                                         │  owns Book   │
                                         └──────┬───────┘
                                                │ trades chan
                        ┌───────────────────────┼───────────────────────┐
                        ▼                       ▼                       ▼
                 Portfolio.Apply            MarketData              (tests / logs)
                 (mutex on maps)            fan-out to N subs
```

* **One symbol** to start (e.g. `"GOOG"`) — do not build a multi-asset router.
* **Prices and quantities:** `int` (e.g. price in cents). Avoid `float64` in matching.
* **Order book:** bids sorted **descending** by price, asks **ascending**. On price ties — FIFO (arrival order).
* **Matching:** best bid ≥ best ask → emit a `Trade` at the **resting** order’s price (maker) — or always at the ask; **pick one rule and document it**.
* **Book owner:** only `Engine` mutates the book. Traders **never** touch `Book` directly.
* **Mutex:** protects portfolios (trader → cash/position maps), not the book. The book is single-writer by design.

Do not mix responsibilities: channel = event flow; mutex = shared state read/written from multiple goroutines.

---

## Phase 1: Orders and trades (Steps 1–4)

### Step 1: Directory and package

Create a `miniexchange` directory (next to `minivm` / `gohooks`) in the same module:

```bash
mkdir -p miniexchange/exchange
```

Implementation package: `exchange`. `main.go` only wires a demo session.

### Step 2: `Side`, `Order`, `Trade`

In `exchange/order.go`:

```go
type Side int

const (
    SideBuy Side = iota
    SideSell
)

type Order struct {
    ID       int64
    TraderID string
    Side     Side
    Symbol   string
    Price    int // e.g. cents
    Qty      int // remaining quantity (shrinks on partial fill)
}

type Trade struct {
    Symbol     string
    Price      int
    Qty        int
    BuyID      int64
    SellID     int64
    BuyTrader  string
    SellTrader string
}
```

Add `String() string` for `Side` (handy in logs).

### Step 3: Order validation

Write:

```go
func ValidateOrder(o Order) error
```

Reject:

* empty `TraderID` / `Symbol`,
* `Price <= 0`, `Qty <= 0`,
* unknown `Side`.

Return sentinel errors (`var ErrInvalidOrder = errors.New(...)`) and assert with `errors.Is`.

### Step 4: Model tests

In `exchange_test.go`:

1. A valid order passes validation.
2. `Qty == 0` → error.
3. `String()` for Buy/Sell is readable.

```bash
go test ./miniexchange/...
```

---

## Phase 2: Order book and matching (Steps 5–9)

### Step 5: `Book` structure

In `exchange/book.go`:

```go
type Book struct {
    bids []Order // descending by Price, then FIFO
    asks []Order // ascending by Price, then FIFO
}
```

A sorted slice plus insert-at-position is fine (binary search optional — linear insert is OK for learning).

Do not export book mutations beyond the package unless needed — methods on `*Book` are enough.

### Step 6: `Add` without immediate matching

Implement internal insertion onto the correct side, preserving price order + FIFO on ties.

Do **not** match yet — test ordering first.

### Step 7: `Match` / `Accept`

Write a method that accepts a new order and:

1. Tries to match it against the opposite side while possible.
2. Produces zero or more `Trade`s (partial fills allowed).
3. Decrements `Qty` on orders; removes fully filled ones from the book.
4. Resting remainder of the aggressive order (if `Qty > 0`) joins the book as maker.

Example:

* book has Sell `@100 qty=5`
* Buy `@100 qty=3` arrives → Trade `qty=3`, ask shrinks to `2`
* Buy `@100 qty=10` arrives → Trade `qty=2`, remaining Buy `qty=8` rests as a bid

### Step 8: Top-of-book reads

Add:

```go
func (b *Book) BestBid() (Order, bool)
func (b *Book) BestAsk() (Order, bool)
func (b *Book) Depth() (bids, asks int)
```

Useful in tests and market data.

### Step 9: Matching tests (table-driven)

Required scenarios:

1. No opposite side → order rests, no trades.
2. Full 1:1 match.
3. Partial fill (aggressive larger than resting).
4. One aggressive order sweeps multiple price levels.
5. Worse price does **not** match (Buy `@90` vs Ask `@100`).
6. FIFO at the same price: older order fills first.

This is the heart of the project — do not move on without green tests.

---

## Phase 3: Channel-driven engine (Steps 10–13)

### Step 10: `Engine` and channels

In `exchange/engine.go`:

```go
type Engine struct {
    orders chan Order
    trades chan Trade
    // book, portfolio, nextID, wg / done — your design
}
```

Constructor example:

```go
func NewEngine(orderBuf, tradeBuf int) *Engine
```

* `orders` — input (traders send),
* `trades` — output (each matched trade),
* buffers > 0 so the demo does not deadlock on bursts; tests may use small buffers on purpose.

### Step 11: `Run` loop

```go
func (e *Engine) Run()
```

Runs in **one goroutine**:

```go
for o := range e.orders {
    // assign ID if needed
    // ValidateOrder
    // trades := book.Accept(o)
    // send each trade on e.trades
    // update portfolio
}
close(e.trades) // when orders closes and the loop exits
```

Rule: **only this goroutine** calls mutating `Book` methods.

### Step 12: Submit and stop

Add API:

```go
func (e *Engine) Submit(o Order) error  // blocking send or non-blocking — pick one and document
func (e *Engine) Close()                // close(orders); wait until Run finishes
```

`Close` should be idempotent or clearly one-shot. After `Close`, further `Submit` → error (e.g. `ErrSessionClosed`).

Stop pattern: `close(orders)` → `Run` finishes matching → `close(trades)` → market-data consumers end their `range`.

### Step 13: Engine test (race-free)

Integration test:

1. `go e.Run()`
2. Submit 2 complementary orders from the test.
3. Receive exactly 1 `Trade` with expected price/qty.
4. `Close()`, assert `trades` closes (`range` ends).

```bash
go test -race ./miniexchange/...
```

---

## Phase 4: Portfolios (mutex) and traders (parallelism) (Steps 14–16)

### Step 14: `Portfolio` under `sync.Mutex`

In `exchange/portfolio.go`:

```go
type Portfolio struct {
    mu       sync.Mutex
    cash     map[string]int // TraderID → cash
    position map[string]int // TraderID → shares
}
```

API:

* `Seed(traderID string, cash, position int)`
* `Apply(trade Trade) error` — buyer loses `price*qty` cash, gains position; seller the opposite
* `Snapshot(traderID string) (cash, position int)` — read under `RLock` if you use `RWMutex`

Reject trades that would make cash negative / create a short (or allow shorts — **decide and test**). Educational recommendation: **reject insufficient funds at Accept/Apply** (simpler invariant).

Design note: if you reject *after* matching, book and portfolio can diverge. Cleaner: check balances before finalizing a trade, or seed generously in the demo and treat `Apply` inconsistency as a test failure. Document the chosen strategy above `Apply`.

### Step 15: Parallel traders

In `exchange/trader.go`:

```go
func RunTrader(e *Engine, id string, orders []Order, wg *sync.WaitGroup)
```

Each trader:

1. `defer wg.Done()`
2. Sends its order list via `e.Submit` (optional small `time.Sleep` / jitter).
3. Does not touch the book or portfolio maps directly.

In a test / demo, start 3–5 traders (`go RunTrader(...)`), `wg.Wait()`, then `e.Close()`.

### Step 16: Portfolio race test

Write a test where portfolios are observed concurrently — e.g. engine calls `Apply` from one goroutine (fine) while parallel `Snapshot`s run during the session.

Run:

```bash
go test -race ./miniexchange/...
```

Goal: zero race reports. If `-race` flags `Book`, you broke the single-owner rule.

---

## Phase 5: Market data (fan-out) and wrap-up (Steps 17–20)

### Step 17: Quote broadcaster

In `exchange/marketdata.go` implement fan-out:

* input: `<-chan Trade` (from the engine),
* subscriber registration: each gets its own buffered `chan Trade`,
* a goroutine reads trades and fans them out to all subs,
* when input closes — close subscriber channels.

```go
type MarketData struct { /* ... */ }

func NewMarketData(in <-chan Trade) *MarketData
func (m *MarketData) Subscribe(buf int) <-chan Trade
func (m *MarketData) Run()
```

Note: a slow subscriber must not block the exchange forever. Baseline: buffered channel + drop/skip when full (`select` + `default`) — **document** the policy you choose.

### Step 18: Session-end signal (`chan struct{}`)

Add a signal channel (0-byte payload):

```go
done := make(chan struct{})
```

After `Close()` / after `Run` exits, close `done` (or send `struct{}{}`) so `main` and subscribers can `select` on session end. This is Day 7 exercise 20 in a realistic setting.

### Step 19: Demo in `main.go`

Wire a readable session (not spaghetti):

1. Seed 3 trader portfolios.
2. Start `Engine.Run` + `MarketData.Run` in goroutines.
3. A market-data subscriber prints trades.
4. Launch traders with orders that **definitely** produce several trades and leave something resting in the book.
5. `wg.Wait` → `Engine.Close` → wait on `done` / market-data shutdown.
6. Print portfolio snapshots + book `Depth()` (you may add `Engine.BookSnapshot()` for debug — watch races: call only after stop, or copy under engine control).

```bash
go run ./miniexchange
go test -race ./miniexchange/...
```

### Step 20: Quality bar and coverage

Add / tighten tests:

1. Table-driven matching (Phase 2).
2. Engine: submit → trade → close.
3. Portfolio: Apply updates cash/position consistently vs Snapshot.
4. MarketData: 2 subscribers receive the same trade.
5. Full suite under `-race`.

```bash
go test -race -cover ./miniexchange/...
```

---

## “Done” checklist

You have a working project when:

* [ ] `Order` / `Trade` / validation work and have tests
* [ ] `Book` sorts bid/ask and matches with partial fill + FIFO
* [ ] `Engine` reads `orders`, writes `trades`, and is the **only** book mutator
* [ ] `Close` ends the session without deadlock (`close` + end of `range`)
* [ ] `Portfolio` is mutex-guarded; balances stay consistent after trades
* [ ] Multiple traders submit in parallel (`WaitGroup`)
* [ ] Market data fans trades out to ≥2 subscribers
* [ ] `go test -race ./miniexchange/...` is clean
* [ ] CLI demo shows a session: orders → trades → snapshots

---

## Tips (no spoilers)

* **Book = single writer.** If you reach for a mutex on the book “just in case”, first check whether one goroutine is enough — that is the main architecture lesson.
* Keep the mutex close to the maps (`Portfolio`), not “around the whole Engine”.
* Only the sender closes a channel. Traders do not close `orders` — `Engine.Close` / `main` does after `WaitGroup`.
* Never send on a closed channel — hence `ErrSessionClosed` or `sync.Once` around `close`.
* `-race` is part of the Definition of Done, not optional.
* Do not build UI, REST, WebSockets, persistence, auth, or margin calls — skip on purpose.
* Prefer `int` prices over `float64`. Float matching is its own nightmare.

---

## Optional extensions (after finishing 1–20)

* order cancel (`Cancel` as another event type on the same channel — `interface{}` / sum type)
* multiple symbols (`symbol → *Book` map, still inside one engine goroutine)
* `context.Context` to cancel a session instead of only `Close`
* per-trader rate limit (`time.Ticker`)
* extra order types (`Limit` vs `Market`)
* periodic order-book snapshots on a separate channel (market depth)
* trade journal to a JSONL file
* hedged “best price” across two books (Day 7 first-responder practice)

Do not start these until the critical path 1–20 is done.

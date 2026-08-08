Here are **20 project exercises: A mini-VM with virtual RAM (in Go)**.

You already know pointers from Day 2, escape analysis (“stack vs heap in Go”), and projects like the parser / gohooks. Here you will build an **educational virtual machine with its own memory**:

* a single byte buffer (virtual RAM),
* pointers as *offsets* into that memory (not native Go `*T`),
* a downward-growing stack + stack pointer (SP),
* a heap with `Alloc` / `Free`,
* call frames (`CALL` / `RET`),
* a small instruction set and an execute loop.

This is not meant to be production-ready (not WASM, not Go’s real GC). The point is to **see memory layout with your own eyes** — ordinary Go hides it behind the compiler.

---

```text
minivm/
├── main.go
└── vm/
    ├── memory.go      # RAM, Addr, word Load/Store
    ├── stack.go       # SP, Push/Pop
    ├── heap.go        # Alloc / Free
    ├── frame.go       # call frames (optional separate file)
    ├── opcode.go      # opcodes
    ├── machine.go     # VM: registers, Execute
    └── vm_test.go
```

You may tweak file names slightly, but keep the split: **memory / stack / heap / instructions / machine**.

---

## Memory model (before you code)

Adopt one convention and stick to it for the whole project:

```text
low address                                               high address
0 ─────────────────────────────────────────────────────► Size-1
[  HEAP grows up → ........ free ........ ← STACK       ]
                     ↑                      ↑
                  HeapBrk                  SP
```

* **Machine word:** `uint32` (4 bytes) or `uint64` (8 bytes) — pick one and use it consistently. Below we assume **`uint32` / 4 bytes** (easier dumps).
* **Endianness:** little-endian when reading/writing words in `[]byte`.
* **`Addr`:** `type Addr uint32` — your “pointer”. `0` = null (must not be dereferenced).
* **Stack:** grows downward (SP decreases on `Push`).
* **Heap:** grows upward from low addresses (after a header / code region — decide in Phase 1).
* **Collision:** if `HeapBrk` meets `SP` → error (stack overflow / out of memory).

Do not mix native Go pointers (`*int`) with VM addresses — the VM API talks in `Addr` and `[]byte`.

---

## Phase 1: Virtual RAM and pointers (Steps 1–5)

### Step 1: Module layout

Create a `minivm` directory (next to `gohooks` / `taskmanager`) in the same module or as a subpackage:

```bash
mkdir -p minivm/vm
```

Implementation package: `vm`. `main.go` only wires up the demo.

### Step 2: `Addr` type and memory

In `vm/memory.go`:

```go
type Addr uint32

const WordSize = 4 // bytes; adjust if you choose uint64

type Memory struct {
    data []byte
}
```

Implement:

* `NewMemory(size int) *Memory` — allocates `make([]byte, size)` (already zeroed).
* `Size() int`
* alignment checks: word operations require addresses divisible by `WordSize`.

### Step 3: Word `Load` / `Store`

Implement:

```go
func (m *Memory) Load(addr Addr) (uint32, error)
func (m *Memory) Store(addr Addr, value uint32) error
```

Rules:

* reject `addr == 0` (null dereference),
* reject out-of-range addresses (`addr+WordSize > size`),
* reject misaligned addresses,
* use `encoding/binary.LittleEndian` (or manual byte packing).

Byte helpers (`Load8` / `Store8`) are optional for later.

### Step 4: Memory dump (debug)

Write `func (m *Memory) Dump(start, end Addr) string` (or `WriteDump(w io.Writer, ...)`) that prints hex in readable lines (e.g. 16 bytes per line + address).

You will need this in every later step — without a dump you will be guessing.

### Step 5: Memory tests

In `vm/vm_test.go` (or `memory_test.go`):

1. `Store` + `Load` round-trips the same value.
2. `Load(0)` → error.
3. Out-of-range address → error.
4. Misaligned address (e.g. `1`) → error.
5. Little-endian: after `Store(addr, 0x01020304)` bytes are `04 03 02 01`.

```bash
go test ./minivm/...
```

---

## Phase 2: Stack and stack pointer (Steps 6–9)

### Step 6: Machine with SP

In `vm/machine.go` (or `stack.go`) define:

```go
type VM struct {
    mem  *Memory
    sp   Addr   // stack pointer
    // later: hp (heap break), ip (instruction pointer), fp (frame pointer)...
}
```

Startup convention:

* `mem` size e.g. `4096` bytes,
* initial `sp = Addr(size)` — “just past the end” / at the end of the region (pick one meaning: whether SP points at free space or the last occupied slot — **document it in a comment** and stay consistent).

Recommendation (classic): **SP points at the top value (last pushed)**; empty stack → `sp == Addr(size)` (sentinel past the region). Alternatively: SP points at the first free word. Pick one model.

### Step 7: `Push` / `Pop`

```go
func (v *VM) Push(value uint32) error
func (v *VM) Pop() (uint32, error)
```

`Push`:

1. Compute the new top address (`sp - WordSize`).
2. Check collision with the heap / lower bound (for now lower bound can be `0` or a fixed `heapStart`).
3. `Store` the value.
4. Update `sp`.

`Pop`: reverse; error on empty stack.

### Step 8: Peek and depth

Add:

* `Peek() (uint32, error)` — read top without popping,
* `StackDepth() int` — how many words are on the stack.

### Step 9: Stack tests

1. Push 1, 2, 3 → Pop → 3, 2, 1 (LIFO).
2. Pop on empty stack → error.
3. Push until nearly full → next Push → error (overflow).
4. Dump shows values at high addresses.

---

## Phase 3: VM pointers — indirect Load/Store (Steps 10–11)

### Step 10: Assembler-style ops before a full ISA

Before a full decoder, add VM operations:

```go
func (v *VM) LoadIndirect() error  // Pop addr; Push mem[addr]
func (v *VM) StoreIndirect() error // Pop value; Pop addr; mem[addr] = value
```

(or with explicit arguments — what matters is that **the address is a value on the stack**, not a native Go pointer).

Manual test scenario:

1. Pick a data-region address, e.g. `0x100`.
2. `Store(0x100, 42)` directly.
3. `Push(0x100)` → `LoadIndirect` → stack top is `42`.

### Step 11: Pointer to pointer (double indirection)

Demonstrate in a test / demo:

1. Address `A` holds value `7`.
2. Address `B` holds value `A` (a pointer).
3. Sequence: Push `B` → LoadIndirect → LoadIndirect yields `7`.

This is the Day 2 `**int` idea, but inside your RAM.

---

## Phase 4: Heap — `Alloc` / `Free` (Steps 12–15)

### Step 12: Bump allocator (simplest heap)

Add to `VM` (or a separate `Heap` struct):

* `heapStart Addr` — fixed heap base (e.g. `0x40`, leave room for “static” data),
* `heapBrk Addr` — current end of used heap (break).

```go
func (v *VM) Alloc(nbytes uint32) (Addr, error)
```

Bump behavior:

1. Round `nbytes` up to a multiple of `WordSize`.
2. Optionally write a block header (size) — **recommended**, it makes `Free` easier later.
3. Return the **payload** address (after the header).
4. Advance `heapBrk`.
5. If `heapBrk` gets too close to `sp` → `ErrOutOfMemory`.

At this stage `Free` may be a no-op or not exist yet.

### Step 13: Block header

Pick a layout, e.g.:

```text
[ size:u32 ][ payload... ]
            ^
            returned Addr
```

Or with a free flag:

```text
[ size:u32 ][ free:u32 ][ payload... ]
```

Document the layout in a comment above `Alloc`. Test: after two `Alloc(8)` calls, addresses differ by aligned size + header.

### Step 14: Free-list (real `Free`)

Extend the allocator:

```go
func (v *VM) Free(ptr Addr) error
```

Minimum requirements:

* `Free` marks the block free and inserts it into a **free-list** — a `next` pointer can live in the payload / a header field,
* the next `Alloc` searches the free-list first (first-fit is enough), then bumps,
* `Free(0)` → error,
* double `Free` of the same block → error (detect via a `free` flag),
* `Free` of a pointer that is not a block payload → error.

You do not need coalescing (merging neighbors) in the base version — that is an optional extension.

### Step 15: Heap tests

1. `Alloc` + `Store`/`Load` through the returned address works.
2. `Free` + another `Alloc` of the same size **reuses** an address (or otherwise proves you are not only bumping forever in an alloc/free loop).
3. Alloc until stack collision → OOM.
4. Double-free → error.
5. Stack vs heap: push a lot onto the stack, then `Alloc` must respect SP.

---

## Phase 5: Call frames and ISA (Steps 16–18)

### Step 16: Frame pointer and `CALL` / `RET`

Add register `fp Addr` (frame pointer) and `ip int` (index into the instruction slice — see Step 17).

Frame convention (simplified — you may simplify further):

```text
high addresses
  ... previous frame ...
  [ arg N ]
  [ arg 1 ]
  [ return addr ]  ← after CALL
  [ saved FP    ]  ← FP points here (or at saved FP)
  [ local 0     ]
  [ local 1     ]
low addresses (stack growth direction)
```

Implement methods (called from Go for now, not from opcodes):

* `Call(returnAddr uint32, nargs int) error` — saves return addr + old FP, sets new FP, optionally reserves local slots,
* `Ret() (returnAddr uint32, error)` — restores FP and SP, returns the return address.

Test 2–3 nested “calls” and check that locals do not overwrite each other.

### Step 17: Opcodes and program

In `opcode.go`:

```go
type Op byte

const (
    OpHalt Op = iota
    OpPush      // operand: u32 literal
    OpPop
    OpAdd       // pop a, pop b, push b+a
    OpLoad      // pop addr, push mem[addr]
    OpStore     // pop val, pop addr, mem[addr]=val
    OpAlloc     // pop nbytes, push ptr
    OpFree      // pop ptr
    OpCall      // operand: target ip; args on stack per your convention
    OpRet
    OpDup       // duplicate top (handy)
    OpPrint     // pop and log — demo only
)
```

```go
type Instr struct {
    Op   Op
    Arg  uint32 // used when the instruction has an operand
}
```

`VM` holds `code []Instr` and `ip int`.

### Step 18: `Execute` loop

```go
func (v *VM) Execute(code []Instr) error
```

Loop:

1. `for { instr := code[ip]; ip++; switch instr.Op ... }`
2. `OpHalt` → return nil.
3. Propagate errors from `Push`/`Pop`/`Alloc` (include `ip` in the message — very helpful).
4. Instruction limit (e.g. 100_000) against infinite loops — even before you have conditional jumps.

Implement first: `Push`, `Pop`, `Add`, `Load`, `Store`, `Halt`. Then wire `Alloc`/`Free`/`Call`/`Ret`.

---

## Phase 6: Demo, safety, wrap-up (Steps 19–20)

### Step 19: Demo program in `main`

Build a readable demo in `main.go` (not one spaghetti dump):

**Program A — stack arithmetic**

* Push 40, Push 2, Add, Print → `42`.

**Program B — heap**

* Alloc 4 bytes, Store `7` at the returned pointer, Load it back, Print, Free.

**Program C — “function” via CALL/RET**

* A tiny procedure: take 1 stack argument, add `1`, return via Push result + `Ret`.
* Main: Push arg, Call, Print.

After each program, print a short dump (SP, FP, heapBrk + a memory slice).

```bash
go run ./minivm
```

### Step 20: Integration tests + quality bar

Add table-driven tests for `Execute`:

1. Program `40 + 2` → top / Print path = 42 (collect output via a callback instead of `fmt` if you prefer).
2. Alloc/Store/Load/Free in bytecode.
3. Nested Call/Ret.
4. Stack overflow from Push in a loop.
5. OOM from Alloc.
6. Null load (`Push 0`, `OpLoad`) → error.

Use `t.Helper()` and clear messages with `ip` / `sp`.

```bash
go test ./minivm/...
go test -cover ./minivm/...
```

---

## “Done” checklist

You have a working project when:

* [ ] Unified RAM (`[]byte`) with word `Load`/`Store` and an `Addr` type
* [ ] Null, OOB, and misaligned access return errors
* [ ] LIFO stack with downward SP + overflow protection
* [ ] `Alloc` (bump) + `Free` (free-list first-fit) sharing memory with the stack
* [ ] Heap↔stack collision is detected
* [ ] Double-free / free(null) are rejected
* [ ] `CALL`/`RET` maintain frames (FP) without clobbering outer locals
* [ ] `Execute` runs bytecode with `Halt`, arithmetic, load/store, alloc/free
* [ ] Unit + integration tests pass
* [ ] CLI demo shows A/B/C and a memory dump

---

## Hints (no solution dump)

* **One truth about SP:** write a 5-line comment on the `sp` field and do not flip semantics mid-project.
* Keep package-level sentinel errors: `var ErrNullDeref = errors.New(...)`, `ErrStackOverflow`, `ErrOutOfMemory`, `ErrDoubleFree` — assert with `errors.Is`.
* Do not `panic` on “bad user program” paths — reserve panic for VM bugs (inconsistent internal state).
* Hex dumps beat ad-hoc `fmt.Printf` when debugging the allocator.
* Do not implement GC, paging, an MMU, or threads — intentionally out of scope.
* Native Go `*T` may hold the VM *struct*; **program data** lives only in `Memory.data`.

---

## Optional extensions (after finishing 1–20)

* coalesce free blocks (merge neighbors on `Free`)
* `OpJmp` / `OpJmpIfZero` + a simple bytecode loop
* separate **code** vs **data** regions (today code is a Go `[]Instr` — you can move immediates into RAM)
* 2–3 general-purpose registers beside the stack (accumulator)
* serialize memory dumps as golden-file tests
* use-after-free detector: fill freed payload with `0xDE` and a generation flag

Do not start these until the critical path 1–20 is done.

Here are **20 exercises for Day 8: Week 1 Wrap-up – Practical CLI Task Manager Project**.

Congratulations on finishing the first week! You already know syntax, pointers, structs, interfaces, error handling, packages, and the basics of concurrency.

Today there are no isolated drills. Instead you will **build a complete, production-minded CLI app (Task Manager) from scratch**. The work is split into 20 architectural steps — each step adds one concrete piece of the system.

---

```text
taskmanager/
├── go.mod
├── main.go
├── task/
│   ├── task.go        # Structs and business logic
│   └── task_test.go   # Unit tests
└── storage/
    ├── storage.go     # Data store interface
    └── json.go        # JSON file-based implementation

```

---

## Phase 1: Core Data Model and Repository (Steps 1–5)

### Step 1: Module init and file layout

Create a `taskmanager` directory and initialize the module (`go mod init taskmanager`). Create the `task/` and `storage/` folders.

### Step 2: Task model (`task/task.go`)

In the `task` package, define a `Task` struct with fields:

* `ID int`
* `Title string`
* `Done bool`
* `CreatedAt time.Time`

Add JSON tags to every field (e.g. `json:"id"`).

### Step 3: Validation and constructor

Write `NewTask(id int, title string) (*Task, error)`. If `title` is empty, return `errors.New("title cannot be empty")`. Otherwise set `CreatedAt` to `time.Now()`.

### Step 4: Storage interface (`storage/storage.go`)

In the `storage` package, define a `Storage` interface:

```go
type Storage interface {
    Save(tasks []task.Task) error
    Load() ([]task.Task, error)
}

```

### Step 5: JSON implementation (`storage/json.go`)

Create a `JSONStorage` struct with a `filename string` field. Implement `Save` and `Load` using the standard `os` and `encoding/json` packages (`json.MarshalIndent` / `json.Unmarshal`).

---

## Phase 2: Task Manager Logic (Steps 6–10)

### Step 6: `TaskManager` struct

In the `task` package, define a `TaskManager` struct that holds:

* `tasks []Task`
* `storage storage.Storage` (depend on the interface, not a concrete type!)

### Step 7: `NewManager` constructor

Write `NewManager(s storage.Storage) (*TaskManager, error)` that calls `s.Load()` on startup and fills the internal `tasks` slice.

### Step 8: `Add` method

Add `Add(title string) error` on `TaskManager`. Generate a unique ID (e.g. `len(tasks) + 1`), create the task, append it to the slice, and call `t.storage.Save(t.tasks)`.

### Step 9: `MarkDone` method

Add `MarkDone(id int) error`. Iterate tasks, find the matching ID, set `Done` to `true`. If the task does not exist, return a custom sentinel error `ErrTaskNotFound`. Persist state to the file.

### Step 10: `List` method

Add `List(showAll bool) []Task`. If `showAll` is `false`, return only incomplete tasks (`Done == false`).

---

## Phase 3: CLI Parser and User Interface (Steps 11–15)

### Step 11: Using `flag` (standard library)

In `main.go`, use the standard `flag` package to parse CLI arguments:

* `-add <title>` (add a task)
* `-done <id>` (mark as done)
* `-list` (list tasks)
* `-all` (modifier for `-list`)

### Step 12: Wire dependencies in `main.go`

Connect everything in `main()`:

1. Create a `JSONStorage` instance pointing at `tasks.json`.
2. Create a `TaskManager` and pass in the storage.
3. Handle startup errors with `log.Fatalf`.

### Step 13: Command handling in `main.go`

Write control flow from the parsed flags (`switch` or `if`). Examples:

```bash
go run main.go -add "Buy milk"
go run main.go -list
go run main.go -done 1

```

### Step 14: Nice console formatting (`tabwriter`)

For listing tasks, use the standard `text/tabwriter` package so columns (ID, Status, Title, Date) align automatically in the terminal.

### Step 15: Graceful shutdown and save on exit

Use `os/signal` with `syscall.SIGINT` / `syscall.SIGTERM` so that on `Ctrl+C` the program can finish any pending file write safely.

---

## Phase 4: Unit Tests and Concurrency (Steps 16–20)

### Step 16: Your first Go test (`task/task_test.go`)

Create `task/task_test.go`. Write `TestNewTask(t *testing.T)` that checks an empty `title` returns an error, and a valid title builds a struct with a non-zero date. Run:

```bash
go test ./...

```

### Step 17: Table-driven tests (idiomatic Go testing)

Rewrite `TestNewTask` as **table-driven tests**:

```go
tests := []struct {
    name    string
    title   string
    wantErr bool
}{
    {"valid title", "Buy milk", false},
    {"empty title", "", true},
}

```

### Step 18: Mocking storage in tests

In the test file, write a simple `MockStorage` that implements `storage.Storage` in memory (no disk I/O). Use it to test `TaskManager.Add()`.

### Step 19: Async notification (goroutine + channel)

Add a `TaskManager` helper that, when a task is marked done, asynchronously (in a goroutine) sends a message on a `chan string` (e.g. `"Task #1 completed!"`). Print the received message in the console.

### Step 20: Build the binary (`go build`) and coverage

1. Run tests with coverage: `go test -cover ./...`.
2. Build a production binary:
```bash
go build -o taskmgr main.go

```

3. Try the binary in the terminal: `./taskmgr -list`.

---

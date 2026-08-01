Here are **exercises for Day 1: Go tooling (CLI)**.

Before syntax and concurrency, get comfortable with the basic toolchain. In Go most work starts with a few terminal commands — counterparts to what you know in JS/TS as `node`, a bundler, or Prettier.

Goal for today: run the project, build a binary, format code, and run a basic checker.

---

## Part 1: Running and building (Exercises 1–4)

### Exercise 1: Run in memory

In the project directory (where `go.mod` and `main.go` live), run:

```bash
go run .
```

or:

```bash
go run main.go
```

See what the program prints. This is compile-and-run in one step — no lasting binary on disk (handy while learning).

### Exercise 2: Compile to a single binary

Build the program:

```bash
go build
```

or with an output name:

```bash
go build -o go-learning
```

Run the resulting file (e.g. `./go-learning`). Compare with `go run`: `build` leaves an artifact on disk.

### Exercise 3: Formatting (`go fmt`)

Intentionally mess up indentation or spacing in some `.go` file, then:

```bash
go fmt ./...
```

This is Go’s Prettier equivalent — one official style for the whole ecosystem.

### Exercise 4: Basic static checks (`go vet`)

Run:

```bash
go vet ./...
```

`vet` catches common mistakes (e.g. suspicious `Printf`, ineffective `Lock`). It does not replace a fuller linter (like `staticcheck`), but it is a solid first CI step.

---

## Part 2: Module and tests (Exercises 5–8)

### Exercise 5: Inspect `go.mod`

Open `go.mod`. Note the module path and language version (`go 1.xx`). This is the project’s dependency manifest.

### Exercise 6: List packages

```bash
go list ./...
```

See which packages Go finds in the repo (`course`, `json-parser`, …).

### Exercise 7: Tests

```bash
go test ./...
```

Even if you are not writing tests today, know how to run the full suite.

### Exercise 8: Built-in help

```bash
go help
go help build
go help test
```

Skim the short descriptions — docs ship with the tool.

---

## Part 3: Learning workflow (Exercises 9–10)

### Exercise 9: Scratch `main.go`

In `main.go`, uncomment / wire a `course` exercise call (from Day 2 onward). Run again with `go run .`.

### Exercise 10: Day checklist

Make sure you can do the following without notes:

1. run code (`go run`),
2. build a binary (`go build`),
3. format (`go fmt ./...`),
4. check with `go vet ./...`.

From Day 2 you enter the language itself: variables, types, and pointers.

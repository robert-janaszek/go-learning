# go-learning

Personal playground for learning Go: a short exercise course plus two small JSON-related packages built from scratch.

## Layout

| Path | What it is |
|------|------------|
| [`course/`](course/) | Day-by-day exercises (tooling → basics → concurrency, …). Notes in [`course/docs/`](course/docs/) (EN + PL). |
| [`course/docs/minivm-pl.md`](course/docs/minivm-pl.md) / [`minivm-en.md`](course/docs/minivm-en.md) | Project brief: mini-VM with virtual RAM (stack, heap, `Alloc`/`Free`, bytecode). |
| [`course/docs/miniexchange-pl.md`](course/docs/miniexchange-pl.md) / [`miniexchange-en.md`](course/docs/miniexchange-en.md) | Project brief: mini exchange / order matching engine (channels, parallelism, mutexes). |
| [`json-parser/`](json-parser/) | Recursive-descent JSON parser → `any` (`Parse`). |
| [`json-fixer/`](json-fixer/) | Best-effort closer for truncated JSON-like input (`Fix`). |
| [`main.go`](main.go) | Scratch entrypoint — uncomment a `course.Day…` or call a package. |

There is also a tiny [`bank/`](bank/) experiment from earlier exercises.

## Run

```bash
go run .
go test ./...
```

Module: `github.com/robert-janaszek/go-learning`

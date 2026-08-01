# go-learning

Personal playground for learning Go: a short exercise course plus two small JSON-related packages built from scratch.

## Layout

| Path | What it is |
|------|------------|
| [`course/`](course/) | Day-by-day exercises (tooling → basics → concurrency, …). Notes in [`course/docs/`](course/docs/) (EN + PL). |
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

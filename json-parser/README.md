# jsonparser

A small learning JSON parser written in Go. It turns a JSON string into Go values via recursive descent — no `encoding/json`, no struct unmarshaling.

```go
v, err := jsonparser.Parse(`{"a":[1,true,null]}`)
// map[string]any{"a": []any{float64(1), true, nil}}
```

**Public API:** `Parse(input string) (any, error)`

| JSON     | Go            |
|----------|---------------|
| object   | `map[string]any` |
| array    | `[]any`       |
| string   | `string`      |
| number   | `float64`     |
| true/false | `bool`      |
| null     | `nil`         |

## How it works

1. **Lexer** — bytes → tokens (`{`, strings, numbers, `true` / `false` / `null`, …), with start offsets on tokens.
2. **Parser** — one function per grammar rule (`parseValue`, `parseObject`, `parseArray`); numbers use a small state machine for RFC-style validation.

```bash
go test ./json-parser/
```

## Not in scope (yet)

- **Strings** — only basic escapes (`\"`, `\\`). Missing full JSON escapes (`\n`, `\t`, `\uXXXX`) and rejection of illegal control characters / bad escapes.
- **Out of scope for this module** — unmarshal into structs (`reflect`), and tolerant / repairing input (see `json-fixer`).

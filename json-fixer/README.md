# jsonfixer

Closes truncated JSON-like input by balancing delimiters: `"`, `{`/`}`, `[`/`]`, and a dangling `\`.

```go
out, err := jsonfixer.Fix(`{"a`)
// out == `{"a"}`
```

## What works

- Missing closing braces and brackets
- Unclosed strings
- Escaped quotes (`\"` does not end a string)
- Trailing backslash (completed as `\\` before closing `"`)
- Rejects clear mismatches (`}`, `{]`, `{[}`)

## What it does not do

- Guarantee valid JSON (e.g. `{{` → `{{}}`)
- Fix truncated literals (`tru`, `nul`, `12.`)
- Fix trailing commas (`[1,`)
- Full escape validation (`\uXXXX`, etc.)

## Test

```bash
go test ./json-fixer/
```

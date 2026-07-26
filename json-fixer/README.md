# jsonfixer

Closes truncated JSON-like input by balancing delimiters: `"`, `{`/`}`, `[`/`]`, and a dangling `\`.
Also strips trailing whitespace/commas and completes partial keywords and numbers outside strings.

```go
out, err := jsonfixer.Fix(`{"a":1e+`)
// out == `{"a":1e+0}`
```

## What works

- Missing closing braces and brackets
- Unclosed strings
- Escaped quotes (`\"` does not end a string)
- Trailing backslash (completed as `\\` before closing `"`)
- Trailing whitespace and comma outside strings (`[1, ` → `[1]`)
- Incomplete `true` / `false` / `null` (`{"a":tru` → `{"a":true}`)
- Truncated numbers: `1.` → `1.0`, `2e` / `1E+` / `1e-` → `…0`, lone `-` → `-0`
  (exponent/`e` only when preceded by a digit, so `true`/`false` stay intact)
- Rejects clear mismatches (`}`, `{]`, `{[}`)

## What it does not do

- Guarantee valid JSON (e.g. `{{` → `{{}}`)
- Fix commas before an existing closer (`{"a":1,}`)
- Full escape validation (`\uXXXX`, etc.)

## Test

```bash
go test ./json-fixer/
```

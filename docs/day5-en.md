Here are **20 exercises for Day 5: Error Handling (`if err != nil`), the `errors` Package, and No Exceptions**.

Your goal today: unlearn `try/catch/throw` thinking, master the idioms of treating errors as values (*errors are values*), and learn the modern error wrapping pattern introduced in newer Go versions.

---

## Part 1: Basics and Returning Errors (Exercises 1–5)

### Exercise 1: Your first error (`errors.New`)

Write a function `ValidateAge(age int) error`. If age is less than 0 or greater than 120, return an error created with `errors.New("invalid age")`. Otherwise return `nil`. Handle the result in `main()`.

### Exercise 2: Formatting an error with context (`fmt.Errorf`)

Modify the function from Exercise 1 so the error message includes the passed value, e.g. `fmt.Errorf("age %d is out of range [0-120]", age)`.

### Exercise 3: Happy Path on the Left (Clean code structure)

Write an order-processing function `ProcessOrder(id int, amount float64) error`. Perform 3 validations: ID > 0, amount > 0, amount < 10000. Write the code so that on error you return immediately (`if err != nil { return err }`), avoiding nested `else` blocks.

### Exercise 4: Sentinel Errors (Error constants/variables)

Declare package-level error variables (naming convention: `Err...`):

* `var ErrNotFound = errors.New("item not found")`
* `var ErrPermissionDenied = errors.New("permission denied")`
Write a function `FindUser(id int) (*User, error)` that returns `ErrNotFound` when ID is 0. Check this error in `main()` with a plain comparison `if err == ErrNotFound`.

### Exercise 5: Ignoring an error (What `_` does)

Call a standard library function (e.g. `strconv.Atoi("123")`) that returns `(int, error)`. Ignore the error with the `_` operator. Consider why linters (e.g. `golangci-lint`) treat this as an anti-pattern in production code.

---

## Part 2: Custom Error Types and Structs (Exercises 6–10)

### Exercise 6: Struct as an error

Create your own error struct `ValidationError`:

```go
type ValidationError struct {
    Field   string
    Message string
}

```

Implement the `error` interface for it (an `Error() string` method that nicely formats the field and message).

### Exercise 7: Returning your own struct as `error`

Write a function `RegisterUser(email, password string) error`. If the email does not contain `@`, return `&ValidationError{Field: "email", Message: "missing @"}`. Notice that the function returns the interface type `error`!

### Exercise 8: Extracting your own error (Type Assertion)

In `main()`, call `RegisterUser` with a bad email. Receive the error as type `error` and use a type assertion (`ve, ok := err.(*ValidationError)`) to access the concrete `Field` and `Message` fields.

### Exercise 9: The `nil` trap with struct types (Key nuance!)

Write a function `badValidate() error` in which you declare a pointer to your struct `var customErr *ValidationError = nil`, then return `customErr`. Check in `main()` why the condition `if err != nil` evaluates to `true`! *(Hint: an interface holding type `*ValidationError` and a `nil` value is itself NOT equal to `nil`)*.

### Exercise 10: Correctly returning `nil` from custom structs

Fix the bug from Exercise 9. Analyze why you should always return an explicit `nil` literal (`return nil`) instead of a nil pointer variable typed as an error.

---

## Part 3: Error Wrapping and the Modern `errors` Package (Exercises 11–15)

### Exercise 11: Wrapping errors with the `%w` verb

Write a function `ReadConfig(path string) error`. Call `os.ReadFile(path)` inside it. If you get an error, return it wrapped with your own context:
`return fmt.Errorf("failed to read config file %s: %w", path, err)`.

### Exercise 12: Checking errors with `errors.Is`

Write a function that calls `ReadConfig` with a nonexistent path. In `main()`, use `errors.Is(err, os.ErrNotExist)` instead of plain `==` to check whether the underlying error (inside the wrap chain) is a file-not-found error.

### Exercise 13: Unwrapping errors with `errors.As`

Use `errors.As` to safely extract an `*os.PathError` from a wrapped error returned by file operations. Print that error's `Path` field.

### Exercise 14: Joining multiple errors (`errors.Join` - Go 1.20+)

Write a form-validation function that collects all validation errors into a `[]error` slice and finally returns a combined error with `errors.Join(err1, err2, err3)`. Check how the printed error looks.

### Exercise 15: Layered error wrapping in a multi-tier architecture

Create a 3-layer call chain: `Repository` → `Service` → `Handler`.

1. `Repository` returns `ErrNotFound`.
2. `Service` wraps it: `fmt.Errorf("user service: %w", err)`.
3. `Handler` checks `errors.Is(err, ErrNotFound)` and returns HTTP status 404.

---

## Part 4: Panic, Recover, and Exceptional Situations (Exercises 16–20)

### Exercise 16: When to use `panic`?

In Go, `panic` is used extremely rarely (equivalent to a critical application abort, e.g. a missing config file at startup). Write a function `MustParseURL(rawURL string)` that does `panic("invalid URL")` if the URL is empty.

### Exercise 17: Catching a panic (`defer` + `recover`)

Write a function `SafeExecute(fn func())` that calls the passed function `fn`. Use a `defer` block and `recover()` inside it to catch a possible panic and prevent the whole application from crashing.

### Exercise 18: Converting `panic` to `error`

Modify `SafeExecute` so it returns an `error`. If a panic occurred inside, turn the recovered value (`recover()`) into a regular `error` and return it.

### Exercise 19: Cleaning up resources with `defer` on errors

Open a file with `file, err := os.Open(...)`. Right after checking the error, add `defer file.Close()`. Call a helper function that returns an error and make sure the file is closed before leaving the main function.

### Exercise 20: Error-logging middleware

Write a simple wrapper function (middleware): `ExecuteWithLogging(fn func() error)`. This function calls `fn()`, checks whether an error was returned, and if so — logs it to the console with a timestamp using the standard `log` package.

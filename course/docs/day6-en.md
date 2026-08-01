Here are **20 exercises for Day 6: Packages, Project Structure, and the Standard Library**.

Today you move from single-file scripts to building **real, structured Go applications**. You will master package organization, encapsulation (field visibility), and key standard library packages (`net/http`, `context`, `slog`, `time`).

---

## Part 1: Packages, Visibility, and Project Structure (Exercises 1–5)

### Exercise 1: Your first sub-package

Create a `config/` folder with a `config.go` file declaring `package config`. Define an `AppConfig` struct with a `Port int` field. Import this package in the main `main.go` file and use the struct.

### Exercise 2: Public vs Private (Capitalization)

In the `config` package, create two functions: `Load()` (public) and `parseEnv()` (private). See what happens when you try to call `config.parseEnv()` from `main.go`.

### Exercise 3: Encapsulation and Getters/Setters

In a new `user/` package, define a `User` struct with a private `email string` field. Expose public methods `SetEmail(e string) error` (with `@` validation) and `Email() string` (getter).

> **Go idiom tip:** Go does not use a `Get` prefix for getters. Instead of `GetEmail()`, just use `Email()`.

### Exercise 4: Import aliases and avoiding conflicts

Imagine you import two packages with the same final name (e.g. `math/rand` and `crypto/rand`). Use an import alias in `main.go` so you can use both at once:

```go
import (
    crand "crypto/rand"
    mrand "math/rand"
)

```

### Exercise 5: Blank imports (Side-effects `_`)

See how imports work when you only need side effects (e.g. registering a database driver): `import _ "[github.com/lib/pq](https://github.com/lib/pq)"`. Learn what the special `init()` function is for in packages.

---

## Part 2: Context (`context.Context`) – A Core Go Concept (Exercises 6–10)

### Exercise 6: Creating a base context

In Go, `context.Context` carries cancellation signals, deadlines, and request metadata. Create a base context with `ctx := context.Background()`.

### Exercise 7: Passing values in context (`context.WithValue`)

Create a function `processRequest(ctx context.Context)`. Add a request ID to the context: `ctx = context.WithValue(ctx, "request_id", "abc-123")`. Inside the function, extract that value and check its type with a type assertion.

### Exercise 8: Cancelling operations (`context.WithCancel`)

Create a context with a cancel function: `ctx, cancel := context.WithCancel(context.Background())`. Run a simulated long operation in a `select` loop listening on `<-ctx.Done()`. Call `cancel()` and observe how the operation stops immediately.

### Exercise 9: Timeouts (`context.WithTimeout`)

Create a context that cancels automatically after 100 milliseconds: `ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)` (always remember `defer cancel()!`). Simulate a 500 ms operation and handle interruption due to timeout.

### Exercise 10: Passing context as the FIRST argument

By Go convention, if a function takes a context, **it must be the first argument**: `func FetchData(ctx context.Context, id string) error`. Coming from TS (where timeouts or options often go at the end), adapt your functions to the Go standard.

---

## Part 3: Building an HTTP Server (`net/http`) (Exercises 11–15)

### Exercise 11: The simplest HTTP server

Create an HTTP server without external frameworks (like Express in Node.js). Use `http.HandleFunc` with the new routing introduced in Go 1.22:

```go
http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
})
http.ListenAndServe(":8080", nil)

```

### Exercise 12: Serving JSON data

Write a `GET /api/user` handler that creates a `User` struct, sets the `w.Header().Set("Content-Type", "application/json")` header, and serializes data directly to the response stream with `json.NewEncoder(w).Encode(user)`.

### Exercise 13: Reading JSON from the body (`POST`)

Write a `POST /api/user` handler that decodes the request body into a struct with `json.NewDecoder(r.Body).Decode(&user)`. Handle errors for invalid JSON.

### Exercise 14: Path parameters (Path Values in Go 1.22+)

Write a `GET /users/{id}` handler that reads a variable from the URL using the built-in method `id := r.PathValue("id")`.

### Exercise 15: Simple HTTP middleware

Middleware in Go is a function that takes an `http.Handler` and returns an `http.Handler`. Write a `LoggingMiddleware` that measures how long each HTTP request takes (use `time.Now()` and `time.Since()`) and prints the method and path.

---

## Part 4: Modern Logging (`slog`) and Time (`time`) (Exercises 16–20)

### Exercise 16: Structured logs with `log/slog`

Since Go 1.21, the standard library includes the `slog` package. Instead of plain `fmt.Println`, use `slog.Info("user logged in", "user_id", 42, "role", "admin")`.

### Exercise 17: Logging in JSON format

Configure `slog` to emit logs in JSON format (ideal for production and Datadog/Grafana):

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)

```

### Exercise 18: Working safely with time (`time.Time`)

In JS, date operations can be a nightmare. In Go you have a powerful `time` package. Create two dates, subtract them (`diff := t2.Sub(t1)`), get a `time.Duration`, and check how many seconds or milliseconds that is.

### Exercise 19: Formatting and parsing dates

In Go, dates are formatted using a **specific reference time**: `Mon Jan 2 15:04:05 MST 2006` (remember the digit layout: 1 2 3 4 5 6 7). Format the current time as `YYYY-MM-DD` using the pattern `"2006-01-02"`.

### Exercise 20: Ticker and Timer

Use `time.NewTicker(1 * time.Second)` to create a loop that runs an action every second (e.g. printing status to the console). Remember to stop the ticker with `defer ticker.Stop()`.

---

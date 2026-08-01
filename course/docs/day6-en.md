Here are **20 exercises for Day 6: Concurrency – Goroutines, Channels, and Select**.

Today you dive into one of Go’s most important and distinctive features. Instead of Node.js’s single-threaded event loop and Promises, you will learn the **CSP (Communicating Sequential Processes)** model.

Go’s golden rule: *"Do not communicate by sharing memory; share memory by communicating."*

---

## Part 1: Lightweight Threads (`Goroutines`) and Synchronization (Exercises 1–5)

### Exercise 1: Your first goroutine

Write a `sayHello(name string)` function that prints a greeting. Call it from `main()` with the `go` keyword: `go sayHello("Gopher")`. Notice that the program exits before the function prints anything. Figure out why that happened.

### Exercise 2: Waiting for goroutines (`sync.WaitGroup`)

Use a `sync.WaitGroup` to fix Exercise 1:

1. Create `var wg sync.WaitGroup`.
2. Add a task to the counter: `wg.Add(1)`.
3. Pass a pointer to `wg` into the goroutine and call `defer wg.Done()` inside it.
4. At the end of `main()`, call `wg.Wait()`.

### Exercise 3: Launching many goroutines in a loop (the loop-variable trap)

Start 5 goroutines in a `for i := 0; i < 5; i++` loop. Pass `i` as an argument into the function inside the goroutine. Check what happens if, instead of passing `i` as a parameter, you use it directly inside an anonymous function: `go func() { fmt.Println(i) }()`.

### Exercise 4: Investigating data races

Write a program where 100 goroutines increment a shared `counter++` variable at the same time with no synchronization. Run it with the race detector:

```bash
go run -race main.go

```

Observe the *Data Race* report.

### Exercise 5: Mutex – protecting shared resources (`sync.Mutex`)

Fix the race from Exercise 4. Use `var mu sync.Mutex` and protect the counter update with `mu.Lock()` and `mu.Unlock()`. Run again with `-race` and confirm the race is gone.

---

## Part 2: Channels – Basics (Exercises 6–10)

### Exercise 6: Unbuffered channel

Create a channel that carries strings: `ch := make(chan string)`. Start a goroutine that sends a message: `ch <- "ping"`. In the main goroutine, receive it into a variable: `msg := <-ch` and print it.

### Exercise 7: Blocking with no receiver (Deadlock)

See what happens when you send to an unbuffered channel with `ch <- "data"` in the same goroutine (`main`), without starting another goroutine to receive. Analyze the error *fatal error: all goroutines are asleep - deadlock!*.

### Exercise 8: Buffered channel

Create a channel with buffer size 2: `ch := make(chan int, 2)`. Send two values back-to-back in the main goroutine (`ch <- 1`, `ch <- 2`). Notice the program does not block. What happens if you try to send a third value?

### Exercise 9: Closing a channel (`close`) and `range`

Write a producer goroutine that sends numbers 1 through 5 into a channel, then **closes the channel**: `close(ch)`. In the main goroutine (consumer), use `for val := range ch` to receive all values.

### Exercise 10: Safely checking whether a channel is open

Receive from a closed channel with the two-value form: `val, ok := <-ch`. Check what `val` and `ok` are for an open vs closed channel.

---

## Part 3: `select` and Advanced Patterns (Exercises 11–15)

### Exercise 11: The `select` statement

Create two channels `ch1` and `ch2`. Start two goroutines that send to those channels after different delays (`time.Sleep`). Use a `select` block to receive from whichever channel responds **first**.

### Exercise 12: Operation timeout with `select`

Write a function that reads from a channel with a time limit. Use `select` combining your channel with `time.After(2 * time.Second)`:

```go
select {
case res := <-dataChan:
    fmt.Println("Received:", res)
case <-time.After(2 * time.Second):
    fmt.Println("Timeout!")
}

```

### Exercise 13: Directional channels

For stronger type safety, functions can restrict channel permissions:

* `func produce(ch chan<- int)` – send-only channel.
* `func consume(ch <-chan int)` – receive-only channel.

Write two functions with those signatures and connect them with a channel in `main()`.

### Exercise 14: Cancellation with `ctx.Done()` and `select`

Write a worker in an infinite `for` loop that uses `select` to watch two channels: a jobs channel and `ctx.Done()` from a context. When the context is cancelled, the worker should exit.

### Exercise 15: Non-blocking channel operations (`default` in `select`)

Use a `default` clause in a `select` block to attempt a receive without blocking when the channel is empty.

---

## Part 4: Practical Concurrency Patterns (Exercises 16–20)

### Exercise 16: Worker pool

This is one of the most common production patterns in Go!

1. Create `jobs := make(chan int, 100)`.
2. Create `results := make(chan int, 100)`.
3. Start 3 worker goroutines; each loops over `jobs`, does work (e.g. `job * 2`), and sends to `results`.
4. Send 10 jobs and close `jobs`. Receive 10 results.

### Exercise 17: Fan-Out, Fan-In

* **Fan-Out:** Split one job source across many goroutines working in parallel.
* **Fan-In:** Merge results from several channels into one shared output channel via a `merge` function.

### Exercise 18: HTTP request race (First Responder / Hedged Requests)

Write a function that sends the same request (e.g. fetching data) to 3 different servers/mirrors at once in separate goroutines. Use a buffered channel to take the **first response** that arrives and ignore the rest.

### Exercise 19: Rate limiting

Use `time.Tick(200 * time.Millisecond)` to build a limiter that allows at most 5 requests per second.

### Exercise 20: Empty struct for signal channels (`chan struct{}`)

Often a channel only carries a signal (“something happened”), not data. In Go you use `chan struct{}` because `struct{}` takes **0 bytes**. Create `done := make(chan struct{})` and use it to notify `main()` that background work finished.

---

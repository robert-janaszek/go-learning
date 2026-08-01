Here are **20 exercises for Day 4: Interfaces, Duck Typing, and Polymorphism**.

Your goal today: understand why interfaces in Go are defined on the *consumer* side, how implicit interfaces work, and how to write clean, loosely coupled code without classes and inheritance.

---

## Part 1: Interface Basics and Implicit Implementation (Exercises 1–5)

### Exercise 1: Your first interface

Define a `Stringer` interface with one method: `String() string`. Create a `User` struct (`Name string`, `Age int`) and implement a `String() string` method for it. Notice that **you don't use any `implements` keyword**.

### Exercise 2: Polymorphism in a function

Write a function `PrintInfo(s Stringer)` that accepts anything satisfying the `Stringer` interface and calls `String()` on it. Call it in `main()` with a `User` instance.

### Exercise 3: Multiple structs satisfying the same interface

Create a second `Book` struct (`Title string`, `Author string`) and also implement a `String() string` method for it. Pass a `Book` instance to the same `PrintInfo` function.

### Exercise 4: Interface zero value (`nil interface`)

Declare an interface-typed variable: `var s Stringer` without assigning any struct to it. Check with `if s == nil` whether it is empty. See what happens if you try to call `s.String()` on an empty interface.

### Exercise 5: Slice / list of interfaces

Create a slice of interface-typed elements: `items := []Stringer{user, book}`. Iterate over it with `for _, item := range items` and call `item.String()` for each.

---

## Part 2: Pointer Receiver vs Value Receiver in Interfaces (Exercises 6–10)

### Exercise 6: Interface with a Pointer Receiver – The assignment trap

Create a `Saver` interface with a `Save() error` method. Create a `Document` struct and implement `Save()` using a **pointer receiver** `(d *Document)`.

### Exercise 7: Trying to assign a value (Check the compiler error!)

Try assigning a plain value to the interface: `var s Saver = Document{}`. See the compiler error! Why must you pass a pointer `var s Saver = &Document{}`? *(This is a key nuance in Go!)*.

### Exercise 8: Interface with a Value Receiver

Create a second `Note` struct and implement `Save()` with a **value receiver** `(n Note)`. Check whether you can assign both `Note{}` (value) and `&Note{}` (pointer) to `Saver`.

### Exercise 9: Standard `io.Reader` and `io.Writer` interfaces

Go has excellent built-in interfaces. Get familiar with `io.Writer` (`Write(p []byte) (n int, err error)`). Write a function `WriteHello(w io.Writer)` that writes the bytes `"Hello Go"` to anything that accepts this interface (e.g. `os.Stdout` or `bytes.Buffer`).

### Exercise 10: Writing to a file and the console with the same code

Use the `WriteHello` function from Exercise 9 twice: once passing `os.Stdout` (console output), and once passing a file created by `os.Create("test.txt")`. Notice the power of polymorphism without building complex class hierarchies!

---

## Part 3: Type Assertion, Type Switch, and Empty Interface (Exercises 11–15)

### Exercise 11: Empty interface (`any` / `interface{}`)

Go 1.18 introduced the `any` alias (equivalent to `interface{}`). It corresponds to TypeScript's `unknown` type. Create a function `Describe(i any)` that accepts any type and prints it.

### Exercise 12: Type Assertion

Given a variable `var val any = "hello Go"`, extract the original `string` type with an assertion: `s := val.(string)`. Print the length of that string (`len(s)`).

### Exercise 13: Safe Type Assertion (`val, ok` pattern)

What happens if you do `num := val.(int)` on a variable holding a string? (The app will panic!). Write a safe version using `n, ok := val.(int)` and handle the case when `ok == false`.

### Exercise 14: Type Switch (Equivalent of `match` / `switch typeof`)

Write a function `ProcessInput(v any)` that uses `switch v.(type)` to check the type of the passed parameter (`int`, `string`, `bool`, `Player`). For each type, print an appropriate message.

### Exercise 15: Extracting a value from a struct behind an interface

Create a `Payer` interface. Create a `CreditCard` struct with a unique `CardNumber string` field. Assign `CreditCard` to a `Payer`-typed variable. Use a type assertion to extract `CardNumber`.

---

## Part 4: Interface Composition and Best Practices (Exercises 16–20)

### Exercise 16: Combining interfaces (Interface Embedding)

Define two small interfaces:

1. `Reader` with a `Read() string` method
2. `Writer` with a `Write(data string)` method
Create a third `ReadWriter` interface that **embeds both of these interfaces**.

### Exercise 17: Go's golden rule: "Small interfaces"

In TypeScript people often create huge interfaces with a dozen methods. In Go the ideal interface has **1 or 2 methods**. Create a `FileHandler` struct that satisfies the `ReadWriter` interface from Exercise 16.

### Exercise 18: Consumer-defined Interface

This is the most important concept in Go. Create a `store/` sub-package with a `PostgresStore` struct that has a `GetUsers() []string` method. In the `main` package, define a `UserGetter` interface and use it in a `UserService`. *(Notice: `PostgresStore` knows nothing about the interface in `main`!)*.

### Exercise 19: Easy mocking in tests (preview)

Thanks to the approach from Exercise 18, create a `MockStore` struct in `main` that implements `GetUsers() []string` with fake data. Swap `PostgresStore` for `MockStore` in `UserService`.

### Exercise 20: Built-in interfaces from the standard library (`error`)

Did you know that Go's built-in `error` type is just an interface with one method?

```go
type error interface {
    Error() string
}

```

Create your own `CustomError` struct (`Code int`, `Message string`), implement an `Error() string` method for it, and return it as a plain `error` from a function.

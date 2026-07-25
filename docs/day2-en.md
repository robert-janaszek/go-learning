Here are **20 exercises for Day 2: Structs, Methods, and No Classes**.

Your goal today: master composition over inheritance and learn to write methods attached to structs with the right receiver (*value* vs *pointer receiver*).

---

## Part 1: Defining and Initializing Structs (Exercises 1–5)

### Exercise 1: Basic struct

Define a `Book` struct with fields: `Title` (string), `Author` (string), `Pages` (int), `IsRead` (bool). Create an instance of this struct in `main()` using field names (*struct literal*) and print it.

### Exercise 2: Different ways to initialize

Create 3 instances of the `Book` struct:

1. Using field names (`Book{Title: "...", ...}`).
2. Without field names (watch the order!).
3. An empty instance (`b := Book{}`) and fill the fields on separate lines (`b.Title = "..."`).

### Exercise 3: Pointer to a struct and factory function (Constructor Pattern)

Go has no `new` keyword in the class sense. You create functions like `NewBook`. Write a function `NewBook(title, author string, pages int) *Book` that returns a **pointer** to a newly created struct.

### Exercise 4: Anonymous struct (Ad-hoc Struct)

Often in Go (e.g. in tests or when extracting JSON) you create one-off structs. Create an anonymous struct with fields `ConfigName` and `Port`, initialize it immediately with values, and print it.

### Exercise 5: Comparing structs

Create two `Book` instances with identical values. Check with `if b1 == b2` whether Go can compare them. Then add a slice field `Tags []string` to the struct and see why the code no longer compiles (structs with reference types are not comparable with `==`).

---

## Part 2: Methods – Pointer vs Value Receiver (Exercises 6–10)

### Exercise 6: First method with a Value Receiver

Add a `Summary() string` method to the `Book` struct. The method should return a string in the format `"Title" by Author (X pages)`. Use a *value receiver* `(b Book)`.

### Exercise 7: Method with a Pointer Receiver (State modification)

Add a `MarkAsRead()` method to `Book`. Think: should the method receiver be a pointer `(b *Book)` or a value `(b Book)`? Test in `main()` by calling this method on a book that had `IsRead = false`.

### Exercise 8: Calling a pointer-receiver method on a value

Create a variable `b := Book{Title: "Go in Action"}` (a plain value, not a pointer). Call the `MarkAsRead()` method from Exercise 7 on it. Notice that Go **automatically takes the address** (`(&b).MarkAsRead()`) — you don't need to turn the variable into a pointer yourself.

### Exercise 9: Methods on custom (basic) types

In Go you can attach methods not only to structs! Define a custom type: `type Celsius float64`. Add a `ToFahrenheit() float64` method to it. Test in `main()`.

### Exercise 10: Method that mutates a custom type

Add an `Add(degrees float64)` method to the `Celsius` type. Choose the appropriate receiver so the method actually modifies the temperature value it was called on.

---

## Part 3: Composition and Embedding Instead of Inheritance (Exercises 11–15)

### Exercise 11: Plain nested structs

Create an `Address` struct (`City string`, `ZipCode string`). Create a `User` struct with fields `Name string` and `Addr Address`. Initialize a `User` and print the user's city (`u.Addr.City`).

### Exercise 12: Anonymous Struct Embedding (Promoted fields)

Modify `User` so that the `Address` field is **anonymous** (an embedded struct):

```go
type User struct {
    Name string
    Address // No field name, just the type!
}

```

Check in `main()` how *field promotion* works — access the city by simply writing `u.City`.

### Exercise 13: Overriding fields and methods (Shadowing in composition)

Add a `FullAddress() string` method to `Address`. Then add your own `FullAddress() string` method to `User` that returns the name and address. Call both in `main()` and see how Go resolves name conflicts.

### Exercise 14: Embedding a pointer to a struct

Create an `Engine` struct (`HorsePower int`). Create a `Car` struct with an embedded `*Engine` pointer. Check what happens when you call a method on `Car` if `Engine` is `nil`.

### Exercise 15: Composition from multiple structs

Create two small structs: `Logger` (method `Log(msg string)`) and `Database` (method `Connect()`). Create a `Server` struct that embeds **both** of these structs. Call `server.Log()` and `server.Connect()`.

---

## Part 4: Practical Patterns, JSON, and Tags (Exercises 16–20)

### Exercise 16: Struct Tags

Define a `Product` struct with fields `ID int`, `Name string`, `Price float64`. Add JSON tags, e.g. `json:"product_id"`. Use `json.Marshal(p)` from the `encoding/json` package to turn the struct into JSON bytes and print the result to the console (`string(bytes)`).

### Exercise 17: Hiding fields in JSON (`json:"-"` and `omitempty`)

Add fields to the `Product` struct:

* `InternalCode string` — should be ignored by JSON (`json:"-"`).
* `Discount float64` — should be omitted from JSON when equal to 0 (`json:"discount,omitempty"`).
Check how `json.Marshal` behaves.

### Exercise 18: Unmarshaling (JSON -> Struct)

Create a variable with a JSON string: `jsonData := []byte('{"name":"Laptop", "price": 2500}')`. Use `json.Unmarshal(jsonData, &p)` to load the data into a `Product` struct. Careful: why must you pass `&p` (a pointer), and not just `p`?

### Exercise 19: Shopping cart with a total method

Warm-up before the project: Create `CartItem` (`Product Product`, `Quantity int`). Create `Cart` with an `Items []CartItem` field. Add methods: `AddItem(p Product, qty int)` and `Total() float64`.

### Exercise 20: Encapsulation and private fields

Create a sub-package in a `bank/` folder. Define an `Account` struct in it with a **private** field `balance float64`. Expose public methods `Deposit(amount float64)`, `Withdraw(amount float64) error`, and `Balance() float64`. Verify in `main.go` that you cannot modify the `balance` field directly.

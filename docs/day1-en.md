Here are **20 short, highly practical exercises** that take you from the basics of variables and types, through memory manipulation, to the nuances of passing data into functions.

They are split into sections so you can gradually build intuition around pointers.

---

## Part 1: Variables, Types, and Declarations (Exercises 1–5)

### Exercise 1: Different ways to declare

Declare a variable `age` of type `int` in 3 ways: using `var` with an explicit type, using `var` with type inference, and using the short declaration `:=`. Print their types with `fmt.Printf("%T\n", age)`.

### Exercise 2: Type conversion (No implicit casting)

In JS you can write `5 + "5"`. In Go there is no automatic type casting. Create a variable `a int = 42` and `b float64 = 3.14`. Convert `a` to `float64`, add it to `b`, and assign the result to a new variable.

### Exercise 3: Zero Values

In JS an uninitialized variable is `undefined`. Declare with `var` without assigning a value: `int`, `float64`, `string`, `bool`, and a pointer `*int`. Print them and check what their default values (*zero values*) are.

### Exercise 4: Constants (`const`) and `iota`

Create a constant block representing the days of the week (from `Monday` to `Sunday`) using the `iota` generator. Print their numeric values to the console.

### Exercise 5: Shadowing

Create a variable `x := 10` in an outer block. Open a new code block `{ ... }`, declare `x := 20` inside it, and print `x`. Outside the block, print `x` again. Analyze what happened.

---

## Part 2: Pointer Basics – Addresses and Dereferencing (Exercises 6–10)

### Exercise 6: Taking an address (`&`)

Create a variable `score := 100`. Create a variable `ptr` that holds the address of `score`. Print the value of `score`, the address of `score`, and the type of `ptr`.

### Exercise 7: Dereferencing (`*`)

Using the pointer `ptr` from Exercise 6, change the value of `score` to `200` **using only the pointer `ptr`** (dereference operator `*ptr = ...`). Print `score`.

### Exercise 8: Two pointers to one variable

Create `x := 50`. Create two pointers `p1` and `p2`, both pointing to `x`. Change the value via `p1` to `100`, then print the value obtained through `*p2`.

### Exercise 9: Nil pointer

Declare a pointer `var p *int` (without assigning an address). Check with `if p == nil` whether the pointer is empty. See what happens (and what error the runtime reports) if you try `*p = 10` without initialization (a *nil pointer dereference panic*).

### Exercise 10: Double pointer (`**int`)

Create a variable `val := 42`. Create a pointer `p` pointing to `val`. Create a pointer to a pointer `pp` (type `**int`) pointing to `p`. Read the value `42` using only `pp`.

---

## Part 3: Pointers in Functions (Exercises 11–15)

### Exercise 11: Swap

Write a function `swap(a, b *int)` that swaps the values of two variables. Test it in `main()` with two variables `x := 1` and `y := 2`.

### Exercise 12: Modifying a string

Write a function `uppercase(s *string)` that takes a pointer to a string and changes its content to uppercase (use `strings.ToUpper`). Check the result in `main()`.

### Exercise 13: Safe division with an optional result

Write a function `safeDivide(a, b float64, result *float64) bool`. If `b == 0`, the function returns `false`. Otherwise, it writes the result at the address `result` and returns `true`.

### Exercise 14: Counter (Incrementer)

Create a structure/variable representing a counter. Write a function `increment(val *int)` that increases the value by 1 on each call. Call it 3 times in a loop.

### Exercise 15: Function returning a pointer

Write a function `createInt(val int) *int` that creates a local variable inside the function and returns its address `&localVal`.
*(Context for C/C++ programmers: In Go this is fully safe! The compiler will perform Escape Analysis and allocate that variable on the heap instead of the stack).*

---

## Part 4: Pointers and Structs (Exercises 16–20)

### Exercise 16: Player struct and state modification

Define a `Player` struct with fields `Name string` and `Health int`. Write a function `takeDamage(p *Player, amount int)` that decreases the player's health points.

### Exercise 17: Value receiver vs Pointer receiver

Add two methods to the `Player` struct:

1. `HealValue(amount int)` with a *value receiver* `(p Player)`
2. `HealPointer(amount int)` with a *pointer receiver* `(p *Player)`
Call both in `main()` and observe which method actually heals the player.

### Exercise 18: Automatic dereferencing with structs

Create a pointer to a struct `p := &Player{Name: "Gopher", Health: 100}`. Change its `Health` field to `90` by simply writing `p.Health = 90`. Notice that in Go you don't need to write `(*p).Health = 90` — the language does this automatically!

### Exercise 19: Optional struct fields (JS/TS concept)

In TypeScript you have optional fields `age?: number` (which can be `undefined`). In Go this is done with pointers!
Create a `User` struct with fields `Name string` and `Age *int`. Create one user without an age (`Age = nil`) and one with an age.

### Exercise 20: Memory leak / Looping over structs

Create a slice of structs `[]Player`. Iterate over it with `for _, player := range players`. Try changing `player.Health = 0` inside the loop. Check why this **doesn't work** (the `range` loop variable is only a copy) and how to fix it using indices `players[i].Health = 0`.

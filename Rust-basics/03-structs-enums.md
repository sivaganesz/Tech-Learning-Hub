# Structs and Enums

## Learning Objectives

By the end of this tutorial, you will be able to:
- Define and instantiate structs
- Use tuple structs and unit structs
- Implement methods and associated functions
- Define enums with and without data
- Use pattern matching with match
- Work with `Option<T>` and `Result<T, E>` types
- Use `if let` for concise pattern matching

---

## Introduction

Now that you understand ownership and borrowing, let's learn how to create custom data types. Structs and enums are Rust's primary tools for structuring data—you'll use them constantly.

---

## Defining Structs

Structs group related data together:

```rust
struct User {
    username: String,
    email: String,
    sign_in_count: u64,
    active: bool,
}

fn main() {
    let user1 = User {
        email: String::from("someone@example.com"),
        username: String::from("someusername123"),
        active: true,
        sign_in_count: 1,
    };

    println!("User: {}", user1.username);
}
```

### Mutable Structs

The entire instance must be mutable to change any field:

```rust
fn main() {
    let mut user1 = User {
        email: String::from("someone@example.com"),
        username: String::from("someusername123"),
        active: true,
        sign_in_count: 1,
    };

    user1.email = String::from("anotheremail@example.com");
    user1.sign_in_count += 1;

    println!("Email: {}", user1.email);
}
```

### Field Init Shorthand

When variable names match field names:

```rust
fn build_user(email: String, username: String) -> User {
    User {
        email,      // Same as email: email
        username,   // Same as username: username
        active: true,
        sign_in_count: 1,
    }
}
```

### Struct Update Syntax

Create a new instance using fields from another:

```rust
fn main() {
    let user1 = User {
        email: String::from("someone@example.com"),
        username: String::from("someusername123"),
        active: true,
        sign_in_count: 1,
    };

    let user2 = User {
        email: String::from("another@example.com"),
        ..user1  // Get remaining fields from user1
    };

    // Note: user1.username moved to user2!
    // user1.active and user1.sign_in_count are still valid (Copy types)
}
```

---

## Tuple Structs

Structs without named fields:

```rust
struct Color(i32, i32, i32);
struct Point(i32, i32, i32);

fn main() {
    let black = Color(0, 0, 0);
    let origin = Point(0, 0, 0);

    // Access by index
    println!("R: {}, G: {}, B: {}", black.0, black.1, black.2);

    // Destructure
    let Point(x, y, z) = origin;
    println!("x: {}, y: {}, z: {}", x, y, z);

    // Note: black and origin are different types!
    // Even though they have the same field types
}
```

---

## Unit Structs

Structs with no fields, useful for traits:

```rust
struct AlwaysEqual;

fn main() {
    let _subject = AlwaysEqual;
    // Useful for implementing traits without storing data
}
```

---

## Impl Blocks and Methods

Methods are functions defined within the context of a struct:

```rust
#[derive(Debug)]
struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    // Method - takes &self
    fn area(&self) -> u32 {
        self.width * self.height
    }

    // Method with additional parameters
    fn can_hold(&self, other: &Rectangle) -> bool {
        self.width > other.width && self.height > other.height
    }

    // Method that mutates
    fn double_size(&mut self) {
        self.width *= 2;
        self.height *= 2;
    }
}

fn main() {
    let mut rect1 = Rectangle {
        width: 30,
        height: 50,
    };

    let rect2 = Rectangle {
        width: 10,
        height: 40,
    };

    println!("Area: {} square pixels", rect1.area());
    println!("Can rect1 hold rect2? {}", rect1.can_hold(&rect2));

    rect1.double_size();
    println!("After doubling: {:?}", rect1);
}
```

### The `self` Parameter

- `&self` - borrows the instance immutably
- `&mut self` - borrows the instance mutably
- `self` - takes ownership (rare, usually for transformations)

```rust
impl Rectangle {
    // Takes ownership - consumes the rectangle
    fn into_square(self) -> Rectangle {
        let side = std::cmp::max(self.width, self.height);
        Rectangle {
            width: side,
            height: side,
        }
    }
}
```

---

## Associated Functions

Functions in `impl` blocks that don't take `self`:

```rust
impl Rectangle {
    // Associated function - no self
    fn new(width: u32, height: u32) -> Rectangle {
        Rectangle { width, height }
    }

    // Another associated function
    fn square(size: u32) -> Rectangle {
        Rectangle {
            width: size,
            height: size,
        }
    }
}

fn main() {
    // Called with ::
    let rect = Rectangle::new(30, 50);
    let square = Rectangle::square(10);

    println!("Rectangle: {:?}", rect);
    println!("Square: {:?}", square);
}
```

### Multiple Impl Blocks

You can have multiple `impl` blocks for the same struct:

```rust
impl Rectangle {
    fn area(&self) -> u32 {
        self.width * self.height
    }
}

impl Rectangle {
    fn perimeter(&self) -> u32 {
        2 * (self.width + self.height)
    }
}
```

---

## Enums

Enums define a type with a fixed set of variants:

```rust
enum IpAddrKind {
    V4,
    V6,
}

fn main() {
    let four = IpAddrKind::V4;
    let six = IpAddrKind::V6;

    route(four);
    route(six);
}

fn route(ip_kind: IpAddrKind) {
    // Process the IP
}
```

### Enums with Data

Each variant can hold different types and amounts of data:

```rust
enum IpAddr {
    V4(u8, u8, u8, u8),
    V6(String),
}

fn main() {
    let home = IpAddr::V4(127, 0, 0, 1);
    let loopback = IpAddr::V6(String::from("::1"));
}
```

### Complex Enum Example

```rust
enum Message {
    Quit,                       // No data
    Move { x: i32, y: i32 },    // Named fields (like struct)
    Write(String),              // Single String
    ChangeColor(i32, i32, i32), // Three i32 values
}

impl Message {
    fn call(&self) {
        match self {
            Message::Quit => println!("Quit"),
            Message::Move { x, y } => println!("Move to ({}, {})", x, y),
            Message::Write(text) => println!("Write: {}", text),
            Message::ChangeColor(r, g, b) => {
                println!("Change color to ({}, {}, {})", r, g, b)
            }
        }
    }
}

fn main() {
    let messages = vec![
        Message::Quit,
        Message::Move { x: 10, y: 20 },
        Message::Write(String::from("hello")),
        Message::ChangeColor(255, 0, 0),
    ];

    for msg in messages {
        msg.call();
    }
}
```

---

## Pattern Matching with `match`

`match` is Rust's powerful control flow construct for pattern matching:

```rust
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter,
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}

fn main() {
    let coin = Coin::Quarter;
    println!("Value: {} cents", value_in_cents(coin));
}
```

### Match with Code Blocks

```rust
fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => {
            println!("Lucky penny!");
            1
        }
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter => 25,
    }
}
```

### Patterns That Bind to Values

```rust
#[derive(Debug)]
enum UsState {
    Alabama,
    Alaska,
    Arizona,
    // ... etc
}

enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(UsState),
}

fn value_in_cents(coin: Coin) -> u8 {
    match coin {
        Coin::Penny => 1,
        Coin::Nickel => 5,
        Coin::Dime => 10,
        Coin::Quarter(state) => {
            println!("State quarter from {:?}!", state);
            25
        }
    }
}

fn main() {
    let coin = Coin::Quarter(UsState::Alaska);
    value_in_cents(coin);
}
```

### Catch-All and Placeholder

```rust
fn main() {
    let dice_roll = 9;

    match dice_roll {
        3 => add_fancy_hat(),
        7 => remove_fancy_hat(),
        other => move_player(other),  // Catch-all that binds
    }

    match dice_roll {
        3 => add_fancy_hat(),
        7 => remove_fancy_hat(),
        _ => reroll(),  // Catch-all that doesn't bind
    }

    match dice_roll {
        3 => add_fancy_hat(),
        7 => remove_fancy_hat(),
        _ => (),  // Do nothing
    }
}

fn add_fancy_hat() {}
fn remove_fancy_hat() {}
fn move_player(num_spaces: u8) {}
fn reroll() {}
```

### Match Must Be Exhaustive

```rust
enum Color {
    Red,
    Green,
    Blue,
}

fn describe(color: Color) -> &'static str {
    match color {
        Color::Red => "red",
        Color::Green => "green",
        // ERROR! Non-exhaustive patterns: Color::Blue not covered
    }
}
```

---

## The `Option` Type

Rust doesn't have null. Instead, it has `Option<T>`:

```rust
enum Option<T> {
    None,
    Some(T),
}
```

`Option` is so common it's included in the prelude—you don't need to import it:

```rust
fn main() {
    let some_number = Some(5);
    let some_string = Some("a string");
    let absent_number: Option<i32> = None;

    println!("{:?}, {:?}, {:?}", some_number, some_string, absent_number);
}
```

### Working with Option

```rust
fn main() {
    let x: i8 = 5;
    let y: Option<i8> = Some(5);

    // Can't add i8 and Option<i8> directly
    // let sum = x + y;  // ERROR!

    // Must handle the Option
    let sum = x + y.unwrap_or(0);
    println!("Sum: {}", sum);
}
```

### Matching on Option

```rust
fn plus_one(x: Option<i32>) -> Option<i32> {
    match x {
        None => None,
        Some(i) => Some(i + 1),
    }
}

fn main() {
    let five = Some(5);
    let six = plus_one(five);
    let none = plus_one(None);

    println!("{:?}, {:?}", six, none);  // Some(6), None
}
```

---

## The `Result` Type

For operations that can fail:

```rust
enum Result<T, E> {
    Ok(T),
    Err(E),
}
```

```rust
use std::fs::File;

fn main() {
    let f = File::open("hello.txt");

    let f = match f {
        Ok(file) => file,
        Err(error) => panic!("Problem opening the file: {:?}", error),
    };
}
```

---

## `if let` for Concise Matching

When you only care about one pattern:

```rust
fn main() {
    let config_max = Some(3u8);

    // Verbose match
    match config_max {
        Some(max) => println!("Maximum is configured to be {}", max),
        _ => (),
    }

    // Concise if let
    if let Some(max) = config_max {
        println!("Maximum is configured to be {}", max);
    }
}
```

### `if let` with `else`

```rust
enum Coin {
    Penny,
    Nickel,
    Dime,
    Quarter(String),
}

fn main() {
    let coin = Coin::Quarter(String::from("Alaska"));

    // With if let
    if let Coin::Quarter(state) = coin {
        println!("State quarter from {}!", state);
    } else {
        println!("Not a quarter");
    }
}
```

### `while let`

```rust
fn main() {
    let mut stack = vec![1, 2, 3];

    while let Some(top) = stack.pop() {
        println!("{}", top);
    }
}
```

---

## Complete Example: A Simple Expression Evaluator

```rust
#[derive(Debug)]
enum Expr {
    Number(f64),
    Add(Box<Expr>, Box<Expr>),
    Sub(Box<Expr>, Box<Expr>),
    Mul(Box<Expr>, Box<Expr>),
    Div(Box<Expr>, Box<Expr>),
}

impl Expr {
    fn evaluate(&self) -> Option<f64> {
        match self {
            Expr::Number(n) => Some(*n),
            Expr::Add(a, b) => {
                let left = a.evaluate()?;
                let right = b.evaluate()?;
                Some(left + right)
            }
            Expr::Sub(a, b) => {
                let left = a.evaluate()?;
                let right = b.evaluate()?;
                Some(left - right)
            }
            Expr::Mul(a, b) => {
                let left = a.evaluate()?;
                let right = b.evaluate()?;
                Some(left * right)
            }
            Expr::Div(a, b) => {
                let left = a.evaluate()?;
                let right = b.evaluate()?;
                if right == 0.0 {
                    None  // Division by zero
                } else {
                    Some(left / right)
                }
            }
        }
    }
}

fn main() {
    // (3 + 4) * 2
    let expr = Expr::Mul(
        Box::new(Expr::Add(
            Box::new(Expr::Number(3.0)),
            Box::new(Expr::Number(4.0)),
        )),
        Box::new(Expr::Number(2.0)),
    );

    match expr.evaluate() {
        Some(result) => println!("Result: {}", result),
        None => println!("Error evaluating expression"),
    }

    // Division by zero
    let bad_expr = Expr::Div(
        Box::new(Expr::Number(10.0)),
        Box::new(Expr::Number(0.0)),
    );

    match bad_expr.evaluate() {
        Some(result) => println!("Result: {}", result),
        None => println!("Error: Division by zero"),
    }
}
```

---

## Common Patterns

### Builder Pattern with Structs

```rust
struct ServerConfig {
    host: String,
    port: u16,
    max_connections: u32,
    timeout: u64,
}

struct ServerConfigBuilder {
    host: String,
    port: u16,
    max_connections: u32,
    timeout: u64,
}

impl ServerConfigBuilder {
    fn new() -> Self {
        ServerConfigBuilder {
            host: String::from("127.0.0.1"),
            port: 8080,
            max_connections: 100,
            timeout: 30,
        }
    }

    fn host(mut self, host: &str) -> Self {
        self.host = host.to_string();
        self
    }

    fn port(mut self, port: u16) -> Self {
        self.port = port;
        self
    }

    fn max_connections(mut self, max: u32) -> Self {
        self.max_connections = max;
        self
    }

    fn build(self) -> ServerConfig {
        ServerConfig {
            host: self.host,
            port: self.port,
            max_connections: self.max_connections,
            timeout: self.timeout,
        }
    }
}

fn main() {
    let config = ServerConfigBuilder::new()
        .host("0.0.0.0")
        .port(3000)
        .max_connections(1000)
        .build();

    println!("Server at {}:{}", config.host, config.port);
}
```

### State Machine with Enums

```rust
enum ConnectionState {
    Disconnected,
    Connecting { attempt: u32 },
    Connected { session_id: String },
    Error { message: String },
}

impl ConnectionState {
    fn connect(&self) -> ConnectionState {
        match self {
            ConnectionState::Disconnected => {
                ConnectionState::Connecting { attempt: 1 }
            }
            ConnectionState::Connecting { attempt } if *attempt < 3 => {
                ConnectionState::Connecting { attempt: attempt + 1 }
            }
            ConnectionState::Connecting { .. } => {
                ConnectionState::Error {
                    message: String::from("Max retries exceeded"),
                }
            }
            _ => self.clone_state(),
        }
    }

    fn clone_state(&self) -> ConnectionState {
        // Simplified clone for this example
        ConnectionState::Disconnected
    }
}
```

---

## Exercises

### Exercise 1: Define a Struct
Create a `Book` struct with:
- title (String)
- author (String)
- pages (u32)
- available (bool)

Implement:
- `new()` associated function
- `borrow()` method (sets available to false if available)
- `return_book()` method
- `summary()` method that prints book info

### Exercise 2: Rectangle Methods
Extend the Rectangle struct with:
- `is_square()` method
- `scale()` method that takes a factor
- `rotate()` method that swaps width and height
- `from_square()` associated function

### Exercise 3: Traffic Light Enum
Create a `TrafficLight` enum with Red, Yellow, Green variants. Implement:
- `duration()` method returning seconds for each light
- `next()` method returning the next state
- A loop that cycles through the lights

### Exercise 4: Shape with Variants
Create a `Shape` enum with:
- Circle(radius: f64)
- Rectangle(width: f64, height: f64)
- Triangle(base: f64, height: f64)

Implement:
- `area()` method
- `perimeter()` method (for Triangle, assume equilateral)

### Exercise 5: Option and Result Practice
Write functions:
1. `divide(a: f64, b: f64) -> Option<f64>` - returns None for division by zero
2. `find_index(vec: &[i32], target: i32) -> Option<usize>` - find element index
3. `parse_and_add(a: &str, b: &str) -> Result<i32, String>` - parse strings and add

Use `match`, `if let`, and the `?` operator.

---

## Summary

You've learned:
- Structs group related data with named fields
- Tuple structs are useful for simple data without named fields
- Methods use `impl` blocks and take `self`
- Associated functions don't take `self` and use `::`
- Enums can hold different types of data in each variant
- `match` provides exhaustive pattern matching
- `Option<T>` replaces null
- `Result<T, E>` handles operations that can fail
- `if let` provides concise pattern matching

---

## Next Steps

Now that you can structure your data, let's learn how to handle errors properly in [Error Handling](./04-error-handling.md).

Error handling is crucial for writing robust Rust applications, and `Option` and `Result` are just the beginning!

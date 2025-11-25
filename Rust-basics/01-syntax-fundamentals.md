# Rust Syntax Fundamentals

## Learning Objectives

By the end of this tutorial, you will be able to:
- Declare variables using `let`, `let mut`, and `const`
- Work with Rust's basic data types
- Write functions with parameters and return types
- Use control flow structures effectively
- Understand the difference between expressions and statements
- Write clear comments and documentation

---

## Introduction

Welcome to Rust! If you're coming from other programming languages, Rust will feel both familiar and different. Don't worry if some concepts take time to click—Rust has a learning curve, but it's worth it. The compiler is your friend and will guide you along the way.

Let's start with the fundamentals that you'll use in every Rust program.

---

## Variables and Mutability

### Immutable Variables with `let`

In Rust, variables are immutable by default. This is a deliberate design choice that helps prevent bugs.

```rust
fn main() {
    let x = 5;
    println!("The value of x is: {}", x);

    // This would cause a compile error:
    // x = 6; // error: cannot assign twice to immutable variable
}
```

### Mutable Variables with `let mut`

When you need a variable that can change, use `let mut`:

```rust
fn main() {
    let mut counter = 0;
    println!("Counter: {}", counter);

    counter = 1;
    println!("Counter: {}", counter);

    counter += 1;
    println!("Counter: {}", counter);
}
```

### Constants with `const`

Constants are always immutable and must have their type annotated:

```rust
const MAX_POINTS: u32 = 100_000;
const PI: f64 = 3.14159;
const APP_NAME: &str = "My Rust App";

fn main() {
    println!("Max points: {}", MAX_POINTS);
    println!("Pi: {}", PI);
    println!("App: {}", APP_NAME);
}
```

**Key differences between `let` and `const`:**
- `const` must have a type annotation
- `const` can be declared in any scope, including global
- `const` must be set to a constant expression (computed at compile time)
- `const` uses SCREAMING_SNAKE_CASE by convention

### Shadowing

You can declare a new variable with the same name as a previous one:

```rust
fn main() {
    let x = 5;
    let x = x + 1;        // x is now 6
    let x = x * 2;        // x is now 12

    println!("The value of x is: {}", x);

    // Shadowing also allows type changes
    let spaces = "   ";           // &str
    let spaces = spaces.len();    // usize

    println!("Number of spaces: {}", spaces);
}
```

---

## Basic Data Types

Rust is statically typed, meaning every variable must have a known type at compile time. The compiler can usually infer types, but sometimes you need to be explicit.

### Integer Types

| Length  | Signed | Unsigned |
|---------|--------|----------|
| 8-bit   | i8     | u8       |
| 16-bit  | i16    | u16      |
| 32-bit  | i32    | u32      |
| 64-bit  | i64    | u64      |
| 128-bit | i128   | u128     |
| arch    | isize  | usize    |

```rust
fn main() {
    // Default integer type is i32
    let a = 42;

    // Explicit type annotations
    let b: i64 = 1_000_000;
    let c: u8 = 255;

    // Different literal forms
    let decimal = 98_222;
    let hex = 0xff;
    let octal = 0o77;
    let binary = 0b1111_0000;
    let byte = b'A';  // u8 only

    println!("a: {}, b: {}, c: {}", a, b, c);
    println!("decimal: {}, hex: {}, octal: {}, binary: {}",
             decimal, hex, octal, binary);
}
```

### Floating-Point Types

```rust
fn main() {
    let x = 2.0;      // f64 (default)
    let y: f32 = 3.0; // f32

    // Numeric operations
    let sum = 5.0 + 10.0;
    let difference = 95.5 - 4.3;
    let product = 4.0 * 30.0;
    let quotient = 56.7 / 32.2;
    let remainder = 43.0 % 5.0;

    println!("sum: {}, diff: {}, prod: {}", sum, difference, product);
    println!("quot: {}, rem: {}", quotient, remainder);
}
```

### Boolean Type

```rust
fn main() {
    let t = true;
    let f: bool = false;

    let is_greater = 10 > 5;
    let is_equal = 5 == 5;

    println!("t: {}, f: {}", t, f);
    println!("is_greater: {}, is_equal: {}", is_greater, is_equal);
}
```

### Character Type

The `char` type represents a Unicode Scalar Value:

```rust
fn main() {
    let c = 'z';
    let z: char = 'Z';
    let heart_eyed_cat = '😻';
    let chinese = '中';

    println!("{} {} {} {}", c, z, heart_eyed_cat, chinese);
}
```

### String Types: `String` and `&str`

This is often confusing for beginners—Rust has two main string types:

```rust
fn main() {
    // &str - string slice (borrowed, immutable)
    let s1: &str = "Hello, world!";

    // String - owned, growable, heap-allocated
    let mut s2 = String::from("Hello");
    s2.push_str(", world!");

    println!("s1: {}", s1);
    println!("s2: {}", s2);

    // Converting between them
    let s3: String = s1.to_string();
    let s4: &str = &s2;

    // Creating strings
    let s5 = "literal".to_string();
    let s6 = String::new();
    let s7 = format!("{} and {}", s1, s2);

    println!("s7: {}", s7);
}
```

### Compound Types: Tuples and Arrays

```rust
fn main() {
    // Tuple - fixed length, different types allowed
    let tup: (i32, f64, u8) = (500, 6.4, 1);
    let (x, y, z) = tup;  // Destructuring

    println!("x: {}, y: {}, z: {}", x, y, z);
    println!("First: {}, Second: {}", tup.0, tup.1);

    // Array - fixed length, same type
    let arr = [1, 2, 3, 4, 5];
    let months = ["Jan", "Feb", "Mar"];

    // Array with type and length
    let zeros: [i32; 5] = [0; 5];  // [0, 0, 0, 0, 0]

    println!("First month: {}", months[0]);
    println!("Array length: {}", arr.len());
}
```

---

## Functions

### Basic Function Syntax

```rust
fn main() {
    println!("Hello from main!");
    another_function();
    greet("Alice");
}

fn another_function() {
    println!("Hello from another function!");
}

fn greet(name: &str) {
    println!("Hello, {}!", name);
}
```

### Parameters

Functions can have multiple parameters:

```rust
fn main() {
    print_labeled_measurement(5, 'h');
    let result = add(5, 3);
    println!("5 + 3 = {}", result);
}

fn print_labeled_measurement(value: i32, unit_label: char) {
    println!("The measurement is: {}{}", value, unit_label);
}

fn add(a: i32, b: i32) -> i32 {
    a + b  // Note: no semicolon = return value
}
```

### Return Values

```rust
fn main() {
    let x = five();
    let y = plus_one(x);

    println!("x: {}, y: {}", x, y);
}

fn five() -> i32 {
    5  // Implicit return (no semicolon)
}

fn plus_one(x: i32) -> i32 {
    x + 1
}

// Multiple return values using tuples
fn swap(a: i32, b: i32) -> (i32, i32) {
    (b, a)
}

// Early return with `return` keyword
fn absolute_value(x: i32) -> i32 {
    if x < 0 {
        return -x;
    }
    x
}
```

---

## Control Flow

### if/else

```rust
fn main() {
    let number = 6;

    if number % 4 == 0 {
        println!("number is divisible by 4");
    } else if number % 3 == 0 {
        println!("number is divisible by 3");
    } else if number % 2 == 0 {
        println!("number is divisible by 2");
    } else {
        println!("number is not divisible by 4, 3, or 2");
    }

    // if as an expression
    let condition = true;
    let number = if condition { 5 } else { 6 };
    println!("The value is: {}", number);
}
```

### loop

The `loop` keyword creates an infinite loop:

```rust
fn main() {
    let mut counter = 0;

    let result = loop {
        counter += 1;

        if counter == 10 {
            break counter * 2;  // Return value from loop
        }
    };

    println!("The result is: {}", result);
}
```

### Loop Labels

```rust
fn main() {
    let mut count = 0;

    'counting_up: loop {
        println!("count = {}", count);
        let mut remaining = 10;

        loop {
            println!("remaining = {}", remaining);
            if remaining == 9 {
                break;
            }
            if count == 2 {
                break 'counting_up;  // Break outer loop
            }
            remaining -= 1;
        }

        count += 1;
    }

    println!("End count = {}", count);
}
```

### while

```rust
fn main() {
    let mut number = 3;

    while number != 0 {
        println!("{}!", number);
        number -= 1;
    }

    println!("LIFTOFF!");
}
```

### for

The `for` loop is the most commonly used loop in Rust:

```rust
fn main() {
    let a = [10, 20, 30, 40, 50];

    // Iterate over array
    for element in a {
        println!("The value is: {}", element);
    }

    // Range
    for number in 1..4 {
        println!("{}", number);  // 1, 2, 3
    }

    // Inclusive range
    for number in 1..=3 {
        println!("{}", number);  // 1, 2, 3
    }

    // Reverse
    for number in (1..4).rev() {
        println!("{}!", number);  // 3, 2, 1
    }

    // With index
    for (index, value) in a.iter().enumerate() {
        println!("Index {}: {}", index, value);
    }
}
```

---

## Expressions vs Statements

This is a crucial concept in Rust:

- **Statements** perform actions but don't return values
- **Expressions** evaluate to a value

```rust
fn main() {
    // Statement - doesn't return a value
    let x = 5;  // `let x = 5` is a statement

    // This won't work:
    // let x = (let y = 6);  // Error!

    // Expression - returns a value
    let y = {
        let x = 3;
        x + 1  // No semicolon = expression returns this value
    };

    println!("y = {}", y);  // y = 4

    // Blocks are expressions
    let z = if true { 5 } else { 6 };

    // Function calls are expressions
    let a = add(1, 2);

    // Math operations are expressions
    let b = 5 + 3;
}

fn add(a: i32, b: i32) -> i32 {
    a + b  // Expression (no semicolon)
}
```

**Important:** Adding a semicolon turns an expression into a statement:

```rust
fn returns_nothing() {
    5;  // This is now a statement, function returns ()
}

fn returns_five() -> i32 {
    5   // Expression, function returns 5
}
```

---

## Comments and Documentation

### Regular Comments

```rust
fn main() {
    // This is a single-line comment

    let lucky_number = 7; // You can also comment at end of line

    /*
     * This is a multi-line comment.
     * Though single-line comments are preferred in Rust.
     */
}
```

### Documentation Comments

Documentation comments generate HTML documentation:

```rust
/// Adds two numbers together.
///
/// # Arguments
///
/// * `a` - The first number
/// * `b` - The second number
///
/// # Returns
///
/// The sum of `a` and `b`
///
/// # Examples
///
/// ```
/// let result = add(2, 3);
/// assert_eq!(result, 5);
/// ```
fn add(a: i32, b: i32) -> i32 {
    a + b
}

//! This is a module-level documentation comment.
//! It documents the entire module or crate.
```

---

## Putting It All Together

Here's a more complete example combining what we've learned:

```rust
//! A simple temperature converter program

/// Temperature scale enumeration
const FREEZING_POINT_F: f64 = 32.0;

/// Converts Fahrenheit to Celsius
///
/// # Examples
///
/// ```
/// let celsius = fahrenheit_to_celsius(32.0);
/// assert_eq!(celsius, 0.0);
/// ```
fn fahrenheit_to_celsius(f: f64) -> f64 {
    (f - 32.0) * 5.0 / 9.0
}

/// Converts Celsius to Fahrenheit
fn celsius_to_fahrenheit(c: f64) -> f64 {
    c * 9.0 / 5.0 + 32.0
}

fn main() {
    let temperatures_f = [32.0, 68.0, 100.0, 212.0];

    println!("Fahrenheit to Celsius Conversion");
    println!("================================");

    for temp in temperatures_f {
        let celsius = fahrenheit_to_celsius(temp);

        let description = if celsius <= 0.0 {
            "freezing"
        } else if celsius < 20.0 {
            "cold"
        } else if celsius < 30.0 {
            "comfortable"
        } else {
            "hot"
        };

        println!("{:.1}°F = {:.1}°C ({})", temp, celsius, description);
    }

    // Using a loop with counter
    println!("\nCounting down from boiling point:");
    let mut temp = 100.0;

    while temp >= 0.0 {
        let f = celsius_to_fahrenheit(temp);
        println!("{:.0}°C = {:.1}°F", temp, f);
        temp -= 25.0;
    }
}
```

---

## Exercises

### Exercise 1: Variable Practice
Create a program that:
1. Declares an immutable variable with your age
2. Declares a mutable counter starting at 0
3. Uses shadowing to transform a string "42" into the integer 42
4. Defines a constant for the number of days in a week

### Exercise 2: Temperature Converter
Extend the temperature converter to:
1. Also support Kelvin conversions
2. Create a function that takes any temperature and its scale, returning all three scales
3. Use tuples for the return value

### Exercise 3: FizzBuzz
Write the classic FizzBuzz program:
- For numbers 1 to 100
- Print "Fizz" for multiples of 3
- Print "Buzz" for multiples of 5
- Print "FizzBuzz" for multiples of both
- Otherwise print the number

### Exercise 4: Fibonacci
Write a function that returns the nth Fibonacci number:
1. First using a loop
2. Then using recursion (bonus: discuss why recursion is less efficient here)

### Exercise 5: Array Statistics
Write a program that:
1. Takes an array of integers
2. Calculates and returns the sum, average, minimum, and maximum
3. Return these as a tuple or struct (if you've read ahead!)

---

## Common Mistakes to Avoid

1. **Forgetting `mut`** for variables you need to change
2. **Adding semicolons** to the last expression in a function that should return a value
3. **Confusing `String` and `&str`** - we'll cover this more in ownership
4. **Using `=` instead of `==`** in conditions (though Rust will catch this!)

---

## Summary

You've learned:
- Variables are immutable by default in Rust
- Rust has strong, static typing with type inference
- Functions use `-> Type` for return types
- The last expression without a semicolon is the return value
- Control flow constructs are expressions in Rust

---

## Next Steps

Now that you understand the basics, it's time to learn about Rust's most distinctive feature: [Ownership and Borrowing](./02-ownership-borrowing.md).

This is where Rust really differs from other languages, and understanding it is key to becoming productive in Rust. Don't worry—take your time with it!

# Ownership and Borrowing

## Learning Objectives

By the end of this tutorial, you will be able to:
- Understand Rust's ownership system and why it exists
- Apply the three rules of ownership
- Distinguish between move, copy, and clone semantics
- Use references and borrowing effectively
- Work with mutable and immutable references
- Understand and use slices
- Write code that satisfies the borrow checker

---

## Introduction

Ownership is Rust's most unique feature and what makes it different from every other programming language. It enables Rust to make memory safety guarantees without needing a garbage collector.

**This is THE most important concept in Rust.** Take your time with this chapter. Re-read sections if needed. Once ownership clicks, everything else becomes much easier.

Don't be discouraged if the borrow checker fights you at first—it's teaching you to think about memory in a way that prevents entire classes of bugs.

---

## Why Ownership?

In other languages, you either:
- Manually manage memory (C/C++) - error-prone, leads to bugs
- Use a garbage collector (Java, Go, Python) - runtime overhead

Rust takes a third approach: the compiler tracks ownership and automatically frees memory when data goes out of scope. Zero runtime cost, maximum safety.

---

## The Three Rules of Ownership

Memorize these rules:

1. **Each value in Rust has a variable that's called its owner.**
2. **There can only be one owner at a time.**
3. **When the owner goes out of scope, the value will be dropped.**

Let's explore each rule with examples.

---

## Rule 1: Each Value Has an Owner

```rust
fn main() {
    let s = String::from("hello");  // s is the owner of this String
    //  ^-- owner                       ^-- value
}
```

The variable `s` owns the String value. When we talk about "the owner," we mean the variable that is responsible for cleaning up that memory.

---

## Rule 2: One Owner at a Time

This is where things get interesting. When you assign a value to another variable, ownership moves:

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1;  // Ownership MOVES from s1 to s2

    // println!("{}", s1);  // ERROR! s1 is no longer valid
    println!("{}", s2);     // OK! s2 is the owner now
}
```

This is called a **move**. After the move, `s1` is no longer valid.

### Why Does This Happen?

Consider what's in memory:

```
Stack:              Heap:
s1: [ptr|len|cap]  -->  "hello"
```

If we simply copied the pointer (like in C), both `s1` and `s2` would point to the same heap data. When both go out of scope, Rust would try to free the same memory twice—a double free error!

Instead, Rust invalidates `s1` after the move:

```
Stack:              Heap:
s1: [invalid]
s2: [ptr|len|cap]  -->  "hello"
```

---

## Rule 3: Drop When Out of Scope

```rust
fn main() {
    {
        let s = String::from("hello");  // s comes into scope
        // do stuff with s
    }  // s goes out of scope, memory is freed automatically

    // s is not accessible here
}
```

Rust automatically calls a special function called `drop` when a variable goes out of scope. This is where the memory is freed.

---

## Move Semantics in Detail

### Moves Happen in Assignments

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1;  // s1 is moved to s2

    // s1 is now invalid
}
```

### Moves Happen When Passing to Functions

```rust
fn main() {
    let s = String::from("hello");
    takes_ownership(s);  // s is moved into the function

    // println!("{}", s);  // ERROR! s is no longer valid
}

fn takes_ownership(some_string: String) {
    println!("{}", some_string);
}  // some_string goes out of scope and is dropped
```

### Moves Happen When Returning from Functions

```rust
fn main() {
    let s1 = gives_ownership();         // Return value moves into s1
    let s2 = String::from("hello");
    let s3 = takes_and_gives_back(s2);  // s2 moves in, return moves into s3

    println!("s1: {}, s3: {}", s1, s3);
    // s2 is invalid here
}

fn gives_ownership() -> String {
    let some_string = String::from("yours");
    some_string  // Returned and moves out
}

fn takes_and_gives_back(a_string: String) -> String {
    a_string  // Returned and moves out
}
```

---

## Clone vs Copy

### Clone: Explicit Deep Copy

If you actually want to copy the heap data, use `clone()`:

```rust
fn main() {
    let s1 = String::from("hello");
    let s2 = s1.clone();  // Deep copy of heap data

    println!("s1: {}, s2: {}", s1, s2);  // Both are valid!
}
```

`clone()` can be expensive for large data structures, so Rust makes you call it explicitly.

### Copy: Implicit Stack Copy

Simple types that live entirely on the stack implement the `Copy` trait:

```rust
fn main() {
    let x = 5;
    let y = x;  // Copy, not move!

    println!("x: {}, y: {}", x, y);  // Both are valid!
}
```

Types that implement `Copy`:
- All integer types (i32, u64, etc.)
- Boolean type (bool)
- Floating point types (f32, f64)
- Character type (char)
- Tuples containing only Copy types

```rust
fn main() {
    // These all implement Copy
    let a: i32 = 5;
    let b = a;  // Copy

    let c: bool = true;
    let d = c;  // Copy

    let e: (i32, f64) = (1, 2.0);
    let f = e;  // Copy

    // This does NOT implement Copy
    let g: (i32, String) = (1, String::from("hello"));
    let h = g;  // Move!
    // println!("{:?}", g);  // ERROR!
}
```

---

## References and Borrowing

Moving ownership everywhere gets tedious. What if we want to use a value without taking ownership? We use **references**.

### Creating References

```rust
fn main() {
    let s1 = String::from("hello");
    let len = calculate_length(&s1);  // &s1 creates a reference

    println!("The length of '{}' is {}.", s1, len);  // s1 is still valid!
}

fn calculate_length(s: &String) -> usize {
    s.len()
}  // s goes out of scope, but doesn't drop what it refers to
```

The `&` symbol creates a reference. This is called **borrowing**—we're borrowing the value without taking ownership.

### Visual Representation

```
Stack:                      Heap:
s1: [ptr|len|cap]  ------>  "hello"
         ^
         |
s:  [ptr]  (reference to s1)
```

### References are Immutable by Default

```rust
fn main() {
    let s = String::from("hello");
    change(&s);
}

fn change(some_string: &String) {
    // some_string.push_str(", world");  // ERROR! Cannot mutate
}
```

---

## Mutable References

To modify borrowed data, use mutable references:

```rust
fn main() {
    let mut s = String::from("hello");  // s must be mut
    change(&mut s);                      // Pass mutable reference

    println!("{}", s);  // "hello, world"
}

fn change(some_string: &mut String) {
    some_string.push_str(", world");
}
```

---

## The Rules of References

Here's where Rust's safety guarantees come from:

### Rule 1: One Mutable Reference OR Any Number of Immutable References

You can have either:
- **One** mutable reference, OR
- **Any number** of immutable references

But not both at the same time!

```rust
fn main() {
    let mut s = String::from("hello");

    // Multiple immutable references - OK!
    let r1 = &s;
    let r2 = &s;
    println!("{} and {}", r1, r2);

    // Now r1 and r2 are no longer used...

    // One mutable reference - OK!
    let r3 = &mut s;
    r3.push_str(" world");
    println!("{}", r3);
}
```

This won't compile:

```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &s;      // immutable borrow
    let r2 = &mut s;  // ERROR! Can't borrow as mutable

    println!("{}, {}", r1, r2);
}
```

And neither will this:

```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &mut s;
    let r2 = &mut s;  // ERROR! Can't have two mutable references

    println!("{}, {}", r1, r2);
}
```

### Why This Rule?

This prevents data races at compile time! A data race occurs when:
- Two or more pointers access the same data
- At least one pointer writes to the data
- No synchronization mechanism

Rust prevents this entirely through the borrow checker.

### Rule 2: References Must Always Be Valid

References must never outlive the data they point to:

```rust
fn main() {
    let reference_to_nothing = dangle();
}

fn dangle() -> &String {  // ERROR!
    let s = String::from("hello");
    &s  // Reference to s
}  // s is dropped here, reference would be invalid!
```

The compiler prevents dangling references. The solution is to return the owned value:

```rust
fn no_dangle() -> String {
    let s = String::from("hello");
    s  // Ownership is moved out
}
```

---

## Non-Lexical Lifetimes (NLL)

Rust is smart about when references are actually used:

```rust
fn main() {
    let mut s = String::from("hello");

    let r1 = &s;
    let r2 = &s;
    println!("{} and {}", r1, r2);
    // r1 and r2 are no longer used after this point

    let r3 = &mut s;  // OK! r1 and r2's scope ended
    println!("{}", r3);
}
```

The compiler sees that `r1` and `r2` aren't used after the `println!`, so their scope ends there.

---

## Slices

Slices let you reference a contiguous sequence of elements without owning them.

### String Slices

```rust
fn main() {
    let s = String::from("hello world");

    let hello = &s[0..5];   // "hello"
    let world = &s[6..11];  // "world"

    println!("{} {}", hello, world);

    // Shorthand syntax
    let hello = &s[..5];    // Start from beginning
    let world = &s[6..];    // Go to end
    let whole = &s[..];     // Entire string

    println!("{}", whole);
}
```

### String Slices are `&str`

```rust
fn main() {
    let s = String::from("hello world");
    let word = first_word(&s);

    println!("First word: {}", word);
}

fn first_word(s: &String) -> &str {
    let bytes = s.as_bytes();

    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[..i];
        }
    }

    &s[..]
}
```

### Better: Accept `&str` in Functions

```rust
fn main() {
    let my_string = String::from("hello world");
    let word = first_word(&my_string);

    let my_string_literal = "hello world";
    let word = first_word(my_string_literal);  // Also works!
}

fn first_word(s: &str) -> &str {  // Accept &str instead of &String
    let bytes = s.as_bytes();

    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[..i];
        }
    }

    &s[..]
}
```

### Array Slices

```rust
fn main() {
    let a = [1, 2, 3, 4, 5];

    let slice = &a[1..3];  // [2, 3]

    assert_eq!(slice, &[2, 3]);

    // Function that takes a slice
    let sum = sum_array(&a);
    let partial_sum = sum_array(&a[0..3]);

    println!("Sum: {}, Partial: {}", sum, partial_sum);
}

fn sum_array(arr: &[i32]) -> i32 {
    let mut sum = 0;
    for &item in arr {
        sum += item;
    }
    sum
}
```

---

## Common Patterns

### Pattern 1: Borrowing for Read-Only Access

```rust
fn main() {
    let data = vec![1, 2, 3, 4, 5];

    // Borrow for reading
    print_vec(&data);
    let sum = sum_vec(&data);

    println!("Sum: {}", sum);
    println!("Data still accessible: {:?}", data);
}

fn print_vec(v: &Vec<i32>) {
    for item in v {
        println!("{}", item);
    }
}

fn sum_vec(v: &Vec<i32>) -> i32 {
    v.iter().sum()
}
```

### Pattern 2: Mutable Borrow for Modification

```rust
fn main() {
    let mut data = vec![1, 2, 3];

    double_values(&mut data);

    println!("{:?}", data);  // [2, 4, 6]
}

fn double_values(v: &mut Vec<i32>) {
    for item in v.iter_mut() {
        *item *= 2;
    }
}
```

### Pattern 3: Returning Owned Data

```rust
fn main() {
    let data = create_data();
    let processed = process_data(data);  // Ownership moves

    println!("{:?}", processed);
}

fn create_data() -> Vec<i32> {
    vec![1, 2, 3, 4, 5]
}

fn process_data(mut v: Vec<i32>) -> Vec<i32> {
    v.push(6);
    v
}
```

### Pattern 4: Builder Pattern

```rust
fn main() {
    let text = String::new()
        .add_word("Hello")
        .add_word("World");

    println!("{}", text);
}

trait StringBuilder {
    fn add_word(self, word: &str) -> Self;
}

impl StringBuilder for String {
    fn add_word(mut self, word: &str) -> Self {
        if !self.is_empty() {
            self.push(' ');
        }
        self.push_str(word);
        self
    }
}
```

---

## Common Borrow Checker Errors and Solutions

### Error: Cannot Move Out of Borrowed Content

```rust
// ERROR
fn main() {
    let v = vec![String::from("hello")];
    let s = v[0];  // ERROR: cannot move
}

// Solution 1: Clone
fn main() {
    let v = vec![String::from("hello")];
    let s = v[0].clone();
}

// Solution 2: Reference
fn main() {
    let v = vec![String::from("hello")];
    let s = &v[0];
}
```

### Error: Cannot Borrow as Mutable More Than Once

```rust
// ERROR
fn main() {
    let mut v = vec![1, 2, 3];
    let first = &mut v[0];
    let second = &mut v[1];  // ERROR!
    *first = 10;
    *second = 20;
}

// Solution: Use indices
fn main() {
    let mut v = vec![1, 2, 3];
    v[0] = 10;
    v[1] = 20;
}

// Or split_at_mut for simultaneous mutable access
fn main() {
    let mut v = vec![1, 2, 3];
    let (first_half, second_half) = v.split_at_mut(1);
    first_half[0] = 10;
    second_half[0] = 20;
}
```

### Error: Borrowed Value Does Not Live Long Enough

```rust
// ERROR
fn main() {
    let r;
    {
        let x = 5;
        r = &x;  // ERROR: x doesn't live long enough
    }
    println!("{}", r);
}

// Solution: Extend lifetime
fn main() {
    let x = 5;
    let r = &x;
    println!("{}", r);
}
```

---

## Real-World Example: A Simple Word Counter

```rust
use std::collections::HashMap;

fn main() {
    let text = "hello world wonderful world";

    let word_counts = count_words(text);

    for (word, count) in &word_counts {
        println!("{}: {}", word, count);
    }
}

fn count_words(text: &str) -> HashMap<&str, u32> {
    let mut counts = HashMap::new();

    for word in text.split_whitespace() {
        let count = counts.entry(word).or_insert(0);
        *count += 1;
    }

    counts
}
```

Note how we:
- Take `&str` as input (borrowing)
- Return `HashMap<&str, u32>` with references to the original text
- Use `entry` API which handles borrowing correctly

---

## Exercises

### Exercise 1: Fix the Borrow Checker Errors
Fix each of these programs:

```rust
// Problem 1
fn main() {
    let s = String::from("hello");
    let s2 = s;
    println!("{}", s);
}

// Problem 2
fn main() {
    let mut s = String::from("hello");
    let r1 = &s;
    let r2 = &mut s;
    println!("{} {}", r1, r2);
}

// Problem 3
fn main() {
    let r;
    {
        let s = String::from("hello");
        r = &s;
    }
    println!("{}", r);
}
```

### Exercise 2: Implement `second_word`
Write a function that returns the second word in a string:

```rust
fn second_word(s: &str) -> &str {
    // Your implementation
}

fn main() {
    let s = String::from("hello world");
    let word = second_word(&s);
    println!("Second word: {}", word);
}
```

### Exercise 3: String Manipulation Without Taking Ownership
Write three functions that:
1. Count vowels in a string (borrow immutably)
2. Remove all spaces from a string (need mutable access)
3. Reverse a string (return a new String)

### Exercise 4: Vector Operations
Implement these functions using proper borrowing:

```rust
fn find_max(v: /* what type? */) -> Option<&i32> {
    // Find maximum value
}

fn find_and_double_max(v: /* what type? */) {
    // Find max and double it in place
}

fn main() {
    let mut numbers = vec![1, 5, 3, 2, 4];

    if let Some(max) = find_max(&numbers) {
        println!("Max: {}", max);
    }

    find_and_double_max(&mut numbers);
    println!("{:?}", numbers);
}
```

### Exercise 5: Struct with References
Create a struct that holds references and ensure lifetimes work:

```rust
struct TextExcerpt<'a> {
    content: &'a str,
    line_number: u32,
}

// Implement:
// 1. A function to create TextExcerpt from a line
// 2. A method to check if excerpt contains a word
// 3. A method to return excerpt length

fn main() {
    let full_text = String::from("Hello, world!");
    let excerpt = TextExcerpt {
        content: &full_text[0..5],
        line_number: 1,
    };
    // Use your implementations
}
```

---

## Tips for Working with the Borrow Checker

1. **Start with owned data** - Use `String` instead of `&str` when learning
2. **Clone when stuck** - It's not always optimal, but it works
3. **Read the error messages** - Rust's compiler errors are excellent
4. **Draw ownership diagrams** - Visualize who owns what
5. **Use smaller scopes** - End borrows early by limiting scope
6. **Consider returning owned data** - Sometimes moving is cleaner than borrowing

---

## Summary

You've learned:
- Ownership is Rust's system for memory management without GC
- Every value has exactly one owner
- When the owner goes out of scope, the value is dropped
- Move semantics transfer ownership
- References allow borrowing without ownership
- You can have either one mutable reference OR many immutable references
- References must always be valid
- Slices are references to portions of data

This is the foundation of Rust. Everything else builds on these concepts.

---

## Next Steps

Now that you understand ownership and borrowing, let's learn how to structure our data with [Structs and Enums](./03-structs-enums.md).

Remember: If the borrow checker is giving you trouble, that's normal! It's catching potential bugs. Take a break, draw out the ownership, and try again.

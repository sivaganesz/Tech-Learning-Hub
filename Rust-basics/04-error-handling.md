# Error Handling

## Learning Objectives

By the end of this tutorial, you will be able to:
- Distinguish between recoverable and unrecoverable errors
- Use `Result<T, E>` for recoverable errors
- Use `Option<T>` for optional values
- Apply the `?` operator for error propagation
- Use `unwrap()` and `expect()` appropriately
- Create custom error types
- Use the `thiserror` crate for cleaner error definitions
- Follow Rust error handling best practices

---

## Introduction

Error handling is one of Rust's strengths. Instead of exceptions that can be thrown anywhere, Rust makes you think about errors upfront. This leads to more robust and reliable programs.

Rust groups errors into two categories:
- **Recoverable errors** (wrong file name, network timeout) - use `Result<T, E>`
- **Unrecoverable errors** (bugs like index out of bounds) - use `panic!`

---

## Panic: Unrecoverable Errors

When something goes terribly wrong and there's no way to continue:

```rust
fn main() {
    panic!("crash and burn");
}
```

Output:
```
thread 'main' panicked at 'crash and burn', src/main.rs:2:5
```

### When Panics Occur

```rust
fn main() {
    let v = vec![1, 2, 3];
    v[99];  // Panic! index out of bounds
}
```

### Viewing the Backtrace

```bash
RUST_BACKTRACE=1 cargo run
```

### When to Use `panic!`

- In tests (failing a test)
- In examples and prototypes
- When a situation is truly unrecoverable
- When you have information the compiler doesn't

```rust
fn main() {
    // We know this won't fail, but the compiler doesn't
    let home: std::net::IpAddr = "127.0.0.1"
        .parse()
        .expect("Hardcoded IP address should be valid");
}
```

---

## Result: Recoverable Errors

Most errors should be handled, not cause a crash:

```rust
enum Result<T, E> {
    Ok(T),
    Err(E),
}
```

### Basic Usage

```rust
use std::fs::File;

fn main() {
    let greeting_file = File::open("hello.txt");

    let greeting_file = match greeting_file {
        Ok(file) => file,
        Err(error) => panic!("Problem opening the file: {:?}", error),
    };
}
```

### Handling Different Errors

```rust
use std::fs::File;
use std::io::ErrorKind;

fn main() {
    let greeting_file = File::open("hello.txt");

    let greeting_file = match greeting_file {
        Ok(file) => file,
        Err(error) => match error.kind() {
            ErrorKind::NotFound => match File::create("hello.txt") {
                Ok(fc) => fc,
                Err(e) => panic!("Problem creating the file: {:?}", e),
            },
            other_error => {
                panic!("Problem opening the file: {:?}", other_error);
            }
        },
    };
}
```

### Cleaner with Closures

```rust
use std::fs::File;
use std::io::ErrorKind;

fn main() {
    let greeting_file = File::open("hello.txt").unwrap_or_else(|error| {
        if error.kind() == ErrorKind::NotFound {
            File::create("hello.txt").unwrap_or_else(|error| {
                panic!("Problem creating the file: {:?}", error);
            })
        } else {
            panic!("Problem opening the file: {:?}", error);
        }
    });
}
```

---

## Option: Representing Absence

`Option<T>` represents a value that might not exist:

```rust
enum Option<T> {
    None,
    Some(T),
}
```

### Common Uses

```rust
fn main() {
    // Finding in a collection
    let numbers = vec![1, 2, 3, 4, 5];
    let third: Option<&i32> = numbers.get(2);

    match third {
        Some(n) => println!("Third number: {}", n),
        None => println!("No third number"),
    }

    // String operations
    let s = "hello";
    let first_char: Option<char> = s.chars().next();

    if let Some(c) = first_char {
        println!("First char: {}", c);
    }
}
```

### Converting Between Option and Result

```rust
fn main() {
    let opt: Option<i32> = Some(5);

    // Option to Result
    let res: Result<i32, &str> = opt.ok_or("No value");

    // Result to Option
    let back: Option<i32> = res.ok();

    println!("{:?}", back);  // Some(5)
}
```

---

## The `?` Operator

The `?` operator is Rust's way of propagating errors elegantly:

### With Result

```rust
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let mut username_file = File::open("hello.txt")?;
    let mut username = String::new();
    username_file.read_to_string(&mut username)?;
    Ok(username)
}
```

The `?` operator:
1. If `Ok`, unwraps and continues
2. If `Err`, returns early from the function with the error

### Chaining with `?`

```rust
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let mut username = String::new();
    File::open("hello.txt")?.read_to_string(&mut username)?;
    Ok(username)
}
```

### Even Shorter

```rust
use std::fs;
use std::io;

fn read_username_from_file() -> Result<String, io::Error> {
    fs::read_to_string("hello.txt")
}
```

### With Option

```rust
fn last_char_of_first_line(text: &str) -> Option<char> {
    text.lines().next()?.chars().last()
}

fn main() {
    assert_eq!(last_char_of_first_line("Hello\nWorld"), Some('o'));
    assert_eq!(last_char_of_first_line(""), None);
    assert_eq!(last_char_of_first_line("\nhi"), None);
}
```

### `?` in `main`

```rust
use std::error::Error;
use std::fs::File;

fn main() -> Result<(), Box<dyn Error>> {
    let greeting_file = File::open("hello.txt")?;
    Ok(())
}
```

---

## Unwrap and Expect

### `unwrap()`: Quick and Dangerous

```rust
fn main() {
    let greeting_file = File::open("hello.txt").unwrap();
    // Panics if Err
}
```

### `expect()`: Unwrap with Context

```rust
fn main() {
    let greeting_file = File::open("hello.txt")
        .expect("hello.txt should be included in this project");
    // Panics with your message if Err
}
```

### When to Use Them

- **Prototyping**: Quick and dirty code
- **Tests**: Expected to panic on failure
- **When you know it won't fail**: And want to document why

```rust
// OK: We control the input
let num: i32 = "42".parse().expect("Static string should parse");

// NOT OK: User input could be anything
let user_input = "not a number";
let num: i32 = user_input.parse().unwrap(); // Don't do this!
```

### Safer Alternatives

```rust
fn main() {
    // unwrap_or - provide default
    let value = Some(5).unwrap_or(0);

    // unwrap_or_default - use type's default
    let value: i32 = None.unwrap_or_default();  // 0

    // unwrap_or_else - compute default lazily
    let value = Some(5).unwrap_or_else(|| expensive_computation());

    // map - transform the inner value
    let value: Option<i32> = Some("5").map(|s| s.len() as i32);
}

fn expensive_computation() -> i32 { 42 }
```

---

## Custom Error Types

### Simple Custom Error

```rust
use std::fmt;

#[derive(Debug)]
struct AppError {
    message: String,
}

impl fmt::Display for AppError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for AppError {}

fn do_something() -> Result<(), AppError> {
    Err(AppError {
        message: String::from("Something went wrong"),
    })
}
```

### Enum-Based Errors

```rust
use std::fmt;
use std::io;
use std::num::ParseIntError;

#[derive(Debug)]
enum ConfigError {
    IoError(io::Error),
    ParseError(ParseIntError),
    MissingField(String),
    InvalidValue { field: String, value: String },
}

impl fmt::Display for ConfigError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            ConfigError::IoError(e) => write!(f, "IO error: {}", e),
            ConfigError::ParseError(e) => write!(f, "Parse error: {}", e),
            ConfigError::MissingField(field) => write!(f, "Missing field: {}", field),
            ConfigError::InvalidValue { field, value } => {
                write!(f, "Invalid value '{}' for field '{}'", value, field)
            }
        }
    }
}

impl std::error::Error for ConfigError {}

// Convert from io::Error
impl From<io::Error> for ConfigError {
    fn from(error: io::Error) -> Self {
        ConfigError::IoError(error)
    }
}

// Convert from ParseIntError
impl From<ParseIntError> for ConfigError {
    fn from(error: ParseIntError) -> Self {
        ConfigError::ParseError(error)
    }
}
```

Now we can use `?` with different error types:

```rust
use std::fs;

fn load_config(path: &str) -> Result<i32, ConfigError> {
    let content = fs::read_to_string(path)?;  // io::Error -> ConfigError
    let value: i32 = content.trim().parse()?; // ParseIntError -> ConfigError
    Ok(value)
}
```

---

## The `thiserror` Crate

The `thiserror` crate makes defining custom errors much cleaner:

```toml
# Cargo.toml
[dependencies]
thiserror = "1.0"
```

```rust
use thiserror::Error;

#[derive(Error, Debug)]
enum DataError {
    #[error("Failed to read file: {0}")]
    IoError(#[from] std::io::Error),

    #[error("Failed to parse data: {0}")]
    ParseError(#[from] std::num::ParseIntError),

    #[error("Missing required field: {0}")]
    MissingField(String),

    #[error("Invalid value '{value}' for field '{field}'")]
    InvalidValue { field: String, value: String },

    #[error("Database error: {0}")]
    DatabaseError(String),
}
```

### Using thiserror

```rust
use std::fs;
use thiserror::Error;

#[derive(Error, Debug)]
enum AppError {
    #[error("Configuration error: {0}")]
    Config(String),

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("Not found: {0}")]
    NotFound(String),
}

fn read_config(path: &str) -> Result<String, AppError> {
    let content = fs::read_to_string(path)?;  // Auto-converts io::Error

    if content.is_empty() {
        return Err(AppError::Config("Config file is empty".to_string()));
    }

    Ok(content)
}

fn main() {
    match read_config("config.txt") {
        Ok(config) => println!("Config: {}", config),
        Err(e) => eprintln!("Error: {}", e),
    }
}
```

---

## Error Handling Patterns

### Pattern 1: Early Return

```rust
fn process_data(data: Option<&str>) -> Result<i32, String> {
    let data = match data {
        Some(d) => d,
        None => return Err("No data provided".to_string()),
    };

    let value: i32 = data
        .parse()
        .map_err(|_| "Invalid number".to_string())?;

    if value < 0 {
        return Err("Value must be positive".to_string());
    }

    Ok(value * 2)
}
```

### Pattern 2: Combining Results

```rust
fn fetch_user_data(id: u32) -> Result<String, String> {
    // Simulated fetches
    let name = fetch_name(id)?;
    let email = fetch_email(id)?;
    let age = fetch_age(id)?;

    Ok(format!("{} ({}) - {} years old", name, email, age))
}

fn fetch_name(id: u32) -> Result<String, String> {
    Ok(format!("User{}", id))
}

fn fetch_email(id: u32) -> Result<String, String> {
    Ok(format!("user{}@example.com", id))
}

fn fetch_age(id: u32) -> Result<u32, String> {
    Ok(20 + id)
}
```

### Pattern 3: Collecting Results

```rust
fn parse_numbers(inputs: Vec<&str>) -> Result<Vec<i32>, std::num::ParseIntError> {
    inputs.iter().map(|s| s.parse::<i32>()).collect()
}

fn main() {
    let good = vec!["1", "2", "3"];
    let bad = vec!["1", "two", "3"];

    println!("{:?}", parse_numbers(good));  // Ok([1, 2, 3])
    println!("{:?}", parse_numbers(bad));   // Err(ParseIntError)
}
```

### Pattern 4: Providing Context

```rust
use std::fs;
use std::io;

fn read_config_file(path: &str) -> Result<String, String> {
    fs::read_to_string(path)
        .map_err(|e| format!("Failed to read config from '{}': {}", path, e))
}
```

Or with the `anyhow` crate for more complex applications:

```rust
// Using anyhow for applications
use anyhow::{Context, Result};
use std::fs;

fn read_config(path: &str) -> Result<String> {
    fs::read_to_string(path)
        .with_context(|| format!("Failed to read config from '{}'", path))
}
```

---

## Best Practices

### 1. Don't Overuse `unwrap()`

```rust
// Bad
let file = File::open("config.txt").unwrap();

// Good
let file = File::open("config.txt")
    .expect("config.txt must exist in the project root");

// Better
let file = File::open("config.txt")?;
```

### 2. Provide Meaningful Error Messages

```rust
// Bad
Err("Error")

// Good
Err(format!("Failed to connect to database at {}: {}", host, e))
```

### 3. Use Type-Safe Error Handling

```rust
// Instead of String errors in production code
fn parse_config(s: &str) -> Result<Config, String>  // Avoid

// Use typed errors
fn parse_config(s: &str) -> Result<Config, ConfigError>  // Better
```

### 4. Handle Errors at the Right Level

```rust
// Low-level function: propagate errors
fn read_file(path: &str) -> Result<String, io::Error> {
    fs::read_to_string(path)
}

// Mid-level function: wrap errors with context
fn load_config(path: &str) -> Result<Config, AppError> {
    let content = read_file(path)?;
    parse_config(&content)
}

// High-level function: handle errors
fn main() {
    match load_config("app.toml") {
        Ok(config) => run_app(config),
        Err(e) => {
            eprintln!("Failed to load config: {}", e);
            std::process::exit(1);
        }
    }
}
```

### 5. Use `Option` for Absence, `Result` for Errors

```rust
// Good: None means "not found"
fn find_user(id: u32) -> Option<User>

// Good: Err means "something went wrong"
fn load_user(id: u32) -> Result<User, DatabaseError>
```

---

## Complete Example: Configuration Parser

```rust
use std::collections::HashMap;
use std::fs;
use thiserror::Error;

#[derive(Error, Debug)]
enum ConfigError {
    #[error("Failed to read config file: {0}")]
    IoError(#[from] std::io::Error),

    #[error("Invalid config format at line {line}: {message}")]
    ParseError { line: usize, message: String },

    #[error("Missing required key: {0}")]
    MissingKey(String),

    #[error("Invalid value for '{key}': expected {expected}, got '{value}'")]
    InvalidValue {
        key: String,
        expected: String,
        value: String,
    },
}

struct Config {
    values: HashMap<String, String>,
}

impl Config {
    fn load(path: &str) -> Result<Self, ConfigError> {
        let content = fs::read_to_string(path)?;
        Self::parse(&content)
    }

    fn parse(content: &str) -> Result<Self, ConfigError> {
        let mut values = HashMap::new();

        for (line_num, line) in content.lines().enumerate() {
            let line = line.trim();

            // Skip empty lines and comments
            if line.is_empty() || line.starts_with('#') {
                continue;
            }

            let parts: Vec<&str> = line.splitn(2, '=').collect();

            if parts.len() != 2 {
                return Err(ConfigError::ParseError {
                    line: line_num + 1,
                    message: "Expected format: key=value".to_string(),
                });
            }

            let key = parts[0].trim().to_string();
            let value = parts[1].trim().to_string();

            values.insert(key, value);
        }

        Ok(Config { values })
    }

    fn get(&self, key: &str) -> Option<&String> {
        self.values.get(key)
    }

    fn get_required(&self, key: &str) -> Result<&String, ConfigError> {
        self.values
            .get(key)
            .ok_or_else(|| ConfigError::MissingKey(key.to_string()))
    }

    fn get_int(&self, key: &str) -> Result<i32, ConfigError> {
        let value = self.get_required(key)?;
        value.parse().map_err(|_| ConfigError::InvalidValue {
            key: key.to_string(),
            expected: "integer".to_string(),
            value: value.clone(),
        })
    }

    fn get_bool(&self, key: &str) -> Result<bool, ConfigError> {
        let value = self.get_required(key)?;
        match value.to_lowercase().as_str() {
            "true" | "yes" | "1" => Ok(true),
            "false" | "no" | "0" => Ok(false),
            _ => Err(ConfigError::InvalidValue {
                key: key.to_string(),
                expected: "boolean (true/false)".to_string(),
                value: value.clone(),
            }),
        }
    }
}

fn main() -> Result<(), ConfigError> {
    // For demonstration, create a config string
    let config_content = r#"
        # Server configuration
        host=localhost
        port=8080
        debug=true
        max_connections=100
    "#;

    let config = Config::parse(config_content)?;

    let host = config.get_required("host")?;
    let port = config.get_int("port")?;
    let debug = config.get_bool("debug")?;
    let max_conn = config.get_int("max_connections")?;

    println!("Server: {}:{}", host, port);
    println!("Debug: {}", debug);
    println!("Max connections: {}", max_conn);

    // Try getting a missing key
    match config.get_required("missing") {
        Ok(value) => println!("Found: {}", value),
        Err(e) => eprintln!("Expected error: {}", e),
    }

    Ok(())
}
```

---

## Exercises

### Exercise 1: File Processing
Write a function that:
1. Opens a file
2. Reads numbers (one per line)
3. Returns the sum
4. Handles all potential errors properly

### Exercise 2: User Input Validation
Create a validation function that checks:
- Email contains '@'
- Age is between 0 and 150
- Username is 3-20 characters
Return custom errors for each validation failure.

### Exercise 3: Custom Error Type
Create a `BankError` enum with variants for:
- `InsufficientFunds { available: f64, requested: f64 }`
- `AccountNotFound(String)`
- `InvalidAmount(f64)`
Use `thiserror` for the implementation.

### Exercise 4: Error Propagation Chain
Create three functions that call each other:
1. `read_file() -> Result<String, FileError>`
2. `parse_data() -> Result<Data, ParseError>`
3. `process() -> Result<Output, AppError>`
Make `AppError` wrap both `FileError` and `ParseError`.

### Exercise 5: Retry Logic
Write a function that retries an operation up to 3 times:
```rust
fn retry<T, E, F>(operation: F) -> Result<T, E>
where
    F: Fn() -> Result<T, E>,
```

---

## Summary

You've learned:
- `panic!` for unrecoverable errors
- `Result<T, E>` for recoverable errors
- `Option<T>` for optional values
- The `?` operator for error propagation
- `unwrap()` and `expect()` for quick (dangerous) unwrapping
- Custom error types for domain-specific errors
- The `thiserror` crate for clean error definitions
- Best practices for error handling

---

## Next Steps

Now let's learn about async programming with [Async and Tokio](./05-async-tokio.md).

Async programming lets you write efficient concurrent code, and Tokio is Rust's most popular async runtime!

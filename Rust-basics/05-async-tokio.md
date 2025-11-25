# Async Programming with Tokio

## Learning Objectives

By the end of this tutorial, you will be able to:
- Understand async/await syntax and concepts
- Work with Futures in Rust
- Set up and use the Tokio runtime
- Spawn concurrent tasks
- Use async channels for communication
- Apply `join!` and `select!` for task coordination
- Build basic TCP and HTTP services

---

## Introduction

Async programming lets you write concurrent code that's efficient and scalable. Instead of blocking while waiting for I/O, async code can do other work. This is crucial for building network services, web servers, and any I/O-heavy applications.

Rust's async model is different from other languages—it's zero-cost and gives you fine-grained control. Let's dive in!

---

## What is Async Programming?

### The Problem with Blocking

```rust
// Blocking code - thread does nothing while waiting
fn fetch_data() -> String {
    std::thread::sleep(std::time::Duration::from_secs(2));
    String::from("data")
}

fn main() {
    let data1 = fetch_data();  // Wait 2 seconds
    let data2 = fetch_data();  // Wait 2 more seconds
    // Total: 4 seconds
}
```

### The Async Solution

```rust
use tokio::time::{sleep, Duration};

async fn fetch_data() -> String {
    sleep(Duration::from_secs(2)).await;
    String::from("data")
}

#[tokio::main]
async fn main() {
    let (data1, data2) = tokio::join!(
        fetch_data(),
        fetch_data()
    );
    // Total: 2 seconds (concurrent)
}
```

---

## Async/Await Syntax

### Defining Async Functions

```rust
// Add `async` before `fn`
async fn say_hello() {
    println!("Hello!");
}

// Async functions can return values
async fn compute_value() -> i32 {
    42
}

// With parameters
async fn greet(name: &str) -> String {
    format!("Hello, {}!", name)
}
```

### Calling Async Functions

Async functions return a `Future`. You must `.await` them:

```rust
#[tokio::main]
async fn main() {
    // This does nothing by itself
    let future = say_hello();

    // `.await` runs the future
    future.await;

    // Usually combined:
    say_hello().await;

    // Getting return values
    let value = compute_value().await;
    println!("Value: {}", value);
}
```

### Important: Futures are Lazy

```rust
#[tokio::main]
async fn main() {
    let future = async {
        println!("This won't print until awaited");
    };

    println!("Future created");

    future.await;  // Now it prints
}
```

---

## Futures Explained

A `Future` is a value that represents a computation that may not be complete yet.

```rust
// Simplified Future trait
pub trait Future {
    type Output;
    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output>;
}

pub enum Poll<T> {
    Ready(T),
    Pending,
}
```

When you `.await` a future:
1. The runtime polls it
2. If `Pending`, the runtime does other work
3. When the future is ready, it resumes
4. If `Ready`, you get the result

You rarely need to implement `Future` directly—`async`/`await` does it for you.

---

## The Tokio Runtime

Tokio is Rust's most popular async runtime. It provides:
- The event loop
- Task scheduling
- Async I/O
- Timers
- Channels
- And much more

### Setting Up Tokio

```toml
# Cargo.toml
[dependencies]
tokio = { version = "1", features = ["full"] }
```

### The `#[tokio::main]` Macro

```rust
#[tokio::main]
async fn main() {
    println!("Hello from async main!");
}

// Expands to roughly:
fn main() {
    tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .unwrap()
        .block_on(async {
            println!("Hello from async main!");
        })
}
```

### Custom Runtime Configuration

```rust
fn main() {
    let runtime = tokio::runtime::Builder::new_multi_thread()
        .worker_threads(4)
        .enable_all()
        .build()
        .unwrap();

    runtime.block_on(async {
        println!("Custom runtime!");
    });
}
```

### Single-Threaded Runtime

```rust
#[tokio::main(flavor = "current_thread")]
async fn main() {
    // Runs on single thread
}
```

---

## Spawning Tasks

`tokio::spawn` runs a future concurrently:

```rust
use tokio::time::{sleep, Duration};

#[tokio::main]
async fn main() {
    // Spawn a task
    let handle = tokio::spawn(async {
        sleep(Duration::from_millis(100)).await;
        "Task complete"
    });

    // Do other work while task runs
    println!("Task spawned");

    // Wait for task to complete
    let result = handle.await.unwrap();
    println!("Result: {}", result);
}
```

### Multiple Tasks

```rust
use tokio::time::{sleep, Duration};

async fn task(id: u32) -> u32 {
    sleep(Duration::from_millis(100)).await;
    println!("Task {} complete", id);
    id * 2
}

#[tokio::main]
async fn main() {
    let mut handles = vec![];

    for i in 0..5 {
        let handle = tokio::spawn(task(i));
        handles.push(handle);
    }

    for handle in handles {
        let result = handle.await.unwrap();
        println!("Got: {}", result);
    }
}
```

### Task Return Types

```rust
#[tokio::main]
async fn main() {
    // JoinHandle<T> contains the task's return value
    let handle: tokio::task::JoinHandle<i32> = tokio::spawn(async {
        42
    });

    let value = handle.await.unwrap();
    println!("Value: {}", value);

    // Handling errors
    let handle = tokio::spawn(async {
        panic!("Oh no!");
    });

    match handle.await {
        Ok(value) => println!("Success: {:?}", value),
        Err(e) => println!("Task panicked: {}", e),
    }
}
```

### `'static` Requirement

Spawned tasks must be `'static`:

```rust
#[tokio::main]
async fn main() {
    let data = String::from("hello");

    // This won't compile:
    // tokio::spawn(async {
    //     println!("{}", data);  // data is borrowed, not owned
    // });

    // Solution: move ownership
    tokio::spawn(async move {
        println!("{}", data);  // data is moved into task
    });

    // data is no longer available here
}
```

---

## Join and Select

### `join!`: Wait for All

```rust
use tokio::time::{sleep, Duration};

async fn fetch_user() -> String {
    sleep(Duration::from_millis(100)).await;
    String::from("User data")
}

async fn fetch_orders() -> Vec<String> {
    sleep(Duration::from_millis(150)).await;
    vec![String::from("Order 1"), String::from("Order 2")]
}

async fn fetch_settings() -> String {
    sleep(Duration::from_millis(50)).await;
    String::from("Settings")
}

#[tokio::main]
async fn main() {
    // Run all concurrently, wait for all to complete
    let (user, orders, settings) = tokio::join!(
        fetch_user(),
        fetch_orders(),
        fetch_settings()
    );

    println!("User: {}", user);
    println!("Orders: {:?}", orders);
    println!("Settings: {}", settings);
}
```

### `try_join!`: Short-Circuit on Error

```rust
use tokio::try_join;

async fn fetch_user() -> Result<String, String> {
    Ok(String::from("User"))
}

async fn fetch_orders() -> Result<Vec<String>, String> {
    Err(String::from("Database error"))
}

#[tokio::main]
async fn main() {
    let result = try_join!(fetch_user(), fetch_orders());

    match result {
        Ok((user, orders)) => println!("Got {} and {:?}", user, orders),
        Err(e) => println!("Error: {}", e),
    }
}
```

### `select!`: Wait for First

```rust
use tokio::time::{sleep, Duration};

async fn slow_operation() -> &'static str {
    sleep(Duration::from_secs(10)).await;
    "Slow"
}

async fn fast_operation() -> &'static str {
    sleep(Duration::from_millis(100)).await;
    "Fast"
}

#[tokio::main]
async fn main() {
    tokio::select! {
        result = slow_operation() => {
            println!("Slow: {}", result);
        }
        result = fast_operation() => {
            println!("Fast: {}", result);
        }
    }
    // Only "Fast: Fast" prints
}
```

### Timeout Pattern with `select!`

```rust
use tokio::time::{sleep, Duration};

async fn long_running_task() -> String {
    sleep(Duration::from_secs(30)).await;
    String::from("Done")
}

#[tokio::main]
async fn main() {
    tokio::select! {
        result = long_running_task() => {
            println!("Task completed: {}", result);
        }
        _ = sleep(Duration::from_secs(5)) => {
            println!("Task timed out!");
        }
    }
}
```

---

## Async Channels

Channels allow communication between tasks.

### `mpsc` (Multi-Producer, Single-Consumer)

```rust
use tokio::sync::mpsc;

#[tokio::main]
async fn main() {
    // Create channel with buffer size 32
    let (tx, mut rx) = mpsc::channel::<String>(32);

    // Spawn producer task
    let tx1 = tx.clone();
    tokio::spawn(async move {
        for i in 0..5 {
            tx1.send(format!("Message {}", i)).await.unwrap();
        }
    });

    // Drop original sender so receiver knows when done
    drop(tx);

    // Receive messages
    while let Some(message) = rx.recv().await {
        println!("Received: {}", message);
    }
}
```

### Bounded vs Unbounded

```rust
use tokio::sync::mpsc;

#[tokio::main]
async fn main() {
    // Bounded: blocks sender when full
    let (tx, mut rx) = mpsc::channel::<i32>(10);

    // Unbounded: never blocks (can use lots of memory!)
    let (tx_unbounded, mut rx_unbounded) = mpsc::unbounded_channel::<i32>();
}
```

### `oneshot` (Single Message)

```rust
use tokio::sync::oneshot;

#[tokio::main]
async fn main() {
    let (tx, rx) = oneshot::channel();

    tokio::spawn(async move {
        // Do some work
        let result = 42;
        tx.send(result).unwrap();
    });

    let result = rx.await.unwrap();
    println!("Result: {}", result);
}
```

### `broadcast` (Multi-Producer, Multi-Consumer)

```rust
use tokio::sync::broadcast;

#[tokio::main]
async fn main() {
    let (tx, mut rx1) = broadcast::channel::<String>(16);
    let mut rx2 = tx.subscribe();

    tokio::spawn(async move {
        while let Ok(msg) = rx1.recv().await {
            println!("Receiver 1: {}", msg);
        }
    });

    tokio::spawn(async move {
        while let Ok(msg) = rx2.recv().await {
            println!("Receiver 2: {}", msg);
        }
    });

    tx.send(String::from("Hello")).unwrap();
    tx.send(String::from("World")).unwrap();

    // Give receivers time to process
    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
}
```

### `watch` (Single Value, Multiple Readers)

```rust
use tokio::sync::watch;

#[tokio::main]
async fn main() {
    let (tx, mut rx) = watch::channel("initial");

    tokio::spawn(async move {
        while rx.changed().await.is_ok() {
            println!("Value changed: {}", *rx.borrow());
        }
    });

    tx.send("updated").unwrap();
    tx.send("final").unwrap();

    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
}
```

---

## Basic TCP with Tokio

### TCP Server

```rust
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::TcpListener;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let listener = TcpListener::bind("127.0.0.1:8080").await?;
    println!("Server listening on port 8080");

    loop {
        let (mut socket, addr) = listener.accept().await?;
        println!("New connection from {}", addr);

        tokio::spawn(async move {
            let mut buf = [0; 1024];

            loop {
                let n = match socket.read(&mut buf).await {
                    Ok(0) => return,  // Connection closed
                    Ok(n) => n,
                    Err(e) => {
                        eprintln!("Failed to read: {}", e);
                        return;
                    }
                };

                // Echo back
                if let Err(e) = socket.write_all(&buf[0..n]).await {
                    eprintln!("Failed to write: {}", e);
                    return;
                }
            }
        });
    }
}
```

### TCP Client

```rust
use tokio::io::{AsyncReadExt, AsyncWriteExt};
use tokio::net::TcpStream;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut stream = TcpStream::connect("127.0.0.1:8080").await?;

    // Write data
    stream.write_all(b"Hello, server!").await?;

    // Read response
    let mut buf = [0; 1024];
    let n = stream.read(&mut buf).await?;

    println!("Response: {}", String::from_utf8_lossy(&buf[..n]));

    Ok(())
}
```

---

## Basic HTTP with Tokio

Using the `reqwest` crate for HTTP client:

```toml
# Cargo.toml
[dependencies]
tokio = { version = "1", features = ["full"] }
reqwest = { version = "0.11", features = ["json"] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"
```

### HTTP Client

```rust
use reqwest;
use serde::{Deserialize, Serialize};

#[derive(Debug, Deserialize)]
struct User {
    id: u32,
    name: String,
    email: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Simple GET
    let body = reqwest::get("https://httpbin.org/get")
        .await?
        .text()
        .await?;
    println!("Response: {}", body);

    // GET with JSON parsing
    let user: User = reqwest::get("https://jsonplaceholder.typicode.com/users/1")
        .await?
        .json()
        .await?;
    println!("User: {:?}", user);

    // POST request
    let client = reqwest::Client::new();
    let response = client
        .post("https://httpbin.org/post")
        .json(&serde_json::json!({
            "name": "John",
            "age": 30
        }))
        .send()
        .await?;

    println!("Status: {}", response.status());

    Ok(())
}
```

### Multiple Concurrent Requests

```rust
use reqwest;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let urls = vec![
        "https://httpbin.org/get",
        "https://httpbin.org/ip",
        "https://httpbin.org/user-agent",
    ];

    let client = reqwest::Client::new();

    let futures: Vec<_> = urls
        .iter()
        .map(|url| {
            let client = client.clone();
            async move {
                let resp = client.get(*url).send().await?;
                let body = resp.text().await?;
                Ok::<_, reqwest::Error>((url, body))
            }
        })
        .collect();

    let results = futures::future::join_all(futures).await;

    for result in results {
        match result {
            Ok((url, body)) => println!("{}: {} bytes", url, body.len()),
            Err(e) => eprintln!("Error: {}", e),
        }
    }

    Ok(())
}
```

---

## Async Patterns

### Pattern 1: Timeout

```rust
use tokio::time::{timeout, Duration};

async fn fetch_data() -> String {
    tokio::time::sleep(Duration::from_secs(10)).await;
    String::from("data")
}

#[tokio::main]
async fn main() {
    match timeout(Duration::from_secs(2), fetch_data()).await {
        Ok(data) => println!("Got: {}", data),
        Err(_) => println!("Operation timed out"),
    }
}
```

### Pattern 2: Retry with Backoff

```rust
use tokio::time::{sleep, Duration};

async fn unreliable_operation() -> Result<String, String> {
    // Simulated failure
    Err(String::from("Failed"))
}

async fn with_retry<F, Fut, T, E>(
    operation: F,
    max_retries: u32,
) -> Result<T, E>
where
    F: Fn() -> Fut,
    Fut: std::future::Future<Output = Result<T, E>>,
{
    let mut attempts = 0;
    loop {
        match operation().await {
            Ok(value) => return Ok(value),
            Err(e) => {
                attempts += 1;
                if attempts >= max_retries {
                    return Err(e);
                }
                let delay = Duration::from_millis(100 * 2_u64.pow(attempts));
                sleep(delay).await;
            }
        }
    }
}

#[tokio::main]
async fn main() {
    let result = with_retry(|| unreliable_operation(), 3).await;
    println!("Result: {:?}", result);
}
```

### Pattern 3: Graceful Shutdown

```rust
use tokio::signal;
use tokio::sync::broadcast;

#[tokio::main]
async fn main() {
    let (shutdown_tx, mut shutdown_rx) = broadcast::channel::<()>(1);

    // Spawn workers
    for i in 0..3 {
        let mut rx = shutdown_tx.subscribe();
        tokio::spawn(async move {
            loop {
                tokio::select! {
                    _ = rx.recv() => {
                        println!("Worker {} shutting down", i);
                        return;
                    }
                    _ = tokio::time::sleep(tokio::time::Duration::from_secs(1)) => {
                        println!("Worker {} tick", i);
                    }
                }
            }
        });
    }

    // Wait for Ctrl+C
    signal::ctrl_c().await.expect("Failed to listen for ctrl+c");
    println!("Shutdown signal received");

    // Signal all workers to shut down
    let _ = shutdown_tx.send(());

    // Give workers time to shut down
    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
}
```

### Pattern 4: Worker Pool

```rust
use tokio::sync::mpsc;

async fn worker(id: u32, mut rx: mpsc::Receiver<String>) {
    while let Some(job) = rx.recv().await {
        println!("Worker {} processing: {}", id, job);
        tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;
    }
    println!("Worker {} done", id);
}

#[tokio::main]
async fn main() {
    let (tx, rx) = mpsc::channel::<String>(100);

    // Spawn workers
    for i in 0..4 {
        let rx = rx.clone();
        tokio::spawn(async move {
            // Share receiver among workers using a mutex wrapper
        });
    }

    // Simple version: one receiver
    tokio::spawn(worker(0, rx));

    // Send jobs
    for i in 0..10 {
        tx.send(format!("Job {}", i)).await.unwrap();
    }

    // Drop sender to signal completion
    drop(tx);

    tokio::time::sleep(tokio::time::Duration::from_secs(2)).await;
}
```

---

## Common Pitfalls

### 1. Blocking in Async Code

```rust
// BAD: This blocks the entire runtime
async fn bad() {
    std::thread::sleep(std::time::Duration::from_secs(1));
}

// GOOD: Use async sleep
async fn good() {
    tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
}

// For blocking operations, use spawn_blocking
async fn blocking_task() -> i32 {
    tokio::task::spawn_blocking(|| {
        // This runs on a blocking thread pool
        std::thread::sleep(std::time::Duration::from_secs(1));
        42
    }).await.unwrap()
}
```

### 2. Forgetting to Await

```rust
#[tokio::main]
async fn main() {
    // This does nothing!
    async_operation();

    // Need to await
    async_operation().await;
}
```

### 3. Holding Locks Across Await

```rust
use std::sync::Mutex;

// BAD: Mutex guard held across await
async fn bad(data: &Mutex<Vec<i32>>) {
    let mut guard = data.lock().unwrap();
    some_async_operation().await;  // Guard still held!
    guard.push(1);
}

// GOOD: Drop guard before await
async fn good(data: &Mutex<Vec<i32>>) {
    {
        let mut guard = data.lock().unwrap();
        guard.push(1);
    }  // Guard dropped
    some_async_operation().await;
}

// BETTER: Use tokio::sync::Mutex for async
use tokio::sync::Mutex as AsyncMutex;

async fn better(data: &AsyncMutex<Vec<i32>>) {
    let mut guard = data.lock().await;
    some_async_operation().await;
    guard.push(1);
}
```

---

## Exercises

### Exercise 1: Concurrent Downloads
Write a program that:
1. Takes a list of URLs
2. Downloads them all concurrently
3. Prints each URL and response size
4. Measures total time

### Exercise 2: Chat Server
Build a simple chat server that:
1. Accepts TCP connections
2. Broadcasts messages to all connected clients
3. Handles client disconnections gracefully

### Exercise 3: Rate Limiter
Implement an async rate limiter:
```rust
struct RateLimiter {
    // ...
}

impl RateLimiter {
    async fn acquire(&self);  // Waits if rate exceeded
}
```

### Exercise 4: Periodic Task
Create a service that:
1. Runs a task every 5 seconds
2. Can be stopped with a shutdown signal
3. Logs each execution

### Exercise 5: Parallel Processing Pipeline
Build a pipeline that:
1. Reads items from one channel
2. Processes them with N workers
3. Sends results to another channel
4. Collects and prints results

---

## Summary

You've learned:
- `async`/`await` syntax for writing async code
- Futures are lazy and need to be polled
- Tokio provides the async runtime
- `tokio::spawn` for concurrent tasks
- `join!` waits for all futures
- `select!` waits for the first future
- Channels for task communication
- Basic TCP and HTTP with async code
- Common patterns and pitfalls

---

## What's Next?

Congratulations! You've completed the Rust basics tutorial series. You now have a solid foundation in:
- Syntax and types
- Ownership and borrowing
- Structs and enums
- Error handling
- Async programming

To continue your Rust journey:
- Read "The Rust Programming Language" book
- Practice on exercism.org or rustlings
- Build small projects
- Explore crates like `serde`, `clap`, `tracing`
- Join the Rust community on Discord or the forums

Remember: Rust has a learning curve, but it's worth it. The compiler is strict because it's protecting you from bugs. Trust it, learn from its messages, and you'll write reliable, efficient code.

Happy coding!

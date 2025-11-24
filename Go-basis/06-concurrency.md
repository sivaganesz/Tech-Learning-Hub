# Concurrency in Go

## Learning Objectives

By the end of this tutorial, you will be able to:
- Create and manage goroutines
- Use channels for communication (buffered and unbuffered)
- Implement select statements for channel multiplexing
- Apply common channel patterns (done, fan-out, fan-in)
- Use sync package primitives (WaitGroup, Mutex)
- Implement cancellation with context

---

## 1. Goroutines

Goroutines are lightweight threads managed by the Go runtime:

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Start goroutines
    go sayHello("Alice")
    go sayHello("Bob")
    go sayHello("Carol")

    // Main goroutine continues
    fmt.Println("Main: Goroutines started")

    // Wait for goroutines to complete
    // (This is not the proper way - use WaitGroup instead)
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Main: Done")
}
```

### Anonymous Goroutines

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Anonymous function as goroutine
    go func() {
        fmt.Println("Anonymous goroutine")
    }()

    // With parameters
    message := "Hello from closure"
    go func(msg string) {
        fmt.Println(msg)
    }(message)

    // Multiple goroutines
    for i := 0; i < 5; i++ {
        go func(n int) {
            fmt.Printf("Goroutine %d\n", n)
        }(i) // Pass i as parameter to avoid closure issue
    }

    time.Sleep(100 * time.Millisecond)
}
```

### Closure Pitfall

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // WRONG - closure captures variable, not value
    fmt.Println("Wrong way:")
    for i := 0; i < 3; i++ {
        go func() {
            fmt.Println(i) // Will likely print 3, 3, 3
        }()
    }
    time.Sleep(50 * time.Millisecond)

    // CORRECT - pass value as parameter
    fmt.Println("\nCorrect way:")
    for i := 0; i < 3; i++ {
        go func(n int) {
            fmt.Println(n) // Will print 0, 1, 2 (in some order)
        }(i)
    }
    time.Sleep(50 * time.Millisecond)
}
```

---

## 2. Channels

Channels are typed conduits for communication between goroutines:

### Unbuffered Channels

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Create unbuffered channel
    ch := make(chan string)

    // Sender goroutine
    go func() {
        fmt.Println("Sender: Sending message...")
        ch <- "Hello, Channel!" // Blocks until receiver is ready
        fmt.Println("Sender: Message sent")
    }()

    // Give sender time to start
    time.Sleep(100 * time.Millisecond)

    // Receiver (main goroutine)
    fmt.Println("Receiver: Waiting for message...")
    msg := <-ch // Blocks until sender sends
    fmt.Println("Receiver: Got", msg)
}
```

### Buffered Channels

```go
package main

import "fmt"

func main() {
    // Buffered channel with capacity 3
    ch := make(chan int, 3)

    // Can send without blocking (up to capacity)
    ch <- 1
    ch <- 2
    ch <- 3
    // ch <- 4 // This would block (buffer full)

    fmt.Println("Buffer length:", len(ch))
    fmt.Println("Buffer capacity:", cap(ch))

    // Receive values
    fmt.Println(<-ch) // 1
    fmt.Println(<-ch) // 2
    fmt.Println(<-ch) // 3
}
```

### Channel Direction

```go
package main

import "fmt"

// Send-only channel parameter
func sender(ch chan<- string) {
    ch <- "Hello"
    ch <- "World"
    close(ch)
}

// Receive-only channel parameter
func receiver(ch <-chan string) {
    for msg := range ch {
        fmt.Println("Received:", msg)
    }
}

func main() {
    ch := make(chan string, 2)

    go sender(ch)
    receiver(ch)
}
```

### Closing Channels

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 5)

    // Send values
    for i := 1; i <= 5; i++ {
        ch <- i
    }
    close(ch) // Signal no more values

    // Range over channel (stops when closed)
    for value := range ch {
        fmt.Println(value)
    }

    // Check if channel is closed
    ch2 := make(chan string, 1)
    ch2 <- "hello"
    close(ch2)

    value, ok := <-ch2
    fmt.Printf("Value: %s, Open: %v\n", value, ok) // hello, true

    value, ok = <-ch2
    fmt.Printf("Value: %s, Open: %v\n", value, ok) // "", false
}
```

---

## 3. Select Statement

Select lets you wait on multiple channel operations:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(100 * time.Millisecond)
        ch1 <- "one"
    }()

    go func() {
        time.Sleep(200 * time.Millisecond)
        ch2 <- "two"
    }()

    // Wait for both
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received", msg2)
        }
    }
}
```

### Non-Blocking Select with Default

```go
package main

import "fmt"

func main() {
    ch := make(chan string, 1)

    // Non-blocking receive
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    default:
        fmt.Println("No message available")
    }

    // Non-blocking send
    ch <- "hello"
    select {
    case ch <- "world":
        fmt.Println("Sent message")
    default:
        fmt.Println("Channel full, message dropped")
    }
}
```

### Timeout with Select

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)

    go func() {
        time.Sleep(2 * time.Second)
        ch <- "result"
    }()

    // Wait with timeout
    select {
    case result := <-ch:
        fmt.Println("Got result:", result)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout!")
    }
}
```

---

## 4. Channel Patterns

### Done Channel (Signaling Completion)

```go
package main

import (
    "fmt"
    "time"
)

func worker(done chan bool) {
    fmt.Println("Working...")
    time.Sleep(time.Second)
    fmt.Println("Done working")

    done <- true
}

func main() {
    done := make(chan bool, 1)

    go worker(done)

    // Wait for worker to finish
    <-done
    fmt.Println("Worker finished")
}
```

### Fan-Out (One to Many)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    jobs := make(chan int, 10)
    var wg sync.WaitGroup

    // Start 3 workers (fan-out)
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, jobs, &wg)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    wg.Wait()
    fmt.Println("All jobs completed")
}
```

### Fan-In (Many to One)

```go
package main

import (
    "fmt"
    "time"
)

func producer(id int, ch chan<- string) {
    for i := 0; i < 3; i++ {
        ch <- fmt.Sprintf("Producer %d: message %d", id, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func fanIn(ch1, ch2, ch3 <-chan string) <-chan string {
    out := make(chan string)

    go func() {
        for {
            select {
            case msg := <-ch1:
                out <- msg
            case msg := <-ch2:
                out <- msg
            case msg := <-ch3:
                out <- msg
            }
        }
    }()

    return out
}

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    ch3 := make(chan string)

    go producer(1, ch1)
    go producer(2, ch2)
    go producer(3, ch3)

    merged := fanIn(ch1, ch2, ch3)

    // Receive merged messages
    for i := 0; i < 9; i++ {
        fmt.Println(<-merged)
    }
}
```

### Pipeline Pattern

```go
package main

import "fmt"

// Stage 1: Generate numbers
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Stage 3: Double numbers
func double(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * 2
        }
        close(out)
    }()
    return out
}

func main() {
    // Build pipeline
    numbers := generate(1, 2, 3, 4, 5)
    squared := square(numbers)
    doubled := double(squared)

    // Consume output
    for result := range doubled {
        fmt.Println(result)
    }
}
```

### Worker Pool

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID  int
    Output string
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        // Process job
        time.Sleep(100 * time.Millisecond)
        result := Result{
            JobID:  job.ID,
            Output: fmt.Sprintf("Worker %d processed: %s", id, job.Data),
        }
        results <- result
    }
}

func main() {
    numWorkers := 3
    numJobs := 10

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    var wg sync.WaitGroup

    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    // Send jobs
    for j := 1; j <= numJobs; j++ {
        jobs <- Job{ID: j, Data: fmt.Sprintf("job-%d", j)}
    }
    close(jobs)

    // Wait for workers and close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    for result := range results {
        fmt.Printf("Job %d: %s\n", result.JobID, result.Output)
    }
}
```

---

## 5. sync Package

### WaitGroup

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Duration(id) * 100 * time.Millisecond)
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }

    wg.Wait()
    fmt.Println("All workers completed")
}
```

### Mutex (Mutual Exclusion)

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := SafeCounter{}
    var wg sync.WaitGroup

    // 1000 goroutines incrementing
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Println("Final count:", counter.Value()) // Always 1000
}
```

### RWMutex (Read-Write Mutex)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock() // Multiple readers allowed
    defer c.mu.RUnlock()
    value, ok := c.data[key]
    return value, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock() // Exclusive access
    defer c.mu.Unlock()
    c.data[key] = value
}

func main() {
    cache := Cache{data: make(map[string]string)}
    var wg sync.WaitGroup

    // Writer
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 5; i++ {
            key := fmt.Sprintf("key%d", i)
            cache.Set(key, fmt.Sprintf("value%d", i))
            time.Sleep(10 * time.Millisecond)
        }
    }()

    // Readers
    for r := 0; r < 3; r++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for i := 0; i < 10; i++ {
                key := fmt.Sprintf("key%d", i%5)
                if value, ok := cache.Get(key); ok {
                    fmt.Printf("Reader %d: %s=%s\n", id, key, value)
                }
                time.Sleep(5 * time.Millisecond)
            }
        }(r)
    }

    wg.Wait()
}
```

### Once (Run Once)

```go
package main

import (
    "fmt"
    "sync"
)

var once sync.Once
var config map[string]string

func loadConfig() {
    fmt.Println("Loading configuration...")
    config = map[string]string{
        "host": "localhost",
        "port": "8080",
    }
}

func getConfig() map[string]string {
    once.Do(loadConfig) // Only runs once
    return config
}

func main() {
    var wg sync.WaitGroup

    // Multiple goroutines trying to get config
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            cfg := getConfig()
            fmt.Printf("Goroutine %d: host=%s\n", id, cfg["host"])
        }(i)
    }

    wg.Wait()
}
```

---

## 6. Context for Cancellation

### Basic Context

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: cancelled\n", id)
            return
        default:
            fmt.Printf("Worker %d: working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // Create cancellable context
    ctx, cancel := context.WithCancel(context.Background())

    // Start workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }

    // Let workers run
    time.Sleep(2 * time.Second)

    // Cancel all workers
    cancel()

    // Wait for workers to finish
    time.Sleep(100 * time.Millisecond)
    fmt.Println("All workers stopped")
}
```

### Context with Timeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func slowOperation(ctx context.Context) error {
    select {
    case <-time.After(5 * time.Second):
        return nil // Operation completed
    case <-ctx.Done():
        return ctx.Err() // Cancelled or timed out
    }
}

func main() {
    // Context with 2 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    fmt.Println("Starting slow operation...")
    err := slowOperation(ctx)
    if err != nil {
        fmt.Println("Operation failed:", err)
    } else {
        fmt.Println("Operation completed")
    }
}
```

### Context with Deadline

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func processRequest(ctx context.Context) {
    deadline, ok := ctx.Deadline()
    if ok {
        fmt.Println("Deadline:", deadline)
    }

    select {
    case <-time.After(3 * time.Second):
        fmt.Println("Request processed")
    case <-ctx.Done():
        fmt.Println("Request cancelled:", ctx.Err())
    }
}

func main() {
    // Set deadline to 2 seconds from now
    deadline := time.Now().Add(2 * time.Second)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()

    processRequest(ctx)
}
```

### Context with Values

```go
package main

import (
    "context"
    "fmt"
)

type contextKey string

const (
    userIDKey   contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

func processRequest(ctx context.Context) {
    userID := ctx.Value(userIDKey)
    requestID := ctx.Value(requestIDKey)

    fmt.Printf("Processing request %v for user %v\n", requestID, userID)
}

func main() {
    // Create context with values
    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey, 123)
    ctx = context.WithValue(ctx, requestIDKey, "req-456")

    processRequest(ctx)
}
```

### HTTP Request with Context

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
)

func fetchURL(ctx context.Context, url string) (string, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return "", err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

func main() {
    // 5 second timeout for HTTP request
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    body, err := fetchURL(ctx, "https://httpbin.org/get")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Response length:", len(body))
}
```

---

## Exercises

### Exercise 1: Concurrent URL Fetcher
Fetch multiple URLs concurrently and collect results.

```go
package main

import (
    "fmt"
    "net/http"
    "sync"
    "time"
)

type Result struct {
    URL        string
    StatusCode int
    Duration   time.Duration
    Error      error
}

// TODO: Implement concurrent URL fetcher
func fetchURLs(urls []string) []Result {
    return nil
}

func main() {
    urls := []string{
        "https://golang.org",
        "https://google.com",
        "https://github.com",
        "https://invalid.url.example",
    }

    results := fetchURLs(urls)
    for _, r := range results {
        if r.Error != nil {
            fmt.Printf("%s: Error - %v\n", r.URL, r.Error)
        } else {
            fmt.Printf("%s: %d (%v)\n", r.URL, r.StatusCode, r.Duration)
        }
    }
}
```

### Exercise 2: Rate Limiter
Implement a rate limiter using channels.

```go
package main

import (
    "fmt"
    "time"
)

type RateLimiter struct {
    // TODO: Add fields
}

// TODO: Implement NewRateLimiter that allows n requests per duration
func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
    return nil
}

// TODO: Implement Wait that blocks until request is allowed
func (rl *RateLimiter) Wait() {
}

func main() {
    // Allow 5 requests per second
    limiter := NewRateLimiter(5, time.Second)

    // Try to make 10 requests
    for i := 1; i <= 10; i++ {
        limiter.Wait()
        fmt.Printf("Request %d at %v\n", i, time.Now().Format("15:04:05.000"))
    }
}
```

### Exercise 3: Producer-Consumer with Multiple Consumers
Implement a producer-consumer pattern with configurable consumers.

```go
package main

import (
    "fmt"
    "time"
)

type Task struct {
    ID   int
    Data string
}

// TODO: Implement producer that generates tasks
func producer(tasks chan<- Task, count int) {
}

// TODO: Implement consumer that processes tasks
func consumer(id int, tasks <-chan Task, done chan<- bool) {
}

func main() {
    numConsumers := 3
    numTasks := 10

    tasks := make(chan Task, numTasks)
    done := make(chan bool, numConsumers)

    // Start consumers
    for i := 1; i <= numConsumers; i++ {
        go consumer(i, tasks, done)
    }

    // Start producer
    go producer(tasks, numTasks)

    // Wait for all consumers to finish
    for i := 0; i < numConsumers; i++ {
        <-done
    }

    fmt.Println("All tasks processed")
}
```

### Exercise 4: Concurrent Map with RWMutex
Implement a thread-safe map.

```go
package main

import (
    "fmt"
    "sync"
)

type ConcurrentMap struct {
    // TODO: Add fields
}

// TODO: Implement NewConcurrentMap
func NewConcurrentMap() *ConcurrentMap {
    return nil
}

// TODO: Implement Get
func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
    return nil, false
}

// TODO: Implement Set
func (m *ConcurrentMap) Set(key string, value interface{}) {
}

// TODO: Implement Delete
func (m *ConcurrentMap) Delete(key string) {
}

// TODO: Implement Len
func (m *ConcurrentMap) Len() int {
    return 0
}

func main() {
    m := NewConcurrentMap()
    var wg sync.WaitGroup

    // Concurrent writes
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            m.Set(fmt.Sprintf("key%d", n), n)
        }(i)
    }

    // Concurrent reads
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            m.Get(fmt.Sprintf("key%d", n))
        }(i)
    }

    wg.Wait()
    fmt.Println("Map size:", m.Len())
}
```

### Exercise 5: Graceful Shutdown
Implement graceful shutdown with context.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Server struct {
    // TODO: Add fields
}

// TODO: Implement Start that runs until context is cancelled
func (s *Server) Start(ctx context.Context) error {
    return nil
}

// TODO: Implement worker goroutines that check context
func (s *Server) worker(ctx context.Context, id int) {
}

func main() {
    // Create context that cancels on interrupt signal
    ctx, cancel := context.WithCancel(context.Background())

    // Handle OS signals
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigCh
        fmt.Println("\nReceived shutdown signal")
        cancel()
    }()

    server := &Server{}
    if err := server.Start(ctx); err != nil {
        fmt.Println("Server error:", err)
    }

    fmt.Println("Server stopped gracefully")
}
```

---

## Summary

In this tutorial, you learned:
- Creating goroutines for concurrent execution
- Using channels for communication (buffered and unbuffered)
- Implementing select for channel multiplexing
- Common patterns: done channels, fan-out, fan-in, pipelines, worker pools
- Using sync primitives: WaitGroup, Mutex, RWMutex, Once
- Implementing cancellation and timeouts with context

---

**Next:** [07-gin-framework.md](07-gin-framework.md) - Learn about the Gin web framework

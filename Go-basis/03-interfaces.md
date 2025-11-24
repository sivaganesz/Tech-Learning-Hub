# Interfaces in Go

## Learning Objectives

By the end of this tutorial, you will be able to:
- Define and implement interfaces
- Understand implicit interface implementation
- Use the empty interface (`any`)
- Work with common interfaces (error, Stringer, io.Reader)
- Perform type assertions and type switches
- Design flexible code using interfaces

---

## 1. Interface Definition

Interfaces define a set of method signatures:

```go
package main

import (
    "fmt"
    "math"
)

// Define an interface
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle implements Shape
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle implements Shape
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Function that accepts any Shape
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 7}

    printShapeInfo(rect)
    printShapeInfo(circle)

    // Slice of shapes
    shapes := []Shape{rect, circle}
    for _, shape := range shapes {
        printShapeInfo(shape)
    }
}
```

---

## 2. Implicit Implementation

Go interfaces are implemented implicitly - no `implements` keyword needed:

```go
package main

import "fmt"

// Interface
type Speaker interface {
    Speak() string
}

// Dog implements Speaker implicitly
type Dog struct {
    Name string
}

func (d Dog) Speak() string {
    return d.Name + " says: Woof!"
}

// Cat implements Speaker implicitly
type Cat struct {
    Name string
}

func (c Cat) Speak() string {
    return c.Name + " says: Meow!"
}

// Robot implements Speaker implicitly
type Robot struct {
    Model string
}

func (r Robot) Speak() string {
    return r.Model + " says: Beep boop!"
}

func makeSpeak(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    dog := Dog{Name: "Buddy"}
    cat := Cat{Name: "Whiskers"}
    robot := Robot{Model: "R2D2"}

    makeSpeak(dog)
    makeSpeak(cat)
    makeSpeak(robot)

    // All speakers in a slice
    speakers := []Speaker{dog, cat, robot}
    for _, s := range speakers {
        makeSpeak(s)
    }
}
```

### Interface Composition

```go
package main

import "fmt"

// Small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

// Composed interfaces
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// Implementation
type File struct {
    name string
}

func (f *File) Read(p []byte) (int, error) {
    fmt.Printf("Reading from %s\n", f.name)
    return len(p), nil
}

func (f *File) Write(p []byte) (int, error) {
    fmt.Printf("Writing to %s: %s\n", f.name, string(p))
    return len(p), nil
}

func (f *File) Close() error {
    fmt.Printf("Closing %s\n", f.name)
    return nil
}

func processData(rw ReadWriter) {
    data := make([]byte, 100)
    rw.Read(data)
    rw.Write([]byte("processed data"))
}

func main() {
    file := &File{name: "data.txt"}
    processData(file)
    file.Close()
}
```

---

## 3. Empty Interface (any)

The empty interface can hold values of any type:

```go
package main

import "fmt"

func main() {
    // Empty interface (any is an alias for interface{})
    var anything interface{}

    anything = 42
    fmt.Printf("Integer: %v (type: %T)\n", anything, anything)

    anything = "hello"
    fmt.Printf("String: %v (type: %T)\n", anything, anything)

    anything = true
    fmt.Printf("Boolean: %v (type: %T)\n", anything, anything)

    anything = []int{1, 2, 3}
    fmt.Printf("Slice: %v (type: %T)\n", anything, anything)

    // Using any (Go 1.18+)
    var value any = "Go is awesome"
    fmt.Println(value)

    // Function that accepts any type
    printValue(42)
    printValue("hello")
    printValue(3.14)
    printValue([]string{"a", "b", "c"})
}

func printValue(v interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", v, v)
}
```

### Working with Empty Interface

```go
package main

import "fmt"

// Map with any values
func main() {
    // Flexible data structure
    config := map[string]interface{}{
        "host":    "localhost",
        "port":    8080,
        "debug":   true,
        "timeout": 30.5,
        "tags":    []string{"api", "web"},
    }

    fmt.Println("Config:")
    for key, value := range config {
        fmt.Printf("  %s: %v (%T)\n", key, value, value)
    }

    // Variadic function with any
    logMessage("Server starting", "port", 8080, "debug", true)
}

func logMessage(message string, keyValues ...interface{}) {
    fmt.Printf("LOG: %s", message)
    for i := 0; i < len(keyValues); i += 2 {
        if i+1 < len(keyValues) {
            fmt.Printf(" %v=%v", keyValues[i], keyValues[i+1])
        }
    }
    fmt.Println()
}
```

---

## 4. Common Interfaces

### The error Interface

```go
package main

import (
    "errors"
    "fmt"
)

// error interface: type error interface { Error() string }

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return ValidationError{
            Field:   "age",
            Message: "cannot be negative",
        }
    }
    if age > 150 {
        return errors.New("age seems unrealistic")
    }
    return nil
}

func main() {
    ages := []int{25, -5, 200}

    for _, age := range ages {
        if err := validateAge(age); err != nil {
            fmt.Printf("Error for age %d: %s\n", age, err)
        } else {
            fmt.Printf("Age %d is valid\n", age)
        }
    }
}
```

### The Stringer Interface

```go
package main

import "fmt"

// fmt.Stringer interface: type Stringer interface { String() string }

type Person struct {
    FirstName string
    LastName  string
    Age       int
}

func (p Person) String() string {
    return fmt.Sprintf("%s %s (age %d)", p.FirstName, p.LastName, p.Age)
}

type Point struct {
    X, Y int
}

func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

type IPAddr [4]byte

func (ip IPAddr) String() string {
    return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func main() {
    person := Person{FirstName: "Alice", LastName: "Smith", Age: 30}
    fmt.Println(person)  // Uses String() method

    point := Point{X: 10, Y: 20}
    fmt.Println(point)  // (10, 20)

    ip := IPAddr{192, 168, 1, 1}
    fmt.Println(ip)  // 192.168.1.1
}
```

### The io.Reader and io.Writer Interfaces

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

func main() {
    // strings.Reader implements io.Reader
    reader := strings.NewReader("Hello, World!")

    // Read into buffer
    buf := make([]byte, 4)
    for {
        n, err := reader.Read(buf)
        if err == io.EOF {
            break
        }
        fmt.Printf("Read %d bytes: %s\n", n, buf[:n])
    }

    // bytes.Buffer implements io.Reader and io.Writer
    var buffer bytes.Buffer

    // Write to buffer
    buffer.WriteString("Hello, ")
    buffer.WriteString("Go!")

    // Read from buffer
    content, _ := io.ReadAll(&buffer)
    fmt.Println("Buffer content:", string(content))

    // Copy from reader to writer
    src := strings.NewReader("Copy this text")
    var dst bytes.Buffer
    written, _ := io.Copy(&dst, src)
    fmt.Printf("Copied %d bytes: %s\n", written, dst.String())
}
```

---

## 5. Type Assertions

Type assertions extract the concrete value from an interface:

```go
package main

import "fmt"

func main() {
    var i interface{} = "hello"

    // Type assertion (will panic if wrong type)
    s := i.(string)
    fmt.Println(s)

    // Safe type assertion (returns ok bool)
    s, ok := i.(string)
    if ok {
        fmt.Println("String value:", s)
    }

    // Check for type without panic
    f, ok := i.(float64)
    if ok {
        fmt.Println("Float value:", f)
    } else {
        fmt.Println("Not a float64")
    }

    // Working with interface values
    processValue("hello")
    processValue(42)
    processValue(3.14)
    processValue(true)
}

func processValue(v interface{}) {
    // Check multiple types
    if str, ok := v.(string); ok {
        fmt.Printf("String: %s (length: %d)\n", str, len(str))
        return
    }

    if num, ok := v.(int); ok {
        fmt.Printf("Integer: %d (squared: %d)\n", num, num*num)
        return
    }

    fmt.Printf("Unknown type: %T\n", v)
}
```

---

## 6. Type Switches

Type switches handle multiple types elegantly:

```go
package main

import "fmt"

func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d (doubled: %d)\n", v, v*2)
    case float64:
        fmt.Printf("Float64: %f (squared: %f)\n", v, v*v)
    case string:
        fmt.Printf("String: %q (length: %d)\n", v, len(v))
    case bool:
        fmt.Printf("Boolean: %t (negated: %t)\n", v, !v)
    case []int:
        fmt.Printf("Int slice with %d elements: %v\n", len(v), v)
    case nil:
        fmt.Println("Nil value")
    default:
        fmt.Printf("Unknown type: %T with value: %v\n", v, v)
    }
}

func main() {
    describe(42)
    describe(3.14)
    describe("hello")
    describe(true)
    describe([]int{1, 2, 3})
    describe(nil)
    describe(map[string]int{"a": 1})
}
```

### Type Switch with Multiple Types

```go
package main

import "fmt"

func printNumeric(v interface{}) {
    switch v := v.(type) {
    case int, int8, int16, int32, int64:
        fmt.Printf("Signed integer: %v\n", v)
    case uint, uint8, uint16, uint32, uint64:
        fmt.Printf("Unsigned integer: %v\n", v)
    case float32, float64:
        fmt.Printf("Floating point: %v\n", v)
    default:
        fmt.Printf("Not a numeric type: %T\n", v)
    }
}

func main() {
    printNumeric(42)
    printNumeric(uint(100))
    printNumeric(3.14)
    printNumeric("hello")
}
```

---

## 7. Interface Design Patterns

### Accept Interfaces, Return Structs

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

// Accept interface
func process(r io.Reader) (string, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

// Return concrete type
type Result struct {
    Data    string
    Success bool
}

func getResult() *Result {
    return &Result{
        Data:    "result data",
        Success: true,
    }
}

func main() {
    // Can pass any io.Reader
    reader := strings.NewReader("Hello, World!")
    data, _ := process(reader)
    fmt.Println(data)

    result := getResult()
    fmt.Printf("%+v\n", result)
}
```

### Interface Segregation

```go
package main

import "fmt"

// Small, focused interfaces
type Saver interface {
    Save() error
}

type Loader interface {
    Load() error
}

type Deleter interface {
    Delete() error
}

// Composed interface
type Repository interface {
    Saver
    Loader
    Deleter
}

// User only needs to save
type UserService struct{}

func (s *UserService) ProcessUser(saver Saver) error {
    // Only need Save capability
    return saver.Save()
}

// Document implementation
type Document struct {
    ID      string
    Content string
}

func (d *Document) Save() error {
    fmt.Printf("Saving document %s\n", d.ID)
    return nil
}

func (d *Document) Load() error {
    fmt.Printf("Loading document %s\n", d.ID)
    return nil
}

func (d *Document) Delete() error {
    fmt.Printf("Deleting document %s\n", d.ID)
    return nil
}

func main() {
    doc := &Document{ID: "123", Content: "Hello"}

    service := &UserService{}
    service.ProcessUser(doc) // Only uses Saver interface
}
```

### Dependency Injection with Interfaces

```go
package main

import "fmt"

// Database interface
type Database interface {
    Query(query string) ([]map[string]interface{}, error)
    Execute(query string) error
}

// UserRepository depends on Database interface
type UserRepository struct {
    db Database
}

func NewUserRepository(db Database) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(id int) (map[string]interface{}, error) {
    results, err := r.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d", id))
    if err != nil {
        return nil, err
    }
    if len(results) == 0 {
        return nil, fmt.Errorf("user not found")
    }
    return results[0], nil
}

// Real implementation
type PostgresDB struct {
    connectionString string
}

func (p *PostgresDB) Query(query string) ([]map[string]interface{}, error) {
    fmt.Printf("PostgreSQL executing: %s\n", query)
    // Real implementation would query PostgreSQL
    return []map[string]interface{}{{"id": 1, "name": "Alice"}}, nil
}

func (p *PostgresDB) Execute(query string) error {
    fmt.Printf("PostgreSQL executing: %s\n", query)
    return nil
}

// Mock implementation for testing
type MockDB struct {
    users []map[string]interface{}
}

func (m *MockDB) Query(query string) ([]map[string]interface{}, error) {
    return m.users, nil
}

func (m *MockDB) Execute(query string) error {
    return nil
}

func main() {
    // Production
    prodDB := &PostgresDB{connectionString: "postgres://localhost/app"}
    prodRepo := NewUserRepository(prodDB)
    user, _ := prodRepo.GetUser(1)
    fmt.Printf("Production user: %v\n", user)

    // Testing
    mockDB := &MockDB{
        users: []map[string]interface{}{
            {"id": 1, "name": "Test User"},
        },
    }
    testRepo := NewUserRepository(mockDB)
    testUser, _ := testRepo.GetUser(1)
    fmt.Printf("Test user: %v\n", testUser)
}
```

---

## Exercises

### Exercise 1: Payment Processor
Implement different payment methods using interfaces.

```go
package main

import "fmt"

// TODO: Define PaymentMethod interface with Process(amount float64) error

// TODO: Implement CreditCard struct with CardNumber, ExpiryDate
// Implement Process method

// TODO: Implement PayPal struct with Email
// Implement Process method

// TODO: Implement BankTransfer struct with AccountNumber, RoutingNumber
// Implement Process method

func processPayment(pm PaymentMethod, amount float64) {
    if err := pm.Process(amount); err != nil {
        fmt.Println("Payment failed:", err)
    } else {
        fmt.Printf("Payment of $%.2f processed successfully\n", amount)
    }
}

func main() {
    cc := CreditCard{CardNumber: "4111111111111111", ExpiryDate: "12/25"}
    pp := PayPal{Email: "user@example.com"}
    bt := BankTransfer{AccountNumber: "123456789", RoutingNumber: "987654321"}

    processPayment(cc, 100.00)
    processPayment(pp, 50.00)
    processPayment(bt, 1000.00)
}
```

### Exercise 2: Logger Interface
Create a flexible logging system.

```go
package main

import (
    "fmt"
    "time"
)

// TODO: Define Logger interface with:
// - Info(message string)
// - Error(message string)
// - Debug(message string)

// TODO: Implement ConsoleLogger that prints to console with timestamps

// TODO: Implement FileLogger that "writes" to a file (just print for now)

// TODO: Implement MultiLogger that sends to multiple loggers

func main() {
    console := &ConsoleLogger{}
    file := &FileLogger{Filename: "app.log"}
    multi := &MultiLogger{Loggers: []Logger{console, file}}

    multi.Info("Application started")
    multi.Debug("Loading configuration")
    multi.Error("Connection timeout")
}
```

### Exercise 3: Type Assertion Practice
Work with a heterogeneous collection.

```go
package main

import "fmt"

func main() {
    // Mixed type collection
    items := []interface{}{
        42,
        "hello",
        3.14,
        true,
        []int{1, 2, 3},
        map[string]int{"a": 1, "b": 2},
    }

    // TODO: Implement processItems that:
    // - Sums all integers
    // - Concatenates all strings
    // - Multiplies all floats
    // - Counts booleans (true count, false count)
    // - Sums all elements in int slices
    // - Counts total keys in maps

    // Expected output:
    // Integer sum: 42
    // String concat: hello
    // Float product: 3.14
    // True count: 1, False count: 0
    // Slice sum: 6
    // Map keys: 2

    processItems(items)
}

func processItems(items []interface{}) {
    // TODO: Implement using type switch
}
```

### Exercise 4: Serializer Interface
Create serializers for different formats.

```go
package main

import "fmt"

type User struct {
    ID       int
    Username string
    Email    string
}

// TODO: Define Serializer interface with Serialize(v interface{}) (string, error)

// TODO: Implement JSONSerializer

// TODO: Implement XMLSerializer

// TODO: Implement YAMLSerializer (simplified)

func main() {
    user := User{ID: 1, Username: "alice", Email: "alice@example.com"}

    serializers := map[string]Serializer{
        "JSON": &JSONSerializer{},
        "XML":  &XMLSerializer{},
        "YAML": &YAMLSerializer{},
    }

    for name, serializer := range serializers {
        output, err := serializer.Serialize(user)
        if err != nil {
            fmt.Printf("%s Error: %v\n", name, err)
            continue
        }
        fmt.Printf("%s:\n%s\n\n", name, output)
    }
}
```

### Exercise 5: Storage Interface
Implement a generic storage system.

```go
package main

import "fmt"

// TODO: Define Storage interface with:
// - Get(key string) (interface{}, error)
// - Set(key string, value interface{}) error
// - Delete(key string) error
// - Exists(key string) bool

// TODO: Implement MemoryStorage using a map

// TODO: Implement CacheStorage that wraps another Storage and adds caching behavior
// (e.g., print "cache hit" or "cache miss")

func main() {
    // Memory storage
    memory := NewMemoryStorage()
    memory.Set("user:1", map[string]string{"name": "Alice"})
    memory.Set("user:2", map[string]string{"name": "Bob"})

    value, _ := memory.Get("user:1")
    fmt.Println("User 1:", value)

    fmt.Println("User 1 exists:", memory.Exists("user:1"))
    fmt.Println("User 3 exists:", memory.Exists("user:3"))

    // With caching
    cached := NewCacheStorage(memory)
    cached.Get("user:1") // cache miss
    cached.Get("user:1") // cache hit
}
```

---

## Summary

In this tutorial, you learned:
- How to define interfaces as method sets
- Implicit implementation of interfaces
- Working with empty interface (`any`)
- Common interfaces: error, Stringer, io.Reader/Writer
- Type assertions and type switches
- Interface design patterns for flexible code

---

**Next:** [04-error-handling.md](04-error-handling.md) - Learn about error handling in Go

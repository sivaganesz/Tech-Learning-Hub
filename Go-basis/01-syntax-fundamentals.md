# Go Syntax Fundamentals

## Learning Objectives

By the end of this tutorial, you will be able to:
- Declare and use variables with `var` and `:=` syntax
- Work with Go's basic types (int, string, bool, float)
- Write functions with various parameter and return patterns
- Use control flow structures (if/else, for, switch)
- Create and manipulate arrays, slices, and maps

---

## 1. Variables

### Variable Declaration with `var`

Go provides multiple ways to declare variables. The `var` keyword is the most explicit:

```go
package main

import "fmt"

func main() {
    // Declare with explicit type
    var name string
    name = "Alice"

    // Declare and initialize
    var age int = 30

    // Type inference (Go infers the type from value)
    var city = "New York"

    // Multiple declarations
    var (
        country string = "USA"
        code    int    = 1
        active  bool   = true
    )

    fmt.Println(name, age, city, country, code, active)
}
```

### Short Variable Declaration with `:=`

Inside functions, you can use the short declaration operator:

```go
package main

import "fmt"

func main() {
    // Short declaration (type inferred)
    name := "Bob"
    age := 25
    height := 5.9
    isStudent := true

    // Multiple short declarations
    first, last := "John", "Doe"

    fmt.Println(name, age, height, isStudent)
    fmt.Println(first, last)
}
```

### Constants

Constants are immutable values declared with `const`:

```go
package main

import "fmt"

const Pi = 3.14159
const AppName = "MyApp"

// Multiple constants
const (
    StatusOK    = 200
    StatusError = 500
    MaxRetries  = 3
)

// iota for enumerations
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
    Wednesday      // 3
    Thursday       // 4
    Friday         // 5
    Saturday       // 6
)

func main() {
    fmt.Println("Pi:", Pi)
    fmt.Println("Sunday is day:", Sunday)
    fmt.Println("Friday is day:", Friday)
}
```

### Zero Values

Uninitialized variables get zero values:

```go
package main

import "fmt"

func main() {
    var i int      // 0
    var f float64  // 0.0
    var b bool     // false
    var s string   // "" (empty string)

    fmt.Printf("int: %d, float: %f, bool: %t, string: %q\n", i, f, b, s)
}
```

---

## 2. Basic Types

### Numeric Types

```go
package main

import "fmt"

func main() {
    // Integers
    var i int = 42           // Platform dependent (32 or 64 bit)
    var i8 int8 = 127        // -128 to 127
    var i16 int16 = 32767    // -32768 to 32767
    var i32 int32 = 2147483647
    var i64 int64 = 9223372036854775807

    // Unsigned integers
    var u uint = 42
    var u8 uint8 = 255       // 0 to 255 (byte is alias)
    var u16 uint16 = 65535
    var u32 uint32 = 4294967295
    var u64 uint64 = 18446744073709551615

    // Floating point
    var f32 float32 = 3.14
    var f64 float64 = 3.141592653589793

    // Complex numbers
    var c64 complex64 = 1 + 2i
    var c128 complex128 = 1 + 2i

    fmt.Println(i, i8, i16, i32, i64)
    fmt.Println(u, u8, u16, u32, u64)
    fmt.Println(f32, f64)
    fmt.Println(c64, c128)
}
```

### Strings

```go
package main

import (
    "fmt"
    "strings"
)

func main() {
    // String declaration
    greeting := "Hello, World!"

    // Multi-line strings (raw strings)
    multiline := `This is a
    multi-line
    string`

    // String operations
    fmt.Println("Length:", len(greeting))
    fmt.Println("Uppercase:", strings.ToUpper(greeting))
    fmt.Println("Contains 'World':", strings.Contains(greeting, "World"))
    fmt.Println("Split:", strings.Split(greeting, ", "))

    // String concatenation
    first := "Go"
    second := "Lang"
    combined := first + " " + second
    fmt.Println(combined)

    // Accessing characters (bytes)
    fmt.Println("First byte:", greeting[0])        // 72 (ASCII for 'H')
    fmt.Println("First char:", string(greeting[0])) // "H"

    fmt.Println(multiline)
}
```

### Booleans

```go
package main

import "fmt"

func main() {
    isActive := true
    isComplete := false

    // Boolean operations
    fmt.Println("AND:", isActive && isComplete)  // false
    fmt.Println("OR:", isActive || isComplete)   // true
    fmt.Println("NOT:", !isActive)               // false

    // Comparison operators return booleans
    a, b := 10, 20
    fmt.Println("a == b:", a == b)   // false
    fmt.Println("a != b:", a != b)   // true
    fmt.Println("a < b:", a < b)     // true
    fmt.Println("a > b:", a > b)     // false
    fmt.Println("a <= b:", a <= b)   // true
    fmt.Println("a >= b:", a >= b)   // false
}
```

---

## 3. Functions

### Basic Functions

```go
package main

import "fmt"

// Function with no parameters and no return
func sayHello() {
    fmt.Println("Hello!")
}

// Function with parameters
func greet(name string) {
    fmt.Println("Hello,", name)
}

// Function with return value
func add(a int, b int) int {
    return a + b
}

// Shorthand for same-type parameters
func multiply(a, b int) int {
    return a * b
}

func main() {
    sayHello()
    greet("Alice")
    result := add(5, 3)
    fmt.Println("5 + 3 =", result)
    fmt.Println("4 * 5 =", multiply(4, 5))
}
```

### Multiple Return Values

```go
package main

import (
    "errors"
    "fmt"
)

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Swap values
func swap(a, b string) (string, string) {
    return b, a
}

func main() {
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("10 / 2 =", result)
    }

    // Ignoring a return value with _
    result2, _ := divide(20, 4)
    fmt.Println("20 / 4 =", result2)

    x, y := swap("first", "second")
    fmt.Println(x, y)
}
```

### Named Return Values

```go
package main

import "fmt"

// Named return values
func rectangle(width, height float64) (area, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return // naked return
}

// Named returns with explicit return
func circle(radius float64) (area, circumference float64) {
    area = 3.14159 * radius * radius
    circumference = 2 * 3.14159 * radius
    return area, circumference
}

func main() {
    a, p := rectangle(5, 3)
    fmt.Printf("Rectangle: area=%.2f, perimeter=%.2f\n", a, p)

    ca, cc := circle(5)
    fmt.Printf("Circle: area=%.2f, circumference=%.2f\n", ca, cc)
}
```

### Variadic Functions

```go
package main

import "fmt"

// Variadic function (variable number of arguments)
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

// Mix regular and variadic parameters
func printf(format string, args ...interface{}) {
    fmt.Printf(format, args...)
}

func main() {
    fmt.Println(sum(1, 2, 3))       // 6
    fmt.Println(sum(1, 2, 3, 4, 5)) // 15

    // Passing slice to variadic function
    nums := []int{10, 20, 30}
    fmt.Println(sum(nums...)) // 60

    printf("Name: %s, Age: %d\n", "Alice", 30)
}
```

### Anonymous Functions and Closures

```go
package main

import "fmt"

func main() {
    // Anonymous function
    add := func(a, b int) int {
        return a + b
    }
    fmt.Println(add(3, 4))

    // Immediately invoked function
    result := func(x int) int {
        return x * x
    }(5)
    fmt.Println("Square of 5:", result)

    // Closure (captures outer variable)
    counter := 0
    increment := func() int {
        counter++
        return counter
    }

    fmt.Println(increment()) // 1
    fmt.Println(increment()) // 2
    fmt.Println(increment()) // 3
}
```

---

## 4. Control Flow

### If/Else Statements

```go
package main

import "fmt"

func main() {
    age := 20

    // Basic if
    if age >= 18 {
        fmt.Println("Adult")
    }

    // If-else
    if age >= 21 {
        fmt.Println("Can drink in US")
    } else {
        fmt.Println("Cannot drink in US")
    }

    // If-else if-else
    score := 85
    if score >= 90 {
        fmt.Println("A")
    } else if score >= 80 {
        fmt.Println("B")
    } else if score >= 70 {
        fmt.Println("C")
    } else {
        fmt.Println("F")
    }

    // If with initialization statement
    if num := 10; num%2 == 0 {
        fmt.Println(num, "is even")
    } else {
        fmt.Println(num, "is odd")
    }
}
```

### For Loops

Go has only one looping construct: `for`

```go
package main

import "fmt"

func main() {
    // Standard for loop
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    // While-style loop
    count := 0
    for count < 3 {
        fmt.Println("count:", count)
        count++
    }

    // Infinite loop (use break to exit)
    sum := 0
    for {
        sum++
        if sum >= 5 {
            break
        }
    }
    fmt.Println("Sum:", sum)

    // Continue statement
    for i := 0; i < 10; i++ {
        if i%2 == 0 {
            continue // skip even numbers
        }
        fmt.Println("Odd:", i)
    }

    // Range over slice
    fruits := []string{"apple", "banana", "cherry"}
    for index, value := range fruits {
        fmt.Printf("%d: %s\n", index, value)
    }

    // Range ignoring index
    for _, fruit := range fruits {
        fmt.Println(fruit)
    }

    // Range over map
    ages := map[string]int{"Alice": 30, "Bob": 25}
    for name, age := range ages {
        fmt.Printf("%s is %d years old\n", name, age)
    }

    // Range over string (runes)
    for i, r := range "Go" {
        fmt.Printf("%d: %c\n", i, r)
    }
}
```

### Switch Statements

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Basic switch
    day := "Monday"
    switch day {
    case "Monday":
        fmt.Println("Start of work week")
    case "Friday":
        fmt.Println("End of work week")
    case "Saturday", "Sunday": // Multiple values
        fmt.Println("Weekend!")
    default:
        fmt.Println("Midweek")
    }

    // Switch with no expression (like if-else chain)
    score := 85
    switch {
    case score >= 90:
        fmt.Println("A")
    case score >= 80:
        fmt.Println("B")
    case score >= 70:
        fmt.Println("C")
    default:
        fmt.Println("F")
    }

    // Switch with initialization
    switch hour := time.Now().Hour(); {
    case hour < 12:
        fmt.Println("Good morning!")
    case hour < 17:
        fmt.Println("Good afternoon!")
    default:
        fmt.Println("Good evening!")
    }

    // Fallthrough
    num := 5
    switch num {
    case 5:
        fmt.Println("Five")
        fallthrough
    case 6:
        fmt.Println("Five or Six")
    }
}
```

---

## 5. Arrays

Arrays have a fixed size defined at compile time:

```go
package main

import "fmt"

func main() {
    // Array declaration
    var numbers [5]int
    fmt.Println(numbers) // [0 0 0 0 0]

    // Array with initialization
    fruits := [3]string{"apple", "banana", "cherry"}
    fmt.Println(fruits)

    // Let compiler count elements
    colors := [...]string{"red", "green", "blue"}
    fmt.Println(colors)

    // Access and modify elements
    numbers[0] = 10
    numbers[1] = 20
    fmt.Println(numbers[0]) // 10

    // Array length
    fmt.Println("Length:", len(fruits))

    // Iterate over array
    for i, fruit := range fruits {
        fmt.Printf("%d: %s\n", i, fruit)
    }

    // Multi-dimensional array
    matrix := [2][3]int{
        {1, 2, 3},
        {4, 5, 6},
    }
    fmt.Println(matrix)
    fmt.Println(matrix[1][2]) // 6
}
```

---

## 6. Slices

Slices are dynamic, flexible views into arrays:

```go
package main

import "fmt"

func main() {
    // Create slice with make
    numbers := make([]int, 3, 5) // length 3, capacity 5
    fmt.Println(numbers, len(numbers), cap(numbers))

    // Slice literal
    fruits := []string{"apple", "banana", "cherry"}
    fmt.Println(fruits)

    // Append elements
    fruits = append(fruits, "date")
    fruits = append(fruits, "elderberry", "fig")
    fmt.Println(fruits)

    // Slice from array or slice
    all := []int{0, 1, 2, 3, 4, 5}
    some := all[1:4] // [1 2 3] (index 1 to 3)
    fmt.Println(some)

    // Slice operations
    first := all[:3]   // [0 1 2]
    last := all[3:]    // [3 4 5]
    copyAll := all[:]  // [0 1 2 3 4 5]
    fmt.Println(first, last, copyAll)

    // Copy slice
    src := []int{1, 2, 3}
    dst := make([]int, len(src))
    copied := copy(dst, src)
    fmt.Println("Copied", copied, "elements:", dst)

    // Iterate over slice
    for i, v := range fruits {
        fmt.Printf("%d: %s\n", i, v)
    }

    // Nil slice
    var nilSlice []int
    fmt.Println(nilSlice == nil)           // true
    fmt.Println(len(nilSlice), cap(nilSlice)) // 0 0

    // Empty slice
    emptySlice := []int{}
    fmt.Println(emptySlice == nil) // false
}
```

### Slice Internals

```go
package main

import "fmt"

func main() {
    // Slices share underlying array
    original := []int{1, 2, 3, 4, 5}
    slice := original[1:4]

    slice[0] = 100 // modifies original too!
    fmt.Println("Original:", original) // [1 100 3 4 5]
    fmt.Println("Slice:", slice)       // [100 3 4]

    // To avoid sharing, use copy
    original2 := []int{1, 2, 3, 4, 5}
    slice2 := make([]int, 3)
    copy(slice2, original2[1:4])

    slice2[0] = 100
    fmt.Println("Original2:", original2) // [1 2 3 4 5]
    fmt.Println("Slice2:", slice2)       // [100 3 4]
}
```

---

## 7. Maps

Maps are key-value pairs:

```go
package main

import "fmt"

func main() {
    // Create map with make
    ages := make(map[string]int)
    ages["Alice"] = 30
    ages["Bob"] = 25
    fmt.Println(ages)

    // Map literal
    scores := map[string]int{
        "Alice": 95,
        "Bob":   87,
        "Carol": 92,
    }
    fmt.Println(scores)

    // Access values
    aliceScore := scores["Alice"]
    fmt.Println("Alice's score:", aliceScore)

    // Check if key exists
    value, exists := scores["Dave"]
    if exists {
        fmt.Println("Dave's score:", value)
    } else {
        fmt.Println("Dave not found")
    }

    // Delete key
    delete(scores, "Bob")
    fmt.Println("After delete:", scores)

    // Length
    fmt.Println("Number of entries:", len(scores))

    // Iterate over map
    for name, score := range scores {
        fmt.Printf("%s: %d\n", name, score)
    }

    // Nested maps
    users := map[string]map[string]string{
        "user1": {
            "name":  "Alice",
            "email": "alice@example.com",
        },
        "user2": {
            "name":  "Bob",
            "email": "bob@example.com",
        },
    }
    fmt.Println(users["user1"]["email"])

    // Nil map (read ok, write panics)
    var nilMap map[string]int
    fmt.Println(nilMap["key"]) // 0 (zero value)
    // nilMap["key"] = 1       // panic!
}
```

---

## Exercises

### Exercise 1: Temperature Converter
Write a program that converts temperatures between Celsius and Fahrenheit.

```go
package main

import "fmt"

// TODO: Implement these functions
func celsiusToFahrenheit(c float64) float64 {
    // Formula: F = C * 9/5 + 32
    return 0
}

func fahrenheitToCelsius(f float64) float64 {
    // Formula: C = (F - 32) * 5/9
    return 0
}

func main() {
    c := 100.0
    f := 212.0

    fmt.Printf("%.2f C = %.2f F\n", c, celsiusToFahrenheit(c))
    fmt.Printf("%.2f F = %.2f C\n", f, fahrenheitToCelsius(f))
}
```

### Exercise 2: FizzBuzz
Write FizzBuzz for numbers 1-100.

```go
package main

import "fmt"

func fizzBuzz(n int) string {
    // TODO: Return "FizzBuzz" if divisible by 3 and 5
    // Return "Fizz" if divisible by 3
    // Return "Buzz" if divisible by 5
    // Return the number as string otherwise
    return ""
}

func main() {
    for i := 1; i <= 100; i++ {
        fmt.Println(fizzBuzz(i))
    }
}
```

### Exercise 3: Word Counter
Count word frequency in a string.

```go
package main

import (
    "fmt"
    "strings"
)

func wordCount(s string) map[string]int {
    // TODO: Split string into words and count each word
    // Hint: Use strings.Fields() to split by whitespace
    return nil
}

func main() {
    text := "the quick brown fox jumps over the lazy dog the fox"
    counts := wordCount(text)

    for word, count := range counts {
        fmt.Printf("%s: %d\n", word, count)
    }
}
```

### Exercise 4: Slice Operations
Implement common slice operations.

```go
package main

import "fmt"

// Remove element at index i
func removeAt(slice []int, i int) []int {
    // TODO: Implement
    return nil
}

// Insert value at index i
func insertAt(slice []int, i int, value int) []int {
    // TODO: Implement
    return nil
}

// Reverse slice in place
func reverse(slice []int) {
    // TODO: Implement
}

func main() {
    nums := []int{1, 2, 3, 4, 5}

    nums = removeAt(nums, 2)
    fmt.Println("After remove:", nums) // [1 2 4 5]

    nums = insertAt(nums, 2, 10)
    fmt.Println("After insert:", nums) // [1 2 10 4 5]

    reverse(nums)
    fmt.Println("After reverse:", nums) // [5 4 10 2 1]
}
```

### Exercise 5: Calculator
Build a simple calculator with functions.

```go
package main

import (
    "errors"
    "fmt"
)

func calculate(a, b float64, op string) (float64, error) {
    // TODO: Implement +, -, *, /
    // Return error for division by zero and unknown operator
    return 0, errors.New("not implemented")
}

func main() {
    operations := []struct {
        a, b float64
        op   string
    }{
        {10, 5, "+"},
        {10, 5, "-"},
        {10, 5, "*"},
        {10, 5, "/"},
        {10, 0, "/"},
        {10, 5, "%"},
    }

    for _, op := range operations {
        result, err := calculate(op.a, op.b, op.op)
        if err != nil {
            fmt.Printf("%.2f %s %.2f = Error: %v\n", op.a, op.op, op.b, err)
        } else {
            fmt.Printf("%.2f %s %.2f = %.2f\n", op.a, op.op, op.b, result)
        }
    }
}
```

---

## Summary

In this tutorial, you learned:
- How to declare variables using `var` and `:=`
- Go's basic types: int, float, string, bool
- Function syntax including multiple returns and named returns
- Control flow with if/else, for loops, and switch
- Working with arrays, slices, and maps

---

**Next:** [02-structs-and-methods.md](02-structs-and-methods.md) - Learn about structs and methods in Go

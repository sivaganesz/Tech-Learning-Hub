# Structs and Methods in Go

## Learning Objectives

By the end of this tutorial, you will be able to:
- Define and instantiate structs
- Create methods with value and pointer receivers
- Use struct embedding for composition
- Apply struct tags for JSON and database mapping
- Understand when to use pointer vs value receivers

---

## 1. Defining Structs

Structs are typed collections of fields used to group data together:

```go
package main

import "fmt"

// Define a struct
type Person struct {
    FirstName string
    LastName  string
    Age       int
    Email     string
}

// Struct with different field types
type Product struct {
    ID          int
    Name        string
    Price       float64
    InStock     bool
    Categories  []string
    Metadata    map[string]string
}

func main() {
    // Create a struct instance
    p := Person{
        FirstName: "Alice",
        LastName:  "Smith",
        Age:       30,
        Email:     "alice@example.com",
    }

    fmt.Println(p)
    fmt.Printf("%+v\n", p) // Print with field names
}
```

---

## 2. Creating Instances

### Different Ways to Create Structs

```go
package main

import "fmt"

type User struct {
    ID       int
    Username string
    Email    string
    Active   bool
}

func main() {
    // Method 1: Named fields (recommended)
    user1 := User{
        ID:       1,
        Username: "alice",
        Email:    "alice@example.com",
        Active:   true,
    }

    // Method 2: Positional (not recommended - fragile)
    user2 := User{2, "bob", "bob@example.com", true}

    // Method 3: Zero value, then assign
    var user3 User
    user3.ID = 3
    user3.Username = "carol"
    user3.Email = "carol@example.com"
    user3.Active = true

    // Method 4: Using new (returns pointer)
    user4 := new(User)
    user4.ID = 4
    user4.Username = "dave"

    // Method 5: Pointer with & (most common)
    user5 := &User{
        ID:       5,
        Username: "eve",
        Email:    "eve@example.com",
        Active:   true,
    }

    fmt.Printf("user1: %+v\n", user1)
    fmt.Printf("user2: %+v\n", user2)
    fmt.Printf("user3: %+v\n", user3)
    fmt.Printf("user4: %+v\n", user4)
    fmt.Printf("user5: %+v\n", user5)
}
```

### Accessing and Modifying Fields

```go
package main

import "fmt"

type Book struct {
    Title    string
    Author   string
    Pages    int
    ISBN     string
}

func main() {
    book := Book{
        Title:  "The Go Programming Language",
        Author: "Alan Donovan",
        Pages:  380,
        ISBN:   "978-0134190440",
    }

    // Access fields
    fmt.Println("Title:", book.Title)
    fmt.Println("Author:", book.Author)

    // Modify fields
    book.Pages = 400
    fmt.Println("Updated pages:", book.Pages)

    // Pointer to struct
    bookPtr := &book
    bookPtr.Title = "New Title" // Automatic dereferencing
    fmt.Println("New title:", book.Title)
}
```

### Anonymous Structs

```go
package main

import "fmt"

func main() {
    // Anonymous struct (useful for one-off data structures)
    person := struct {
        Name string
        Age  int
    }{
        Name: "Alice",
        Age:  30,
    }

    fmt.Printf("%+v\n", person)

    // Common use: test cases
    testCases := []struct {
        input    int
        expected int
    }{
        {1, 2},
        {2, 4},
        {3, 6},
    }

    for _, tc := range testCases {
        result := tc.input * 2
        fmt.Printf("Input: %d, Expected: %d, Got: %d\n",
            tc.input, tc.expected, result)
    }
}
```

---

## 3. Methods

Methods are functions with a receiver argument:

### Value Receivers

```go
package main

import (
    "fmt"
    "math"
)

type Rectangle struct {
    Width  float64
    Height float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Circumference() float64 {
    return 2 * math.Pi * c.Radius
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    fmt.Printf("Rectangle: Area=%.2f, Perimeter=%.2f\n",
        rect.Area(), rect.Perimeter())

    circle := Circle{Radius: 5}
    fmt.Printf("Circle: Area=%.2f, Circumference=%.2f\n",
        circle.Area(), circle.Circumference())
}
```

### Pointer Receivers

Use pointer receivers when you need to modify the struct:

```go
package main

import "fmt"

type Counter struct {
    value int
}

// Pointer receiver - can modify the struct
func (c *Counter) Increment() {
    c.value++
}

func (c *Counter) Add(n int) {
    c.value += n
}

func (c *Counter) Reset() {
    c.value = 0
}

// Value receiver - for read-only operations
func (c Counter) Value() int {
    return c.value
}

type BankAccount struct {
    Owner   string
    Balance float64
}

func (a *BankAccount) Deposit(amount float64) {
    if amount > 0 {
        a.Balance += amount
    }
}

func (a *BankAccount) Withdraw(amount float64) bool {
    if amount > 0 && amount <= a.Balance {
        a.Balance -= amount
        return true
    }
    return false
}

func (a BankAccount) GetBalance() float64 {
    return a.Balance
}

func main() {
    counter := Counter{}
    counter.Increment()
    counter.Increment()
    counter.Add(10)
    fmt.Println("Counter:", counter.Value()) // 12

    counter.Reset()
    fmt.Println("After reset:", counter.Value()) // 0

    account := BankAccount{Owner: "Alice", Balance: 100}
    account.Deposit(50)
    fmt.Printf("After deposit: $%.2f\n", account.GetBalance())

    if account.Withdraw(30) {
        fmt.Printf("After withdrawal: $%.2f\n", account.GetBalance())
    }
}
```

### When to Use Pointer vs Value Receivers

```go
package main

import "fmt"

type SmallStruct struct {
    X, Y int
}

type LargeStruct struct {
    Data [1000]int
    Name string
}

// Value receiver - small struct, no modification needed
func (s SmallStruct) Sum() int {
    return s.X + s.Y
}

// Pointer receiver - large struct (avoid copying)
func (l *LargeStruct) Process() {
    // Process data
}

// Pointer receiver - needs to modify
func (l *LargeStruct) SetName(name string) {
    l.Name = name
}

/*
Guidelines for choosing receiver type:

1. Use POINTER receiver when:
   - Method needs to modify the receiver
   - Struct is large (avoid copying)
   - Consistency (if any method uses pointer, all should)

2. Use VALUE receiver when:
   - Method doesn't modify receiver
   - Struct is small (like time.Time)
   - You want immutability
*/

func main() {
    small := SmallStruct{X: 10, Y: 20}
    fmt.Println("Sum:", small.Sum())

    large := &LargeStruct{}
    large.SetName("Test")
    fmt.Println("Name:", large.Name)
}
```

---

## 4. Struct Embedding (Composition)

Go uses composition instead of inheritance:

```go
package main

import "fmt"

// Base struct
type Address struct {
    Street  string
    City    string
    Country string
    ZipCode string
}

func (a Address) FullAddress() string {
    return fmt.Sprintf("%s, %s, %s %s",
        a.Street, a.City, a.Country, a.ZipCode)
}

// Embedding Address in Person
type Person struct {
    Name    string
    Age     int
    Address // Embedded (anonymous field)
}

// Embedding Address in Company
type Company struct {
    Name    string
    Address // Embedded
}

func main() {
    person := Person{
        Name: "Alice",
        Age:  30,
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
            ZipCode: "10001",
        },
    }

    // Access embedded fields directly
    fmt.Println("Name:", person.Name)
    fmt.Println("City:", person.City)         // Promoted field
    fmt.Println("Full:", person.FullAddress()) // Promoted method

    // Or access through the embedded type
    fmt.Println("Street:", person.Address.Street)

    company := Company{
        Name: "Acme Inc",
        Address: Address{
            Street:  "456 Business Ave",
            City:    "San Francisco",
            Country: "USA",
            ZipCode: "94102",
        },
    }

    fmt.Println("\nCompany:", company.Name)
    fmt.Println("Location:", company.FullAddress())
}
```

### Multiple Embedding and Field Shadowing

```go
package main

import "fmt"

type Engine struct {
    Power    int
    FuelType string
}

func (e Engine) Start() {
    fmt.Println("Engine starting...")
}

type Wheels struct {
    Count int
    Size  int
}

// Multiple embedding
type Car struct {
    Brand string
    Model string
    Engine
    Wheels
}

// Override embedded method
func (c Car) Start() {
    fmt.Println("Car starting with key...")
    c.Engine.Start() // Call embedded method explicitly
}

type ContactInfo struct {
    Email string
    Phone string
}

type Employee struct {
    Name  string
    Email string // Same field name as in ContactInfo
    ContactInfo
}

func main() {
    car := Car{
        Brand:  "Toyota",
        Model:  "Camry",
        Engine: Engine{Power: 200, FuelType: "Gasoline"},
        Wheels: Wheels{Count: 4, Size: 17},
    }

    fmt.Printf("Car: %s %s\n", car.Brand, car.Model)
    fmt.Printf("Power: %d HP, Fuel: %s\n", car.Power, car.FuelType)
    fmt.Printf("Wheels: %d x %d inch\n", car.Count, car.Size)
    car.Start()

    // Field shadowing
    emp := Employee{
        Name:  "Alice",
        Email: "alice@work.com", // This shadows ContactInfo.Email
        ContactInfo: ContactInfo{
            Email: "alice@personal.com",
            Phone: "555-1234",
        },
    }

    fmt.Println("\nEmployee:", emp.Name)
    fmt.Println("Work Email:", emp.Email)                  // alice@work.com
    fmt.Println("Personal Email:", emp.ContactInfo.Email)  // alice@personal.com
    fmt.Println("Phone:", emp.Phone)                       // Promoted
}
```

---

## 5. Struct Tags

Struct tags provide metadata for fields:

### JSON Tags

```go
package main

import (
    "encoding/json"
    "fmt"
)

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"-"`                    // Ignored in JSON
    CreatedAt string `json:"created_at,omitempty"` // Omit if empty
    Age       int    `json:"age,string"`           // Output as string
}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func main() {
    user := User{
        ID:       1,
        Username: "alice",
        Email:    "alice@example.com",
        Password: "secret123",
        Age:      30,
    }

    // Struct to JSON
    jsonData, err := json.Marshal(user)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("JSON:", string(jsonData))

    // Pretty print JSON
    jsonPretty, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("Pretty JSON:")
    fmt.Println(string(jsonPretty))

    // JSON to struct
    jsonStr := `{"id": 2, "username": "bob", "email": "bob@example.com"}`
    var newUser User
    err = json.Unmarshal([]byte(jsonStr), &newUser)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("\nParsed user: %+v\n", newUser)
}
```

### Database Tags

```go
package main

import "fmt"

// GORM-style tags
type Product struct {
    ID          uint    `gorm:"primaryKey;autoIncrement"`
    Name        string  `gorm:"size:255;not null"`
    Description string  `gorm:"type:text"`
    Price       float64 `gorm:"not null;default:0"`
    SKU         string  `gorm:"uniqueIndex;size:50"`
    CategoryID  uint    `gorm:"index"`
    CreatedAt   string  `gorm:"autoCreateTime"`
    UpdatedAt   string  `gorm:"autoUpdateTime"`
    DeletedAt   string  `gorm:"index"`
}

// SQL-style tags
type Order struct {
    ID         int     `db:"id" json:"id"`
    CustomerID int     `db:"customer_id" json:"customer_id"`
    Total      float64 `db:"total" json:"total"`
    Status     string  `db:"status" json:"status"`
}

func main() {
    product := Product{
        ID:    1,
        Name:  "Widget",
        Price: 29.99,
        SKU:   "WGT-001",
    }

    fmt.Printf("%+v\n", product)

    order := Order{
        ID:         1,
        CustomerID: 100,
        Total:      99.99,
        Status:     "pending",
    }

    fmt.Printf("%+v\n", order)
}
```

### Validation Tags

```go
package main

import "fmt"

// Common validation library tags (go-playground/validator)
type Registration struct {
    Username  string `json:"username" validate:"required,min=3,max=20"`
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    Age       int    `json:"age" validate:"required,gte=18,lte=120"`
    Phone     string `json:"phone" validate:"required,e164"`
    Website   string `json:"website" validate:"omitempty,url"`
}

// Custom struct with multiple tags
type Config struct {
    Host     string `json:"host" yaml:"host" env:"APP_HOST" default:"localhost"`
    Port     int    `json:"port" yaml:"port" env:"APP_PORT" default:"8080"`
    Debug    bool   `json:"debug" yaml:"debug" env:"APP_DEBUG" default:"false"`
    Database string `json:"database" yaml:"database" env:"DATABASE_URL" required:"true"`
}

func main() {
    reg := Registration{
        Username: "alice",
        Email:    "alice@example.com",
        Password: "securepass123",
        Age:      25,
        Phone:    "+1234567890",
    }

    fmt.Printf("Registration: %+v\n", reg)

    config := Config{
        Host:     "0.0.0.0",
        Port:     3000,
        Debug:    true,
        Database: "postgres://localhost/mydb",
    }

    fmt.Printf("Config: %+v\n", config)
}
```

### Reading Struct Tags with Reflection

```go
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    ID    int    `json:"id" db:"user_id"`
    Name  string `json:"name" db:"user_name"`
    Email string `json:"email" db:"email_address"`
}

func main() {
    user := User{}
    t := reflect.TypeOf(user)

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        jsonTag := field.Tag.Get("json")
        dbTag := field.Tag.Get("db")

        fmt.Printf("Field: %s, JSON: %s, DB: %s\n",
            field.Name, jsonTag, dbTag)
    }
}
```

---

## 6. Constructor Functions

Go doesn't have constructors, but we use factory functions:

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

type User struct {
    ID        int
    Username  string
    Email     string
    CreatedAt time.Time
    settings  map[string]string // unexported
}

// Constructor function
func NewUser(username, email string) *User {
    return &User{
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
        settings:  make(map[string]string),
    }
}

// Constructor with validation
func NewUserWithID(id int, username, email string) (*User, error) {
    if id <= 0 {
        return nil, errors.New("id must be positive")
    }
    if username == "" {
        return nil, errors.New("username is required")
    }
    if email == "" {
        return nil, errors.New("email is required")
    }

    return &User{
        ID:        id,
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
        settings:  make(map[string]string),
    }, nil
}

// Method to set settings
func (u *User) SetSetting(key, value string) {
    u.settings[key] = value
}

func (u *User) GetSetting(key string) string {
    return u.settings[key]
}

func main() {
    // Using simple constructor
    user1 := NewUser("alice", "alice@example.com")
    fmt.Printf("User1: %+v\n", user1)

    // Using constructor with validation
    user2, err := NewUserWithID(1, "bob", "bob@example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("User2: %+v\n", user2)

    // Using settings
    user1.SetSetting("theme", "dark")
    fmt.Println("Theme:", user1.GetSetting("theme"))
}
```

---

## Exercises

### Exercise 1: Shape Calculator
Create structs for different shapes with area and perimeter methods.

```go
package main

import (
    "fmt"
    "math"
)

// TODO: Define Rectangle struct with Width and Height
// Add methods: Area(), Perimeter()

// TODO: Define Circle struct with Radius
// Add methods: Area(), Circumference()

// TODO: Define Triangle struct with Base and Height
// Add methods: Area()

func main() {
    // Test your implementations
    rect := Rectangle{Width: 10, Height: 5}
    fmt.Printf("Rectangle: Area=%.2f, Perimeter=%.2f\n",
        rect.Area(), rect.Perimeter())

    circle := Circle{Radius: 7}
    fmt.Printf("Circle: Area=%.2f, Circumference=%.2f\n",
        circle.Area(), circle.Circumference())

    triangle := Triangle{Base: 10, Height: 8}
    fmt.Printf("Triangle: Area=%.2f\n", triangle.Area())
}
```

### Exercise 2: Bank Account System
Implement a bank account with transactions.

```go
package main

import (
    "fmt"
    "time"
)

type Transaction struct {
    Type      string    // "deposit" or "withdrawal"
    Amount    float64
    Timestamp time.Time
    Balance   float64   // Balance after transaction
}

type BankAccount struct {
    AccountNumber string
    Owner         string
    Balance       float64
    Transactions  []Transaction
}

// TODO: Implement NewBankAccount constructor

// TODO: Implement Deposit method (add transaction record)

// TODO: Implement Withdraw method (check balance, add transaction record)

// TODO: Implement GetStatement method (return formatted transaction history)

func main() {
    account := NewBankAccount("123456", "Alice")

    account.Deposit(1000)
    account.Deposit(500)
    account.Withdraw(200)
    account.Withdraw(100)

    fmt.Println("Account:", account.Owner)
    fmt.Printf("Balance: $%.2f\n", account.Balance)
    fmt.Println("\nStatement:")
    fmt.Println(account.GetStatement())
}
```

### Exercise 3: Library System
Create a library system with books and members.

```go
package main

import "fmt"

type Book struct {
    ISBN      string
    Title     string
    Author    string
    Available bool
}

type Member struct {
    ID            int
    Name          string
    BorrowedBooks []string // ISBNs
}

type Library struct {
    Name    string
    Books   map[string]*Book   // ISBN -> Book
    Members map[int]*Member    // ID -> Member
}

// TODO: Implement NewLibrary constructor

// TODO: Implement AddBook method

// TODO: Implement RegisterMember method

// TODO: Implement BorrowBook(memberID int, isbn string) error

// TODO: Implement ReturnBook(memberID int, isbn string) error

// TODO: Implement ListAvailableBooks() []*Book

func main() {
    library := NewLibrary("City Library")

    // Add books
    library.AddBook(&Book{ISBN: "978-0134190440", Title: "The Go Programming Language", Author: "Donovan & Kernighan", Available: true})
    library.AddBook(&Book{ISBN: "978-1491941959", Title: "Introducing Go", Author: "Caleb Doxsey", Available: true})

    // Register members
    library.RegisterMember(&Member{ID: 1, Name: "Alice"})
    library.RegisterMember(&Member{ID: 2, Name: "Bob"})

    // Borrow books
    err := library.BorrowBook(1, "978-0134190440")
    if err != nil {
        fmt.Println("Error:", err)
    }

    // List available books
    available := library.ListAvailableBooks()
    fmt.Println("Available books:")
    for _, book := range available {
        fmt.Printf("- %s by %s\n", book.Title, book.Author)
    }
}
```

### Exercise 4: JSON API Response
Create structs for API responses with proper JSON tags.

```go
package main

import (
    "encoding/json"
    "fmt"
)

// TODO: Create User struct with JSON tags
// Fields: ID, Username, Email, Password (hidden), CreatedAt (omitempty)

// TODO: Create APIResponse struct
// Fields: Success, Data (omitempty), Error (omitempty), Timestamp

// TODO: Create PaginatedResponse struct
// Fields: Data, Page, PerPage, Total, TotalPages

func main() {
    // Test User serialization
    user := User{
        ID:       1,
        Username: "alice",
        Email:    "alice@example.com",
        Password: "secret",
    }

    userJSON, _ := json.MarshalIndent(user, "", "  ")
    fmt.Println("User JSON:")
    fmt.Println(string(userJSON))

    // Test API response
    response := APIResponse{
        Success: true,
        Data:    user,
    }

    respJSON, _ := json.MarshalIndent(response, "", "  ")
    fmt.Println("\nAPI Response:")
    fmt.Println(string(respJSON))
}
```

### Exercise 5: Embedded Structs
Create a vehicle hierarchy using embedding.

```go
package main

import "fmt"

// TODO: Create Vehicle struct with Make, Model, Year

// TODO: Create Car struct embedding Vehicle
// Add fields: Doors, FuelType
// Add method: Drive()

// TODO: Create Motorcycle struct embedding Vehicle
// Add field: EngineCC
// Add method: Ride()

// TODO: Create Truck struct embedding Vehicle
// Add field: PayloadCapacity
// Add method: Haul()

func main() {
    car := Car{
        Vehicle: Vehicle{Make: "Toyota", Model: "Camry", Year: 2022},
        Doors:   4,
        FuelType: "Gasoline",
    }
    fmt.Printf("Car: %s %s (%d)\n", car.Make, car.Model, car.Year)
    car.Drive()

    motorcycle := Motorcycle{
        Vehicle:  Vehicle{Make: "Honda", Model: "CBR600RR", Year: 2021},
        EngineCC: 599,
    }
    fmt.Printf("\nMotorcycle: %s %s (%d)\n", motorcycle.Make, motorcycle.Model, motorcycle.Year)
    motorcycle.Ride()

    truck := Truck{
        Vehicle:         Vehicle{Make: "Ford", Model: "F-150", Year: 2023},
        PayloadCapacity: 3300,
    }
    fmt.Printf("\nTruck: %s %s (%d)\n", truck.Make, truck.Model, truck.Year)
    truck.Haul()
}
```

---

## Summary

In this tutorial, you learned:
- How to define structs to group related data
- Different ways to create struct instances
- Methods with value and pointer receivers
- When to use pointer vs value receivers
- Struct embedding for composition
- Struct tags for JSON, database, and validation

---

**Next:** [03-interfaces.md](03-interfaces.md) - Learn about interfaces and polymorphism in Go

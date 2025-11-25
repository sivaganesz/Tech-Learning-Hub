# Packages and Modules in Go

## Learning Objectives

By the end of this tutorial, you will be able to:
- Understand package declaration and naming
- Import standard library and third-party packages
- Distinguish between exported and unexported identifiers
- Work with go.mod and go.sum files
- Use go mod commands effectively
- Follow project structure conventions

---

## 1. Package Declaration

Every Go file starts with a package declaration:

```go
// main package - executable programs
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Package Naming Conventions

```go
// Package names should be:
// - lowercase
// - single word (avoid underscores or camelCase)
// - short but descriptive

// Good package names:
package http
package json
package user
package auth
package database

// Avoid:
package httpHandler    // Use httphandler or separate packages
package user_service   // Use userservice or separate packages
package MyPackage      // Use mypackage
```

### Creating a Package

```go
// File: calculator/calculator.go
package calculator

// Add returns the sum of two integers
func Add(a, b int) int {
    return a + b
}

// Subtract returns the difference of two integers
func Subtract(a, b int) int {
    return a - b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
    return a * b
}

// Divide returns the quotient of two integers
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}
```

### Using Your Package

```go
// File: main.go
package main

import (
    "fmt"
    "myproject/calculator"
)

func main() {
    sum := calculator.Add(10, 5)
    fmt.Println("10 + 5 =", sum)

    diff := calculator.Subtract(10, 5)
    fmt.Println("10 - 5 =", diff)
}
```

---

## 2. Imports

### Standard Library Imports

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

func main() {
    // fmt - formatted I/O
    fmt.Println("Hello, World!")

    // strings - string manipulation
    upper := strings.ToUpper("hello")
    fmt.Println(upper)

    // strconv - string conversion
    num, _ := strconv.Atoi("42")
    fmt.Println(num)

    // time - time and duration
    now := time.Now()
    fmt.Println(now)

    // os - operating system interface
    hostname, _ := os.Hostname()
    fmt.Println("Hostname:", hostname)

    // path/filepath - file path manipulation
    dir := filepath.Dir("/path/to/file.txt")
    fmt.Println("Directory:", dir)

    // encoding/json - JSON encoding/decoding
    data, _ := json.Marshal(map[string]int{"a": 1})
    fmt.Println(string(data))

    // net/http - HTTP client and server
    // io - I/O primitives
    _, _ = http.Get, io.EOF // Just to use the imports
}
```

### Third-Party Imports

```go
package main

import (
    "fmt"

    // Third-party packages (need go get or go mod tidy)
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go.uber.org/zap"
)

func main() {
    // Generate UUID
    id := uuid.New()
    fmt.Println("UUID:", id)

    // Create Gin router
    router := gin.Default()

    // Create Zap logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    logger.Info("Application started")

    _ = router // Use router
}
```

### Import Aliases and Dot Imports

```go
package main

import (
    "fmt"

    // Alias import
    myjson "encoding/json"

    // Blank import (for side effects only)
    _ "image/png"

    // Multiple packages with same name
    "crypto/rand"
    mathrand "math/rand"

    // Dot import (not recommended - pollutes namespace)
    // . "strings"
)

func main() {
    // Using alias
    data, _ := myjson.Marshal(map[string]int{"a": 1})
    fmt.Println(string(data))

    // Using both rand packages
    cryptoBytes := make([]byte, 16)
    rand.Read(cryptoBytes)

    mathNum := mathrand.Intn(100)
    fmt.Println("Random number:", mathNum)
}
```

### Import Organization

```go
package main

import (
    // Standard library (alphabetical)
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    // Third-party packages (alphabetical)
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go.uber.org/zap"

    // Internal packages (alphabetical)
    "myproject/internal/auth"
    "myproject/internal/database"
    "myproject/internal/models"
)

// goimports tool will format these automatically
```

---

## 3. Exported vs Unexported

### Visibility Rules

```go
// File: user/user.go
package user

import "time"

// Exported (capital letter) - accessible from other packages
type User struct {
    ID        int       // Exported field
    Username  string    // Exported field
    Email     string    // Exported field
    password  string    // Unexported field (lowercase)
    createdAt time.Time // Unexported field
}

// Exported constant
const MaxUsers = 1000

// Unexported constant
const defaultTimeout = 30

// Exported variable
var AdminEmail = "admin@example.com"

// Unexported variable
var connectionPool = make(map[string]interface{})

// Exported function
func NewUser(username, email, password string) *User {
    return &User{
        Username:  username,
        Email:     email,
        password:  hashPassword(password),
        createdAt: time.Now(),
    }
}

// Unexported function (helper)
func hashPassword(password string) string {
    // Implementation
    return "hashed_" + password
}

// Exported method
func (u *User) GetUsername() string {
    return u.Username
}

// Unexported method
func (u *User) validateEmail() bool {
    // Implementation
    return true
}
```

### Using Exported Members

```go
// File: main.go
package main

import (
    "fmt"
    "myproject/user"
)

func main() {
    // Can access exported members
    u := user.NewUser("alice", "alice@example.com", "secret123")
    fmt.Println("Username:", u.GetUsername())
    fmt.Println("Email:", u.Email)
    fmt.Println("Max users:", user.MaxUsers)

    // Cannot access unexported members
    // fmt.Println(u.password)      // Error: unexported field
    // fmt.Println(u.createdAt)     // Error: unexported field
    // user.hashPassword("test")    // Error: unexported function
}
```

### Package-Level Access

```go
// File: user/admin.go
package user

// Same package can access unexported members

func CreateAdmin() *User {
    admin := &User{
        Username: "admin",
        Email:    AdminEmail,
        password: hashPassword("admin123"), // Can access unexported
    }
    admin.validateEmail() // Can call unexported method
    return admin
}

func GetUserPassword(u *User) string {
    // Can access unexported field within same package
    return u.password
}
```

---

## 4. go.mod and go.sum

### Creating a Module

```bash
# Create a new module
mkdir myproject
cd myproject
go mod init github.com/username/myproject
```

### go.mod File

```go
// go.mod
module github.com/username/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/google/uuid v1.4.0
    go.uber.org/zap v1.26.0
)

require (
    // Indirect dependencies (transitive)
    github.com/bytedance/sonic v1.9.1 // indirect
    github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
    github.com/gabriel-vasile/mimetype v1.4.2 // indirect
    // ... more indirect dependencies
)
```

### go.sum File

```
// go.sum - checksums for dependencies
github.com/gin-gonic/gin v1.9.1 h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=
github.com/gin-gonic/gin v1.9.1/go.mod h1:hPrL7YrpYKXt5YId3A/Dn+qS8P1Uk3M3VpY4/yKHPWo=
github.com/google/uuid v1.4.0 h1:MtMxsa51/r9yyhkyLsVeVt0B+BGQZzpQiTQ4eHZ8bc4=
github.com/google/uuid v1.4.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
// ... checksums for all dependencies
```

### Understanding go.mod Directives

```go
// go.mod with all directives
module github.com/username/myproject

go 1.21

// Direct dependencies
require (
    github.com/gin-gonic/gin v1.9.1
    github.com/google/uuid v1.4.0
)

// Replace a module with local version or fork
replace github.com/original/package => ../local/package
replace github.com/original/package => github.com/fork/package v1.0.0

// Exclude a specific version
exclude github.com/bad/package v1.2.3

// Retract versions (module author declares versions as bad)
retract (
    v1.0.0 // Contains critical bug
    [v1.1.0, v1.2.0] // Range retraction
)
```

---

## 5. go mod Commands

### Essential Commands

```bash
# Initialize a new module
go mod init github.com/username/project

# Add missing dependencies, remove unused ones
go mod tidy

# Download dependencies to local cache
go mod download

# Verify dependencies haven't been modified
go mod verify

# Print module dependency graph
go mod graph

# Show why a package is needed
go mod why github.com/some/package

# Create a vendor directory
go mod vendor

# Edit go.mod programmatically
go mod edit -require github.com/pkg/errors@v0.9.1
go mod edit -droprequire github.com/old/package
go mod edit -replace github.com/old=github.com/new@v1.0.0
```

### Common Workflows

```bash
# Adding a new dependency
go get github.com/gin-gonic/gin
go get github.com/gin-gonic/gin@v1.9.1  # Specific version
go get github.com/gin-gonic/gin@latest  # Latest version

# Updating dependencies
go get -u github.com/gin-gonic/gin      # Update to latest
go get -u ./...                          # Update all dependencies

# Downgrading
go get github.com/gin-gonic/gin@v1.8.0

# Removing a dependency
# 1. Remove import from code
# 2. Run:
go mod tidy

# Checking for updates
go list -m -u all                        # List all with updates
go list -m -u github.com/gin-gonic/gin   # Check specific package
```

### Working with Private Repositories

```bash
# Set GOPRIVATE for private repos
export GOPRIVATE=github.com/mycompany/*

# Or in .gitconfig
git config --global url."git@github.com:".insteadOf "https://github.com/"

# Configure authentication
# Option 1: SSH keys (recommended)
# Option 2: Personal access token in .netrc
machine github.com
    login USERNAME
    password TOKEN
```

---

## 6. Project Structure Conventions

### Simple Project

```
myproject/
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Package-Oriented Layout

```
myproject/
├── go.mod
├── go.sum
├── main.go
├── auth/
│   ├── auth.go
│   ├── auth_test.go
│   └── token.go
├── database/
│   ├── database.go
│   ├── migrations/
│   │   └── 001_initial.sql
│   └── queries.go
├── models/
│   ├── user.go
│   └── product.go
└── handlers/
    ├── user_handler.go
    └── product_handler.go
```

### Standard Go Project Layout

```
myproject/
├── cmd/
│   ├── api/
│   │   └── main.go           # API server entry point
│   └── worker/
│       └── main.go           # Background worker entry point
├── internal/
│   ├── auth/                 # Internal auth package
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── database/             # Database access
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # HTTP middleware
│   ├── models/               # Data models
│   ├── repositories/         # Data access layer
│   └── services/             # Business logic
├── pkg/
│   ├── logger/               # Reusable logger package
│   └── validator/            # Reusable validation
├── api/
│   └── openapi.yaml          # API specification
├── configs/
│   ├── config.yaml
│   └── config.go
├── scripts/
│   ├── build.sh
│   └── deploy.sh
├── migrations/
│   └── 001_initial.sql
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### internal vs pkg

```go
// internal/ - Private to your module
// Only code within your module can import these
myproject/
└── internal/
    └── auth/
        └── auth.go  // Only myproject can import

// pkg/ - Public packages
// Other projects can import these
myproject/
└── pkg/
    └── logger/
        └── logger.go  // Anyone can import
```

### Example: Complete Project

```go
// cmd/api/main.go
package main

import (
    "log"
    "myproject/internal/config"
    "myproject/internal/database"
    "myproject/internal/handlers"
    "myproject/internal/services"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Initialize database
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Initialize services
    userService := services.NewUserService(db)
    authService := services.NewAuthService(cfg.JWTSecret)

    // Initialize handlers
    handler := handlers.NewHandler(userService, authService)

    // Start server
    log.Println("Server starting on", cfg.Port)
    handler.Start(cfg.Port)
}

// internal/config/config.go
package config

import "os"

type Config struct {
    Port        string
    DatabaseURL string
    JWTSecret   string
}

func Load() (*Config, error) {
    return &Config{
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/app"),
        JWTSecret:   getEnv("JWT_SECRET", "secret"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// internal/models/user.go
package models

import "time"

type User struct {
    ID        int       `json:"id" db:"id"`
    Username  string    `json:"username" db:"username"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"-" db:"password"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// internal/services/user_service.go
package services

import "myproject/internal/models"

type UserService struct {
    db Database
}

func NewUserService(db Database) *UserService {
    return &UserService{db: db}
}

func (s *UserService) GetUser(id int) (*models.User, error) {
    // Implementation
    return nil, nil
}

func (s *UserService) CreateUser(user *models.User) error {
    // Implementation
    return nil
}
```

---

## Exercises

### Exercise 1: Create a Calculator Package
Build a calculator package with proper exports.

```
calculator/
├── go.mod
├── main.go
├── calc/
│   ├── basic.go      # Add, Subtract, Multiply, Divide
│   ├── advanced.go   # Power, Sqrt, Factorial
│   └── helpers.go    # unexported helper functions
└── calc_test.go
```

Implement:
- Exported functions for all operations
- Unexported helper functions
- Proper error handling for invalid inputs

### Exercise 2: Multi-Package Project
Create a user management system:

```
usersystem/
├── go.mod
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── models/
│   │   └── user.go
│   ├── repository/
│   │   └── user_repo.go
│   └── service/
│       └── user_service.go
└── pkg/
    └── validator/
        └── validator.go
```

### Exercise 3: go mod Practice
Practice module commands:

```bash
# 1. Create a new module
mkdir modpractice && cd modpractice
go mod init github.com/yourname/modpractice

# 2. Add dependencies
# Create main.go that uses gin and uuid

# 3. Run go mod tidy

# 4. Check the dependency graph
go mod graph

# 5. Update a dependency
go get -u github.com/gin-gonic/gin

# 6. Vendor dependencies
go mod vendor
```

### Exercise 4: Import Refactoring
Refactor this code with proper import organization:

```go
package main

import "github.com/gin-gonic/gin"
import "fmt"
import "myproject/internal/auth"
import "github.com/google/uuid"
import "net/http"
import "myproject/internal/models"
import "encoding/json"
import "go.uber.org/zap"
import "time"
import "myproject/internal/handlers"

// Refactor to use proper import groups and formatting
```

### Exercise 5: Package Design
Design packages for an e-commerce system:

```
ecommerce/
├── cmd/
│   └── server/
├── internal/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── handlers/
└── pkg/
```

Create the following packages with appropriate exports:
- models: Product, Order, Customer
- repository: ProductRepo, OrderRepo, CustomerRepo
- service: ProductService, OrderService, PaymentService
- handlers: HTTP handlers for all services

---

## Summary

In this tutorial, you learned:
- Package declaration and naming conventions
- How to import standard and third-party packages
- Exported vs unexported identifiers
- Working with go.mod and go.sum
- Essential go mod commands
- Project structure best practices

---

**Next:** [06-concurrency.md](06-concurrency.md) - Learn about goroutines and channels in Go

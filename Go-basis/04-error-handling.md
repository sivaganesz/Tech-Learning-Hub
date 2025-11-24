# Error Handling in Go

## Learning Objectives

By the end of this tutorial, you will be able to:
- Understand the error type and how it works
- Create errors with `errors.New` and `fmt.Errorf`
- Check and handle errors properly
- Wrap errors with context using `%w`
- Use `errors.Is` and `errors.As` for error inspection
- Apply error handling best practices

---

## 1. Error Type Basics

In Go, errors are values that implement the `error` interface:

```go
package main

import (
    "errors"
    "fmt"
)

// The error interface
// type error interface {
//     Error() string
// }

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    // Successful operation
    result, err := divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("10 / 2 =", result)

    // Failed operation
    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

### The Error Pattern

```go
package main

import (
    "errors"
    "fmt"
    "strconv"
)

func parseAge(s string) (int, error) {
    age, err := strconv.Atoi(s)
    if err != nil {
        return 0, err
    }
    if age < 0 {
        return 0, errors.New("age cannot be negative")
    }
    if age > 150 {
        return 0, errors.New("age seems unrealistic")
    }
    return age, nil
}

func main() {
    ages := []string{"25", "-5", "abc", "200", "30"}

    for _, ageStr := range ages {
        age, err := parseAge(ageStr)
        if err != nil {
            fmt.Printf("Failed to parse '%s': %v\n", ageStr, err)
            continue
        }
        fmt.Printf("Parsed age: %d\n", age)
    }
}
```

---

## 2. Creating Errors

### Using errors.New

```go
package main

import (
    "errors"
    "fmt"
)

// Package-level error variables (sentinel errors)
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrInvalidInput = errors.New("invalid input")
)

type User struct {
    ID   int
    Name string
}

var users = map[int]User{
    1: {ID: 1, Name: "Alice"},
    2: {ID: 2, Name: "Bob"},
}

func getUser(id int) (User, error) {
    if id <= 0 {
        return User{}, ErrInvalidInput
    }
    user, ok := users[id]
    if !ok {
        return User{}, ErrNotFound
    }
    return user, nil
}

func main() {
    // Test different scenarios
    ids := []int{1, 2, 3, -1}

    for _, id := range ids {
        user, err := getUser(id)
        if err != nil {
            fmt.Printf("Error getting user %d: %v\n", id, err)
            continue
        }
        fmt.Printf("Found user: %+v\n", user)
    }
}
```

### Using fmt.Errorf

```go
package main

import (
    "fmt"
)

func validateUsername(username string) error {
    if len(username) == 0 {
        return fmt.Errorf("username cannot be empty")
    }
    if len(username) < 3 {
        return fmt.Errorf("username '%s' is too short (minimum 3 characters)", username)
    }
    if len(username) > 20 {
        return fmt.Errorf("username '%s' is too long (maximum 20 characters)", username)
    }
    return nil
}

func validateEmail(email string) error {
    if len(email) == 0 {
        return fmt.Errorf("email cannot be empty")
    }
    // Simple check - real validation would be more complex
    for _, char := range email {
        if char == '@' {
            return nil
        }
    }
    return fmt.Errorf("email '%s' is invalid: missing @ symbol", email)
}

func main() {
    usernames := []string{"", "ab", "alice", "verylongusernamethatexceeds20chars"}
    emails := []string{"", "invalid", "user@example.com"}

    fmt.Println("Username validation:")
    for _, u := range usernames {
        if err := validateUsername(u); err != nil {
            fmt.Printf("  Invalid: %v\n", err)
        } else {
            fmt.Printf("  Valid: %s\n", u)
        }
    }

    fmt.Println("\nEmail validation:")
    for _, e := range emails {
        if err := validateEmail(e); err != nil {
            fmt.Printf("  Invalid: %v\n", err)
        } else {
            fmt.Printf("  Valid: %s\n", e)
        }
    }
}
```

### Custom Error Types

```go
package main

import (
    "fmt"
    "time"
)

// Custom error type
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// Custom error with more context
type APIError struct {
    StatusCode int
    Message    string
    Timestamp  time.Time
    RequestID  string
}

func (e APIError) Error() string {
    return fmt.Sprintf("[%d] %s (request: %s)", e.StatusCode, e.Message, e.RequestID)
}

// Temporary error interface (for retry logic)
type TemporaryError struct {
    Message string
}

func (e TemporaryError) Error() string {
    return e.Message
}

func (e TemporaryError) Temporary() bool {
    return true
}

func validateUser(name, email string) error {
    if name == "" {
        return ValidationError{Field: "name", Message: "cannot be empty"}
    }
    if email == "" {
        return ValidationError{Field: "email", Message: "cannot be empty"}
    }
    return nil
}

func callAPI() error {
    // Simulate API error
    return APIError{
        StatusCode: 429,
        Message:    "rate limit exceeded",
        Timestamp:  time.Now(),
        RequestID:  "req-123",
    }
}

func main() {
    // Validation error
    if err := validateUser("", ""); err != nil {
        fmt.Println("Validation failed:", err)
    }

    // API error
    if err := callAPI(); err != nil {
        fmt.Println("API call failed:", err)
    }
}
```

---

## 3. Checking Errors

### Basic Error Checking

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

func main() {
    // Pattern 1: Check immediately
    file, err := os.Open("nonexistent.txt")
    if err != nil {
        fmt.Println("Failed to open file:", err)
        // Handle error and return/continue
    } else {
        defer file.Close()
        // Use file
    }

    // Pattern 2: Inline error check
    if _, err := os.Stat("config.json"); err != nil {
        if os.IsNotExist(err) {
            fmt.Println("Config file does not exist")
        } else {
            fmt.Println("Error checking config:", err)
        }
    }

    // Pattern 3: Multiple operations
    data, err := readConfig()
    if err != nil {
        fmt.Println("Config error:", err)
        return
    }

    processed, err := processData(data)
    if err != nil {
        fmt.Println("Processing error:", err)
        return
    }

    err = saveResult(processed)
    if err != nil {
        fmt.Println("Save error:", err)
        return
    }

    fmt.Println("All operations completed successfully")
}

func readConfig() (string, error) {
    return "config data", nil
}

func processData(data string) (string, error) {
    return "processed " + data, nil
}

func saveResult(data string) error {
    return nil
}
```

### Comparing Errors

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
)

func getResource(id int, isAuthenticated, isAdmin bool) (string, error) {
    if !isAuthenticated {
        return "", ErrUnauthorized
    }
    if id == 999 && !isAdmin {
        return "", ErrForbidden
    }
    if id > 100 {
        return "", ErrNotFound
    }
    return fmt.Sprintf("Resource %d", id), nil
}

func main() {
    testCases := []struct {
        id              int
        isAuthenticated bool
        isAdmin         bool
    }{
        {1, true, false},
        {1, false, false},
        {999, true, false},
        {999, true, true},
        {200, true, false},
    }

    for _, tc := range testCases {
        resource, err := getResource(tc.id, tc.isAuthenticated, tc.isAdmin)

        if err == nil {
            fmt.Printf("Got resource: %s\n", resource)
            continue
        }

        // Compare with sentinel errors
        switch err {
        case ErrNotFound:
            fmt.Printf("Resource %d not found\n", tc.id)
        case ErrUnauthorized:
            fmt.Println("Please log in first")
        case ErrForbidden:
            fmt.Println("Admin access required")
        default:
            fmt.Println("Unknown error:", err)
        }
    }
}
```

---

## 4. Wrapping Errors

### Using %w for Error Wrapping

```go
package main

import (
    "errors"
    "fmt"
)

var ErrDatabase = errors.New("database error")

func queryDatabase() error {
    // Simulate a low-level database error
    return ErrDatabase
}

func getUser(id int) error {
    err := queryDatabase()
    if err != nil {
        // Wrap with context using %w
        return fmt.Errorf("failed to get user %d: %w", id, err)
    }
    return nil
}

func handleRequest() error {
    err := getUser(123)
    if err != nil {
        // Add more context
        return fmt.Errorf("request handler error: %w", err)
    }
    return nil
}

func main() {
    err := handleRequest()
    if err != nil {
        fmt.Println("Error:", err)
        // Output: Error: request handler error: failed to get user 123: database error
    }
}
```

### Building Error Context

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrConnection = errors.New("connection failed")
    ErrTimeout    = errors.New("operation timed out")
    ErrInvalidData = errors.New("invalid data")
)

type Repository struct {
    name string
}

func (r *Repository) connect() error {
    return ErrConnection
}

func (r *Repository) fetchData(query string) ([]byte, error) {
    if err := r.connect(); err != nil {
        return nil, fmt.Errorf("repository %s: %w", r.name, err)
    }
    return nil, nil
}

type Service struct {
    repo *Repository
}

func (s *Service) GetUserData(userID int) ([]byte, error) {
    data, err := s.repo.fetchData(fmt.Sprintf("SELECT * FROM users WHERE id = %d", userID))
    if err != nil {
        return nil, fmt.Errorf("GetUserData(%d): %w", userID, err)
    }
    return data, nil
}

type Handler struct {
    service *Service
}

func (h *Handler) HandleGetUser(userID int) error {
    _, err := h.service.GetUserData(userID)
    if err != nil {
        return fmt.Errorf("HandleGetUser: %w", err)
    }
    return nil
}

func main() {
    repo := &Repository{name: "users_db"}
    service := &Service{repo: repo}
    handler := &Handler{service: service}

    err := handler.HandleGetUser(123)
    if err != nil {
        fmt.Println("Full error chain:")
        fmt.Println(err)
        // Output: HandleGetUser: GetUserData(123): repository users_db: connection failed

        // Unwrap to check original error
        fmt.Println("\nIs connection error:", errors.Is(err, ErrConnection))
    }
}
```

---

## 5. errors.Is and errors.As

### Using errors.Is

```go
package main

import (
    "errors"
    "fmt"
    "os"
)

var (
    ErrNotFound = errors.New("not found")
    ErrTimeout  = errors.New("timeout")
)

func findUser(id int) error {
    // Simulate wrapped error
    return fmt.Errorf("database query failed: %w", ErrNotFound)
}

func fetchData() error {
    return fmt.Errorf("HTTP request failed: %w", ErrTimeout)
}

func main() {
    // errors.Is checks if any error in the chain matches
    err := findUser(123)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("User not found - showing 404 page")
    }

    err = fetchData()
    if errors.Is(err, ErrTimeout) {
        fmt.Println("Request timed out - will retry")
    }

    // Works with standard library errors too
    _, err = os.Open("nonexistent.txt")
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("File does not exist")
    }

    // Compare the difference
    fmt.Println("\nDirect comparison vs errors.Is:")
    wrappedErr := fmt.Errorf("wrapped: %w", ErrNotFound)

    fmt.Println("Direct ==:", wrappedErr == ErrNotFound)        // false
    fmt.Println("errors.Is:", errors.Is(wrappedErr, ErrNotFound)) // true
}
```

### Using errors.As

```go
package main

import (
    "errors"
    "fmt"
)

// Custom error type
type ValidationError struct {
    Field string
    Value interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid value %v for field %s", e.Value, e.Field)
}

// Another custom error type
type HTTPError struct {
    StatusCode int
    Message    string
}

func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func validateAge(age int) error {
    if age < 0 || age > 150 {
        return &ValidationError{Field: "age", Value: age}
    }
    return nil
}

func processRequest() error {
    err := validateAge(-5)
    if err != nil {
        return fmt.Errorf("request processing failed: %w", err)
    }
    return nil
}

func fetchFromAPI() error {
    return fmt.Errorf("API call failed: %w", &HTTPError{
        StatusCode: 429,
        Message:    "rate limit exceeded",
    })
}

func main() {
    // Extract ValidationError from wrapped error
    err := processRequest()
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation failed on field '%s' with value '%v'\n",
            validationErr.Field, validationErr.Value)
    }

    // Extract HTTPError
    err = fetchFromAPI()
    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        fmt.Printf("HTTP Error: status=%d, message=%s\n",
            httpErr.StatusCode, httpErr.Message)

        if httpErr.StatusCode == 429 {
            fmt.Println("Will retry after delay...")
        }
    }
}
```

---

## 6. Best Practices

### Handle Errors Once

```go
package main

import (
    "errors"
    "fmt"
    "log"
)

var ErrInvalidInput = errors.New("invalid input")

// BAD: Logging and returning
func badExample(input string) error {
    if input == "" {
        err := ErrInvalidInput
        log.Println("Error:", err) // Logged here
        return err                  // And returned - will be logged again!
    }
    return nil
}

// GOOD: Return error, let caller decide
func goodExample(input string) error {
    if input == "" {
        return fmt.Errorf("processing input: %w", ErrInvalidInput)
    }
    return nil
}

func main() {
    // Caller handles the error appropriately
    if err := goodExample(""); err != nil {
        // Log once at the appropriate level
        log.Printf("Operation failed: %v", err)
    }
}
```

### Add Context When Wrapping

```go
package main

import (
    "errors"
    "fmt"
)

var ErrDatabase = errors.New("database error")

// BAD: No context added
func badWrap() error {
    err := queryDB()
    if err != nil {
        return err // No context - where did this come from?
    }
    return nil
}

// GOOD: Add meaningful context
func goodWrap(userID int, action string) error {
    err := queryDB()
    if err != nil {
        return fmt.Errorf("%s for user %d: %w", action, userID, err)
    }
    return nil
}

func queryDB() error {
    return ErrDatabase
}

func main() {
    err := goodWrap(123, "fetching profile")
    if err != nil {
        fmt.Println(err)
        // Output: fetching profile for user 123: database error
    }
}
```

### Use Sentinel Errors for Expected Conditions

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors for expected conditions
var (
    ErrNotFound       = errors.New("resource not found")
    ErrAlreadyExists  = errors.New("resource already exists")
    ErrInvalidInput   = errors.New("invalid input")
    ErrUnauthorized   = errors.New("unauthorized")
    ErrRateLimited    = errors.New("rate limited")
)

type UserService struct {
    users map[int]string
}

func (s *UserService) GetUser(id int) (string, error) {
    if id <= 0 {
        return "", ErrInvalidInput
    }
    user, ok := s.users[id]
    if !ok {
        return "", ErrNotFound
    }
    return user, nil
}

func (s *UserService) CreateUser(id int, name string) error {
    if id <= 0 || name == "" {
        return ErrInvalidInput
    }
    if _, ok := s.users[id]; ok {
        return ErrAlreadyExists
    }
    s.users[id] = name
    return nil
}

func main() {
    service := &UserService{users: make(map[int]string)}

    // Create user
    err := service.CreateUser(1, "Alice")
    if err != nil {
        handleError(err)
        return
    }

    // Try to create duplicate
    err = service.CreateUser(1, "Alice2")
    if err != nil {
        handleError(err)
    }

    // Get existing user
    user, err := service.GetUser(1)
    if err != nil {
        handleError(err)
        return
    }
    fmt.Println("Found user:", user)

    // Get non-existing user
    _, err = service.GetUser(999)
    if err != nil {
        handleError(err)
    }
}

func handleError(err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        fmt.Println("404 - Not Found")
    case errors.Is(err, ErrAlreadyExists):
        fmt.Println("409 - Conflict")
    case errors.Is(err, ErrInvalidInput):
        fmt.Println("400 - Bad Request")
    case errors.Is(err, ErrUnauthorized):
        fmt.Println("401 - Unauthorized")
    default:
        fmt.Println("500 - Internal Server Error:", err)
    }
}
```

### Don't Panic

```go
package main

import (
    "errors"
    "fmt"
)

// BAD: Panicking on recoverable errors
func badDivide(a, b int) int {
    if b == 0 {
        panic("division by zero") // Don't do this!
    }
    return a / b
}

// GOOD: Return an error
func goodDivide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    // Good approach
    result, err := goodDivide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        // Handle gracefully
    } else {
        fmt.Println("Result:", result)
    }

    // Panic is appropriate for:
    // - Programming errors (bugs that should never happen)
    // - Initialization failures that prevent the program from running
}
```

---

## Exercises

### Exercise 1: File Reader with Error Handling
Implement a file reader with proper error handling.

```go
package main

import (
    "fmt"
    "os"
)

// TODO: Define custom errors
// var ErrFileNotFound = ...
// var ErrFileEmpty = ...
// var ErrReadFailed = ...

// TODO: Implement ReadFileContents that:
// - Returns ErrFileNotFound if file doesn't exist
// - Returns ErrFileEmpty if file is empty
// - Wraps any other errors with context
func ReadFileContents(filename string) (string, error) {
    return "", nil
}

func main() {
    files := []string{"existing.txt", "nonexistent.txt", "empty.txt"}

    for _, file := range files {
        content, err := ReadFileContents(file)
        if err != nil {
            // Handle different error types appropriately
            fmt.Printf("Error reading %s: %v\n", file, err)
            continue
        }
        fmt.Printf("Content of %s: %s\n", file, content)
    }
}
```

### Exercise 2: User Registration with Validation
Create a registration system with detailed error handling.

```go
package main

import "fmt"

// TODO: Define ValidationError struct with Field and Message

// TODO: Define RegistrationError struct that can hold multiple ValidationErrors

// TODO: Implement ValidateRegistration that checks:
// - Username: 3-20 characters, alphanumeric only
// - Email: contains @
// - Password: minimum 8 characters, at least one number
// - Age: between 13 and 120

type Registration struct {
    Username string
    Email    string
    Password string
    Age      int
}

func main() {
    registrations := []Registration{
        {"alice", "alice@example.com", "password123", 25},
        {"ab", "invalid-email", "short", 10},
        {"valid_user", "user@test.com", "securepass1", 30},
    }

    for _, reg := range registrations {
        if err := ValidateRegistration(reg); err != nil {
            fmt.Printf("Registration failed: %v\n\n", err)
        } else {
            fmt.Printf("Registration valid for %s\n\n", reg.Username)
        }
    }
}
```

### Exercise 3: Error Wrapping Chain
Implement a service with proper error wrapping through layers.

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors
var (
    ErrNotFound    = errors.New("not found")
    ErrDatabase    = errors.New("database error")
    ErrPermission  = errors.New("permission denied")
)

// TODO: Implement Repository layer
type Repository struct{}

func (r *Repository) FindByID(id int) (string, error) {
    // Return ErrNotFound for id > 100
    // Return ErrDatabase for id == 0
    // Return data otherwise
    return "", nil
}

// TODO: Implement Service layer that wraps Repository errors
type Service struct {
    repo *Repository
}

func (s *Service) GetUser(id int) (string, error) {
    // Call repository and wrap errors with context
    return "", nil
}

// TODO: Implement Handler layer that wraps Service errors
type Handler struct {
    service *Service
}

func (h *Handler) HandleGetUser(id int) error {
    // Call service and wrap errors with context
    return nil
}

func main() {
    handler := &Handler{
        service: &Service{
            repo: &Repository{},
        },
    }

    testIDs := []int{1, 50, 0, 150}

    for _, id := range testIDs {
        err := handler.HandleGetUser(id)
        if err != nil {
            fmt.Printf("Error for ID %d: %v\n", id, err)

            // Check the underlying error
            if errors.Is(err, ErrNotFound) {
                fmt.Println("  -> Return 404")
            } else if errors.Is(err, ErrDatabase) {
                fmt.Println("  -> Return 500, alert ops team")
            }
        } else {
            fmt.Printf("Success for ID %d\n", id)
        }
        fmt.Println()
    }
}
```

### Exercise 4: Retry with Error Types
Implement retry logic based on error types.

```go
package main

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

// Custom error types
type TemporaryError struct {
    Message string
}

func (e TemporaryError) Error() string { return e.Message }
func (e TemporaryError) Temporary() bool { return true }

type PermanentError struct {
    Message string
}

func (e PermanentError) Error() string { return e.Message }

// TODO: Implement operation that randomly fails
func unreliableOperation() error {
    // Return TemporaryError 60% of time
    // Return PermanentError 20% of time
    // Return nil 20% of time
    return nil
}

// TODO: Implement retry logic
// - Retry up to maxRetries times for temporary errors
// - Return immediately for permanent errors
// - Return nil on success
func withRetry(maxRetries int, operation func() error) error {
    return nil
}

func main() {
    rand.Seed(time.Now().UnixNano())

    for i := 1; i <= 5; i++ {
        fmt.Printf("Attempt %d:\n", i)
        err := withRetry(3, unreliableOperation)
        if err != nil {
            fmt.Printf("  Failed: %v\n", err)
        } else {
            fmt.Println("  Success!")
        }
        fmt.Println()
    }
}
```

### Exercise 5: API Error Handler
Create a comprehensive API error handler.

```go
package main

import (
    "errors"
    "fmt"
)

// HTTP-style errors
type HTTPError struct {
    Code    int
    Message string
    Details map[string]string
}

func (e HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// Application errors
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrBadRequest   = errors.New("bad request")
    ErrInternal     = errors.New("internal error")
)

// TODO: Implement toHTTPError that converts application errors to HTTPError
func toHTTPError(err error) HTTPError {
    // Map application errors to HTTP status codes
    // Add appropriate messages
    return HTTPError{}
}

// TODO: Implement error handler middleware
func errorHandler(err error) {
    httpErr := toHTTPError(err)
    // Print the response that would be sent
    fmt.Printf("Response: %d - %s\n", httpErr.Code, httpErr.Message)
}

func main() {
    errors := []error{
        ErrNotFound,
        ErrUnauthorized,
        fmt.Errorf("user lookup: %w", ErrNotFound),
        fmt.Errorf("database: %w", ErrInternal),
        ErrBadRequest,
    }

    for _, err := range errors {
        fmt.Printf("Error: %v\n", err)
        errorHandler(err)
        fmt.Println()
    }
}
```

---

## Summary

In this tutorial, you learned:
- How Go's error interface works
- Creating errors with `errors.New` and `fmt.Errorf`
- Proper error checking patterns
- Wrapping errors with context using `%w`
- Inspecting errors with `errors.Is` and `errors.As`
- Best practices for error handling

---

**Next:** [05-packages-modules.md](05-packages-modules.md) - Learn about packages and modules in Go

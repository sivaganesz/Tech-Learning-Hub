# Gin Web Framework

## Learning Objectives

By the end of this tutorial, you will be able to:
- Install and set up the Gin framework
- Create routes for different HTTP methods
- Handle route parameters and query strings
- Bind request data (JSON, form, query)
- Return responses in various formats
- Create and use middleware
- Organize routes with groups
- Handle errors properly

---

## 1. Installation and Setup

### Installing Gin

```bash
# Create a new project
mkdir gin-tutorial
cd gin-tutorial
go mod init gin-tutorial

# Install Gin
go get -u github.com/gin-gonic/gin
```

### Basic Server

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    // Create default Gin router with logging and recovery middleware
    router := gin.Default()

    // Define a route
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello, World!",
        })
    })

    // Start server on port 8080
    router.Run(":8080")
}
```

### Minimal Setup

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    // Create router without default middleware
    router := gin.New()

    // Add middleware manually
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    router.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })

    router.Run(":8080")
}
```

### Different Modes

```go
package main

import (
    "os"

    "github.com/gin-gonic/gin"
)

func main() {
    // Set mode before creating router
    // Options: gin.DebugMode, gin.ReleaseMode, gin.TestMode

    // From environment variable
    gin.SetMode(os.Getenv("GIN_MODE"))

    // Or explicitly
    gin.SetMode(gin.ReleaseMode)

    router := gin.Default()
    // ... routes
    router.Run(":8080")
}
```

---

## 2. Basic Routing

### HTTP Methods

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // GET - retrieve resources
    router.GET("/users", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"action": "list users"})
    })

    // POST - create resources
    router.POST("/users", func(c *gin.Context) {
        c.JSON(http.StatusCreated, gin.H{"action": "create user"})
    })

    // PUT - update resources (full update)
    router.PUT("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{"action": "update user", "id": id})
    })

    // PATCH - partial update
    router.PATCH("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{"action": "partial update", "id": id})
    })

    // DELETE - remove resources
    router.DELETE("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{"action": "delete user", "id": id})
    })

    // Handle multiple methods
    router.Any("/any", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"method": c.Request.Method})
    })

    router.Run(":8080")
}
```

### Handler Functions

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Inline handler
func main() {
    router := gin.Default()

    // Inline handler
    router.GET("/inline", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"type": "inline"})
    })

    // Named handler function
    router.GET("/users", getUsers)
    router.GET("/users/:id", getUserByID)
    router.POST("/users", createUser)

    router.Run(":8080")
}

func getUsers(c *gin.Context) {
    users := []gin.H{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
    }
    c.JSON(http.StatusOK, users)
}

func getUserByID(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{
        "id":   id,
        "name": "User " + id,
    })
}

func createUser(c *gin.Context) {
    // Process request...
    c.JSON(http.StatusCreated, gin.H{
        "message": "User created",
    })
}
```

---

## 3. Route Parameters and Query Strings

### Path Parameters

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // Single parameter
    router.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{"user_id": id})
    })

    // Multiple parameters
    router.GET("/posts/:year/:month/:day", func(c *gin.Context) {
        year := c.Param("year")
        month := c.Param("month")
        day := c.Param("day")
        c.JSON(http.StatusOK, gin.H{
            "year":  year,
            "month": month,
            "day":   day,
        })
    })

    // Wildcard parameter (catches rest of path)
    router.GET("/files/*filepath", func(c *gin.Context) {
        filepath := c.Param("filepath")
        c.JSON(http.StatusOK, gin.H{"filepath": filepath})
    })

    router.Run(":8080")
}
```

### Query Strings

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // /search?q=golang&page=1
    router.GET("/search", func(c *gin.Context) {
        query := c.Query("q")              // Returns "" if not present
        page := c.DefaultQuery("page", "1") // With default value

        c.JSON(http.StatusOK, gin.H{
            "query": query,
            "page":  page,
        })
    })

    // Get as array: /tags?tag=go&tag=web&tag=api
    router.GET("/tags", func(c *gin.Context) {
        tags := c.QueryArray("tag")
        c.JSON(http.StatusOK, gin.H{"tags": tags})
    })

    // Get as map: /filter?filters[status]=active&filters[role]=admin
    router.GET("/filter", func(c *gin.Context) {
        filters := c.QueryMap("filters")
        c.JSON(http.StatusOK, gin.H{"filters": filters})
    })

    router.Run(":8080")
}
```

---

## 4. Request Binding

### JSON Binding

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"required,gte=0,lte=130"`
    Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
    Username string `json:"username" binding:"omitempty,min=3,max=20"`
    Email    string `json:"email" binding:"omitempty,email"`
    Age      int    `json:"age" binding:"omitempty,gte=0,lte=130"`
}

func main() {
    router := gin.Default()

    // Bind JSON body
    router.POST("/users", func(c *gin.Context) {
        var req CreateUserRequest

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{
            "message":  "User created",
            "username": req.Username,
            "email":    req.Email,
        })
    })

    // Bind with validation
    router.PUT("/users/:id", func(c *gin.Context) {
        var req UpdateUserRequest

        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        id := c.Param("id")
        c.JSON(http.StatusOK, gin.H{
            "id":      id,
            "updated": req,
        })
    })

    router.Run(":8080")
}
```

### Form Binding

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type LoginForm struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required"`
    Remember bool   `form:"remember"`
}

func main() {
    router := gin.Default()

    // Form data (application/x-www-form-urlencoded)
    router.POST("/login", func(c *gin.Context) {
        var form LoginForm

        if err := c.ShouldBind(&form); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "username": form.Username,
            "remember": form.Remember,
        })
    })

    // Multipart form (file upload)
    router.POST("/upload", func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Save file
        dst := "./" + file.Filename
        if err := c.SaveUploadedFile(file, dst); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "filename": file.Filename,
            "size":     file.Size,
        })
    })

    router.Run(":8080")
}
```

### Query Binding

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type SearchQuery struct {
    Query   string   `form:"q" binding:"required"`
    Page    int      `form:"page,default=1" binding:"min=1"`
    PerPage int      `form:"per_page,default=20" binding:"min=1,max=100"`
    Tags    []string `form:"tags"`
    Sort    string   `form:"sort,default=created_at"`
    Order   string   `form:"order,default=desc" binding:"oneof=asc desc"`
}

func main() {
    router := gin.Default()

    // Bind query parameters
    router.GET("/search", func(c *gin.Context) {
        var query SearchQuery

        if err := c.ShouldBindQuery(&query); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "query":    query.Query,
            "page":     query.Page,
            "per_page": query.PerPage,
            "tags":     query.Tags,
            "sort":     query.Sort,
            "order":    query.Order,
        })
    })

    router.Run(":8080")
}
```

### URI Binding

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type UserURI struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

type PostURI struct {
    UserID int64 `uri:"user_id" binding:"required,min=1"`
    PostID int64 `uri:"post_id" binding:"required,min=1"`
}

func main() {
    router := gin.Default()

    // Bind URI parameters
    router.GET("/users/:id", func(c *gin.Context) {
        var uri UserURI

        if err := c.ShouldBindUri(&uri); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"user_id": uri.ID})
    })

    router.GET("/users/:user_id/posts/:post_id", func(c *gin.Context) {
        var uri PostURI

        if err := c.ShouldBindUri(&uri); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "user_id": uri.UserID,
            "post_id": uri.PostID,
        })
    })

    router.Run(":8080")
}
```

---

## 5. Response Formats

### JSON Response

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}

func main() {
    router := gin.Default()

    // Using gin.H (map shorthand)
    router.GET("/simple", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello",
            "status":  "success",
        })
    })

    // Using struct
    router.GET("/user", func(c *gin.Context) {
        user := User{
            ID:        1,
            Username:  "alice",
            Email:     "alice@example.com",
            CreatedAt: "2024-01-15T10:30:00Z",
        }
        c.JSON(http.StatusOK, user)
    })

    // Slice of structs
    router.GET("/users", func(c *gin.Context) {
        users := []User{
            {ID: 1, Username: "alice", Email: "alice@example.com"},
            {ID: 2, Username: "bob", Email: "bob@example.com"},
        }
        c.JSON(http.StatusOK, users)
    })

    // Pretty JSON (indented)
    router.GET("/pretty", func(c *gin.Context) {
        data := gin.H{
            "nested": gin.H{
                "key": "value",
            },
        }
        c.IndentedJSON(http.StatusOK, data)
    })

    // Secure JSON (prevents JSON hijacking)
    router.GET("/secure", func(c *gin.Context) {
        c.SecureJSON(http.StatusOK, []string{"a", "b", "c"})
    })

    router.Run(":8080")
}
```

### Other Response Formats

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type XMLUser struct {
    ID       int    `xml:"id,attr"`
    Username string `xml:"username"`
    Email    string `xml:"email"`
}

func main() {
    router := gin.Default()

    // XML response
    router.GET("/xml", func(c *gin.Context) {
        user := XMLUser{
            ID:       1,
            Username: "alice",
            Email:    "alice@example.com",
        }
        c.XML(http.StatusOK, user)
    })

    // YAML response
    router.GET("/yaml", func(c *gin.Context) {
        c.YAML(http.StatusOK, gin.H{
            "name": "Alice",
            "age":  30,
        })
    })

    // String response
    router.GET("/string", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello, %s!", "World")
    })

    // HTML response
    router.GET("/html", func(c *gin.Context) {
        c.Data(http.StatusOK, "text/html; charset=utf-8",
            []byte("<h1>Hello World</h1>"))
    })

    // Redirect
    router.GET("/redirect", func(c *gin.Context) {
        c.Redirect(http.StatusMovedPermanently, "https://example.com")
    })

    // File download
    router.GET("/download", func(c *gin.Context) {
        c.File("./file.pdf")
    })

    // Attachment (forces download)
    router.GET("/attachment", func(c *gin.Context) {
        c.FileAttachment("./file.pdf", "downloaded-file.pdf")
    })

    router.Run(":8080")
}
```

---

## 6. Middleware

### Built-in Middleware

```go
package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    // With default middleware (Logger + Recovery)
    router := gin.Default()

    // Or add manually
    router = gin.New()
    router.Use(gin.Logger())   // Request logging
    router.Use(gin.Recovery()) // Panic recovery

    router.GET("/", func(c *gin.Context) {
        c.String(200, "OK")
    })

    router.Run(":8080")
}
```

### Custom Middleware

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// Simple middleware
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // Process request
        c.Next()

        // After request
        duration := time.Since(start)
        status := c.Writer.Status()

        fmt.Printf("[%s] %s %s %d %v\n",
            time.Now().Format("2006-01-02 15:04:05"),
            c.Request.Method,
            c.Request.URL.Path,
            status,
            duration,
        )
    }
}

// Authentication middleware
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            c.Abort()
            return
        }

        // Validate token (simplified)
        if token != "Bearer valid-token" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token",
            })
            c.Abort()
            return
        }

        // Set user info in context
        c.Set("user_id", 123)
        c.Set("username", "alice")

        c.Next()
    }
}

// CORS middleware
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

func main() {
    router := gin.New()

    // Global middleware
    router.Use(Logger())
    router.Use(CORS())
    router.Use(gin.Recovery())

    // Public routes
    router.GET("/public", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "public"})
    })

    // Protected routes
    protected := router.Group("/api")
    protected.Use(AuthRequired())
    {
        protected.GET("/profile", func(c *gin.Context) {
            userID := c.GetInt("user_id")
            username := c.GetString("username")
            c.JSON(http.StatusOK, gin.H{
                "user_id":  userID,
                "username": username,
            })
        })
    }

    router.Run(":8080")
}
```

### Request Timeout Middleware

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()

        c.Request = c.Request.WithContext(ctx)

        finished := make(chan struct{})

        go func() {
            c.Next()
            close(finished)
        }()

        select {
        case <-finished:
            // Request completed
        case <-ctx.Done():
            c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
                "error": "Request timeout",
            })
        }
    }
}

func main() {
    router := gin.Default()

    // Apply timeout to specific routes
    router.GET("/slow", TimeoutMiddleware(2*time.Second), func(c *gin.Context) {
        time.Sleep(3 * time.Second) // Simulates slow operation
        c.JSON(http.StatusOK, gin.H{"message": "done"})
    })

    router.GET("/fast", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "fast"})
    })

    router.Run(":8080")
}
```

---

## 7. Route Groups

### Basic Grouping

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // API v1 group
    v1 := router.Group("/api/v1")
    {
        v1.GET("/users", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"version": "v1", "users": []string{}})
        })
        v1.GET("/products", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"version": "v1", "products": []string{}})
        })
    }

    // API v2 group
    v2 := router.Group("/api/v2")
    {
        v2.GET("/users", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"version": "v2", "users": []string{}})
        })
        v2.GET("/products", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"version": "v2", "products": []string{}})
        })
    }

    router.Run(":8080")
}
```

### Nested Groups with Middleware

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if user is admin
        isAdmin := c.GetHeader("X-Admin") == "true"
        if !isAdmin {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin required"})
            c.Abort()
            return
        }
        c.Next()
    }
}

func main() {
    router := gin.Default()

    // Public routes
    public := router.Group("/")
    {
        public.GET("/health", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"status": "healthy"})
        })
        public.POST("/login", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"token": "xxx"})
        })
    }

    // Authenticated routes
    api := router.Group("/api")
    api.Use(AuthMiddleware())
    {
        // User routes
        users := api.Group("/users")
        {
            users.GET("", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"users": []string{}})
            })
            users.GET("/:id", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"user": c.Param("id")})
            })
            users.POST("", func(c *gin.Context) {
                c.JSON(http.StatusCreated, gin.H{"created": true})
            })
        }

        // Admin routes (nested middleware)
        admin := api.Group("/admin")
        admin.Use(AdminMiddleware())
        {
            admin.GET("/stats", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"stats": "admin only"})
            })
            admin.DELETE("/users/:id", func(c *gin.Context) {
                c.JSON(http.StatusOK, gin.H{"deleted": c.Param("id")})
            })
        }
    }

    router.Run(":8080")
}
```

---

## 8. Error Handling

### Basic Error Handling

```go
package main

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
)

var (
    ErrNotFound     = errors.New("resource not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrBadRequest   = errors.New("bad request")
)

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func main() {
    router := gin.Default()

    router.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")

        // Simulate different errors
        switch id {
        case "0":
            c.JSON(http.StatusBadRequest, APIError{
                Code:    http.StatusBadRequest,
                Message: "Invalid user ID",
            })
            return
        case "999":
            c.JSON(http.StatusNotFound, APIError{
                Code:    http.StatusNotFound,
                Message: "User not found",
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "id":   id,
            "name": "User " + id,
        })
    })

    router.Run(":8080")
}
```

### Custom Error Handler Middleware

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // Check if there are any errors
        if len(c.Errors) > 0 {
            err := c.Errors.Last()

            // Check for AppError
            var appErr *AppError
            if errors.As(err.Err, &appErr) {
                c.JSON(appErr.Code, gin.H{
                    "error": appErr.Message,
                })
                return
            }

            // Default error response
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
    }
}

func main() {
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())
    router.Use(ErrorHandler())

    router.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")

        if id == "0" {
            c.Error(&AppError{
                Code:    http.StatusBadRequest,
                Message: "Invalid user ID",
            })
            return
        }

        if id == "999" {
            c.Error(&AppError{
                Code:    http.StatusNotFound,
                Message: "User not found",
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{"id": id, "name": "User " + id})
    })

    router.Run(":8080")
}

import "errors"
```

### Panic Recovery

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func CustomRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the error
                // Send to error tracking service

                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                    "code":  "PANIC_RECOVERED",
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}

func main() {
    router := gin.New()
    router.Use(gin.Logger())
    router.Use(CustomRecovery())

    router.GET("/panic", func(c *gin.Context) {
        panic("Something went wrong!")
    })

    router.GET("/ok", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    router.Run(":8080")
}
```

---

## Exercises

### Exercise 1: RESTful API for Books
Create a complete CRUD API for books.

```go
package main

import "github.com/gin-gonic/gin"

type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title" binding:"required"`
    Author string `json:"author" binding:"required"`
    Year   int    `json:"year" binding:"required,min=1000,max=2100"`
}

// TODO: Implement in-memory storage and handlers
// GET    /books         - List all books
// GET    /books/:id     - Get book by ID
// POST   /books         - Create book
// PUT    /books/:id     - Update book
// DELETE /books/:id     - Delete book

func main() {
    router := gin.Default()

    // TODO: Add routes

    router.Run(":8080")
}
```

### Exercise 2: Authentication Middleware
Implement JWT-style authentication middleware.

```go
package main

import "github.com/gin-gonic/gin"

// TODO: Implement auth middleware that:
// - Checks for Authorization header
// - Validates token format "Bearer <token>"
// - Extracts user info from token
// - Sets user info in context

// TODO: Implement protected routes that use the middleware

func main() {
    router := gin.Default()

    // Public
    router.POST("/login", loginHandler)

    // Protected
    protected := router.Group("/api")
    // TODO: Add auth middleware
    {
        protected.GET("/profile", getProfile)
        protected.GET("/settings", getSettings)
    }

    router.Run(":8080")
}
```

### Exercise 3: Rate Limiting Middleware
Create middleware that limits requests per IP.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "sync"
    "time"
)

// TODO: Implement rate limiter that:
// - Tracks requests per IP address
// - Allows X requests per minute
// - Returns 429 Too Many Requests when exceeded

type RateLimiter struct {
    // TODO: Add fields
    mu sync.Mutex
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
    // TODO: Implement
    return nil
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    // TODO: Implement
    return nil
}

func main() {
    router := gin.Default()

    limiter := NewRateLimiter(10) // 10 requests per minute
    router.Use(limiter.Middleware())

    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "ok"})
    })

    router.Run(":8080")
}
```

### Exercise 4: File Upload API
Create an API for uploading and serving files.

```go
package main

import "github.com/gin-gonic/gin"

// TODO: Implement:
// POST /upload       - Upload single file
// POST /upload/multi - Upload multiple files
// GET  /files/:name  - Download file
// GET  /files        - List all files

func main() {
    router := gin.Default()

    // Set max upload size
    router.MaxMultipartMemory = 8 << 20 // 8 MB

    // TODO: Add routes

    router.Run(":8080")
}
```

### Exercise 5: Complete API with All Features
Build a user management API with all learned features.

```go
package main

import "github.com/gin-gonic/gin"

// TODO: Implement complete API with:
// - User CRUD operations
// - Authentication middleware
// - Request validation
// - Error handling
// - Route groups (public, authenticated, admin)
// - Logging middleware
// - CORS middleware

// Models
type User struct {
    ID       string `json:"id"`
    Username string `json:"username" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Role     string `json:"role" binding:"oneof=user admin"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func main() {
    router := gin.New()

    // Global middleware
    // TODO: Add Logger, Recovery, CORS

    // Public routes
    // TODO: POST /login, POST /register

    // Authenticated routes
    // TODO: GET /profile, PUT /profile

    // Admin routes
    // TODO: GET /admin/users, DELETE /admin/users/:id

    router.Run(":8080")
}
```

---

## Summary

In this tutorial, you learned:
- Installing and setting up Gin
- Creating routes for different HTTP methods
- Handling path parameters and query strings
- Binding request data with validation
- Returning responses in various formats
- Creating custom middleware
- Organizing routes with groups
- Handling errors properly

---

## What's Next?

With these Go basics under your belt, you're ready to:
- Build production-ready APIs
- Implement database integration
- Add authentication and authorization
- Deploy your applications

Happy coding with Go and Gin!

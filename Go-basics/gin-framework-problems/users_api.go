package main

import (
	"net/http"
	"os"
	"sync"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"oneof=user admin"`
	Password string `json:"-"` // not returned
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var (
	users   = map[string]User{} // id -> user
	usersMu sync.Mutex
	idSeq   = 1

	// tokens to userID
	tokens   = map[string]string{}
	tokensMu sync.Mutex
)

func nextID() string {
	id := idSeq
	idSeq++
	return fmt.Sprintf("%d", id)
}

// ---- Helpers ----
func findUserByUsername(username string) (User, bool) {
	usersMu.Lock()
	defer usersMu.Unlock()
	for _, u := range users {
		if u.Username == username {
			return u, true
		}
	}
	return User{}, false
}

func registerHandler(c *gin.Context) {
	var u User
	var raw struct {
		Username string `json:"username" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// ensure unique username/email
	if _, ok := findUserByUsername(raw.Username); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
		return
	}
	u = User{
		ID:       nextID(),
		Username: raw.Username,
		Email:    raw.Email,
		Role:     "user",
		Password: raw.Password,
	}

	usersMu.Lock()
	users[u.ID] = u
	usersMu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"role":     u.Role,
	})
}

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, ok := findUserByUsername(req.Username)
	if !ok || u.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// create token
	token := fmt.Sprintf("tk_%s_%d", u.ID, len(tokens)+1)
	tokensMu.Lock()
	tokens[token] = u.ID
	tokensMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "bad auth header"})
			return
		}
		token := parts[1]

		tokensMu.Lock()
		uid, ok := tokens[token]
		tokensMu.Unlock()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		usersMu.Lock()
		user, ok := users[uid]
		usersMu.Unlock()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func requireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		user := v.(User)
		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

// ---- Handlers ----
func getProfile(c *gin.Context) {
	u := c.MustGet("user").(User)
	// hide password
	c.JSON(http.StatusOK, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
		"role":     u.Role,
	})
}

func updateProfile(c *gin.Context) {
	u := c.MustGet("user").(User)
	var req struct {
		Email string `json:"email" binding:"omitempty,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	usersMu.Lock()
	stored := users[u.ID]
	if req.Email != "" {
		stored.Email = req.Email
	}
	users[u.ID] = stored
	usersMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func adminListUsers(c *gin.Context) {
	usersMu.Lock()
	defer usersMu.Unlock()
	out := []gin.H{}
	for _, u := range users {
		out = append(out, gin.H{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
			"role":     u.Role,
		})
	}
	c.JSON(http.StatusOK, out)
}

func adminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	usersMu.Lock()
	defer usersMu.Unlock()
	if _, ok := users[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	delete(users, id)
	c.Status(http.StatusNoContent)
}

// ---- Middleware: simple CORS ----
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	// create a default admin user
	admin := User{
		ID:       nextID(),
		Username: "admin",
		Email:    "admin@example.com",
		Role:     "admin",
		Password: "admin123",
	}
	users[admin.ID] = admin

	router := gin.New()
	// Logging and recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Public
	public := router.Group("/api")
	{
		public.POST("/register", registerHandler)
		public.POST("/login", loginHandler)
	}

	// Authenticated
	private := router.Group("/api")
	private.Use(authMiddleware())
	{
		private.GET("/profile", getProfile)
		private.PUT("/profile", updateProfile)
	}

	// Admin
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(authMiddleware(), requireAdmin())
	{
		adminRoutes.GET("/users", adminListUsers)
		adminRoutes.DELETE("/users/:id", adminDeleteUser)
	}

	// make sure uploads dir exists for potential file endpoints
	_ = os.MkdirAll("./uploads", 0755)

	router.Run(":8080")
}

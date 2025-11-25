package main

import (
	"net/http"
	"strings"
	"sync"
	"fmt"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

var (
	// in-memory users (username->password,role)
	users = map[string]struct {
		Password string
		Role     string
	}{
		"alice": {Password: "password1", Role: "user"},
		"bob":   {Password: "adminpass", Role: "admin"},
	}

	// token -> UserInfo
	tokens   = map[string]UserInfo{}
	tokensMu sync.Mutex
)

func loginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, ok := users[req.Username]
	if !ok || u.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// create a simple token: username + ":" + role + ":" + counter
	token := createTokenForUser(req.Username, u.Role)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func createTokenForUser(username, role string) string {
	tokensMu.Lock()
	defer tokensMu.Unlock()
	// simple token generation (not secure for production)
	token := fmt.Sprintf("tok_%s_%d", username, len(tokens)+1)
	tokens[token] = UserInfo{Username: username, Role: role}
	return token
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization format"})
			return
		}
		token := parts[1]

		tokensMu.Lock()
		user, ok := tokens[token]
		tokensMu.Unlock()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// attach user info to context
		c.Set("user", user)
		c.Next()
	}
}

func getProfile(c *gin.Context) {
	u, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"profile": u})
}

func getSettings(c *gin.Context) {
	u, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"settings_for": u})
}

func main() {
	router := gin.Default()

	// Public
	router.POST("/login", loginHandler)

	// Protected
	protected := router.Group("/api")
	protected.Use(authMiddleware())
	{
		protected.GET("/profile", getProfile)
		protected.GET("/settings", getSettings)
	}

	router.Run(":8080")
}


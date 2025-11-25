package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requestsPerMinute int
	mu                sync.Mutex
	clients           map[string][]time.Time
	window            time.Duration
}

func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		requestsPerMinute: requestsPerMinute,
		clients:           make(map[string][]time.Time),
		window:            time.Minute,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		rl.mu.Lock()
		timestamps := rl.clients[ip]

		// prune older than window
		pruned := make([]time.Time, 0, len(timestamps))
		for _, t := range timestamps {
			if now.Sub(t) <= rl.window {
				pruned = append(pruned, t)
			}
		}

		if len(pruned) >= rl.requestsPerMinute {
			// exceeded
			rl.clients[ip] = pruned
			rl.mu.Unlock()
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		// allow and record
		pruned = append(pruned, now)
		rl.clients[ip] = pruned
		rl.mu.Unlock()

		c.Next()
	}
}

func main() {
	router := gin.Default()

	limiter := NewRateLimiter(10) // 10 requests per minute per IP
	router.Use(limiter.Middleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	router.Run(":8080")
}

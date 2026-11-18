package middleware

import (
	"net/http"
	"sync"
	"time"

	"app/internal/response"

	"github.com/gin-gonic/gin"
)

type clientLimit struct {
	count     int
	resetTime time.Time
}

type MemoryRateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientLimit
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *MemoryRateLimiter {
	return &MemoryRateLimiter{
		clients: make(map[string]*clientLimit),
		limit:   limit,
		window:  window,
	}
}

func (rl *MemoryRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.mu.Lock()
		now := time.Now()

		client, exists := rl.clients[ip]
		if !exists || now.After(client.resetTime) {
			rl.clients[ip] = &clientLimit{
				count:     1,
				resetTime: now.Add(rl.window),
			}
			rl.mu.Unlock()
			c.Next()
			return
		}

		if client.count >= rl.limit {
			rl.mu.Unlock()
			response.Error(c, http.StatusTooManyRequests, "Terlalu banyak percobaan. Coba lagi dalam beberapa saat.", nil)
			c.Abort()
			return
		}

		client.count++
		rl.mu.Unlock()
		c.Next()
	}
}

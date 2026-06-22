package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SimpleRateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int
	window   time.Duration
}

type visitor struct {
	count       int
	windowStart time.Time
}

func NewSimpleRateLimiter(requestsPerSecond int) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerSecond,
		window:   time.Second,
	}
}

func (rl *SimpleRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}

func (rl *SimpleRateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	now := time.Now()

	if !exists || now.Sub(v.windowStart) > rl.window {
		rl.visitors[ip] = &visitor{count: 1, windowStart: now}
		return true
	}

	if v.count >= rl.rate {
		return false
	}

	v.count++
	return true
}

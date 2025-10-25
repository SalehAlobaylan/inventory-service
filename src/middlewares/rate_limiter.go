package middlewares

import (
	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
	
)

var limiter = rate.NewLimiter(1, 5) // 1 req/sec, burst of 5

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
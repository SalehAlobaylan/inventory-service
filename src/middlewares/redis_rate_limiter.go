package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	rateLimiter *redis_rate.Limiter
)

// InitRedisRateLimiter initializes Redis connection for rate limiting
func InitRedisRateLimiter(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	redisClient = redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	rateLimiter = redis_rate.NewLimiter(redisClient)
	return nil
}

// RedisRateLimiter creates a Redis-based rate limiting middleware
func RedisRateLimiter(requestsPerSecond int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Use client IP as the key for per-IP rate limiting
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		// Check rate limit
		limit := redis_rate.PerSecond(requestsPerSecond)
		limit.Burst = burst

		result, err := rateLimiter.Allow(ctx, key, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter error",
			})
			c.Abort()
			return
		}

		// Set rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerSecond))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))

		// Check if rate limit exceeded
		if result.Allowed == 0 {
			retryAfterSeconds := result.RetryAfter.Seconds()
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfterSeconds))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": fmt.Sprintf("%.0f seconds", retryAfterSeconds),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RedisRateLimiterByUser creates rate limiting based on authenticated user
func RedisRateLimiterByUser(requestsPerSecond int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract user ID from context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			// Fall back to IP-based limiting for unauthenticated requests
			userID = c.ClientIP()
		}

		key := fmt.Sprintf("rate_limit:user:%v", userID)

		limit := redis_rate.PerSecond(requestsPerSecond)
		limit.Burst = burst

		result, err := rateLimiter.Allow(ctx, key, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter error",
			})
			c.Abort()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerSecond))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))

		if result.Allowed == 0 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests",
				"retry_after": result.RetryAfter.String(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

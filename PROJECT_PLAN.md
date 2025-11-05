# Inventory Management System - Project Plan

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture & Technology Stack](#architecture--technology-stack)
3. [Development Phases](#development-phases)
4. [Implementation Details](#implementation-details)
5. [Concerns & Considerations](#concerns--considerations)
6. [Testing Strategy](#testing-strategy)
7. [Deployment Guide](#deployment-guide)
8. [Stand-Out Features](#stand-out-features)

---

## Project Overview

### Goal

Build a RESTful API backend server for an inventory management system using Go, Gin framework, PostgreSQL, and GORM.

### Core Features

- **CRUD Operations**: Create, Read, Update, Delete inventory items
- **Data Model**: Item (ID, Name, Stock, Price)
- **Advanced Features**:
  - Rate limiting (1 req/sec, burst of 5)
  - Pagination for large datasets
  - Sorting (by name, stock, price - asc/desc)
  - Filtering (by name, minimum stock)
  - UUID-based primary keys

### Success Criteria

- ✅ All CRUD endpoints functional
- ✅ PostgreSQL database with GORM integration
- ✅ Proper HTTP status codes and error handling
- ✅ JSON responses for all endpoints
- ✅ Rate limiting middleware
- ✅ Pagination, sorting, and filtering support
- ✅ Clean, documented, maintainable code

---

## Architecture & Technology Stack

### Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin (HTTP routing, middleware)
- **Database**: PostgreSQL 15+
- **ORM**: GORM (database interactions)
- **ID Generation**: google/uuid
- **Rate Limiting**: golang.org/x/time/rate

### Project Structure

```
inventory-service/
├── src/
│   ├── main.go                 # Application entry point
│   ├── models/
│   │   └── item.go             # Item struct definition
│   ├── controllers/
│   │   └── item_controller.go  # CRUD handlers
│   ├── routes/
│   │   └── router.go           # Route registration
│   ├── middlewares/
│   │   ├── logger.go           # Logging middleware
│   │   └── rate_limiter.go     # Rate limiting middleware
│   ├── utils/
│   │   └── database.go         # DB connection & config
│   └── seeds/
│       └── seeder.go           # Database seeding
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yaml
├── .env.example
├── .gitignore
└── README.md
```

---

## Development Phases

### Phase 1: Environment Setup ✅

**Objective**: Set up PostgreSQL, Go environment, and project dependencies

**Tasks**:

1. Install PostgreSQL locally or use Docker
2. Create database: `inventory_db`
3. Initialize Go module: `go mod init inventory-service`
4. Install dependencies:
   ```bash
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   go get gorm.io/driver/postgres
   go get github.com/google/uuid
   go get golang.org/x/time/rate
   ```

**Concerns**:

- ⚠️ PostgreSQL version compatibility (use 12+)
- ⚠️ Connection string format: `postgres://user:password@host:port/dbname?sslmode=disable`
- ⚠️ Environment variable management (use `.env` file)

---

### Phase 2: Database & Data Models ✅

**Objective**: Define Item model, set up GORM, create migrations

**Tasks**:

1. **Define Item struct** (`src/models/item.go`):

   ```go
   type Item struct {
       ID        string    `json:"id" gorm:"type:uuid;primary_key"`
       Name      string    `json:"name" gorm:"type:varchar(255);not null"`
       Stock     int       `json:"stock" gorm:"not null"`
       Price     float64   `json:"price" gorm:"not null"`
       CreatedAt time.Time `json:"created_at"`
       UpdatedAt time.Time `json:"updated_at"`
   }
   ```

2. **Database connection** (`src/utils/database.go`):

   - Singleton pattern for DB instance
   - Read `DATABASE_URL` from environment
   - Error handling for connection failures

3. **Auto-migration**:

   - Run `db.AutoMigrate(&models.Item{})` on startup

4. **Seed database** (`src/seeds/seeder.go`):
   - Check if table is empty
   - Insert 5+ sample items
   - Generate UUIDs using `uuid.NewString()`

**Concerns**:

- ⚠️ UUID generation: Use GORM hooks (`BeforeCreate`) to auto-generate IDs
- ⚠️ Data types: Stock (int), Price (float64), ensure no negative values
- ⚠️ Timestamps: GORM handles `CreatedAt`/`UpdatedAt` automatically
- ⚠️ Idempotent seeding: Don't re-seed if data exists

---

### Phase 3: Core CRUD Endpoints 🔄

**Objective**: Implement 5 RESTful endpoints with proper handlers

**Endpoints**:

| Method | Path           | Handler     | Status Codes       |
| ------ | -------------- | ----------- | ------------------ |
| GET    | /inventory     | GetItems    | 200, 500           |
| GET    | /inventory/:id | GetItemByID | 200, 404, 500      |
| POST   | /inventory     | CreateItem  | 201, 400, 500      |
| PUT    | /inventory/:id | UpdateItem  | 200, 400, 404, 500 |
| DELETE | /inventory/:id | DeleteItem  | 204, 404, 500      |

**Handler Details**:

1. **GetItems** (GET /inventory):

   - Returns all items as JSON array
   - Status: 200 OK

2. **GetItemByID** (GET /inventory/:id):

   - Extract `:id` from URL params
   - Query: `db.First(&item, "id = ?", id)`
   - Status: 200 (found), 404 (not found)

3. **CreateItem** (POST /inventory):

   - Bind JSON request body to Item struct
   - Validate required fields (Name, Stock, Price)
   - UUID auto-generated in `BeforeCreate` hook
   - Status: 201 (created), 400 (bad input)

4. **UpdateItem** (PUT /inventory/:id):

   - Find existing item by ID (404 if not found)
   - Support partial updates (only update provided fields)
   - Use pointer fields to detect "not provided" vs "zero value"
   - Status: 200 (updated), 404 (not found), 400 (bad input)

5. **DeleteItem** (DELETE /inventory/:id):
   - Find item by ID (404 if not found)
   - Soft delete with GORM: `db.Delete(&item)`
   - Status: 204 (no content), 404 (not found)

**Concerns**:

- ⚠️ Error handling: Always return JSON errors, not plain text
- ⚠️ Input validation: Check for empty names, negative stock/price
- ⚠️ 404 responses: Return `{"error": "item not found"}` with 404 status
- ⚠️ Partial updates: Use struct with pointer fields for PUT requests
- ⚠️ Idempotency: DELETE should return 404 if already deleted

---

### Phase 4: Rate Limiting Middleware

**Objective**: Prevent API abuse with token bucket rate limiter

**Implementation** (`src/middlewares/rate_limiter.go`):

```go
import "golang.org/x/time/rate"

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
```

**Apply globally**:

```go
router.Use(middlewares.RateLimiter())
```

**Concerns**:

- ⚠️ Global vs per-IP limiting: Current implementation is global (shared across all clients)
- ⚠️ For production: Use per-IP limiting with a map of limiters
- ⚠️ Redis alternative: For distributed systems, use Redis-based rate limiting
- ⚠️ Status code: Return 429 (Too Many Requests)

#### 4.1 Redis-Based Rate Limiting (Production-Ready)

**Why Redis?**

- ✅ **Distributed**: Works across multiple server instances
- ✅ **Persistent**: Rate limits survive application restarts
- ✅ **Scalable**: Handles high-traffic applications
- ✅ **Per-IP**: Easy to implement per-client rate limiting
- ✅ **Flexible**: Support for multiple rate limit tiers (per-user, per-endpoint)

**Dependencies**:

```bash
go get github.com/go-redis/redis/v8
go get github.com/go-redis/redis_rate/v10
```

**Implementation** (`src/middlewares/redis_rate_limiter.go`):

```go
package middlewares

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "github.com/go-redis/redis_rate/v10"
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
            retryAfter := time.Until(result.RetryAfter).Seconds()
            c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))

            c.JSON(http.StatusTooManyRequests, gin.H{
                "error":       "too many requests",
                "retry_after": fmt.Sprintf("%.0f seconds", retryAfter),
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
                "retry_after": time.Until(result.RetryAfter).String(),
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
```

**Usage in main.go**:

```go
package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"

    "inventory-service/src/middlewares"
    "inventory-service/src/routes"
)

func main() {
    // Initialize Redis rate limiter
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        redisURL = "redis://localhost:6379/0"
    }

    if err := middlewares.InitRedisRateLimiter(redisURL); err != nil {
        log.Fatalf("Failed to initialize Redis rate limiter: %v", err)
    }
    defer middlewares.CloseRedis()

    router := gin.Default()

    // Apply Redis-based rate limiting globally
    router.Use(middlewares.RedisRateLimiter(1, 5)) // 1 req/sec, burst of 5

    routes.RegisterRoutes(router)

    if err := router.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

**Advanced: Per-Endpoint Rate Limiting**:

```go
// Different rate limits for different endpoints
func setupRoutes(router *gin.Engine) {
    // Public endpoints: 10 req/sec
    public := router.Group("/inventory")
    public.Use(middlewares.RedisRateLimiter(10, 20))
    {
        public.GET("", controllers.GetItems)
        public.GET("/:id", controllers.GetItemByID)
    }

    // Write endpoints: 5 req/sec
    protected := router.Group("/inventory")
    protected.Use(middlewares.RedisRateLimiter(5, 10))
    {
        protected.POST("", controllers.CreateItem)
        protected.PUT("/:id", controllers.UpdateItem)
        protected.DELETE("/:id", controllers.DeleteItem)
    }
}
```

**Docker Compose with Redis**:

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:15-alpine
    container_name: inventory-postgres
    environment:
      POSTGRES_DB: inventory_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    container_name: inventory-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  app:
    build: .
    container_name: inventory-api
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/inventory_db?sslmode=disable
      REDIS_URL: redis://redis:6379/0
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

**Environment Variables**:

```bash
# .env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable
REDIS_URL=redis://localhost:6379/0
```

**Benefits of Redis-Based Rate Limiting**:

1. **Distributed Systems**: Multiple API servers share the same rate limit state
2. **Persistence**: Rate limit counters survive application restarts
3. **Accurate**: No race conditions in distributed environments
4. **Flexible Keys**: Can rate limit by IP, user ID, API key, etc.
5. **Multiple Strategies**: Support for different algorithms (token bucket, sliding window, fixed window)
6. **Headers**: Automatic rate limit headers (X-RateLimit-\*)
7. **Monitoring**: Easy to monitor rate limit usage via Redis

**Testing Redis Rate Limiter**:

```bash
# Test with multiple rapid requests
for i in {1..20}; do
  curl -w "\nStatus: %{http_code}\n" http://localhost:8080/inventory
  sleep 0.1
done

# Check rate limit headers
curl -I http://localhost:8080/inventory

# Monitor Redis keys
redis-cli
> KEYS rate_limit:*
> TTL rate_limit:127.0.0.1
> GET rate_limit:127.0.0.1
```

**Production Considerations**:

- ⚠️ **Redis Availability**: Use Redis Sentinel or Redis Cluster for high availability
- ⚠️ **Fallback Strategy**: Implement fallback to in-memory rate limiting if Redis is down
- ⚠️ **Key Expiration**: Set appropriate TTL on rate limit keys
- ⚠️ **Monitoring**: Monitor Redis memory usage and connection pool
- ⚠️ **Security**: Enable Redis authentication with `requirepass`
- ⚠️ **Network**: Use TLS for Redis connections in production

---

### Phase 5: Pagination, Sorting, Filtering

**Objective**: Handle large datasets efficiently

#### 5.1 Pagination

**Query Parameters**:

- `limit`: Number of items per page (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)

**Implementation**:

```go
func GetItems(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

    if limit > 100 { limit = 100 }
    if limit < 1 { limit = 10 }

    var items []models.Item
    db.Limit(limit).Offset(offset).Find(&items)

    c.JSON(200, items)
}
```

**Response format**:

```json
{
  "items": [...],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 150
  }
}
```

**Concerns**:

- ⚠️ Performance: Offset pagination is slow for large offsets (use cursor-based for production)
- ⚠️ Validation: Limit max page size to prevent memory issues
- ⚠️ Count query: `db.Model(&Item{}).Count(&total)` can be expensive

#### 5.2 Sorting

**Query Parameters**:

- `sort_by`: Field to sort by (name, stock, price)
- `order`: Sort direction (asc, desc)

**Implementation**:

```go
sortBy := c.DefaultQuery("sort_by", "created_at")
order := c.DefaultQuery("order", "desc")

allowedFields := map[string]bool{"name": true, "stock": true, "price": true, "created_at": true}
if !allowedFields[sortBy] {
    sortBy = "created_at"
}

if order != "asc" && order != "desc" {
    order = "desc"
}

db.Order(sortBy + " " + order).Find(&items)
```

**Concerns**:

- ⚠️ SQL injection: Whitelist allowed sort fields (don't directly interpolate user input)
- ⚠️ Case sensitivity: PostgreSQL sorting is case-sensitive
- ⚠️ Index support: Add database indexes on frequently sorted columns

#### 5.3 Filtering

**Query Parameters**:

- `name`: Filter by name (case-insensitive partial match)
- `min_stock`: Minimum stock level

**Implementation**:

```go
query := db.Model(&models.Item{})

if name := c.Query("name"); name != "" {
    query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
}

if minStock := c.Query("min_stock"); minStock != "" {
    query = query.Where("stock >= ?", minStock)
}

query.Find(&items)
```

**Concerns**:

- ⚠️ LIKE queries: Use indexes or full-text search for better performance
- ⚠️ Multiple filters: Combine with AND logic by default
- ⚠️ Empty filters: Validate that filter values are not empty strings

---

### Phase 6: Logging & Error Handling

**Objective**: Comprehensive logging and user-friendly errors

**Logging Middleware** (`src/middlewares/logger.go`):

```go
router.Use(gin.Logger())
router.Use(gin.Recovery())
```

**Structured Error Responses**:

```json
{
  "error": "item not found",
  "code": "ITEM_NOT_FOUND",
  "timestamp": "2025-10-24T10:30:00Z"
}
```

**Concerns**:

- ⚠️ Don't expose internal errors (database errors, stack traces)
- ⚠️ Log to file in production (use logrus or zap)
- ⚠️ Request ID tracking for debugging

---

### Phase 7: Testing

**Objective**: Ensure reliability with unit and integration tests

**Test Coverage**:

1. **Unit Tests** (`*_test.go`):

   - Model validation
   - Handler logic (mock database)
   - Middleware functionality

2. **Integration Tests**:
   - End-to-end API tests with test database
   - Use `httptest` package

**Example**:

```go
func TestCreateItem(t *testing.T) {
    router := setupTestRouter()

    body := `{"name":"Test Item","stock":10,"price":99.99}`
    req, _ := http.NewRequest("POST", "/inventory", strings.NewReader(body))
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)
}
```

**Test Database**:

- Use separate test database
- Clean up after each test

**Concerns**:

- ⚠️ Don't test against production database
- ⚠️ Use table-driven tests for multiple scenarios
- ⚠️ Test edge cases (empty input, invalid IDs, etc.)

---

### Phase 8: Docker & Deployment

**Objective**: Containerize application for easy deployment

**Dockerfile**:

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./src

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/inventory-service ./
EXPOSE 8080
ENTRYPOINT ["./inventory-service"]
```

**docker-compose.yaml**:

```yaml
version: "3.8"

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: inventory_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/inventory_db?sslmode=disable
    depends_on:
      - postgres

volumes:
  postgres_data:
```

**Concerns**:

- ⚠️ Health checks: Add `/health` endpoint for container orchestration
- ⚠️ Secrets management: Don't hardcode passwords (use Docker secrets)
- ⚠️ Database readiness: Wait for PostgreSQL to be ready before starting app

---

## Concerns & Considerations

### Security

- ❗ **SQL Injection**: Use parameterized queries (GORM handles this)
- ❗ **Input Validation**: Validate all user inputs (length, type, range)
- ❗ **CORS**: Configure allowed origins in production
- ❗ **Rate Limiting**: Implement per-IP limiting for production
- ❗ **Authentication**: Add JWT authentication for production use

### Performance

- ❗ **Database Indexes**: Add indexes on frequently queried fields
- ❗ **Connection Pooling**: Configure GORM connection pool
- ❗ **Pagination**: Use cursor-based pagination for large datasets
- ❗ **Caching**: Consider Redis for frequently accessed items

### Scalability

- ❗ **Horizontal Scaling**: Design for stateless instances
- ❗ **Database Replication**: Use read replicas for read-heavy workloads
- ❗ **Load Balancing**: Use nginx or cloud load balancer

### Monitoring

- ❗ **Metrics**: Expose Prometheus metrics
- ❗ **Logging**: Centralized logging (ELK stack)
- ❗ **Tracing**: Distributed tracing with OpenTelemetry

### Data Integrity

- ❗ **Transactions**: Use GORM transactions for multi-step operations
- ❗ **Soft Deletes**: GORM soft delete preserves data
- ❗ **Backups**: Regular PostgreSQL backups

---

## Testing Strategy

### Manual Testing with cURL

**Create Item**:

```bash
curl -X POST http://localhost:8080/inventory \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","stock":10,"price":999.99}'
```

**Get All Items**:

```bash
curl http://localhost:8080/inventory?limit=10&offset=0&sort_by=name&order=asc
```

**Get Single Item**:

```bash
curl http://localhost:8080/inventory/{id}
```

**Update Item**:

```bash
curl -X PUT http://localhost:8080/inventory/{id} \
  -H "Content-Type: application/json" \
  -d '{"stock":15}'
```

**Delete Item**:

```bash
curl -X DELETE http://localhost:8080/inventory/{id}
```

**Test Rate Limiting**:

```bash
for i in {1..10}; do curl http://localhost:8080/inventory; done
```

### Postman Collection

Create a Postman collection with:

- All CRUD operations
- Pagination examples
- Sorting/filtering examples
- Error scenarios (404, 400, 429)

---

## Deployment Guide

### Local Development

1. **Start PostgreSQL**:

   ```bash
   docker run -d --name postgres \
     -e POSTGRES_DB=inventory_db \
     -e POSTGRES_PASSWORD=postgres \
     -p 5432:5432 postgres:15-alpine
   ```

2. **Set environment variable**:

   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
   ```

3. **Run application**:
   ```bash
   go run src/main.go
   ```

### Docker Deployment

```bash
docker-compose up --build
```

### Production Checklist

- [ ] Use environment variables for all secrets
- [ ] Enable SSL/TLS for database connections
- [ ] Set up HTTPS with valid certificates
- [ ] Configure CORS properly
- [ ] Enable request logging
- [ ] Set up monitoring and alerts
- [ ] Configure automated backups
- [ ] Use managed PostgreSQL (AWS RDS, Google Cloud SQL)

---

## Stand-Out Features

### 1. Swagger API Documentation 🌟

**Implementation**:

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

**Add annotations to handlers**:

```go
// @Summary Get all items
// @Description Get all items with pagination, sorting, and filtering
// @Tags inventory
// @Accept json
// @Produce json
// @Param limit query int false "Items per page"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Item
// @Router /inventory [get]
func GetItems(c *gin.Context) { ... }
```

**Generate docs**:

```bash
swag init -g src/main.go
```

**Access**: http://localhost:8080/swagger/index.html

### 2. Redis Caching 🌟

**Use Cases**:

- Cache frequently accessed items
- Cache paginated results
- Reduce database load

**Implementation**:

```go
import "github.com/go-redis/redis/v8"

// Check cache first
cacheKey := "item:" + id
cached, err := redisClient.Get(ctx, cacheKey).Result()
if err == nil {
    // Return cached data
}

// If not in cache, fetch from DB and cache it
redisClient.Set(ctx, cacheKey, itemJSON, 5*time.Minute)
```

### 3. JWT Authentication 🌟

**Implementation**:

```go
import "github.com/golang-jwt/jwt/v5"

// Protected routes
authorized := router.Group("/")
authorized.Use(middlewares.JWTAuth())
{
    authorized.POST("/inventory", controllers.CreateItem)
    authorized.PUT("/inventory/:id", controllers.UpdateItem)
    authorized.DELETE("/inventory/:id", controllers.DeleteItem)
}

// Public routes (read-only)
router.GET("/inventory", controllers.GetItems)
router.GET("/inventory/:id", controllers.GetItemByID)
```

### 4. GraphQL API 🌟

**Alternative to REST**:

- Single endpoint
- Client-specified queries
- Reduced over-fetching

### 5. Audit Logging 🌟

Track all changes:

- Who made the change
- What was changed
- When it was changed

### 6. Export to CSV/Excel 🌟

```go
router.GET("/inventory/export", controllers.ExportToCSV)
```

### 7. Bulk Operations 🌟

```go
router.POST("/inventory/bulk", controllers.BulkCreateItems)
router.DELETE("/inventory/bulk", controllers.BulkDeleteItems)
```

---

## Timeline Estimate

| Phase                           | Tasks                                   | Estimated Time |
| ------------------------------- | --------------------------------------- | -------------- |
| 1. Environment Setup            | PostgreSQL, Go, dependencies            | 1-2 hours      |
| 2. Database & Models            | Struct, connection, migrations, seeding | 2-3 hours      |
| 3. CRUD Endpoints               | 5 handlers, routes, error handling      | 4-6 hours      |
| 4. Rate Limiting                | Middleware implementation               | 1-2 hours      |
| 5. Pagination/Sorting/Filtering | Query parameter handling                | 3-4 hours      |
| 6. Logging & Error Handling     | Middleware, structured errors           | 1-2 hours      |
| 7. Testing                      | Unit tests, integration tests           | 3-4 hours      |
| 8. Docker & Deployment          | Dockerfile, docker-compose              | 2-3 hours      |
| 9. Documentation                | README, API docs                        | 2-3 hours      |

**Total**: 19-29 hours (2-4 days of focused work)

---

## Success Metrics

### Functional Requirements

- ✅ All 5 CRUD endpoints working
- ✅ PostgreSQL database connected
- ✅ GORM handling all DB operations
- ✅ UUID primary keys
- ✅ Database seeded with sample data
- ✅ Rate limiting active (1 req/sec, burst 5)
- ✅ Pagination working (limit/offset)
- ✅ Sorting by name/stock/price
- ✅ Filtering by name/min_stock

### Non-Functional Requirements

- ✅ Proper HTTP status codes (200, 201, 204, 400, 404, 429, 500)
- ✅ JSON responses for all endpoints
- ✅ Error messages in JSON format
- ✅ Code organized in logical packages
- ✅ README with setup/usage instructions
- ✅ Simple `go run` to start server

### Code Quality

- ✅ No syntax errors or warnings
- ✅ Consistent naming conventions
- ✅ Comments on exported functions
- ✅ Error handling on all DB operations
- ✅ Input validation

---

## Resources

### Documentation

- [Gin Web Framework](https://gin-gonic.com/docs/)
- [GORM Guide](https://gorm.io/docs/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Go Standard Library](https://pkg.go.dev/std)

### Tutorials

- [Building REST APIs with Gin](https://blog.logrocket.com/building-rest-api-go-gin/)
- [GORM CRUD Operations](https://gorm.io/docs/create.html)
- [Rate Limiting in Go](https://www.alexedwards.net/blog/how-to-rate-limit-http-requests)

### Tools

- [Postman](https://www.postman.com/) - API testing
- [TablePlus](https://tableplus.com/) - Database GUI
- [Docker Desktop](https://www.docker.com/products/docker-desktop) - Containerization

---

## Conclusion

This project plan provides a comprehensive roadmap for building a production-ready inventory management REST API. By following the phased approach and addressing the concerns outlined, you'll create a scalable, maintainable, and well-documented application that meets all rubric requirements and incorporates industry best practices.

**Next Steps**:

1. Review this plan and adjust timeline based on your availability
2. Set up development environment (Phase 1)
3. Follow phases sequentially, testing each component before moving forward
4. Document issues and solutions as you encounter them
5. Iterate and refine based on testing feedback

Good luck with your implementation! 🚀

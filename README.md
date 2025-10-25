# inventory-service

## Tech Stack

<p align="left">
  <img src="https://skillicons.dev/icons?i=go" alt="Go" height="50" />
  <img src="https://skillicons.dev/icons?i=postgres" alt="PostgreSQL" height="50" />
  <img src="https://skillicons.dev/icons?i=docker" alt="Docker" height="50" />
</p>

A RESTful API backend for inventory management built with Go, Gin, PostgreSQL, and GORM.

## Features

- ✅ **CRUD Operations**: Create, Read, Update, Delete inventory items
- ✅ **RESTful API**: Clean and intuitive endpoints
- ✅ **PostgreSQL Database**: Reliable data persistence with GORM ORM
- ✅ **UUID Primary Keys**: Unique identifiers for all items
- ✅ **Rate Limiting**: Prevent API abuse (1 req/sec, burst of 5)
- ✅ **Pagination**: Handle large datasets efficiently
- ✅ **Sorting & Filtering**: Sort by name/stock/price, filter by criteria
- ✅ **Docker Support**: Containerized deployment

## Project Structure

```
inventory-service/
├── src/
│   ├── main.go                 # Application entry point
│   ├── models/
│   │   └── item.go             # Item data model
│   ├── controllers/
│   │   └── item_controller.go  # Request handlers
│   ├── routes/
│   │   └── router.go           # Route definitions
│   ├── middlewares/
│   │   ├── logger.go           # Logging middleware
│   │   └── rate_limiter.go     # Rate limiting (TODO)
│   ├── utils/
│   │   └── database.go         # Database connection
│   └── seeds/
│       └── seeder.go           # Initial data seeding
├── Dockerfile
├── docker-compose.yaml
└── README.md
```

## Prerequisites

- **Go**: 1.22 or higher ([Download](https://go.dev/dl/))
- **PostgreSQL**: 12 or higher ([Download](https://www.postgresql.org/download/))
- **Docker** (optional): For containerized deployment ([Download](https://www.docker.com/))

## Installation

### Option 1: Local Development (Manual Setup)

#### 1. Clone the Repository

```bash
git clone <repository-url>
cd inventory-service
```

#### 2. Set Up PostgreSQL Database

**Using Docker** (Recommended):

```bash
docker run -d --name inventory-postgres \
  -e POSTGRES_DB=inventory_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15-alpine
```

**Using Local PostgreSQL**:

```sql
CREATE DATABASE inventory_db;
```

#### 3. Configure Environment Variables

Create a `.env` file in the project root:

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable
```

Or export directly:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable"
```

#### 4. Install Go Dependencies

```bash
go mod download
```

#### 5. Run the Application

```bash
go run src/main.go
```

The server will start on `http://localhost:8080`

### Option 2: Docker Compose (Recommended for Quick Start)

```bash
docker-compose up --build
```

This will automatically:

- Set up PostgreSQL database
- Build and run the Go application
- Expose the API on `http://localhost:8080`

## Database Connection String Format

The `DATABASE_URL` environment variable should follow this format:

```
postgres://[user]:[password]@[host]:[port]/[database]?sslmode=[mode]
```

**Example**:

```
postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable
```

**Parameters**:

- `user`: PostgreSQL username (default: `postgres`)
- `password`: PostgreSQL password
- `host`: Database host (default: `localhost`)
- `port`: PostgreSQL port (default: `5432`)
- `database`: Database name (default: `inventory_db`)
- `sslmode`: SSL mode (`disable` for local, `require` for production)

## API Endpoints

### Base URL

```
http://localhost:8080
```

### Inventory Operations

| Method | Endpoint         | Description                                             | Status Codes       |
| ------ | ---------------- | ------------------------------------------------------- | ------------------ |
| GET    | `/inventory`     | Get all items (supports pagination, sorting, filtering) | 200, 500           |
| GET    | `/inventory/:id` | Get a single item by ID                                 | 200, 404, 500      |
| POST   | `/inventory`     | Create a new item                                       | 201, 400, 500      |
| PUT    | `/inventory/:id` | Update an existing item                                 | 200, 400, 404, 500 |
| DELETE | `/inventory/:id` | Delete an item                                          | 204, 404, 500      |

### Query Parameters for GET /inventory

- `limit`: Number of items per page (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)
- `sort_by`: Field to sort by (`name`, `stock`, `price`, `created_at`)
- `order`: Sort order (`asc` or `desc`)
- `name`: Filter by item name (partial match, case-insensitive)
- `min_stock`: Filter by minimum stock level

## Usage Examples

### Using cURL

#### 1. Get All Items

```bash
curl http://localhost:8080/inventory
```

#### 2. Get All Items with Pagination

```bash
curl "http://localhost:8080/inventory?limit=10&offset=0"
```

#### 3. Get All Items with Sorting

```bash
curl "http://localhost:8080/inventory?sort_by=price&order=desc"
```

#### 4. Filter Items by Name

```bash
curl "http://localhost:8080/inventory?name=laptop"
```

#### 5. Filter Items by Minimum Stock

```bash
curl "http://localhost:8080/inventory?min_stock=10"
```

#### 6. Get a Single Item

```bash
curl http://localhost:8080/inventory/{item-id}
```

#### 7. Create a New Item

```bash
curl -X POST http://localhost:8080/inventory \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "stock": 10,
    "price": 999.99
  }'
```

#### 8. Update an Item

```bash
curl -X PUT http://localhost:8080/inventory/{item-id} \
  -H "Content-Type: application/json" \
  -d '{
    "stock": 15,
    "price": 899.99
  }'
```

#### 9. Delete an Item

```bash
curl -X DELETE http://localhost:8080/inventory/{item-id}
```

### Using Postman

1. Import the API endpoints into Postman
2. Set the base URL to `http://localhost:8080`
3. For POST/PUT requests, set `Content-Type: application/json` header
4. Add request body in JSON format

## Request/Response Examples

### Create Item Request

```json
POST /inventory
{
  "name": "Wireless Mouse",
  "stock": 50,
  "price": 29.99
}
```

### Create Item Response (201 Created)

```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "name": "Wireless Mouse",
  "stock": 50,
  "price": 29.99,
  "created_at": "2025-10-24T10:30:00Z",
  "updated_at": "2025-10-24T10:30:00Z"
}
```

### Get All Items Response (200 OK)

```json
[
  {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "name": "Laptop",
    "stock": 10,
    "price": 999.99,
    "created_at": "2025-10-24T09:00:00Z",
    "updated_at": "2025-10-24T09:00:00Z"
  },
  {
    "id": "b2c3d4e5-f6g7-8901-bcde-f12345678901",
    "name": "Smartphone",
    "stock": 25,
    "price": 699.99,
    "created_at": "2025-10-24T09:01:00Z",
    "updated_at": "2025-10-24T09:01:00Z"
  }
]
```

### Error Response (404 Not Found)

```json
{
  "error": "item not found"
}
```

## Database Schema

### Items Table

| Column     | Type         | Constraints |
| ---------- | ------------ | ----------- |
| id         | UUID         | PRIMARY KEY |
| name       | VARCHAR(255) | NOT NULL    |
| stock      | INTEGER      | NOT NULL    |
| price      | NUMERIC      | NOT NULL    |
| created_at | TIMESTAMP    | AUTO        |
| updated_at | TIMESTAMP    | AUTO        |

## Initial Data Seeding

The application automatically seeds the database with sample inventory items on first run:

- Laptop (Stock: 10, Price: $999.99)
- Smartphone (Stock: 25, Price: $699.99)
- Headphones (Stock: 15, Price: $199.99)
- Keyboard (Stock: 30, Price: $89.99)
- Monitor (Stock: 12, Price: $299.99)

To re-seed the database, drop all items and restart the application.

## Development

### Project Dependencies

```bash
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/google/uuid
```

### Running Tests

```bash
go test ./... -v
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./src
```

## Troubleshooting

### Common Issues

**1. Database Connection Failed**

- Ensure PostgreSQL is running: `docker ps` or `pg_isready`
- Verify `DATABASE_URL` environment variable is set correctly
- Check database credentials and network connectivity

**2. Port 8080 Already in Use**

- Change port in `main.go`: `srv := &http.Server{Addr: ":3000", Handler: router}`
- Or kill the process using port 8080

**3. Module Download Errors**

- Run `go mod tidy` to clean up dependencies
- Clear module cache: `go clean -modcache`

**4. Rate Limit Errors (429 Too Many Requests)**

- Wait 1 second between requests
- Current limit: 1 request/second with burst of 5

## Environment Variables

| Variable     | Required | Default | Description                  |
| ------------ | -------- | ------- | ---------------------------- |
| DATABASE_URL | Yes      | -       | PostgreSQL connection string |
| PORT         | No       | 8080    | Server port                  |

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on GitHub.

---

**Built with ❤️ using Go, Gin, PostgreSQL, and GORM**

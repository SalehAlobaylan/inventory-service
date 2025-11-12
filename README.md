# inventory-service

## Tech Stack

<p align="left">
  <img src="https://skillicons.dev/icons?i=go" alt="Go" height="50" />
  <img src="https://skillicons.dev/icons?i=postgres" alt="PostgreSQL" height="50" />
  <img src="https://skillicons.dev/icons?i=docker" alt="Docker" height="50" />
</p>

## Key points and Features

- ✅ **CRUD Operations**: Create, Read, Update, Delete inventory items
- ✅ **RESTful API**: Clean and intuitive endpoints
- ✅ **PostgreSQL Database**: Reliable data persistence with GORM ORM
- ✅ **UUID Primary Keys**: Unique identifiers for all items
- ✅ **Rate Limiting**: Prevent API abuse (1 req/sec, burst of 5)
- ✅ **Pagination**: Handle large datasets efficiently
- ✅ **Sorting & Filtering**: Sort by name/stock/price, filter by criteria

### Docker (recommended)

```bash
docker-compose up --build
```

Runs PostgreSQL, Redis, and the API on `http://localhost:8080`.

### Manual run

1. Create `.env` in the project root:
   ```
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/inventory_db?sslmode=disable
   REDIS_URL=redis://localhost:6379/0
   ```
2. Start services (PostgreSQL + Redis)
3. Run the API:
   ```bash
   go run ./src/main.go
   ```

## API summary

Base URL: `http://localhost:8080`

| Method | Path             | Notes                                     |
| ------ | ---------------- | ----------------------------------------- |
| GET    | `/inventory`     | List items, supports filters + pagination |
| GET    | `/inventory/:id` | Fetch single item                         |
| POST   | `/inventory`     | Create new item                           |
| PUT    | `/inventory/:id` | Partial update                            |
| DELETE | `/inventory/:id` | Remove item                               |

Query params for `GET /inventory`: `limit`, `offset`, `sort_by`, `order`, `name`, `min_stock`.

## Ready-to-use cURL calls

- List items

  ```bash
  curl "http://localhost:8080/inventory?limit=10&offset=0&sort_by=created_at&order=desc" \
    -H "Accept: application/json"
  ```

- Get item (replace `{id}`)

  ```bash
  curl "http://localhost:8080/inventory/{id}" \
    -H "Accept: application/json"
  ```

- Create item

  ```bash
  curl -X POST "http://localhost:8080/inventory" \
    -H "Content-Type: application/json" \
    -d '{"name":"Wireless Mouse","stock":25,"price":29.99}'
  ```

- Update item

  ```bash
  curl -X PUT "http://localhost:8080/inventory/{id}" \
    -H "Content-Type: application/json" \
    -d '{"stock":35,"price":24.99}'
  ```

- Delete item

  ```bash
  curl -X DELETE "http://localhost:8080/inventory/{id}"
  ```

- Download Swagger spec
  ```bash
  curl "http://localhost:8080/swagger/doc.json" -o swagger.json
  ```

## Docs & tooling

- Swagger UI: `http://localhost:8080/swagger/index.html`
- Regenerate docs after handler changes: `swag init -g src/main.go -o docs`

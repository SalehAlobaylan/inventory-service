FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-service ./src

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/inventory-service ./inventory-service

EXPOSE 8080

ENTRYPOINT ["./inventory-service"]

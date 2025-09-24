package utils

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

// ConnectDatabase establishes a PostgreSQL connection using the provided DSN
// and returns a singleton GORM database instance.
func ConnectDatabase() *gorm.DB {
	if dbInstance != nil {
		return dbInstance
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	dbInstance = connection
	return dbInstance
}

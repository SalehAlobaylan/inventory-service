package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Item represents a single inventory item stored in the database.
type Item struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Stock     int       `json:"stock" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate ensures that each item receives a unique identifier when inserted.
func (item *Item) BeforeCreate(tx *gorm.DB) error {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	return nil
}

package seeds

import (
	"log"

	"gorm.io/gorm"

	"inventory-service/src/models"
)

// just for testing purposes to make sure the database contains some data
func SeedDatabase(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Item{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Database already seeded, skipping initial data load")
		return nil
	}

	items := []models.Item{
		{Name: "Laptop", Stock: 10, Price: 999.99},
		{Name: "Smartphone", Stock: 25, Price: 699.99},
		{Name: "Headphones", Stock: 15, Price: 199.99},
		{Name: "Keyboard", Stock: 30, Price: 89.99},
		{Name: "Monitor", Stock: 12, Price: 299.99},
	}

	if err := db.Create(&items).Error; err != nil {
		return err
	}

	log.Printf("Seeded database with %d inventory items", len(items))
	return nil
}

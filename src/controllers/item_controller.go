package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"inventory-service/src/models"
	"inventory-service/src/utils"
)

// GetItems handles GET /inventory requests and returns all inventory items with pagination.
func GetItems(c *gin.Context) {
	db := utils.ConnectDatabase()

	// Extract pagination parameters
	paginationParams := utils.ExtractPaginationParams(c)

	// Build query with filters and sorting using QueryBuilder
	query := utils.NewQueryBuilder(db.Model(&models.Item{}), c).
		ApplySorting(utils.SortConfig{
			AllowedFields: map[string]bool{
				"name":       true,
				"stock":      true,
				"price":      true,
				"created_at": true,
			},
			DefaultField: "created_at",
			DefaultOrder: "desc",
		}).
		ApplyFilters([]utils.FilterConfig{
			{QueryParam: "name", DBColumn: "name", Operator: "ILIKE"},
			{QueryParam: "min_stock", DBColumn: "stock", Operator: ">="},
		}).
		Build()

	// Apply pagination
	result, err := utils.PaginateWithQuery[models.Item](query, paginationParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetItemByID handles GET /inventory/:id requests and returns the matching item.
func GetItemByID(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	db := utils.ConnectDatabase()
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateItem handles POST /inventory requests to add a new inventory item.
func CreateItem(c *gin.Context) {
	var input models.Item
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.ConnectDatabase()
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateItem handles PUT /inventory/:id requests to modify an existing inventory item.
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	db := utils.ConnectDatabase()
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	var payload struct {
		Name  *string  `json:"name"`
		Stock *int     `json:"stock"`
		Price *float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Name != nil {
		item.Name = *payload.Name
	}
	if payload.Stock != nil {
		item.Stock = *payload.Stock
	}
	if payload.Price != nil {
		item.Price = *payload.Price
	}

	if err := db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeleteItem handles DELETE /inventory/:id requests to remove an item from inventory.
func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	db := utils.ConnectDatabase()

	var item models.Item
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	if err := db.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

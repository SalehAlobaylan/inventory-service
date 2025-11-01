package controllers

import (
	"net/http"
	"strconv"
	"fmt"

	"github.com/gin-gonic/gin"

	"inventory-service/src/models"
	"inventory-service/src/utils"
)

// GetItems handles GET /inventory requests and returns all inventory items.
func GetItems(c *gin.Context) {
	var items []models.Item
	

	db := utils.ConnectDatabase()
	
	// for pagination
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	// for sorting
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	// for filtering
	// name := c.Query("name")
	// minStock := c.Query("min_stock")
	
	
	limitInt ,_ := strconv.Atoi(limit)
	offsetInt ,_ := strconv.Atoi(offset)
	// sortByInt ,_ := strconv.Atoi(sortBy)
	// orderInt ,_ := strconv.Atoi(order)

	if order != "asc" && order != "desc" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameter"})
        return
    }


	if err := db.Limit(limitInt).Offset(offsetInt).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	orderClause := fmt.Sprintf("%s %s", sortBy, order)

	if err := db.Order(orderClause).Find(&items).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
        return
    }

	c.JSON(http.StatusOK, items)
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

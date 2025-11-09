package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inventory-service/src/models"
	"inventory-service/src/utils"
)

// CreateItemRequest defines the payload required to create a new inventory item.
type CreateItemRequest struct {
	Name  string  `json:"name" binding:"required" example:"Laptop"`
	Stock int     `json:"stock" binding:"required" example:"10"`
	Price float64 `json:"price" binding:"required" example:"999.99"`
}

// UpdateItemRequest defines the fields that can be updated on an inventory item.
type UpdateItemRequest struct {
	Name  *string  `json:"name" example:"Laptop Pro"`
	Stock *int     `json:"stock" example:"15"`
	Price *float64 `json:"price" example:"849.99"`
}

// GetItems handles GET /inventory requests and returns all inventory items.
// @Summary List inventory items
// @Description Retrieve inventory items with optional filtering, sorting, and pagination.
// @Tags inventory
// @Accept json
// @Produce json
// @Param name query string false "Filter by item name (case-insensitive)"
// @Param min_stock query int false "Minimum stock filter"
// @Param limit query int false "Items per page (default 10, max 100)"
// @Param offset query int false "Offset for pagination"
// @Param sort_by query string false "Sort field (name|stock|price|created_at)"
// @Param order query string false "Sort order (asc|desc)"
// @Success 200 {array} models.Item
// @Failure 500 {object} map[string]string
// @Router /inventory [get]
func GetItems(c *gin.Context) {
	var items []models.Item

	db := utils.ConnectDatabase()
	query := db.Model(&models.Item{})

	// Filters
	if name := c.Query("name"); name != "" {
		// Case-insensitive match (PostgreSQL)
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if minStockStr := c.Query("min_stock"); minStockStr != "" {
		if minStock, err := strconv.Atoi(minStockStr); err == nil {
			query = query.Where("stock >= ?", minStock)
		}
	}

	// Sorting (whitelist fields) to prevent SQL injection
	allowedFields := map[string]bool{
		"name":       true,
		"stock":      true,
		"price":      true,
		"created_at": true,
	}
	sortBy := c.DefaultQuery("sort_by", "created_at")
	if !allowedFields[sortBy] {
		sortBy = "created_at"
	}
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, order)
	query = query.Order(orderClause)

	// Pagination (with sane bounds)
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	if err := query.Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetItemByID handles GET /inventory/:id requests and returns the matching item.
// @Summary Get an inventory item
// @Description Retrieve a single inventory item by its identifier.
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} models.Item
// @Failure 404 {object} map[string]string
// @Router /inventory/{id} [get]
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
// @Summary Create a new inventory item
// @Description Create a new inventory item by providing its core attributes.
// @Tags inventory
// @Accept json
// @Produce json
// @Param item body CreateItemRequest true "Item to create"
// @Success 201 {object} models.Item
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory [post]
func CreateItem(c *gin.Context) {
	var input CreateItemRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.ConnectDatabase()
	item := models.Item{
		Name:  input.Name,
		Stock: input.Stock,
		Price: input.Price,
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// UpdateItem handles PUT /inventory/:id requests to modify an existing inventory item.
// @Summary Update an inventory item
// @Description Update the mutable fields of an existing inventory item.
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body UpdateItemRequest true "Fields to update"
// @Success 200 {object} models.Item
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory/{id} [put]
func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item

	db := utils.ConnectDatabase()
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	var payload UpdateItemRequest
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
// @Summary Delete an inventory item
// @Description Remove an inventory item by its identifier.
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventory/{id} [delete]
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

package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SortConfig defines allowed sorting fields and their mappings
type SortConfig struct {
	AllowedFields map[string]bool
	DefaultField  string
	DefaultOrder  string
}

// FilterConfig defines a single filter configuration
type FilterConfig struct {
	QueryParam string
	DBColumn   string
	Operator   string // "=", "LIKE", ">=", "<=", etc.
}

// QueryBuilder helps build complex queries with sorting and filtering
type QueryBuilder struct {
	db      *gorm.DB
	context *gin.Context
}

// NewQueryBuilder creates a new QueryBuilder instance
func NewQueryBuilder(db *gorm.DB, c *gin.Context) *QueryBuilder {
	return &QueryBuilder{
		db:      db,
		context: c,
	}
}

// ApplySorting applies sorting based on query parameters
func (qb *QueryBuilder) ApplySorting(config SortConfig) *QueryBuilder {
	sortBy := qb.context.DefaultQuery("sort_by", config.DefaultField)
	order := qb.context.DefaultQuery("order", config.DefaultOrder)

	// Validate sort field
	if !config.AllowedFields[sortBy] {
		sortBy = config.DefaultField
	}

	// Validate order
	if order != "asc" && order != "desc" {
		order = config.DefaultOrder
	}

	qb.db = qb.db.Order(sortBy + " " + order)
	return qb
}

// ApplyFilter applies a single filter
func (qb *QueryBuilder) ApplyFilter(config FilterConfig) *QueryBuilder {
	value := qb.context.Query(config.QueryParam)
	if value == "" {
		return qb
	}

	switch config.Operator {
	case "LIKE":
		qb.db = qb.db.Where(config.DBColumn+" LIKE ?", "%"+value+"%")
	case "ILIKE":
		qb.db = qb.db.Where("LOWER("+config.DBColumn+") LIKE ?", "%"+value+"%")
	case ">=":
		qb.db = qb.db.Where(config.DBColumn+" >= ?", value)
	case "<=":
		qb.db = qb.db.Where(config.DBColumn+" <= ?", value)
	case ">":
		qb.db = qb.db.Where(config.DBColumn+" > ?", value)
	case "<":
		qb.db = qb.db.Where(config.DBColumn+" < ?", value)
	default: // "="
		qb.db = qb.db.Where(config.DBColumn+" = ?", value)
	}

	return qb
}

// ApplyFilters applies multiple filters
func (qb *QueryBuilder) ApplyFilters(configs []FilterConfig) *QueryBuilder {
	for _, config := range configs {
		qb.ApplyFilter(config)
	}
	return qb
}

// GetQuery returns the built query
func (qb *QueryBuilder) GetQuery() *gorm.DB {
	return qb.db
}

// Build returns the built query (alias for GetQuery)
func (qb *QueryBuilder) Build() *gorm.DB {
	return qb.GetQuery()
}

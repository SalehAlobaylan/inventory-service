package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams holds pagination query parameters
type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Offset   int `json:"-"`
}

// PaginationMeta contains metadata about the pagination
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

// PaginatedResponse is a generic response structure for paginated data
type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// ExtractPaginationParams extracts and validates pagination parameters from request
func ExtractPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Validate and set bounds
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// Paginate applies pagination to a GORM query and returns paginated results
func Paginate[T any](db *gorm.DB, params PaginationParams) (PaginatedResponse[T], error) {
	var data []T
	var total int64

	// Get total count
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	// Get paginated data
	if err := db.Offset(params.Offset).Limit(params.PageSize).Find(&data).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	return PaginatedResponse[T]{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage:  params.Page,
			PageSize:     params.PageSize,
			TotalPages:   totalPages,
			TotalRecords: total,
			HasNext:      hasNext,
			HasPrev:      hasPrev,
		},
	}, nil
}

// PaginateWithQuery applies pagination to a pre-filtered GORM query
func PaginateWithQuery[T any](query *gorm.DB, params PaginationParams) (PaginatedResponse[T], error) {
	var data []T
	var total int64

	// Get total count from the filtered query
	countQuery := query.Session(&gorm.Session{})
	if err := countQuery.Model(new(T)).Count(&total).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	// Get paginated data
	if err := query.Offset(params.Offset).Limit(params.PageSize).Find(&data).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	return PaginatedResponse[T]{
		Data: data,
		Pagination: PaginationMeta{
			CurrentPage:  params.Page,
			PageSize:     params.PageSize,
			TotalPages:   totalPages,
			TotalRecords: total,
			HasNext:      hasNext,
			HasPrev:      hasPrev,
		},
	}, nil
}

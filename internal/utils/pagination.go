package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination holds pagination and sorting info. You can add filter fields if needed.
type Pagination struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`   // e.g., "date asc" or "amount desc"
	Search string `json:"search"` // optional: for searching across title/description, etc.
}

// NewPaginationFromRequest parses query params from Gin context
// and returns a Pagination struct with defaults if not provided.
func NewPaginationFromRequest(c *gin.Context) Pagination {
	// Default values
	p := Pagination{
		Page:  1,
		Limit: 10, // change to your preference
		Sort:  "created_at desc",
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			p.Page = page
		}
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			p.Limit = limit
		}
	}

	// Parse sort (e.g. sort=title asc, or sort=amount desc)
	if sortStr := c.Query("sort"); sortStr != "" {
		// You could validate that it matches <column> <asc/desc>, etc.
		p.Sort = sortStr
	}

	// Optional: parse a "search" param for global text search
	if searchStr := c.Query("search"); searchStr != "" {
		p.Search = searchStr
	}

	return p
}

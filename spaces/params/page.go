package params

import (
	"github.com/fox-gonic/fox/database"
	"github.com/miclle/space/models"
)

// CreatePage create page params
type CreatePage struct {
	SpaceID         int64
	CreatorID       int64
	ParentID        int64
	PageID          int64
	Lang            string
	Version         string
	Status          models.PageStatus
	Title           string
	ShortTitle      string
	Body            string
	MetaKeywords    []string
	MetaDescription string
}

// DescribePages describe page detail params
type DescribePages struct {
	SpaceID int64
	Lang    string
	Version string
	Depth   string
}

// DescribePage describe page detail params
type DescribePage struct {
	SpaceID int64
	PageID  int64
	Lang    string
	Version string
}

// UpdatePage update page params
type UpdatePage struct {
	ID              int64
	Lang            *string
	Version         *string
	Status          *models.PageStatus
	Title           *string
	ShortTitle      *string
	Body            *string
	MetaKeywords    *[]string
	MetaDescription *string
}

// Search page params
type Search struct {
	database.Pagination[*models.Page]
	Lang string
	Q    string
}

package params

import "github.com/miclle/space/models"

// CreatePage create page params
type CreatePage struct {
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
	Lang    string
	Version string
	Depth   string
}

// DescribePage describe page detail params
type DescribePage struct {
	ID      string
	Lang    string
	Version string
}

// UpdatePage update page params
type UpdatePage struct {
	Lang            *string
	Version         *string
	Status          *models.PageStatus
	Title           *string
	ShortTitle      *string
	Body            *string
	MetaKeywords    *[]string
	MetaDescription *string
}

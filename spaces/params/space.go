package params

import (
	"github.com/fox-gonic/fox/database"

	"github.com/miclle/space/models"
)

// CreateSpace create space params
type CreateSpace struct {
	Name         string
	Key          string
	Lang         string
	FallbackLang string
	Description  string
	Avatar       string
	Status       models.SpaceStatus
}

// DescribeSpaces describe spaces params
type DescribeSpaces struct {
	database.Pagination[*models.Space]
	Q string
}

// DescribeSpace describe space detail params
type DescribeSpace struct {
	Key string
}

// UpdateSpace update space params
type UpdateSpace struct {
	Key          string
	Name         *string
	Lang         *string
	FallbackLang *string
	HomepageID   *int64
	Description  *string
	Avatar       *string
	Status       models.SpaceStatus
}

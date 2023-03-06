package spaces

import (
	"context"

	"github.com/fox-gonic/fox/database"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
)

// Service for spaces interface
type Service interface {
	CreateSpace(context.Context, *params.CreateSpace) (*models.Space, error)
	DescribeSpaces(context.Context, *params.DescribeSpaces) (*database.Pagination[*models.Space], error)
	DescribeSpace(context.Context, *params.DescribeSpace) (*models.Space, error)
	UpdateSpace(context.Context, *params.UpdateSpace) (*models.Space, error)

	CreatePage(context.Context, *params.CreatePage) (*models.Page, error)
	DescribePages(context.Context, *params.DescribePages) ([]*models.Page, error)
	DescribePage(context.Context, *params.DescribePage) (*models.Page, error)
	UpdatePage(context.Context, *params.UpdatePage) (*models.Page, error)
}

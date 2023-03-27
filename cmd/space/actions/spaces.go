package actions

import (
	"errors"

	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/engine"
	"gorm.io/gorm"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
)

// -----------------------------------------------------------------------------

// CreateSpaceArgs create space args
type CreateSpaceArgs struct {
	Name         string             `json:"name"`
	Key          string             `json:"key"`
	Lang         string             `json:"lang"`
	FallbackLang string             `json:"fallback_lang"`
	Description  string             `json:"description"`
	Avatar       string             `json:"avatar"`
	Status       models.SpaceStatus `json:"status"`
}

// CreateSpace create space
// POST /api/spaces
func (actions *Actions) CreateSpace(c *engine.Context, args *CreateSpaceArgs) (*models.Space, error) {

	var account = c.MustGet("account").(*models.Account)

	var params = &params.CreateSpace{
		Name:         args.Name,
		Key:          args.Key,
		Lang:         args.Lang,
		FallbackLang: args.FallbackLang,
		Description:  args.Description,
		Avatar:       args.Avatar,
		Status:       args.Status,
		CreatorID:    account.ID,
	}

	return actions.Spacer.CreateSpace(c, params)
}

// DescribeSpacesArgs describe spaces args
type DescribeSpacesArgs struct {
	database.Pagination[*models.Space]
	Q string `query:"q"`
}

// DescribeSpaces describe spaces
// GET /api/spaces
func (actions *Actions) DescribeSpaces(c *engine.Context, args *DescribeSpacesArgs) (*database.Pagination[*models.Space], error) {

	var params = &params.DescribeSpaces{
		Pagination: args.Pagination,
		Q:          args.Q,
	}

	return actions.Spacer.DescribeSpaces(c, params)
}

// -----------------------------------------------------------------------------

// DescribeSpaceArgs describe space detail args
type DescribeSpaceArgs struct {
	Key string `uri:"key"`
}

// DescribeSpace describe space detail
// GET /api/spaces/:key
func (actions *Actions) DescribeSpace(c *engine.Context, args *DescribeSpaceArgs) (*models.Space, error) {

	var space = c.MustGet("space").(*models.Space)

	return space, nil
}

// -----------------------------------------------------------------------------

// UpdateSpaceArgs update space args
type UpdateSpaceArgs struct {
	Key          string             `uri:"key"`
	Name         *string            `json:"name"`
	Lang         *string            `json:"lang"`
	FallbackLang *string            `json:"fallback_lang"`
	HomepageID   *int64             `json:"homepage_id"`
	Description  *string            `json:"description"`
	Avatar       *string            `json:"avatar"`
	Status       models.SpaceStatus `json:"status"`
}

// UpdateSpace update space
// PATCH /api/spaces/:key
func (actions *Actions) UpdateSpace(c *engine.Context, args *UpdateSpaceArgs) (*models.Space, error) {

	var params = &params.UpdateSpace{
		Key:          args.Key,
		Name:         args.Name,
		Lang:         args.Lang,
		FallbackLang: args.FallbackLang,
		HomepageID:   args.HomepageID,
		Description:  args.Description,
		Avatar:       args.Avatar,
		Status:       args.Status,
	}

	return actions.Spacer.UpdateSpace(c, params)
}

// -----------------------------------------------------------------------------

// SetSpaceArgs describe space detail args
type SetSpaceArgs struct {
	Key string `uri:"key"`
}

// SetSpace describe space detail
// MATCH route `/api/spaces/:key`
func (actions *Actions) SetSpace(c *engine.Context, args *SetSpaceArgs) error {

	space, err := actions.Spacer.DescribeSpace(c, &params.DescribeSpace{
		Key: args.Key,
	})

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// TODO(m) custom errors.New
		return gorm.ErrRecordNotFound
	}

	if err != nil {
		c.Logger.Error("find space failed, err: %+v", err)
		return err
	}

	c.Set("space", space)

	return nil
}

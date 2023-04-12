package actions

import (
	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/engine"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
)

// ----------------------------------------------------------------------------

// CreatePageArgs create page args
type CreatePageArgs struct {
	ParentID        int64             `json:"parent_id"`
	PageID          int64             `json:"page_id"`
	Lang            string            `json:"lang"`
	Version         string            `json:"version"`
	Status          models.PageStatus `json:"status"`
	Title           string            `json:"title"`
	ShortTitle      string            `json:"short_title"`
	Body            string            `json:"body"`
	MetaKeywords    []string          `json:"meta_keywords"`
	MetaDescription string            `json:"meta_description"`
}

// CreatePage create page
// POST /api/spaces/:key/pages
func (actions *Actions) CreatePage(c *engine.Context, args *CreatePageArgs) (*models.Page, error) {

	var (
		space   = c.MustGet("space").(*models.Space)
		account = c.MustGet("account").(*models.Account)
		params  = &params.CreatePage{
			SpaceID:         space.ID,
			CreatorID:       account.ID,
			ParentID:        args.ParentID,
			Lang:            args.Lang,
			Version:         args.Version,
			Status:          args.Status,
			Title:           args.Title,
			ShortTitle:      args.ShortTitle,
			Body:            args.Body,
			MetaKeywords:    args.MetaKeywords,
			MetaDescription: args.MetaDescription,
		}
	)

	return actions.Spacer.CreatePage(c, params)
}

// ----------------------------------------------------------------------------

// DescribePagesArgs describe page detail args
type DescribePagesArgs struct {
	Lang    string `query:"lang"`
	Version string `query:"version"`
	Depth   string `query:"depth"` // all, root, default: all
}

// DescribePages describe pages
// GET /api/spaces/:key/pages
func (actions *Actions) DescribePages(c *engine.Context, args *DescribePagesArgs) (*database.Pagination[*models.Page], error) {

	var (
		space  = c.MustGet("space").(*models.Space)
		params = &params.DescribePages{
			SpaceID: space.ID,
		}
	)

	lang := args.Lang
	if lang == "" {
		params.Lang = space.Lang
	}

	if len(args.Version) > 0 {
		params.Version = args.Version
	}

	return actions.Spacer.DescribePages(c, params)
}

// ----------------------------------------------------------------------------

// DescribePageArgs describe page detail args
type DescribePageArgs struct {
	PageID  string `uri:"id"`
	Lang    string `query:"lang"`
	Version string `query:"version"`
}

// DescribePage describe page detail
// GET /api/spaces/:key/pages/:id
func (actions *Actions) DescribePage(c *engine.Context, args *DescribePageArgs) (*models.Page, error) {

	var page = c.MustGet("page").(*models.Page)

	return page, nil
}

// ----------------------------------------------------------------------------

// UpdatePageArgs update page args
type UpdatePageArgs struct {
	Lang            *string            `json:"lang"`
	Version         *string            `json:"version"`
	Status          *models.PageStatus `json:"status"`
	Title           *string            `json:"title"`
	ShortTitle      *string            `json:"short_title"`
	Body            *string            `json:"body"`
	MetaKeywords    *[]string          `json:"meta_keywords"`
	MetaDescription *string            `json:"meta_description"`
}

// UpdatePage update page
// PATCH /api/spaces/:key/pages/:id
func (actions *Actions) UpdatePage(c *engine.Context, args *UpdatePageArgs) (*models.Page, error) {
	var (
		page   = c.MustGet("page").(*models.PageContent)
		params = &params.UpdatePage{
			ID:              page.ID,
			Lang:            args.Lang,
			Version:         args.Version,
			Status:          args.Status,
			Title:           args.Title,
			ShortTitle:      args.ShortTitle,
			Body:            args.Body,
			MetaKeywords:    args.MetaKeywords,
			MetaDescription: args.MetaDescription,
		}
	)

	return actions.Spacer.UpdatePage(c, params)
}

// middleware
// ----------------------------------------------------------------------------

// SetPageArgs describe space detail args
type SetPageArgs struct {
	PageID  int64  `uri:"id"`
	Lang    string `query:"lang"    json:"-"`
	Version string `query:"version" json:"-"`
}

// SetPage describe space detail
// MATCH route `/api/spaces/:key/pages/:id`
func (actions *Actions) SetPage(c *engine.Context, args *SetPageArgs) error {
	var (
		space  = c.MustGet("space").(*models.Space)
		params = &params.DescribePage{
			SpaceID: space.ID,
			PageID:  args.PageID,
		}
	)

	lang := args.Lang
	if lang == "" {
		params.Lang = space.Lang
	}

	if len(args.Version) > 0 {
		params.Version = args.Version
	}

	page, err := actions.Spacer.DescribePage(c, params)
	if err != nil {
		c.Logger.Error("find page failed, err: %+v", err)
		return err
	}

	c.Set("page", page)

	return nil
}

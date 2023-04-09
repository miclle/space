package website

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fox-gonic/fox/engine"
	"gorm.io/gorm"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
	"github.com/miclle/space/ui"
)

// DescribePageArgs describe page detail args
type DescribePageArgs struct {
	SpaceKey string `uri:"space_key"`
	PageID   int64  `uri:"page_id"`
	Lang     string `uri:"lang"`
	Version  string `query:"version"`
}

// DescribePage describe page detail
// GET /docs/:id
func (actions *Actions) DescribePage(c *engine.Context, args *DescribePageArgs) {

	var (
		spaces = c.MustGet("spaces").([]*models.Space)
		space  = c.MustGet("space").(*models.Space)
		// pageTree = c.MustGet("pageTree").(models.PageTree)
		page *models.Page
	)

	data := ui.PageData{
		Lang:   args.Lang,
		Spaces: spaces,
		Space:  space,
		// PageTree: pageTree,
	}

	page, err := actions.Spacer.DescribePage(c, &params.DescribePage{
		SpaceID: space.ID,
		PageID:  args.PageID,
		Lang:    args.Lang,
		Version: args.Version,
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(404, "404.html", data)
			return
		}

		c.Logger.Error("get spaces failed", err)
		c.HTML(500, "500.html", data)
		return
	}

	// redirect to space homepage
	if page.ID == space.Homepage.ID {
		c.Redirect(http.StatusFound, fmt.Sprintf("/%s/docs/%s", args.Lang, space.Key))
		return
	}

	// data.Title = page.Title
	// data.Page = page

	c.HTML(200, "page.html", data)
}

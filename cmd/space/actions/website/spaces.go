package website

import (
	"errors"

	"github.com/fox-gonic/fox/engine"
	"gorm.io/gorm"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
	"github.com/miclle/space/ui"
)

// DescribeSpaceArgs describe space args
type DescribeSpaceArgs struct {
	Lang     string `uri:"lang"`
	SpaceKey string `uri:"space_key"`
}

// DescribeSpace describe docs
// GET /:lang/docs/:space
func (actions *Actions) DescribeSpace(c *engine.Context, args *DescribeSpaceArgs) {

	var (
		spaces = c.MustGet("spaces").([]*models.Space)
		space  = c.MustGet("space").(*models.Space)
		// pageTree = c.MustGet("pageTree").(models.PageTree)
	)

	data := ui.PageData{
		Lang:   args.Lang,
		Title:  space.Name,
		Spaces: spaces,
		Space:  space,
		// PageTree: pageTree,
	}

	c.HTML(200, "space.html", data)
}

// ----------------------------------------------------------------------------

// SetSpaceArgs describe space args
type SetSpaceArgs struct {
	Lang     string `uri:"lang"`
	SpaceKey string `uri:"space_key"`
	Version  string `query:"version"`
}

// SetSpace describe docs
// match route: `/:lang/docs/:space`
func (actions *Actions) SetSpace(c *engine.Context, args *SetSpaceArgs) {

	var (
		space *models.Space
		// pageTree models.PageTree
	)

	// find space
	space, err := actions.Spacer.DescribeSpace(c, &params.DescribeSpace{
		Key:     args.SpaceKey,
		Lang:    args.Lang,
		Version: args.Version,
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.HTML(404, "404.html", map[string]interface{}{
				"Lang":  args.Lang,
				"Title": "Page not found",
			})
			return
		}

		c.Logger.Error("get spaces failed", err)
		c.HTML(500, "500.html", map[string]interface{}{})
		return
	}

	// pageTree, err = actions.Spacer.DescribePageTree(c, &params.DescribePages{
	// 	SpaceID: space.ID,
	// 	Lang:    args.Lang,
	// 	Version: args.Version,
	// })

	if err != nil {
		c.Logger.Error("get space pages failed", err)
		c.HTML(500, "500.html", map[string]interface{}{})
		return
	}

	c.Set("space", space)
	// c.Set("pageTree", pageTree)
}

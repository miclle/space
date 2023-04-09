package website

import (
	"github.com/fox-gonic/fox/database"
	"github.com/fox-gonic/fox/engine"

	"github.com/miclle/space/models"
	"github.com/miclle/space/spaces/params"
)

// SearchArgs search args
type SearchArgs struct {
	database.Pagination[*models.PageContent]
	Lang string `uri:"lang"`
	Q    string `query:"q"`
}

// Search website homepage
// GET /:lang/search
func (actions *Actions) Search(c *engine.Context, args *SearchArgs) {

	var spaces = c.MustGet("spaces").([]*models.Space)

	pagination, err := actions.Spacer.Serach(c, &params.Search{
		Pagination: args.Pagination,
		Lang:       args.Lang,
		Q:          args.Q,
	})

	if err != nil {
		c.Logger.Error("get pages failed", err)
		c.HTML(500, "500.html", map[string]interface{}{})
		return
	}

	c.HTML(200, "search.html", map[string]interface{}{
		"Lang":       args.Lang,
		"Title":      "Search",
		"Spaces":     spaces,
		"Q":          args.Q,
		"Pagination": pagination,
	})
}

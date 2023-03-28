package website

import (
	"github.com/fox-gonic/fox/engine"

	"github.com/miclle/space/models"
)

// Homepage website homepage
// GET /
func (actions *Actions) Homepage(c *engine.Context) {
	var (
		lang   = c.MustGet("lang").(string)
		spaces = c.MustGet("spaces").([]*models.Space)
	)

	c.HTML(200, "homepage.html", map[string]interface{}{
		"Lang":   lang,
		"Spaces": spaces,
	})
}

// NotFound 404 page
// GET /
func (actions *Actions) NotFound(c *engine.Context) {

	c.HTML(404, "404.html", map[string]interface{}{
		"Title": "Page not found",
	})
}

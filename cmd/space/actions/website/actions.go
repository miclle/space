package website

import (
	"net/http"
	"path"

	"github.com/fox-gonic/fox/engine"
	"github.com/gin-gonic/gin/render"

	"github.com/miclle/space/accounts"
	"github.com/miclle/space/config"
	"github.com/miclle/space/spaces"
	"github.com/miclle/space/spaces/params"
)

// Actions type
type Actions struct {
	Configuration config.Configuration
	Accounter     accounts.Service
	Spacer        spaces.Service
}

// SetLangArgs set lang middleware
type SetLangArgs struct {
	Lang string `uri:"lang"`
}

// SetLang set lang middleware
func (actions *Actions) SetLang(c *engine.Context, args *SetLangArgs) (res interface{}) {

	// accept := c.GetHeader("Accept-Language")
	// c.Logger.Debug("accept", accept)

	// TODO(m): check args.Lang is valid

	if args.Lang == "" {
		return render.Redirect{
			Code:     http.StatusFound,
			Location: path.Join("/en-US", c.Request.URL.String()),
		}
	}

	c.Set("lang", args.Lang)

	return
}

// SetGlobal global middleware
func (actions *Actions) SetGlobal(c *engine.Context) {

	lang := c.MustGet("lang").(string)

	var params = &params.DescribeSpaces{
		Lang: lang,
	}

	params.PageSize = 1000

	pagination, err := actions.Spacer.DescribeSpaces(c, params)

	if err != nil {
		c.Logger.Error("get spaces failed", err)
		c.HTML(500, "500.html", map[string]interface{}{})
		return
	}

	c.Set("spaces", pagination.Items)
}

package actions

import (
	"strings"

	"github.com/fox-gonic/fox/engine"
	"github.com/fox-gonic/fox/errors"
	"github.com/gin-contrib/sessions"
	"github.com/samber/lo"

	"github.com/miclle/space/accounts/params"
)

// SessionStaffLoginKey session account context key
const SessionStaffLoginKey = "staff-login"

var skipPaths = []string{
	"/static",
	"/sso",
	"/logout",
	"/capacity",
	"/ping",
	"/forbidden",
	"/500",
	"/open",

	"/favicon.ico",
	"/robots.txt",
	"/manifest.json",
}

// AuthMiddleware 登录验证中间件
func (actions *Actions) AuthMiddleware(c *engine.Context) (res interface{}) {

	path := c.Request.URL.Path

	contain := lo.ContainsBy(skipPaths, func(prefix string) bool {
		return strings.HasPrefix(path, prefix)
	})
	if contain {
		return
	}

	var (
		session = sessions.Default(c.Context)
		login   = session.Get(SessionStaffLoginKey)
	)

	if login != nil {
		if l, ok := login.(string); ok {
			account, err := actions.Accounter.DescribeAccount(c, &params.DescribeAccount{
				Login: l,
			})

			if err != nil {
				c.Logger.Errorf("database.First() failed, err: %+v", err)
				session.Delete("account")

				if err := session.Save(); err != nil {
					c.Logger.Error("session.Save() failed, err: %+v", err)
				}
			} else {
				c.Set("account", account)
			}
		}
	}

	if strings.HasPrefix(path, "/api") {
		return errors.ErrUnauthorized
	}

	// return render.Redirect{
	// 	Code:     302,
	// 	Location: "/signin",
	// }

	return
}

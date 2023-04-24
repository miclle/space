package main

import (
	"net/http"

	"github.com/fox-gonic/fox/engine"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/miclle/space/accounts"
	"github.com/miclle/space/cmd/space/actions"
	"github.com/miclle/space/cmd/space/actions/website"
	"github.com/miclle/space/config"
	"github.com/miclle/space/spaces"
	"github.com/miclle/space/ui"
)

func router(
	configuration config.Configuration,
	accounter accounts.Service,
	spacer spaces.Service,
) http.Handler {

	// --------------------------------------------------------------------------
	store := cookie.NewStore([]byte(configuration.Secret))
	options := sessions.Options{
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}
	if engine.Mode() == engine.DebugMode {
		options.Secure = true
	}
	store.Options(options)

	// --------------------------------------------------------------------------

	router := engine.New()
	router.SetHTMLTemplate(ui.Template)

	router.Use(sessions.Sessions("SPACE", store))

	// Embed website assets
	ui.EmbedAssets(router)

	router.GET("ping", func(c *engine.Context) string {
		c.Logger.Debug("ping/pong testing")
		return "pong"
	})

	// website
	// --------------------------------------------------------------------------
	{
		website := &website.Actions{
			Configuration: configuration,
			Accounter:     accounter,
			Spacer:        spacer,
		}

		group := router.Group("", website.SetLang, website.SetGlobal)
		group.GET("/", website.Homepage)
		group.GET("/:lang", website.Homepage)

		group.GET("/search", website.Search)
		group.GET("/:lang/search", website.Search)

		{
			spaceGroup := group.Group("", website.SetSpace)

			spaceGroup.GET("/docs/:space_key", website.DescribeSpace)
			spaceGroup.GET("/:lang/docs/:space_key", website.DescribeSpace)

			spaceGroup.GET("/docs/:space_key/:page_id", website.DescribePage)
			spaceGroup.GET("/:lang/docs/:space_key/:page_id", website.DescribePage)
		}
	}

	// api
	// --------------------------------------------------------------------------
	{
		api := &actions.Actions{
			Configuration: configuration,
			Accounter:     accounter,
			Spacer:        spacer,
		}

		router.GET("/logout", api.Logout)
		router.POST("/api/accounts/signup", api.Signup)
		router.POST("/api/accounts/signin", api.Signin)

		group := router.Group("/api", api.AuthMiddleware)

		group.GET("/accounts/overview", api.Overview)
		group.GET("/accounts", api.DescribeAccounts)
		group.PATCH("/accounts/:id", api.UpdateAccount)

		group.POST("/spaces", api.CreateSpace)
		group.GET("/spaces", api.DescribeSpaces)

		space := group.Group("/spaces/:key", api.SetSpace)
		space.GET("", api.DescribeSpace)
		space.PATCH("", api.UpdateSpace)
		space.POST("/pages", api.CreatePage)
		space.GET("/pages", api.DescribePages)

		page := space.Group("/pages/:id", api.SetPage)
		page.GET("", api.DescribePage)
		page.PATCH("", api.UpdatePage)

		group.POST("/markdown/preview", api.PreviewMarkdown)
	}

	return router
}

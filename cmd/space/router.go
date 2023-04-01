package main

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/fox-gonic/fox/engine"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin/render"

	"github.com/miclle/space/accounts"
	"github.com/miclle/space/cmd/space/actions"
	"github.com/miclle/space/cmd/space/actions/website"
	"github.com/miclle/space/config"
	"github.com/miclle/space/spaces"
	"github.com/miclle/space/ui"
)

var spaRoutes = []string{
	"/signin",
	"/signup",
	"/spaces",
	"/spaces/*filepath",
	"/accounts",
	"/accounts/*filepath",
}

// embedPublicAssets embed fs from `public` dir
func embedPublicAssets(router *engine.Engine) {

	var (
		embedFS = ui.EmbedFS()
		tmpl    = template.Must(template.New("").ParseFS(embedFS, "build/*.html"))
	)

	homepage := render.HTML{
		Template: tmpl,
		Name:     "index.html",
		Data:     map[string]string{},
	}

	// handle home page
	router.GET("/", func(c *engine.Context) (res interface{}) {
		return homepage
	})

	router.NotFound(func(c *engine.Context) (res interface{}) {

		if c.Request.Method != http.MethodGet {
			return http.StatusNotFound
		}

		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			return http.StatusNotFound
		}

		filepath := "public" + c.Request.URL.Path

		file, err := embedFS.Open(filepath)
		if errors.Is(err, fs.ErrNotExist) {
			return homepage
		}

		info, err := file.Stat()
		if err != nil {
			return err
		}

		return render.Reader{
			ContentLength: info.Size(),
			Reader:        file,
		}
	})

	// handle static assets
	router.GET("/static/*filepath", func(c *engine.Context) (res interface{}) {
		c.FileFromFS("public/static/"+c.Param("filepath"), http.FS(embedFS))
		c.Abort()
		return
	})
}

// reverseProxyWebsiteApp proxy development website
func reverseProxyWebsiteApp(router *engine.Engine, website string) {

	origin, _ := url.Parse(website)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	for _, path := range spaRoutes {
		router.GET(path, func(c *engine.Context) (res interface{}) {
			c.Logger.Debug("/spaces/*filepath", c.Request.URL)
			proxy.ServeHTTP(c.Writer, c.Request)
			return
		})
	}

	router.NotFound(func(c *engine.Context) {
		c.Logger.Debug("NotFound", c.Request.URL)
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

// ----------------------------------------------------------------------------

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
	// !IMPORTANT:
	if engine.Mode() == engine.DebugMode {
		reverseProxyWebsiteApp(router, os.Getenv("PROXY_WEBSITE"))
		router.StaticFS("/static", ui.StaticFS("public/static/"))
	} else {
		embedPublicAssets(router)
	}

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

		group := router.Group("/api", api.AuthMiddleware)

		group.GET("/accounts/overview", api.Overview)
		group.POST("/accounts/signup", api.Signup)
		group.POST("/accounts/signin", api.Signin)

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

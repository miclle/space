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
	"github.com/miclle/space/config"
	"github.com/miclle/space/spaces"
	"github.com/miclle/space/ui"
)

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

	router.GET("/", func(c *engine.Context) (res interface{}) {
		proxy.ServeHTTP(c.Writer, c.Request)
		return
	})

	router.NotFound(func(c *engine.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}

// ----------------------------------------------------------------------------

func router(
	configuration config.Configuration,
	accounter accounts.Service,
	spacer spaces.Service,
) http.Handler {

	actions := &actions.Actions{
		Configuration: configuration,
		Accounter:     accounter,
		Spacer:        spacer,
	}

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
	router.Use(sessions.Sessions("SPACE", store))

	// auth middleware
	router.Use(actions.AuthMiddleware)

	// Embed website assets
	// !IMPORTANT:
	if engine.Mode() == engine.DebugMode {
		reverseProxyWebsiteApp(router, os.Getenv("PROXY_WEBSITE"))
	} else {
		embedPublicAssets(router)
	}

	router.GET("ping", func(c *engine.Context) string {
		c.Logger.Debug("ping/pong testing")
		return "pong"
	})

	// staff admin panel
	// --------------------------------------------------------------------------
	{
		// router.GET("/logout", actions.Logout)         // 登出
		// router.GET("/api/overview", actions.Overview) // 用户信息概览

		api := router.Group("/api")
		api.POST("/spaces", actions.CreateSpace)
		api.GET("/spaces", actions.DescribeSpaces)

		space := api.Group("/spaces/:key", actions.SetSpace)
		space.GET("", actions.DescribeSpace)
		space.PATCH("", actions.UpdateSpace)
		space.POST("/pages", actions.CreatePage)
		space.GET("/pages", actions.DescribePages)

		page := space.Group("/pages/:id", actions.SetPage)
		page.GET("", actions.DescribePage)
		page.PATCH("", actions.UpdatePage)

		api.POST("/markdown/preview", actions.PreviewMarkdown)
	}

	return router
}

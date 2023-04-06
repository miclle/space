//go:build development
// +build development

package ui

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/fox-gonic/fox/engine"
)

// EmbedAssets proxy development website
func EmbedAssets(router *engine.Engine) {

	website := os.Getenv("PROXY_WEBSITE")
	origin, _ := url.Parse(website)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	for _, path := range StaticAssets {
		router.GET(path, func(c *engine.Context) (res interface{}) {
			proxy.ServeHTTP(c.Writer, c.Request)
			return
		})
	}

	for _, path := range SPARoutes {
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

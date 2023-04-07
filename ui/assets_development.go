//go:build development
// +build development

package ui

import (
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/fox-gonic/fox/engine"
	"gopkg.in/ini.v1"
)

//go:embed .env.development
var env string

var origin *url.URL

func init() {
	cfg, err := ini.Load([]byte(env))
	if err != nil {
		fmt.Printf("Fail to read file: %+v", err)
		os.Exit(1)
	}

	var (
		section = cfg.Section("")
		scheme  = "http"
		host    = section.Key("HOST").MustString("localhost")
		port    = section.Key("PORT").MustString("3000")
	)

	if section.Key("HTTPS").MustBool() {
		scheme = "https"
	}

	origin, err = url.Parse(fmt.Sprintf("%s://%s:%s", scheme, host, port))
	if err != nil {
		fmt.Printf("Fail to parse url: %+v", err)
		os.Exit(1)
	}
}

// EmbedAssets proxy development website
func EmbedAssets(router *engine.Engine) {

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{
		Director: director,
	}

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

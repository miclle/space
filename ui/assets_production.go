//go:build !development
// +build !development

package ui

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/fox-gonic/fox/engine"
	"github.com/gin-gonic/gin/render"
)

// EmbedAssets embed fs from `public` dir
func EmbedAssets(router *engine.Engine) {

	var (
		embedFS = EmbedFS()
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

package ui

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"time"

	"github.com/Masterminds/sprig/v3"

	"github.com/miclle/space/models"
)

//go:embed build/* templates/*.html
var embedFS embed.FS

// StaticAssets static assets
var StaticAssets = []string{
	"/android-chrome-192x192.png",
	"/android-chrome-512x512.png",
	"/apple-touch-icon.png",
	"/favicon-16x16.png",
	"/favicon-32x32.png",
	"/favicon.ico",
	"/index.html",
	"/manifest.json",
	"/robots.txt",

	"/static",
	"/static/*filepath",
}

// SPARoutes front SPA app routes
var SPARoutes = []string{
	"/signin",
	"/signup",
	"/spaces",
	"/spaces/*filepath",
	"/accounts",
	"/accounts/*filepath",
}

// Template for all
var Template *template.Template

// PageData template obj
type PageData struct {
	Lang   string
	Title  string
	Spaces []*models.Space
	Space  *models.Space
	// PageTree models.PageTree
	Page *models.Page
}

func init() {
	funcMap := sprig.FuncMap()
	funcMap["timeUnix"] = time.Unix
	funcMap["unescapeHTML"] = unescapeHTML

	Template = template.Must(template.New("").Funcs(funcMap).ParseFS(embedFS, "templates/*.html"))
}

// resource is an interface that provides static file
type resource struct {
	prefix string
	fs     embed.FS
}

// Open to implement the interface by http.FS required
func (r *resource) Open(name string) (fs.File, error) {
	name = path.Join(r.prefix, name)
	return r.fs.Open(name)
}

// EmbedFS return embed FS
func EmbedFS() embed.FS {
	return embedFS
}

// StaticFS static http file system
func StaticFS(prefix ...string) http.FileSystem {
	var p = "static"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	return http.FS(&resource{prefix: p, fs: embedFS})
}

// unescapeHTML unescape HTML content
func unescapeHTML(x string) template.HTML {
	return template.HTML(x)
}

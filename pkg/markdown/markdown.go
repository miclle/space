package markdown

import (
	"bytes"
	"fmt"

	"github.com/longbridgeapp/autocorrect"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			// don’t using html.WithHardWraps
		),
	)
}

// Options with markdown parser
type Options struct {
	Format bool
}

// Parse markdown convert to html
func Parse(content string, option ...Options) (string, error) {

	opt := Options{
		Format: true,
	}

	if len(option) > 0 {
		opt = option[0]
	}

	buf := new(bytes.Buffer)

	err := md.Convert([]byte(content), buf)
	if err != nil {
		return "", fmt.Errorf("markdown convert failed, err: %+v", err)
	}

	html := buf.String()

	if opt.Format {
		html, err = autocorrect.FormatHTML(html)
		if err != nil {
			return "", fmt.Errorf("format markdown html failed, err: %+v", err)
		}
	}

	// scrub content of XSS
	html = bluemonday.UGCPolicy().Sanitize(html)

	return html, nil
}

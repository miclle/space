package actions

import (
	"html/template"

	"github.com/fox-gonic/fox/engine"
	"github.com/gin-gonic/gin/render"

	"github.com/miclle/space/pkg/markdown"
)

// PreviewMarkdownArgs create space args
type PreviewMarkdownArgs struct {
	Content string `json:"content"`
}

// PreviewMarkdown markdown to html
func (actions *Actions) PreviewMarkdown(c *engine.Context, args *PreviewMarkdownArgs) (render.HTML, error) {

	result, err := markdown.Parse(args.Content)
	if err != nil {
		return render.HTML{}, err
	}

	tmpl, err := template.New("markdown").Parse(result)
	if err != nil {
		return render.HTML{}, err
	}

	var html = render.HTML{
		Template: tmpl,
	}

	return html, nil
}

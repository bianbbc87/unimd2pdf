package parser

import (
	"bytes"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	goldmarkparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"unimd2pdf/config"
)

// GoldmarkParser implements convert.MarkdownParser.
type GoldmarkParser struct {
	Extra []*GoldmarkExtension
}

// Use registers a custom extension plugin.
func (g *GoldmarkParser) Use(ext *GoldmarkExtension) {
	g.Extra = append(g.Extra, ext)
}

func (g *GoldmarkParser) Parse(source []byte, cfg *config.Config) (string, error) {
	exts := []goldmark.Extender{
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle(cfg.Code.HighlightStyle),
		),
	}

	if cfg.Markdown.Footnote {
		exts = append(exts, extension.NewFootnote())
	}
	if cfg.Markdown.DefinitionList {
		exts = append(exts, extension.DefinitionList)
	}
	if cfg.Markdown.Typographer {
		exts = append(exts, extension.NewTypographer())
	}
	if cfg.Markdown.CJK {
		exts = append(exts, extension.NewCJK())
	}

	for _, ext := range g.Extra {
		exts = append(exts, ext.Extender())
	}

	md := goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(
			goldmarkparser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

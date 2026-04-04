package convert

import (
	"encoding/base64"
	"fmt"
	htmlpkg "html"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"unimd2pdf/config"
)

// Pipeline orchestrates: Markdown → HTML → diagrams → images → CSS → PDF.
// All fields are interfaces defined in this package (convert.go).
// Concrete types are injected by main.go.
type Pipeline struct {
	Parser           MarkdownParser
	Renderer         Renderer
	Theme            ThemeProvider
	DiagramRenderers []DiagramRenderer
}

// RegisterDiagram adds a DiagramRenderer to the pipeline.
func (p *Pipeline) RegisterDiagram(r DiagramRenderer) {
	p.DiagramRenderers = append(p.DiagramRenderers, r)
}

// Convert runs the full pipeline.
func (p *Pipeline) Convert(mdBytes []byte, cfg *config.Config) ([]byte, error) {
	// 1. Markdown → HTML
	html, err := p.Parser.Parse(mdBytes, cfg)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	// 2. Diagram code blocks → inline SVG
	html = p.renderDiagrams(html, cfg)

	// 3. Local images → base64 data URIs
	html = inlineImages(html, cfg.BaseDir)

	// 4. Wrap with theme CSS
	html = p.wrapHTML(html, cfg)

	// 5. HTML → PDF
	return p.Renderer.Render(html, cfg)
}

func (p *Pipeline) renderDiagrams(htmlStr string, cfg *config.Config) string {
	for _, r := range p.DiagramRenderers {
		if !r.Available() {
			fmt.Fprintf(os.Stderr, "Warning: %s not available, keeping as code block\n", r.Name())
			continue
		}
		pattern := fmt.Sprintf(`(?s)<pre><code class="language-%s">(.*?)</code></pre>`, regexp.QuoteMeta(r.Name()))
		re := regexp.MustCompile(pattern)
		htmlStr = re.ReplaceAllStringFunc(htmlStr, func(match string) string {
			parts := re.FindStringSubmatch(match)
			if len(parts) < 2 {
				return match
			}
			result, err := r.Render(htmlpkg.UnescapeString(parts[1]), cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s failed: %v\n", r.Name(), err)
				return match
			}
			return result
		})
	}
	return htmlStr
}

func (p *Pipeline) wrapHTML(body string, cfg *config.Config) string {
	css := p.Theme.CSS(cfg)
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
%s
</style>
</head>
<body>
%s
</body>
</html>`, css, body)
}

func inlineImages(htmlStr string, baseDir string) string {
	re := regexp.MustCompile(`(<img\s[^>]*?src=")([^"]+)(")`)
	return re.ReplaceAllStringFunc(htmlStr, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) < 4 {
			return match
		}
		src := parts[2]
		if strings.HasPrefix(src, "data:") || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			return match
		}
		imgPath := src
		if !filepath.IsAbs(imgPath) {
			imgPath = filepath.Join(baseDir, imgPath)
		}
		data, err := os.ReadFile(imgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot read image %s: %v\n", imgPath, err)
			return match
		}
		mimeType := http.DetectContentType(data)
		if mimeType == "application/octet-stream" {
			mimeType = mime.TypeByExtension(filepath.Ext(imgPath))
			if mimeType == "" {
				mimeType = "image/png"
			}
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		return parts[1] + "data:" + mimeType + ";base64," + encoded + parts[3]
	})
}

package theme

import (
	"fmt"
	"os"
	"strings"

	"unimd2pdf/config"
)

// Theme is any type with a CSS method. Satisfied by Light, Dark, Custom.
// (Matches convert.ThemeProvider interface via Go structural typing.)
type Theme = interface{ CSS(cfg *config.Config) string }

// Resolve returns the appropriate theme for the config.
func Resolve(cfg *config.Config) Theme {
	switch cfg.Theme {
	case "dark":
		return &Dark{}
	case "light", "":
		return &Light{}
	default:
		// Treat as custom CSS file path
		return &Custom{Path: cfg.Theme}
	}
}

// Light is the default light theme.
type Light struct{}

func (l *Light) CSS(cfg *config.Config) string {
	return baseCSS(cfg, lightColors())
}

// Dark theme.
type Dark struct{}

func (d *Dark) CSS(cfg *config.Config) string {
	return baseCSS(cfg, darkColors())
}

// Custom loads CSS from a file path.
type Custom struct{ Path string }

func (c *Custom) CSS(cfg *config.Config) string {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cannot read custom theme %s: %v, falling back to light\n", c.Path, err)
		return (&Light{}).CSS(cfg)
	}
	// Prepend page/font config, then custom CSS
	return pageCSS(cfg) + "\n" + string(data)
}

type colors struct {
	bg, fg, codeBg, codeBorder    string
	tableBorder, tableHeaderBg    string
	tableStripeBg                 string
	blockquoteBorder, blockquoteBg string
	blockquoteFg                  string
	linkColor                     string
	hrColor                       string
	h1Border, h2Border            string
}

func lightColors() colors {
	return colors{
		bg: "#ffffff", fg: "#1a1a1a",
		codeBg: "#f6f8fa", codeBorder: "#e1e4e8",
		tableBorder: "#d0d7de", tableHeaderBg: "#f6f8fa",
		tableStripeBg: "#f9fafb",
		blockquoteBorder: "#dfe2e5", blockquoteBg: "#f9fafb",
		blockquoteFg: "#57606a",
		linkColor: "#0969da",
		hrColor: "#d0d7de",
		h1Border: "#e1e4e8", h2Border: "#eaecef",
	}
}

func darkColors() colors {
	return colors{
		bg: "#0d1117", fg: "#e6edf3",
		codeBg: "#161b22", codeBorder: "#30363d",
		tableBorder: "#30363d", tableHeaderBg: "#161b22",
		tableStripeBg: "#0d1117",
		blockquoteBorder: "#3b434b", blockquoteBg: "#161b22",
		blockquoteFg: "#8b949e",
		linkColor: "#58a6ff",
		hrColor: "#30363d",
		h1Border: "#30363d", h2Border: "#21262d",
	}
}

func pageCSS(cfg *config.Config) string {
	margin := parseMargin(cfg.Page.Margin)
	size := cfg.Page.Size
	if cfg.IsLandscape() {
		size += " landscape"
	}

	return fmt.Sprintf(`@page {
  size: %s;
  margin: %s;
}`, size, margin)
}

func baseCSS(cfg *config.Config, c colors) string {
	return fmt.Sprintf(`%s

body {
  font-family: %s;
  font-size: %s;
  line-height: 1.6;
  color: %s;
  background-color: %s;
  max-width: 100%%;
}

h1 {
  font-size: 22pt;
  font-weight: 700;
  margin-top: 0;
  margin-bottom: 0.6em;
  padding-bottom: 0.3em;
  border-bottom: 1px solid %s;
}

h2 {
  font-size: 17pt;
  font-weight: 600;
  margin-top: 1.4em;
  margin-bottom: 0.5em;
  padding-bottom: 0.2em;
  border-bottom: 1px solid %s;
}

h3 {
  font-size: 13pt;
  font-weight: 600;
  margin-top: 1.2em;
  margin-bottom: 0.4em;
}

p { margin: 0.6em 0; }

ul, ol { padding-left: 2em; margin: 0.5em 0; }
li { margin: 0.25em 0; }

code {
  font-family: "SF Mono", "Fira Code", "Fira Mono", Menlo, Consolas, monospace;
  font-size: 0.9em;
  background-color: %s;
  padding: 0.15em 0.4em;
  border-radius: 3px;
}

pre {
  background-color: %s;
  border: 1px solid %s;
  border-radius: 6px;
  padding: 12px 16px;
  overflow-x: auto;
  font-size: 0.85em;
  line-height: 1.5;
}

pre code { background: none; padding: 0; border-radius: 0; }

table {
  border-collapse: collapse;
  width: 100%%;
  margin: 0.8em 0;
  font-size: 0.92em;
}

th, td {
  border: 1px solid %s;
  padding: 8px 12px;
  text-align: left;
}

th {
  background-color: %s;
  font-weight: 600;
}

tr:nth-child(even) { background-color: %s; }

blockquote {
  margin: 0.8em 0;
  padding: 0.4em 1em;
  border-left: 4px solid %s;
  color: %s;
  background-color: %s;
}

blockquote p { margin: 0.3em 0; }

hr {
  border: none;
  border-top: 2px solid %s;
  margin: 1.5em 0;
}

img { max-width: 100%%; height: auto; }

a { color: %s; text-decoration: none; }

strong { font-weight: 600; }

.mermaid-diagram { text-align: center; margin: 1em 0; }
.mermaid-diagram svg { max-width: 100%%; height: auto; }`,
		pageCSS(cfg),
		cfg.Font.Family, cfg.Font.Size,
		c.fg, c.bg,
		c.h1Border,
		c.h2Border,
		c.codeBg,       // inline code bg
		c.codeBg,       // pre bg
		c.codeBorder,   // pre border
		c.tableBorder,  // th/td border
		c.tableHeaderBg,
		c.tableStripeBg,
		c.blockquoteBorder, c.blockquoteFg, c.blockquoteBg,
		c.hrColor,
		c.linkColor,
	)
}

// parseMargin normalizes margin input.
// "20mm" → "20mm", "20mm 18mm" → "20mm 18mm", etc.
func parseMargin(m string) string {
	m = strings.TrimSpace(m)
	if m == "" {
		return "20mm"
	}
	return m
}

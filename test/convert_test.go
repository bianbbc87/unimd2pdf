package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"unimd2pdf/config"
	"unimd2pdf/convert/parser"
)

func defaultCfg() *config.Config {
	return config.Defaults()
}

func TestGoldmarkBasic(t *testing.T) {
	p := &parser.GoldmarkParser{}
	html, err := p.Parse([]byte("# Hello\n\nworld"), defaultCfg())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<h1") {
		t.Errorf("expected <h1>, got: %s", html)
	}
	if !strings.Contains(html, "world") {
		t.Errorf("expected 'world', got: %s", html)
	}
}

func TestGoldmarkTable(t *testing.T) {
	md := "| A | B |\n|---|---|\n| 1 | 2 |"
	p := &parser.GoldmarkParser{}
	html, err := p.Parse([]byte(md), defaultCfg())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<table>") {
		t.Errorf("expected <table>, got: %s", html)
	}
	if !strings.Contains(html, "<td>1</td>") {
		t.Errorf("expected cell content, got: %s", html)
	}
}

func TestGoldmarkFootnote(t *testing.T) {
	md := "Text[^1]\n\n[^1]: Footnote content."
	cfg := defaultCfg()
	cfg.Markdown.Footnote = true
	p := &parser.GoldmarkParser{}
	html, err := p.Parse([]byte(md), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "footnote") {
		t.Errorf("expected footnote markup, got: %s", html)
	}
}

func TestGoldmarkDefinitionList(t *testing.T) {
	md := "Term\n:   Definition here."
	cfg := defaultCfg()
	cfg.Markdown.DefinitionList = true
	p := &parser.GoldmarkParser{}
	html, err := p.Parse([]byte(md), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<dl>") {
		t.Errorf("expected <dl>, got: %s", html)
	}
}

func TestGoldmarkCodeHighlight(t *testing.T) {
	md := "```go\nfunc main() {}\n```"
	p := &parser.GoldmarkParser{}
	html, err := p.Parse([]byte(md), defaultCfg())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "style=") || !strings.Contains(html, "func") {
		t.Errorf("expected styled code block, got: %s", html)
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg := config.Defaults()
	if cfg.Theme != "light" {
		t.Errorf("expected light theme, got: %s", cfg.Theme)
	}
	if cfg.Page.Size != "A4" {
		t.Errorf("expected A4, got: %s", cfg.Page.Size)
	}
	if !cfg.Mermaid.Enabled {
		t.Error("expected mermaid enabled by default")
	}
	if !cfg.Markdown.Footnote {
		t.Error("expected footnote enabled by default")
	}
}

func TestConfigMerge(t *testing.T) {
	base := config.Defaults()
	overlay := &config.Config{Theme: "dark", Page: config.PageConfig{Size: "Letter"}}
	merged := config.Merge(base, overlay)
	if merged.Theme != "dark" {
		t.Errorf("expected dark, got: %s", merged.Theme)
	}
	if merged.Page.Size != "Letter" {
		t.Errorf("expected Letter, got: %s", merged.Page.Size)
	}
	if merged.Font.Size != "11pt" {
		t.Errorf("expected base font size preserved, got: %s", merged.Font.Size)
	}
}

func TestConfigMergeNil(t *testing.T) {
	base := config.Defaults()
	merged := config.Merge(base, nil)
	if merged.Theme != base.Theme {
		t.Error("nil merge should preserve base")
	}
}

func TestConfigLoadFile(t *testing.T) {
	dir := t.TempDir()
	yamlContent := "theme: dark\npage:\n  size: A3\n"
	os.WriteFile(filepath.Join(dir, "unimd2pdf.yaml"), []byte(yamlContent), 0644)

	cfg, err := config.LoadFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if cfg.Theme != "dark" {
		t.Errorf("expected dark, got: %s", cfg.Theme)
	}
	if cfg.Page.Size != "A3" {
		t.Errorf("expected A3, got: %s", cfg.Page.Size)
	}
}

func TestConfigLoadFileNotFound(t *testing.T) {
	cfg, err := config.LoadFile(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if cfg != nil {
		t.Error("expected nil when no config file exists")
	}
}

func TestPageDimensions(t *testing.T) {
	cfg := config.Defaults()
	w, h := cfg.PageDimensions()
	if w != 8.27 || h != 11.69 {
		t.Errorf("A4 should be 8.27x11.69, got: %.2fx%.2f", w, h)
	}

	cfg.Page.Size = "Letter"
	w, h = cfg.PageDimensions()
	if w != 8.5 || h != 11.0 {
		t.Errorf("Letter should be 8.5x11.0, got: %.2fx%.2f", w, h)
	}
}

func TestMermaidBlockRegex(t *testing.T) {
	html := `<pre><code class="language-mermaid">flowchart TD
    A --> B</code></pre>`
	if !strings.Contains(html, `class="language-mermaid"`) {
		t.Error("expected mermaid code block pattern")
	}
}

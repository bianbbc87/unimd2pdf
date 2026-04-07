package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"unimd2pdf/config"
	"unimd2pdf/convert"
	"unimd2pdf/convert/diagram"
	"unimd2pdf/convert/parser"
	"unimd2pdf/convert/renderer"
	"unimd2pdf/convert/theme"
)

var version = "dev"

func main() {
	input := flag.String("i", "", "Input markdown file (required)")
	output := flag.String("o", "", "Output PDF path (default: input with .pdf extension)")
	themeName := flag.String("theme", "", "Theme: light, dark, or path to custom CSS")
	pageSize := flag.String("page-size", "", "Page size: A3, A4, A5, Letter")
	orientation := flag.String("orientation", "", "Orientation: portrait, landscape")
	fontFamily := flag.String("font", "", "Font family")
	fontSize := flag.String("font-size", "", "Font size (e.g. 14pt)")
	margin := flag.String("margin", "", `Margin (e.g. "20mm", "20mm 18mm")`)
	highlightStyle := flag.String("highlight-style", "", "Code highlight style (e.g. github, monokai)")
	mermaidTheme := flag.String("mermaid-theme", "", "Mermaid theme: default, dark, forest, neutral")
	noMermaid := flag.Bool("no-mermaid", false, "Disable mermaid rendering (code block fallback)")
	showVersion := flag.Bool("version", false, "Print version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "unimd2pdf - Universal Markdown to PDF converter\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  unimd2pdf -i input.md [-o output.pdf] [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nConfig file:\n")
		fmt.Fprintf(os.Stderr, "  Place unimd2pdf.yaml in your project root for declarative config.\n")
		fmt.Fprintf(os.Stderr, "  CLI flags override config file values.\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("unimd2pdf %s\n", version)
		os.Exit(0)
	}

	if *input == "" && flag.NArg() > 0 {
		*input = flag.Arg(0)
	}
	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: input file required (-i flag or positional argument)")
		flag.Usage()
		os.Exit(1)
	}
	if *output == "" {
		*output = strings.TrimSuffix(*input, filepath.Ext(*input)) + ".pdf"
	}

	// --- Config: defaults → file → CLI ---
	cfg := config.Defaults()
	inputDir, _ := filepath.Abs(filepath.Dir(*input))
	if fileCfg, err := config.LoadFile(inputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: config file error: %v\n", err)
	} else {
		cfg = config.Merge(cfg, fileCfg)
	}
	cliCfg := &config.Config{}
	cliCfg.Theme = *themeName
	cliCfg.Page.Size = *pageSize
	cliCfg.Page.Orientation = *orientation
	cliCfg.Page.Margin = *margin
	cliCfg.Font.Family = *fontFamily
	cliCfg.Font.Size = *fontSize
	cliCfg.Code.HighlightStyle = *highlightStyle
	cliCfg.Mermaid.Theme = *mermaidTheme
	cfg = config.Merge(cfg, cliCfg)
	if *noMermaid {
		cfg.Mermaid.Enabled = false
	}
	cfg.Input = *input
	cfg.Output = *output
	cfg.BaseDir = inputDir

	// Auto-select dark highlight style when theme is dark and user didn't override
	if cfg.Theme == "dark" && cfg.Code.HighlightStyle == "github" {
		cfg.Code.HighlightStyle = "monokai"
	}

	// --- Read ---
	mdBytes, err := os.ReadFile(cfg.Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot read %s: %v\n", cfg.Input, err)
		os.Exit(1)
	}

	// --- Assemble pipeline (Layer 3: wiring) ---
	// main.go is the ONLY place that imports concrete types.
	// Pipeline fields are all convert.* interfaces.
	p := &convert.Pipeline{
		Parser:   &parser.GoldmarkParser{},
		Renderer: &renderer.ChromeRenderer{},
		Theme:    theme.Resolve(cfg),
	}
	if cfg.Mermaid.Enabled {
		p.RegisterDiagram(diagram.NewMermaidRenderer())
	}

	// --- Convert ---
	pdfBytes, err := p.Convert(mdBytes, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(cfg.Output, pdfBytes, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot write %s: %v\n", cfg.Output, err)
		os.Exit(1)
	}
	fmt.Printf("PDF generated: %s (%d bytes)\n", cfg.Output, len(pdfBytes))
}

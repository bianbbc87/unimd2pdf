package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config is the central configuration struct.
// Populated by: defaults → config file → CLI flags (each layer overrides the previous).
type Config struct {
	Theme string     `yaml:"theme"` // "light" | "dark" | path to custom CSS
	Page  PageConfig `yaml:"page"`
	Font  FontConfig `yaml:"font"`
	Code  CodeConfig `yaml:"code"`

	Markdown MarkdownConfig `yaml:"markdown"`
	Mermaid  MermaidConfig  `yaml:"mermaid"`

	// Runtime fields (not from YAML)
	Input   string `yaml:"-"`
	Output  string `yaml:"-"`
	BaseDir string `yaml:"-"` // resolved from input path
}

type PageConfig struct {
	Size        string `yaml:"size"`        // A3, A4, A5, Letter
	Orientation string `yaml:"orientation"` // portrait, landscape
	Margin      string `yaml:"margin"`      // e.g. "20mm", "20mm 18mm", "20mm 18mm 20mm 18mm"
}

type FontConfig struct {
	Family string `yaml:"family"` // e.g. "Apple SD Gothic Neo"
	Size   string `yaml:"size"`   // e.g. "11pt"
}

type CodeConfig struct {
	HighlightStyle string `yaml:"highlight-style"` // chroma style name
}

type MarkdownConfig struct {
	Footnote       bool `yaml:"footnote"`       // [^1] footnote syntax
	DefinitionList bool `yaml:"definitionlist"` // term/definition list
	Typographer    bool `yaml:"typographer"`    // smart quotes, dashes, ellipses
	CJK            bool `yaml:"cjk"`            // CJK line break handling
	// Math: reserved for future implementation (KaTeX/MathJax)
}

type MermaidConfig struct {
	Enabled bool   `yaml:"enabled"`
	Theme   string `yaml:"theme"` // default, dark, forest, neutral
}

// PageDimensions returns (width, height) in inches for the configured page size.
func (c *Config) PageDimensions() (float64, float64) {
	switch c.Page.Size {
	case "A3":
		return 11.69, 16.54
	case "A5":
		return 5.83, 8.27
	case "Letter":
		return 8.5, 11.0
	default: // A4
		return 8.27, 11.69
	}
}

// IsLandscape returns true if orientation is landscape.
func (c *Config) IsLandscape() bool {
	return c.Page.Orientation == "landscape"
}

// Defaults returns a Config with sensible default values.
func Defaults() *Config {
	return &Config{
		Theme: "light",
		Page: PageConfig{
			Size:        "A4",
			Orientation: "portrait",
			Margin:      "20mm",
		},
		Font: FontConfig{
			Family: `-apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif`,
			Size:   "11pt",
		},
		Code: CodeConfig{
			HighlightStyle: "github",
		},
		Markdown: MarkdownConfig{
			Footnote:       true,
			DefinitionList: true,
			Typographer:    true,
			CJK:            true,
		},
		Mermaid: MermaidConfig{
			Enabled: true,
			Theme:   "default",
		},
	}
}

// LoadFile reads a unimd2pdf.yaml from the given directory (or parents).
// Returns nil if no config file is found (not an error).
func LoadFile(dir string) (*Config, error) {
	path := findConfigFile(dir)
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return &cfg, nil
}

// findConfigFile walks up from dir looking for unimd2pdf.yaml.
func findConfigFile(dir string) string {
	dir, _ = filepath.Abs(dir)
	for {
		candidate := filepath.Join(dir, "unimd2pdf.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		candidate = filepath.Join(dir, ".unimd2pdf.yaml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// Merge overlays non-zero values from overlay onto base, returning a new Config.
func Merge(base, overlay *Config) *Config {
	result := *base
	if overlay == nil {
		return &result
	}
	if overlay.Theme != "" {
		result.Theme = overlay.Theme
	}
	if overlay.Page.Size != "" {
		result.Page.Size = overlay.Page.Size
	}
	if overlay.Page.Orientation != "" {
		result.Page.Orientation = overlay.Page.Orientation
	}
	if overlay.Page.Margin != "" {
		result.Page.Margin = overlay.Page.Margin
	}
	if overlay.Font.Family != "" {
		result.Font.Family = overlay.Font.Family
	}
	if overlay.Font.Size != "" {
		result.Font.Size = overlay.Font.Size
	}
	if overlay.Code.HighlightStyle != "" {
		result.Code.HighlightStyle = overlay.Code.HighlightStyle
	}
	if overlay.Mermaid.Theme != "" {
		result.Mermaid.Theme = overlay.Mermaid.Theme
	}
	// Mermaid.Enabled: only override if explicitly set in overlay file
	// (YAML zero value is false, so we can't distinguish "not set" from "set to false" here.
	//  CLI --no-mermaid handles the explicit disable.)
	return &result
}

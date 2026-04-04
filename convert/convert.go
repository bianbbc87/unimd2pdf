// Package convert defines the interfaces for the unimd2pdf conversion pipeline.
//
// Layer rules:
//
//	Layer 0: config/       — data structs only, imports nothing internal
//	Layer 1: convert/      — interfaces only, imports config/ only
//	Layer 2: convert/*/    — implementations, import config/ and convert/ only
//	Layer 3: main.go       — wiring, imports everything
//
// Implementations MUST NOT import each other (no parser → renderer, etc.).
// main.go is the only place that knows about concrete types.
package convert

import "unimd2pdf/config"

// MarkdownParser converts markdown bytes to HTML string.
type MarkdownParser interface {
	Parse(source []byte, cfg *config.Config) (string, error)
}

// Renderer converts a complete HTML document to PDF bytes.
type Renderer interface {
	Render(html string, cfg *config.Config) ([]byte, error)
}

// DiagramRenderer converts a fenced code block to inline HTML (typically SVG).
type DiagramRenderer interface {
	Name() string
	Available() bool
	Render(source string, cfg *config.Config) (string, error)
}

// ThemeProvider generates CSS from configuration.
type ThemeProvider interface {
	CSS(cfg *config.Config) string
}

// MarkdownExtension is a pluggable goldmark extension.
// Kept intentionally minimal — implementations provide the goldmark.Extender.
type MarkdownExtension interface {
	Name() string
}

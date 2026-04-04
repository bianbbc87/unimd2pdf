# Architecture

## Layer Rules

unimd2pdf enforces a strict four-layer dependency hierarchy. Each layer may only import from layers below it.

```
Layer 3: main.go            (wiring — imports everything)
   |
Layer 2: convert/*/         (implementations — import config/ only)
   |
Layer 1: convert/           (interfaces — imports config/ only)
   |
Layer 0: config/            (data structs — imports nothing internal)
```

Implementations in Layer 2 MUST NOT import each other. There is no `parser -> renderer` or `diagram -> theme` dependency. All cross-cutting coordination happens through the `Pipeline` in Layer 1, wired together by `main.go` in Layer 3.

## Package Responsibilities

### `config/` (Layer 0)

Data-only package. Defines the `Config` struct and its sub-structs (`PageConfig`, `FontConfig`, `CodeConfig`, `MarkdownConfig`, `MermaidConfig`).

Provides:
- `Defaults()` — sensible default configuration
- `LoadFile(dir)` — reads `unimd2pdf.yaml` from the given directory or parents
- `Merge(base, overlay)` — overlays non-zero values onto a base config

### `convert/` (Layer 1)

Interface definitions and pipeline orchestration. Defines:
- `MarkdownParser` — converts markdown bytes to HTML string
- `Renderer` — converts HTML to PDF bytes
- `DiagramRenderer` — converts fenced code blocks to inline HTML/SVG
- `ThemeProvider` — generates CSS from configuration
- `MarkdownExtension` — pluggable goldmark extension

The `Pipeline` struct orchestrates the full conversion flow:

```
Markdown → HTML → Diagrams → Images → CSS → PDF
```

### `convert/parser/` (Layer 2)

`GoldmarkParser` implements `MarkdownParser` using the goldmark library. Supports GFM, syntax highlighting, footnotes, definition lists, typographer, and CJK extensions (all configurable).

### `convert/renderer/` (Layer 2)

`ChromeRenderer` implements `Renderer` using headless Chrome via chromedp. Handles page size, orientation, and margin configuration.

### `convert/diagram/` (Layer 2)

`MermaidRenderer` implements `DiagramRenderer` using the mmdc CLI (Mermaid CLI). Converts mermaid code blocks to inline SVG.

### `convert/theme/` (Layer 2)

Provides `Light`, `Dark`, and `Custom` theme implementations. `Resolve(cfg)` returns the appropriate theme based on config. Custom themes load CSS from a file path.

### `main.go` (Layer 3)

The only file that imports concrete types. Handles:
- CLI flag parsing
- Config precedence: defaults -> YAML file -> CLI flags
- Pipeline assembly (injecting concrete implementations)
- File I/O

## Dependency Flow

```
main.go
  ├── config/              (Config, Defaults, LoadFile, Merge)
  ├── convert/             (Pipeline, interfaces)
  ├── convert/parser/      (GoldmarkParser)
  ├── convert/renderer/    (ChromeRenderer)
  ├── convert/diagram/     (MermaidRenderer)
  └── convert/theme/       (Resolve, Light, Dark, Custom)

convert/pipeline.go
  └── config/              (Config)

convert/parser/goldmark.go
  └── config/              (Config)

convert/renderer/chrome.go
  └── config/              (Config)

convert/diagram/mermaid.go
  └── config/              (Config)

convert/theme/theme.go
  └── config/              (Config)
```

## Adding New Implementations

### New Parser

1. Create `convert/parser/myparser.go`
2. Implement the `convert.MarkdownParser` interface:
   ```go
   type MyParser struct{}
   func (p *MyParser) Parse(source []byte, cfg *config.Config) (string, error) { ... }
   ```
3. Wire it in `main.go`: `p.Parser = &parser.MyParser{}`

### New Renderer

1. Create `convert/renderer/myrenderer.go`
2. Implement the `convert.Renderer` interface:
   ```go
   type MyRenderer struct{}
   func (r *MyRenderer) Render(html string, cfg *config.Config) ([]byte, error) { ... }
   ```
3. Wire it in `main.go`: `p.Renderer = &renderer.MyRenderer{}`

### New Diagram Renderer

1. Create `convert/diagram/mydiagram.go`
2. Implement the `convert.DiagramRenderer` interface:
   ```go
   type MyDiagramRenderer struct{}
   func (d *MyDiagramRenderer) Name() string { ... }
   func (d *MyDiagramRenderer) Available() bool { ... }
   func (d *MyDiagramRenderer) Render(source string, cfg *config.Config) (string, error) { ... }
   ```
3. Register in `main.go`: `p.RegisterDiagram(&diagram.MyDiagramRenderer{})`

### New Theme

1. Create `convert/theme/mytheme.go`
2. Implement the `convert.ThemeProvider` interface:
   ```go
   type MyTheme struct{}
   func (t *MyTheme) CSS(cfg *config.Config) string { ... }
   ```
3. Add a case in `theme.Resolve()` or set it directly in `main.go`: `p.Theme = &theme.MyTheme{}`

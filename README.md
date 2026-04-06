# unimd2pdf

Universal Markdown to PDF converter. Single binary, zero config, batteries included.

## Features

- GFM (tables, strikethrough, task lists, autolinks)
- Syntax highlighting (github, monokai, dracula, ...)
- Mermaid diagrams (flowchart, sequence, class, state) via mmdc
- Footnotes and definition lists
- Smart typography (curly quotes, em/en dashes)
- CJK text support
- Light / Dark / Custom CSS themes
- Declarative config via `unimd2pdf.yaml`
- Local images auto-embedded as base64

## Architecture

```
main.go              CLI + wiring (Layer 3)
  |
convert/             interfaces + pipeline (Layer 1)
  |
convert/parser/      GoldmarkParser (Layer 2)
convert/renderer/    ChromeRenderer (Layer 2)
convert/diagram/     MermaidRenderer (Layer 2)
convert/theme/       Light/Dark/Custom (Layer 2)
  |
config/              data structs (Layer 0)
```

Layer 2 packages never import each other. See [docs/architecture.md](docs/architecture.md).

## Quick Start

### Install

```bash
brew install bianbbc87/tap/unimd2pdf
```

Or with Go:

```bash
go install github.com/bianbbc87/unimd2pdf@latest
```

### Usage

```bash
unimd2pdf -i README.md                          # basic
unimd2pdf -i doc.md --theme dark                 # dark theme
unimd2pdf -i doc.md --page-size Letter           # letter size
unimd2pdf -i doc.md --font "Noto Sans KR"        # custom font
```

## Samples

Generate all sample PDFs:

```bash
cd samples
for f in *.md; do unimd2pdf -i "$f"; done
unimd2pdf -i 08-dark-theme.md --theme dark       # dark theme sample
```

| Sample | Description | PDF |
|--------|-------------|-----|
| [01-basic.md](samples/01-basic.md) | Headings, lists, quotes, links | [Download](samples/01-basic.pdf) |
| [02-gfm.md](samples/02-gfm.md) | Tables, strikethrough, task lists | [Download](samples/02-gfm.pdf) |
| [03-code.md](samples/03-code.md) | Syntax highlighting (6 languages) | [Download](samples/03-code.pdf) |
| [04-mermaid.md](samples/04-mermaid.md) | Flowchart, sequence, class, state | [Download](samples/04-mermaid.pdf) |
| [05-images.md](samples/05-images.md) | Local image embedding | [Download](samples/05-images.pdf) |
| [06-footnotes.md](samples/06-footnotes.md) | Footnotes, definition lists | [Download](samples/06-footnotes.pdf) |
| [07-typography.md](samples/07-typography.md) | Smart quotes, dashes, CJK | [Download](samples/07-typography.pdf) |
| [08-dark-theme.md](samples/08-dark-theme.md) | Dark theme showcase | [Download](samples/08-dark-theme.pdf) |
| [09-full-proposal.md](samples/09-full-proposal.md) | Full document with all features | [Download](samples/09-full-proposal.pdf) |

## Config File

Place `unimd2pdf.yaml` in your project root:

```yaml
theme: light
page:
  size: A4
  margin: "20mm 18mm"
code:
  highlight-style: github
markdown:
  footnote: true
  definitionlist: true
  typographer: true
  cjk: true
mermaid:
  enabled: true
  theme: default
```

Precedence: defaults < yaml config < CLI flags.

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i` | (required) | Input markdown file |
| `-o` | `<input>.pdf` | Output PDF path |
| `-theme` | `light` | `light`, `dark`, or CSS file path |
| `-page-size` | `A4` | `A3`, `A4`, `A5`, `Letter` |
| `-orientation` | `portrait` | `portrait`, `landscape` |
| `-font` | system | Font family |
| `-font-size` | `11pt` | Font size |
| `-margin` | `20mm` | Page margin |
| `-highlight-style` | `github` | Chroma style name |
| `-mermaid-theme` | `default` | `default`, `dark`, `forest`, `neutral` |
| `-no-mermaid` | `false` | Disable mermaid rendering |

## Prerequisites

- **Chrome/Chromium** (required) - for PDF rendering
- **mmdc** (optional) - for Mermaid diagrams: `npm i -g @mermaid-js/mermaid-cli`

## Contributing

See [docs/contributing.md](docs/contributing.md).

## License

MIT - see [LICENSE](LICENSE).

# Usage

## Basic Usage

```bash
unimd2pdf -i input.md -o output.pdf
```

If `-o` is omitted, the output file is the input file with a `.pdf` extension:

```bash
unimd2pdf -i README.md
# produces README.pdf
```

You can also pass the input file as a positional argument:

```bash
unimd2pdf README.md
```

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i` | (required) | Input markdown file |
| `-o` | `<input>.pdf` | Output PDF path |
| `-theme` | `light` | Theme: `light`, `dark`, or path to custom CSS file |
| `-page-size` | `A4` | Page size: `A3`, `A4`, `A5`, `Letter` |
| `-orientation` | `portrait` | Orientation: `portrait`, `landscape` |
| `-font` | system default | Font family |
| `-font-size` | `11pt` | Font size (e.g., `14pt`) |
| `-margin` | `20mm` | Margin (e.g., `20mm`, `20mm 18mm`, `20mm 18mm 20mm 18mm`) |
| `-highlight-style` | `github` | Code highlight style (chroma style name, e.g., `monokai`) |
| `-mermaid-theme` | `default` | Mermaid theme: `default`, `dark`, `forest`, `neutral` |
| `-no-mermaid` | `false` | Disable mermaid rendering (code block fallback) |
| `-version` | | Print version and exit |

## Config File

Place a `unimd2pdf.yaml` file in your project directory for declarative configuration. The tool searches the input file's directory and walks up to parent directories.

### Full Example

```yaml
# Theme: light, dark, or path to custom CSS
theme: light

# Page settings
page:
  size: A4                  # A3, A4, A5, Letter
  orientation: portrait     # portrait, landscape
  margin: "20mm"            # "20mm", "20mm 18mm", "20mm 18mm 20mm 18mm"

# Font settings
font:
  family: '-apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif'
  size: 11pt

# Code highlighting
code:
  highlight-style: github   # chroma style: github, monokai, dracula, etc.

# Markdown extensions
markdown:
  footnote: true            # [^1] footnote syntax
  definitionlist: true      # term/definition lists
  typographer: true         # smart quotes, em/en dashes, ellipses
  cjk: true                 # CJK line break handling

# Mermaid diagram rendering
mermaid:
  enabled: true
  theme: default            # default, dark, forest, neutral
```

## Config Precedence

Configuration is resolved in order, with each layer overriding the previous:

1. **Defaults** — built-in sensible defaults
2. **YAML config** — `unimd2pdf.yaml` found in input directory (or parent)
3. **CLI flags** — command-line flags override everything

For example, if the YAML file sets `theme: dark` but you pass `-theme light` on the command line, the light theme is used.

## Examples

Convert with dark theme:

```bash
unimd2pdf -i doc.md -theme dark
```

A4 landscape with custom margins:

```bash
unimd2pdf -i slides.md -page-size A4 -orientation landscape -margin "15mm 20mm"
```

Custom font and highlight style:

```bash
unimd2pdf -i code-review.md -font "Fira Sans" -font-size 12pt -highlight-style monokai
```

Disable Mermaid (faster if you have no diagrams):

```bash
unimd2pdf -i notes.md -no-mermaid
```

Use a custom CSS theme:

```bash
unimd2pdf -i report.md -theme ./my-theme.css
```

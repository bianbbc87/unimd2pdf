# Installation

## Homebrew (macOS / Linux)

```bash
brew install bianbbc87/tap/unimd2pdf
```

## Go Install

```bash
go install github.com/bianbbc87/unimd2pdf@latest
```

Requires Go 1.23 or later.

## Binary Download

Download pre-built binaries from [GitHub Releases](https://github.com/bianbbc87/unimd2pdf/releases).

Available platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

After downloading, extract the archive and place the `unimd2pdf` binary in your `$PATH`.

## Prerequisites

### Chrome / Chromium (required)

unimd2pdf uses headless Chrome for PDF rendering. You need Chrome or Chromium installed:

- **macOS**: `brew install --cask google-chrome` or use the bundled Chromium
- **Linux**: `apt install chromium-browser` or `snap install chromium`
- **Windows**: Install [Google Chrome](https://www.google.com/chrome/)

chromedp will automatically detect the Chrome installation.

### Mermaid CLI (optional)

To render Mermaid diagrams, install the `mmdc` CLI:

```bash
npm install -g @mermaid-js/mermaid-cli
```

If `mmdc` is not available, Mermaid code blocks are rendered as plain code blocks. You can also explicitly disable Mermaid rendering with `--no-mermaid`.

## Verify Installation

```bash
unimd2pdf --version
```

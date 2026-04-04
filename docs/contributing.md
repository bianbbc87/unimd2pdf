# Contributing

## Prerequisites

- **Go 1.23+**
- **Chrome or Chromium** (required for PDF rendering)
- **mmdc** (Mermaid CLI, optional) — install via `npm install -g @mermaid-js/mermaid-cli`
- **golangci-lint** (for linting) — install via `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Building

```bash
go build -o unimd2pdf ./main.go
```

## Testing

```bash
go test -v ./...
```

With race detection and coverage:

```bash
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Linting

```bash
golangci-lint run ./...
go vet ./...
```

## Code Style

### Layer Rules

This project enforces a strict layer architecture. Before submitting code, verify that your imports follow these rules:

| Layer | Package | May Import |
|-------|---------|------------|
| 0 | `config/` | Nothing internal |
| 1 | `convert/` | `config/` only |
| 2 | `convert/*/` | `config/` only |
| 3 | `main.go` | Everything |

- Implementations in `convert/*/` MUST NOT import each other.
- Only `main.go` imports concrete types.
- See [docs/architecture.md](architecture.md) for full details.

### General Guidelines

- Follow standard Go conventions (`gofmt`, `go vet`).
- Keep functions short and focused.
- Use meaningful variable names.
- Add comments for exported types and functions.
- Write tests for new functionality.

## Pull Request Process

1. Fork the repository and create a feature branch from `main`.
2. Make your changes, following the layer rules and code style.
3. Add or update tests as needed.
4. Run `go test ./...` and `golangci-lint run ./...` locally.
5. Commit with a clear message describing the change.
6. Open a pull request against `main`.
7. Ensure CI passes (tests, vet, lint, build).

## Reporting Issues

Open an issue on GitHub with:
- A clear title and description.
- Steps to reproduce (if applicable).
- Expected vs. actual behavior.
- Go version, OS, and Chrome version.

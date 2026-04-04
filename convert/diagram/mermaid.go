package diagram

import (
	"fmt"
	"os"
	"os/exec"

	"unimd2pdf/config"
)

// MermaidRenderer implements DiagramRenderer using mmdc CLI.
type MermaidRenderer struct {
	mmdc string // resolved path, empty if not available
}

func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{mmdc: findMmdc()}
}

func (m *MermaidRenderer) Name() string { return "mermaid" }

func (m *MermaidRenderer) Available() bool { return m.mmdc != "" }

func (m *MermaidRenderer) Render(source string, cfg *config.Config) (string, error) {
	if !m.Available() {
		return "", fmt.Errorf("mmdc not found")
	}

	inFile, err := os.CreateTemp("", "mermaid-*.mmd")
	if err != nil {
		return "", err
	}
	defer os.Remove(inFile.Name())
	inFile.WriteString(source)
	inFile.Close()

	outFile := inFile.Name() + ".svg"
	defer os.Remove(outFile)

	args := []string{"-i", inFile.Name(), "-o", outFile, "-b", "transparent"}
	if cfg.Mermaid.Theme != "" {
		args = append(args, "-t", cfg.Mermaid.Theme)
	}

	var cmd *exec.Cmd
	if m.mmdc == "npx:mmdc" {
		cmd = exec.Command("npx", append([]string{"mmdc"}, args...)...)
	} else {
		cmd = exec.Command(m.mmdc, args...)
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mmdc: %w", err)
	}

	svgBytes, err := os.ReadFile(outFile)
	if err != nil {
		return "", err
	}

	return `<div class="mermaid-diagram">` + string(svgBytes) + `</div>`, nil
}

func findMmdc() string {
	if path, err := exec.LookPath("mmdc"); err == nil {
		return path
	}
	if _, err := exec.LookPath("npx"); err == nil {
		return "npx:mmdc"
	}
	return ""
}

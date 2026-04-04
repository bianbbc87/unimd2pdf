package renderer

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"unimd2pdf/config"
)

// ChromeRenderer implements Renderer using headless Chrome via chromedp.
type ChromeRenderer struct{}

func (r *ChromeRenderer) Render(htmlContent string, cfg *config.Config) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "unimd2pdf-*.html")
	if err != nil {
		return nil, fmt.Errorf("temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(htmlContent); err != nil {
		return nil, fmt.Errorf("write html: %w", err)
	}
	tmpFile.Close()

	fileURL := "file://" + tmpFile.Name()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	w, h := cfg.PageDimensions()
	if cfg.IsLandscape() {
		w, h = h, w
	}

	margins := parseMargins(cfg.Page.Margin)

	var pdfBuf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithLandscape(false). // we already swapped w/h
				WithDisplayHeaderFooter(false).
				WithPrintBackground(true).
				WithScale(1.0).
				WithPaperWidth(w).
				WithPaperHeight(h).
				WithMarginTop(margins[0]).
				WithMarginRight(margins[1]).
				WithMarginBottom(margins[2]).
				WithMarginLeft(margins[3]).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp: %w", err)
	}

	return pdfBuf, nil
}

// parseMargins converts margin string to [top, right, bottom, left] in inches.
func parseMargins(m string) [4]float64 {
	m = strings.TrimSpace(m)
	if m == "" {
		m = "20mm"
	}

	parts := strings.Fields(m)
	vals := make([]float64, len(parts))
	for i, p := range parts {
		vals[i] = parseLengthToInches(p)
	}

	switch len(vals) {
	case 1:
		return [4]float64{vals[0], vals[0], vals[0], vals[0]}
	case 2:
		return [4]float64{vals[0], vals[1], vals[0], vals[1]}
	case 3:
		return [4]float64{vals[0], vals[1], vals[2], vals[1]}
	case 4:
		return [4]float64{vals[0], vals[1], vals[2], vals[3]}
	default:
		d := mmToInch(20)
		return [4]float64{d, d, d, d}
	}
}

func parseLengthToInches(s string) float64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "mm") {
		return mmToInch(parseFloat(strings.TrimSuffix(s, "mm")))
	}
	if strings.HasSuffix(s, "cm") {
		return mmToInch(parseFloat(strings.TrimSuffix(s, "cm")) * 10)
	}
	if strings.HasSuffix(s, "in") {
		return parseFloat(strings.TrimSuffix(s, "in"))
	}
	return mmToInch(parseFloat(s))
}

func mmToInch(mm float64) float64 {
	return math.Round(mm/25.4*1000) / 1000
}

func parseFloat(s string) float64 {
	var f float64
	_, _ = fmt.Sscanf(s, "%f", &f)
	return f
}

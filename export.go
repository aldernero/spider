package spider

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

// Save saves the chart to a file, automatically detecting the format from the file extension
// Supports PNG (.png) and SVG (.svg) formats
func (c *Chart) Save(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return c.SavePNG(filename)
	case ".svg":
		return c.SaveSVG(filename)
	default:
		return fmt.Errorf("unsupported file format: %s (supported formats: .png, .svg)", ext)
	}
}

// SavePNG saves the chart as a PNG image
func (c *Chart) SavePNG(filename string) error {
	canvasWidth := c.Width()
	canvasHeight := c.Height()

	// Create canvas
	canv := canvas.New(canvasWidth, canvasHeight)
	ctx := canvas.NewContext(canv)

	// Draw chart
	if err := c.Draw(ctx); err != nil {
		return fmt.Errorf("failed to draw chart: %w", err)
	}

	// Save as PNG
	if err := renderers.Write(filename, canv, canvas.DefaultResolution); err != nil {
		return fmt.Errorf("failed to save PNG: %w", err)
	}

	return nil
}

// SaveSVG saves the chart as an SVG image
func (c *Chart) SaveSVG(filename string) error {
	canvasWidth := c.Width()
	canvasHeight := c.Height()

	// Create canvas
	canv := canvas.New(canvasWidth, canvasHeight)
	ctx := canvas.NewContext(canv)

	// Draw chart
	if err := c.Draw(ctx); err != nil {
		return fmt.Errorf("failed to draw chart: %w", err)
	}

	// Save as SVG
	if err := renderers.Write(filename, canv); err != nil {
		return fmt.Errorf("failed to save SVG: %w", err)
	}

	return nil
}

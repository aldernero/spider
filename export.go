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
	if err := renderers.Write(filename, canv, canvas.DPMM(1)); err != nil {
		return fmt.Errorf("failed to save SVG: %w", err)
	}

	return nil
}

func (c *Chart) SaveDebug() error {
	// Validate chart before drawing
	if err := c.validate(); err != nil {
		return err
	}
	canvasWidth := c.Width()
	canvasHeight := c.Height()

	// Create canvas
	canv := canvas.New(canvasWidth, canvasHeight)
	ctx := canvas.NewContext(canv)

	c.calcRects()

	c.drawBackground(ctx)
	ctx.SetStrokeWidth(1)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetFillColor(canvas.White)
	ctx.DrawPath(0, 0, canvas.Rectangle(canvasWidth, canvasHeight))
	ctx.FillStroke()
	ctx.SetFillColor(canvas.Transparent)
	ctx.SetStrokeColor(canvas.Gray)
	ctx.DrawPath(c.pageRect.X0, c.pageRect.Y0, canvas.Rectangle(c.pageRect.W(), c.pageRect.H()))
	ctx.Stroke()
	ctx.SetStrokeColor(canvas.Red)
	ctx.DrawPath(c.titleRect.X0, c.titleRect.Y0, canvas.Rectangle(c.titleRect.W(), c.titleRect.H()))
	ctx.Stroke()
	ctx.SetStrokeColor(canvas.Green)
	ctx.DrawPath(c.subtitleRect.X0, c.subtitleRect.Y0, canvas.Rectangle(c.subtitleRect.W(), c.subtitleRect.H()))
	ctx.Stroke()
	ctx.SetStrokeColor(canvas.Blue)
	ctx.DrawPath(c.plotRect.X0, c.plotRect.Y0, canvas.Rectangle(c.plotRect.W(), c.plotRect.H()))
	ctx.Stroke()
	ctx.SetStrokeColor(canvas.Yellow)
	ctx.DrawPath(c.legendRect.X0, c.legendRect.Y0, canvas.Rectangle(c.legendRect.W(), c.legendRect.H()))
	ctx.Stroke()

	c.drawTitle(ctx)
	c.drawSubtitle(ctx)
	c.drawPlotBackground(ctx)
	c.drawAxes(ctx)
	c.drawSeries(ctx)
	c.drawLegend(ctx)
	// Save as PNG
	if err := renderers.Write("debug.png", canv, canvas.DefaultResolution); err != nil {
		return fmt.Errorf("failed to save PNG: %w", err)
	}

	return nil
}

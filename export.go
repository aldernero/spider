package spider

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	svgrenderer "github.com/tdewolff/canvas/renderers/svg"
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

	// Render to memory so the paint declarations can be normalized before writing
	var buf bytes.Buffer
	svg := svgrenderer.New(&buf, canvasWidth, canvasHeight, nil)
	canv.RenderTo(svg)
	if err := svg.Close(); err != nil {
		return fmt.Errorf("failed to render SVG: %w", err)
	}

	if err := os.WriteFile(filename, svgToSVG11Colors(buf.Bytes()), 0644); err != nil {
		return fmt.Errorf("failed to save SVG: %w", err)
	}

	return nil
}

// rgbaPaint matches a fill or stroke declaration written with CSS rgba()
// functional notation.
var rgbaPaint = regexp.MustCompile(`(fill|stroke):rgba\((\d+),(\d+),(\d+),([0-9.]+)\)`)

// svgToSVG11Colors rewrites translucent paints into their SVG 1.1 form.
//
// The renderer emits `fill:rgba(r,g,b,a)`, which is CSS Color 4 syntax. Browsers
// accept it, but SVG 1.1 renderers — Inkscape among them — do not, and fall back
// to the initial value, painting every translucent shape solid black. Splitting
// the alpha into a separate fill-opacity/stroke-opacity property is understood
// everywhere and renders identically.
func svgToSVG11Colors(svg []byte) []byte {
	return rgbaPaint.ReplaceAllFunc(svg, func(match []byte) []byte {
		m := rgbaPaint.FindSubmatch(match)
		property, alpha := m[1], m[5]

		rgb := make([]byte, 3)
		for i := range rgb {
			v, err := strconv.Atoi(string(m[i+2]))
			if err != nil || v < 0 || v > 255 {
				return match // not something we understand; leave it alone
			}
			rgb[i] = byte(v)
		}

		return fmt.Appendf(nil, "%s:#%02x%02x%02x;%s-opacity:%s",
			property, rgb[0], rgb[1], rgb[2], property, alpha)
	})
}

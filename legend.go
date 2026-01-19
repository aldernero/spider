package spider

import (
	"github.com/tdewolff/canvas"
)

// LegendConfig represents configuration for the legend
type LegendOptions struct {
	Show             bool            `json:"show" yaml:"show"`                           // Whether to show the legend
	Placement        LegendPlacement `json:"placement" yaml:"placement"`                 // Legend placement
	OutlineThickness float64         `json:"outline_thickness" yaml:"outline_thickness"` // Outline thickness in millimeters
	OutlineColor     Color           `json:"outline_color" yaml:"outline_color"`         // Outline color
	LegendStyle      Font            `json:"style" yaml:"style"`                         // Legend text style
	Padding          float64         `json:"padding,omitempty" yaml:"padding,omitempty"` // Padding around legend items
	ShowOutline      bool            `json:"show_outline" yaml:"show_outline"`           // Whether to show the outline
}

// DefaultLegendOptions returns default legend options
func DefaultLegendOptions() LegendOptions {
	return LegendOptions{
		Show:             true,
		Placement:        LegendPlacementRight,
		OutlineThickness: 0.5,
		OutlineColor:     Color("#000000"),
	}
}

// drawLegend draws the legend on the canvas
func (c *Chart) drawLegend(ctx *canvas.Context) {
	legend := c.Options.LegendOptions
	if !legend.Show || len(c.Data.Series) == 0 {
		return
	}

	// Draw legend background (white) and border
	rect := canvas.Rectangle(c.legendRect.W(), c.legendRect.H())
	ctx.SetFillColor(canvas.Transparent)
	ctx.DrawPath(c.legendRect.X0, c.legendRect.Y0, rect)
	ctx.Fill()

	if legend.OutlineThickness > 0 {
		ctx.SetStrokeColor(legend.OutlineColor.ToCanvasColor())
		ctx.SetStrokeWidth(legend.OutlineThickness)
		ctx.DrawPath(c.legendRect.X0, c.legendRect.Y0, rect)
		ctx.Stroke()
	}
}

// drawLegendItem draws a single legend item
func (c *Chart) drawLegendItem(ctx *canvas.Context, series Series, x, y, width float64) {
	legend := c.Options.LegendOptions
	// Draw sample line and point
	sampleLength := 20.0
	sampleY := y
	sampleX1 := x
	sampleX2 := x + sampleLength

	// Draw line sample
	style := c.Options.SeriesOptions
	ctx.SetStrokeColor(style.LineColor.ToCanvasColor())
	ctx.SetStrokeWidth(style.LineThickness)
	ctx.MoveTo(sampleX1, sampleY)
	ctx.LineTo(sampleX2, sampleY)
	ctx.Stroke()

	// Draw point sample
	if style.ShowPoints && style.PointShape != PointShapeNone {
		pointSize := style.PointSize
		if pointSize == 0 {
			pointSize = DefaultPointSize
		}
		pointColor := style.PointStrokeColor
		if pointColor == "" {
			pointColor = style.LineColor
		}
		drawPoint(ctx, sampleX2, sampleY, pointSize, style.PointShape, pointColor.ToCanvasColor())
	}

	// Draw series name
	textX := sampleX2 + 5.0
	textY := y
	fontSize := legend.LegendStyle.Size
	if fontSize == 0 {
		fontSize = DefaultFontSize
	}
	textColor := legend.LegendStyle.Color
	if textColor == "" {
		textColor = Color("#000000")
	}
	ctx.DrawText(textX, textY, canvas.NewTextLine(c.fonts["legend_label"], series.Name, canvas.Left))
}

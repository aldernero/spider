package spider

import (
	"fmt"
	"math"

	"github.com/tdewolff/canvas"
)

// drawBackground draws the chart background
func (c *Chart) drawBackground(ctx *canvas.Context) {
	// Check if background is transparent before converting (case-insensitive)
	ctx.SetStrokeColor(canvas.Transparent)
	ctx.SetStrokeWidth(0)
	ctx.SetFillColor(c.Options.Background.ToCanvasColor())
	ctx.DrawPath(0, 0, canvas.Rectangle(c.Width(), c.Height()))
	ctx.Fill()
}

// drawTitle draws the chart title
func (c *Chart) drawTitle(ctx *canvas.Context) {
	if !c.Options.ShowTitle || c.Options.Title == "" {
		return
	}
	ctx.DrawText(c.titleRect.X0, c.titleRect.Y1, canvas.NewTextBox(c.fonts["title"], c.Options.Title, c.titleRect.W(), c.titleRect.H(), canvas.Center, canvas.Bottom, nil))
}

// drawSubtitle draws the chart subtitle
func (c *Chart) drawSubtitle(ctx *canvas.Context) {
	if !c.Options.ShowSubtitle || c.Options.Subtitle == "" {
		return
	}
	ctx.DrawText(c.subtitleRect.X0, c.subtitleRect.Y1, canvas.NewTextBox(c.fonts["subtitle"], c.Options.Subtitle, c.subtitleRect.W(), c.subtitleRect.H(), canvas.Center, canvas.Top, nil))
}

// foreground returns the color to use for chart furniture, falling back to the
// chart-wide foreground and finally to black.
func (c *Chart) foreground(col Color) Color {
	if col != "" {
		return col
	}
	if c.Options.Foreground != "" {
		return c.Options.Foreground
	}
	return Color("black")
}

// drawPlotBackground draws the plot background shape (circle or polygon)
func (c *Chart) drawPlotBackground(ctx *canvas.Context) {
	ctx.SetFillColor(canvas.Transparent)
	ctx.SetStrokeColor(c.foreground(c.Options.PlotOptions.OutlineColor).ToCanvasColor())
	ctx.SetStrokeWidth(c.Options.PlotOptions.OutlineThickness)

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	switch c.Options.PlotOptions.ConnectType {
	case ConnectTypeCircle:
		circle := canvas.Circle(c.Radius())
		ctx.DrawPath(centerX, centerY, circle)
		ctx.Stroke()
	case ConnectTypePolygon:
		nAxes := len(c.Data.Axes)
		if nAxes > 0 {
			polygon := canvas.RegularPolygon(nAxes, c.Radius(), true)
			// For odd number of sides, RegularPolygon has a vertex at the top by default
			ctx.DrawPath(centerX, centerY, polygon)
			ctx.Stroke()
		}
	}
}

// drawAxes draws all axes and their labels
func (c *Chart) drawAxes(ctx *canvas.Context, maxima map[string]float64) {
	opts := c.Options.AxisOptions
	if !opts.ShowAxis {
		return
	}

	nAxes := len(c.Data.Axes)

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	radius := c.Radius()

	dt := 360.0 / float64(nAxes)
	// Start at π/2 (top) to align with polygon vertex at top
	// For odd n, RegularPolygon has a vertex at the top by default
	theta := 90.0

	showName := c.Options.ShowAxisNames && opts.ShowName
	showTicks := c.Options.ShowTicks && opts.ShowTicks
	showTickLabels := showTicks && c.Options.ShowTickLabels && opts.ShowTickLabels

	axisColor := c.foreground(opts.LineColor).ToCanvasColor()
	ctx.SetStrokeColor(axisColor)
	labelOffset := opts.LabelOffset
	for _, axis := range c.Data.Axes {
		// Draw axis line using Push/Pop with transformations
		ctx.Push()
		ctx.Translate(centerX, centerY) // Move origin to center
		ctx.Rotate(theta)               // Rotate around the new origin
		ctx.SetStrokeWidth(opts.LineThickness)
		ctx.MoveTo(0, 0)      // Start at origin (center)
		ctx.LineTo(radius, 0) // Draw line along rotated x-axis
		ctx.Stroke()
		// Draw axis name
		if showName {
			ctx.Push()
			ctx.Translate(radius+labelOffset, 0)
			if theta > 180 && theta < 360 {
				ctx.Rotate(90)
			} else {
				ctx.Rotate(-90)
			}
			ctx.DrawText(0, 0, canvas.NewTextLine(c.fonts["axis_label"], axis.Name, canvas.Center))
			ctx.Pop()
		}
		if showTicks {
			max := maxima[axis.Name]
			// major ticks
			majorTicks := linspace(0, radius, opts.MajorTicks+2)
			majorLen := opts.MajorTickLength / 2
			minorLen := opts.MinorTickLength / 2
			for i := 1; i < len(majorTicks)-1; i++ {
				ctx.SetStrokeWidth(opts.MajorTickLineThickness)
				ctx.MoveTo(majorTicks[i], -majorLen)
				ctx.LineTo(majorTicks[i], majorLen)
				ctx.Stroke()
				// tick labels
				if showTickLabels {
					label := val2String(linmap(0, radius, 0, max, majorTicks[i]))
					ctx.Push()
					ctx.Translate(majorTicks[i], -majorLen-opts.LabelOffset)
					ctx.Rotate(-theta)
					ctx.DrawText(0, 0, canvas.NewTextLine(c.fonts["tick_label"], label, canvas.Center))
					ctx.Pop()
				}
			}
			// minor ticks, spaced within each major interval
			ctx.SetStrokeWidth(opts.MinorTickLineThickness)
			for i := 0; i+1 < len(majorTicks); i++ {
				step := (majorTicks[i+1] - majorTicks[i]) / float64(opts.MinorTicks+1)
				for j := 1; j <= opts.MinorTicks; j++ {
					x := majorTicks[i] + float64(j)*step
					ctx.MoveTo(x, -minorLen)
					ctx.LineTo(x, minorLen)
					ctx.Stroke()
				}
			}
		}
		ctx.Pop() // Restore previous transformation state
		theta += dt
	}
}

// resolveSeriesOptions returns the options for a series with the chart-level
// palette and defaults filled in wherever the series does not override them.
func (c *Chart) resolveSeriesOptions(i int) SeriesOptions {
	opts := c.Data.Series[i].Options
	seriesColor := c.Options.Colors[i%len(c.Options.Colors)]

	// An unset color is not "transparent": it means "use the palette". Leaving
	// FillColor empty here would render it as black at FillOpacity.
	for _, col := range []*Color{&opts.LineColor, &opts.FillColor, &opts.PointStrokeColor, &opts.PointFillColor} {
		if *col == "" {
			*col = seriesColor
		}
	}
	if opts.PointShape == "" {
		opts.PointShape = c.Options.PointMarkers[i%len(c.Options.PointMarkers)]
	}
	if opts.PointSize == 0 {
		opts.PointSize = DefaultPointSize
	}
	if opts.PointLineThickness == 0 {
		opts.PointLineThickness = DefaultSeriesLineThickness
	}
	if opts.LineThickness == 0 {
		opts.LineThickness = DefaultSeriesLineThickness
	}
	return opts
}

// drawSeries draws all series on the chart
func (c *Chart) drawSeries(ctx *canvas.Context, maxima map[string]float64) {
	nAxes := len(c.Data.Axes)

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	radius := c.Radius()

	// Calculate starting angle (same as in drawAxes) - start at top (π/2)
	startTheta := math.Pi / 2
	dt := Tau / float64(nAxes)

	points := make([]canvas.Point, nAxes)
	for i := range c.Data.Series {
		series := &c.Data.Series[i]
		seriesOpts := c.resolveSeriesOptions(i)
		// Calculate points for this series
		theta := startTheta
		for j, axis := range c.Data.Axes {
			value := series.GetDataValue(axis.Name)
			scaledRadius := linmap(0, maxima[axis.Name], 0, radius, value)
			points[j].X = centerX + scaledRadius*math.Cos(theta)
			points[j].Y = centerY + scaledRadius*math.Sin(theta)
			theta += dt
		}
		// draw series line
		ctx.SetFillColor(seriesOpts.FillColor.ToCanvasColorWithOpacity(seriesOpts.FillOpacity))
		ctx.SetStrokeColor(seriesOpts.LineColor.ToCanvasColor())
		ctx.SetStrokeWidth(seriesOpts.LineThickness)
		ctx.MoveTo(points[0].X, points[0].Y)
		for j := 1; j < len(points); j++ {
			ctx.LineTo(points[j].X, points[j].Y)
		}
		ctx.Close()
		ctx.FillStroke()
		// draw series points
		if c.Options.ShowPointMarkers {
			for _, point := range points {
				c.drawSeriesPoint(ctx, point, seriesOpts)
			}
		}
	}
}

// drawSeriesPoint draws a point for a series
func (c *Chart) drawSeriesPoint(ctx *canvas.Context, point canvas.Point, seriesOpts SeriesOptions) {
	ctx.SetFillColor(seriesOpts.PointFillColor.ToCanvasColorWithOpacity(seriesOpts.PointFillOpacity))
	ctx.SetStrokeColor(seriesOpts.PointStrokeColor.ToCanvasColor())
	ctx.SetStrokeWidth(seriesOpts.PointLineThickness)

	switch seriesOpts.PointShape {
	case PointShapeCircle:
		circle := canvas.Circle(seriesOpts.PointSize / 2)
		ctx.DrawPath(point.X, point.Y, circle)
		ctx.FillStroke()
	case PointShapeSquare:
		rect := canvas.Rectangle(seriesOpts.PointSize, seriesOpts.PointSize)
		ctx.DrawPath(point.X-seriesOpts.PointSize/2, point.Y-seriesOpts.PointSize/2, rect)
		ctx.FillStroke()
	case PointShapeTriangle:
		triangle := canvas.RegularPolygon(3, seriesOpts.PointSize/2, true)
		ctx.DrawPath(point.X, point.Y, triangle)
		ctx.FillStroke()
	case PointShapeDiamond:
		diamond := canvas.RegularPolygon(4, seriesOpts.PointSize/2, true)
		ctx.DrawPath(point.X, point.Y, diamond)
		ctx.FillStroke()
	}
}

// drawLegend draws the legend on the canvas
func (c *Chart) drawLegend(ctx *canvas.Context) {
	legend := c.Options.LegendOptions
	if !c.Options.ShowLegend || !legend.Show || legend.Placement == LegendPlacementNone || len(c.Data.Series) == 0 {
		return
	}

	var legendHorizontalTextAlignment canvas.TextAlign
	var legendVerticalTextAlignment canvas.TextAlign
	switch legend.Placement {
	case LegendPlacementTop:
		legendHorizontalTextAlignment = canvas.Center
		legendVerticalTextAlignment = canvas.Bottom
	case LegendPlacementBottom:
		legendHorizontalTextAlignment = canvas.Center
		legendVerticalTextAlignment = canvas.Top
	case LegendPlacementLeft:
		legendHorizontalTextAlignment = canvas.Right
		legendVerticalTextAlignment = canvas.Center
	case LegendPlacementRight:
		legendHorizontalTextAlignment = canvas.Left
		legendVerticalTextAlignment = canvas.Center
	default:
		legendHorizontalTextAlignment = canvas.Center
		legendVerticalTextAlignment = canvas.Center
	}

	// Draw legend items
	rt := canvas.NewRichText(c.fonts["legend_label"])
	textOptions := &canvas.TextOptions{
		Linebreaker: canvas.KnuthLinebreaker{},
	}
	w := c.legendRect.W()
	dw := 0.0
	spaceCanvas, spaceWidth := c.canvasString(" ")
	for i, series := range c.Data.Series {
		cnvs, width := c.drawLegendSeriesPath(i)
		dw += width
		rt.WriteCanvas(cnvs, canvas.FontMiddle)
		cnvs, width = c.canvasString(series.Name)
		dw += width
		rt.WriteCanvas(cnvs, canvas.FontMiddle)
		if (legend.Placement == LegendPlacementRight || legend.Placement == LegendPlacementLeft) && i < len(c.Data.Series)-1 {
			rt.WriteString("\n")
			dw = 0.0
		} else {
			dw += spaceWidth
			rt.WriteCanvas(spaceCanvas, canvas.FontMiddle)
		}
		if dw > 0.85*w {
			rt.WriteString("\n")
			dw = 0.0
		}
	}
	text := rt.ToText(c.legendRect.W(), c.legendRect.H(), legendHorizontalTextAlignment, legendVerticalTextAlignment, textOptions)
	ctx.DrawText(c.legendRect.X0, c.legendRect.Y1, text)
}

func (c *Chart) canvasString(s string) (*canvas.Canvas, float64) {
	cnvs := canvas.New(10, 10)
	ctx := canvas.NewContext(cnvs)
	ctx.DrawText(0, 0, canvas.NewTextLine(c.fonts["legend_label"], s, canvas.Center))
	cnvs.Fit(0)
	width, _ := cnvs.Size()
	return cnvs, width
}

func (c *Chart) drawLegendSeriesPath(seriesIndex int) (*canvas.Canvas, float64) {
	seriesOpts := c.resolveSeriesOptions(seriesIndex)
	cnvs := canvas.New(10, 10)
	ctx := canvas.NewContext(cnvs)
	ctx.SetStrokeColor(seriesOpts.LineColor.ToCanvasColor())
	ctx.SetStrokeWidth(c.Options.LegendOptions.LineThickness)
	ctx.MoveTo(0, 0)
	ctx.LineTo(c.Options.LegendOptions.LineLength, 0)
	ctx.Stroke()
	if c.Options.ShowPointMarkers {
		center := canvas.Point{X: c.Options.LegendOptions.LineLength / 2, Y: 0}
		c.drawSeriesPoint(ctx, center, seriesOpts)
	}
	cnvs.Fit(0)
	width, _ := cnvs.Size()
	return cnvs, width
}

// Linspace creates a slice of linearly distributed values in a range, inclusive of the end value
func linspace(i float64, j float64, n int) []float64 {
	// Fewer than two points has no well-defined spacing, and n-1 would divide by zero
	if n < 1 {
		return nil
	}
	if n == 1 {
		return []float64{i}
	}
	result := make([]float64, n)
	d := (j - i) / float64(n-1)
	for k := range result {
		result[k] = i + float64(k)*d
	}
	return result
}

// clamp constrains v to [min, max]
func clamp(v, min, max float64) float64 {
	return math.Min(math.Max(v, min), max)
}

// Lerp calculates the linear interpolation between two values
func lerp(a float64, b float64, i float64) float64 {
	return a + i*(b-a)
}

// Map calculates the linear interpolation from one range to another
func linmap(a float64, b float64, c float64, d float64, i float64) float64 {
	p := (i - a) / (b - a)
	return lerp(c, d, p)
}

func val2String(val float64) string {
	if val < 1.0 {
		return fmt.Sprintf("%.3f", val)
	} else if val < 10.0 {
		return fmt.Sprintf("%.2f", val)
	} else if val < 100.0 {
		return fmt.Sprintf("%.1f", val)
	} else if val < 1000.0 {
		return fmt.Sprintf("%.0f", val)
	} else if val < 1000000.0 {
		return fmt.Sprintf("%.0fk", val/1000.0)
	} else if val < 1000000000.0 {
		return fmt.Sprintf("%.0fM", val/1000000.0)
	} else if val < 1000000000000.0 {
		return fmt.Sprintf("%.0fG", val/1000000000.0)
	} else if val < 1000000000000000.0 {
		return fmt.Sprintf("%.0fT", val/1000000000000.0)
	} else {
		return fmt.Sprintf("%.0fP", val/1000000000000000.0)
	}
}

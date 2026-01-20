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
	if c.Options.Title == "" {
		return
	}

	style := c.Options.TitleStyle
	if style.Size == 0 {
		style.Size = DefaultTitleFontSize
	}
	if style.Color == "" {
		style.Color = Color("#000000")
	}

	ctx.DrawText(c.titleRect.X0, c.titleRect.Y1, canvas.NewTextBox(c.fonts["title"], c.Options.Title, c.titleRect.W(), c.titleRect.H(), canvas.Center, canvas.Bottom, nil))
}

// drawSubtitle draws the chart subtitle
func (c *Chart) drawSubtitle(ctx *canvas.Context) {
	if c.Options.Subtitle == "" {
		return
	}

	style := c.Options.SubtitleStyle
	if style.Size == 0 {
		style.Size = DefaultSubtitleFontSize
	}
	if style.Color == "" {
		style.Color = Color("#000000")
	}

	ctx.DrawText(c.subtitleRect.X0, c.subtitleRect.Y1, canvas.NewTextBox(c.fonts["subtitle"], c.Options.Subtitle, c.subtitleRect.W(), c.subtitleRect.H(), canvas.Center, canvas.Top, nil))
}

// drawPlotBackground draws the plot background shape (circle or polygon)
func (c *Chart) drawPlotBackground(ctx *canvas.Context) {
	ctx.SetFillColor(canvas.Transparent)
	ctx.SetStrokeColor(canvas.Black)
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
func (c *Chart) drawAxes(ctx *canvas.Context) {
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(DefaultAxisLineThickness)

	nAxes := len(c.Data.Axes)

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	radius := c.Radius()

	dt := 360.0 / float64(nAxes)
	// Start at π/2 (top) to align with polygon vertex at top
	// For odd n, RegularPolygon has a vertex at the top by default
	theta := 90.0

	// Get all series data for max calculation
	seriesData := getAllSeriesData(c.Data.Series)
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(c.Options.AxisOptions.LineThickness)
	labelOffset := c.Options.AxisOptions.LabelOffset
	for _, axis := range c.Data.Axes {
		// Draw axis line using Push/Pop with transformations
		ctx.Push()
		ctx.Translate(centerX, centerY) // Move origin to center
		ctx.Rotate(theta)               // Rotate around the new origin
		ctx.MoveTo(0, 0)                // Start at origin (center)
		ctx.LineTo(radius, 0)           // Draw line along rotated x-axis
		ctx.Stroke()
		// Draw axis name
		if c.Options.ShowAxisNames {
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
		// Draw ticks
		max := axis.GetMax(seriesData)
		// major ticks
		majorTicks := linspace(0, radius, c.Options.AxisOptions.MajorTicks+2)
		ctx.SetStrokeWidth(c.Options.AxisOptions.MajorTickLineThickness)
		for i := 1; i < len(majorTicks)-1; i++ {
			l := c.Options.AxisOptions.MajorTickLength / 2
			o := c.Options.AxisOptions.LabelOffset
			val := linmap(0, radius, 0, max, majorTicks[i])
			label := val2String(val)
			ctx.SetStrokeWidth(c.Options.AxisOptions.MajorTickLineThickness)
			ctx.MoveTo(majorTicks[i], -l)
			ctx.LineTo(majorTicks[i], l)
			ctx.Stroke()
			// tick labels
			if c.Options.ShowTickLabels {
				ctx.Push()
				ctx.Translate(majorTicks[i], -l-o)
				ctx.Rotate(-theta)
				ctx.DrawText(0, 0, canvas.NewTextLine(c.fonts["tick_label"], label, canvas.Center))
				ctx.Pop()
			}
		}
		// minor ticks
		ctx.SetStrokeWidth(c.Options.AxisOptions.MinorTickLineThickness)
		for i := 0; i < len(majorTicks)-1; i++ {
			minorTicks := linspace(majorTicks[i], majorTicks[i+1], c.Options.AxisOptions.MinorTicks+2)
			for j := 1; j < len(minorTicks)-1; j++ {
				ctx.MoveTo(minorTicks[j], -c.Options.AxisOptions.MinorTickLength/2)
				ctx.LineTo(minorTicks[j], c.Options.AxisOptions.MinorTickLength/2)
				ctx.Stroke()
			}
		}
		ctx.Pop() // Restore previous transformation state
		theta += dt
	}
}

// drawSeries draws all series on the chart
func (c *Chart) drawSeries(ctx *canvas.Context) {
	nAxes := len(c.Data.Axes)

	seriesData := getAllSeriesData(c.Data.Series)

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	radius := c.Radius()

	// Calculate starting angle (same as in drawAxes) - start at top (π/2)
	startTheta := math.Pi / 2
	dt := Tau / float64(nAxes)

	nColors := len(c.Options.Colors)
	nPointMarkers := len(c.Options.PointMarkers)
	for i := range c.Data.Series {
		series := &c.Data.Series[i]
		// Apply default colors and point markers if not set
		seriesOpts := SeriesOptions{
			LineColor:          series.Options.LineColor,
			FillColor:          series.Options.FillColor,
			PointStrokeColor:   series.Options.PointStrokeColor,
			PointFillColor:     series.Options.PointFillColor,
			PointFillOpacity:   series.Options.PointFillOpacity,
			PointShape:         series.Options.PointShape,
			PointSize:          series.Options.PointSize,
			PointLineThickness: series.Options.PointLineThickness,
			LineThickness:      series.Options.LineThickness,
			FillOpacity:        series.Options.FillOpacity,
		}
		if seriesOpts.LineColor == "" {
			seriesOpts.LineColor = c.Options.Colors[i%nColors]
		}
		if seriesOpts.PointStrokeColor == "" {
			seriesOpts.PointStrokeColor = c.Options.Colors[i%nColors]
		}
		if seriesOpts.PointFillColor == "" {
			seriesOpts.PointFillColor = c.Options.Colors[i%nColors]
		}
		if series.Options.LineColor == "" {
			seriesOpts.LineColor = c.Options.Colors[i%nColors]
		}
		if series.Options.PointStrokeColor == "" {
			seriesOpts.PointStrokeColor = c.Options.Colors[i%nColors]
		}
		if series.Options.PointFillColor == "" {
			seriesOpts.PointFillColor = c.Options.Colors[i%nColors]
		}
		if series.Options.PointShape == "" {
			seriesOpts.PointShape = c.Options.PointMarkers[i%nPointMarkers]
		}
		if series.Options.PointSize == 0 {
			seriesOpts.PointSize = DefaultPointSize
		}
		if series.Options.PointLineThickness == 0 {
			seriesOpts.PointLineThickness = DefaultSeriesLineThickness
		}
		// Calculate points for this series
		points := make([]canvas.Point, nAxes)
		theta := startTheta
		for j, axis := range c.Data.Axes {
			max := axis.GetMax(seriesData)
			value := series.GetDataValue(axis.Name)
			scaledRadius := linmap(0, max, 0, radius, value)
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
	if !legend.Show || len(c.Data.Series) == 0 {
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
			cnvs, width := c.canvasString(" ")
			dw += width
			rt.WriteCanvas(cnvs, canvas.FontMiddle)
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
	seriesOpts := SeriesOptions{
		LineColor:          c.Data.Series[seriesIndex].Options.LineColor,
		FillColor:          c.Data.Series[seriesIndex].Options.FillColor,
		PointStrokeColor:   c.Data.Series[seriesIndex].Options.PointStrokeColor,
		PointFillColor:     c.Data.Series[seriesIndex].Options.PointFillColor,
		PointFillOpacity:   c.Data.Series[seriesIndex].Options.PointFillOpacity,
		PointShape:         c.Data.Series[seriesIndex].Options.PointShape,
		PointSize:          c.Data.Series[seriesIndex].Options.PointSize,
		PointLineThickness: c.Data.Series[seriesIndex].Options.PointLineThickness,
		LineThickness:      c.Data.Series[seriesIndex].Options.LineThickness,
		FillOpacity:        c.Data.Series[seriesIndex].Options.FillOpacity,
	}
	nColors := len(c.Options.Colors)
	nPointMarkers := len(c.Options.PointMarkers)
	if seriesOpts.LineColor == "" {
		seriesOpts.LineColor = c.Options.Colors[seriesIndex%nColors]
	}
	if seriesOpts.PointStrokeColor == "" {
		seriesOpts.PointStrokeColor = c.Options.Colors[seriesIndex%nColors]
	}
	if seriesOpts.PointFillColor == "" {
		seriesOpts.PointFillColor = c.Options.Colors[seriesIndex%nColors]
	}
	if seriesOpts.LineColor == "" {
		seriesOpts.LineColor = c.Options.Colors[seriesIndex%nColors]
	}
	if seriesOpts.PointStrokeColor == "" {
		seriesOpts.PointStrokeColor = c.Options.Colors[seriesIndex%nColors]
	}
	if seriesOpts.PointFillColor == "" {
		seriesOpts.PointFillColor = c.Options.Colors[seriesIndex%nColors]
	}
	if seriesOpts.PointShape == "" {
		seriesOpts.PointShape = c.Options.PointMarkers[seriesIndex%nPointMarkers]
	}
	if seriesOpts.PointSize == 0 {
		seriesOpts.PointSize = DefaultPointSize
	}
	if seriesOpts.PointLineThickness == 0 {
		seriesOpts.PointLineThickness = DefaultSeriesLineThickness
	}
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
	var result []float64
	d := (j - i) / float64(n-1)
	for k := 0; k < n; k++ {
		result = append(result, i+float64(k)*d)
	}
	return result
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

package spider

import (
	"fmt"
	"image/color"
	"math"

	"github.com/tdewolff/canvas"
)

// drawBackground draws the chart background
func (c *Chart) drawBackground(ctx *canvas.Context) {
	// Check if background is transparent before converting (case-insensitive)
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
	ctx.SetStrokeWidth(0.5)

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
	fmt.Println("drawing axes", len(c.Data.Axes))
	ctx.SetStrokeColor(canvas.Black)
	ctx.SetStrokeWidth(DefaultAxisLineThickness)

	nAxes := len(c.Data.Axes)
	if nAxes == 0 {
		return
	}

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

	dt := Tau / float64(nAxes)
	seriesData := getAllSeriesData(c.Data.Series)

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	radius := c.Radius()

	// Calculate starting angle (same as in drawAxes) - start at top (π/2)
	startTheta := math.Pi / 2

	for _, series := range c.Data.Series {
		// Calculate points for this series
		points := make([]canvas.Point, nAxes)
		theta := startTheta
		for i, axis := range c.Data.Axes {
			max := axis.GetMax(seriesData)
			value := series.GetDataValue(axis.Name)
			scaledRadius := linmap(0, radius, 0, max, value)
			points[i].X = centerX + scaledRadius*math.Cos(theta)
			points[i].Y = centerY + scaledRadius*math.Sin(theta)
			theta += dt
		}
		//c.drawSingleSeries(ctx, series, dt, startTheta, seriesData)
	}
}

// drawSingleSeries draws a single series
func (c *Chart) drawSingleSeries(ctx *canvas.Context, series Series, dt float64, startTheta float64, seriesData []map[string]float64) {
	nAxes := len(c.Data.Axes)
	if nAxes == 0 {
		return
	}

	centerX := c.plotRect.X0 + c.plotRect.W()/2
	centerY := c.plotRect.Y0 + c.plotRect.H()/2
	radius := c.Radius()

	// Get series style with defaults
	style := series.Options
	if style.LineThickness == 0 {
		style.LineThickness = DefaultSeriesLineThickness
	}
	if style.LineColor == "" {
		style.LineColor = Color("#000000")
	}

	// Calculate points for this series
	points := make([]canvas.Point, nAxes)
	theta := startTheta

	for i, axis := range c.Data.Axes {
		max := axis.GetMax(seriesData)
		value := series.GetDataValue(axis.Name)
		scaledRadius := radius * axis.ScaleValue(value, max)

		points[i].x = centerX + scaledRadius*math.Cos(theta)
		points[i].y = centerY + scaledRadius*math.Sin(theta)
		theta += dt
	}

	// Draw fill (if opacity > 0)
	if style.FillOpacity > 0 {
		c.drawSeriesFill(ctx, points, style.FillColor, style.FillOpacity)
	}

	// Draw line
	drawSeriesLine(ctx, points, style.Line)

	// Draw points
	if style.ShowPoints && style.PointShape != PointShapeNone {
		for _, point := range points {
			// Skip drawing points at the center (check if distance from center is very small)
			dx := point.x - centerX
			dy := point.y - centerY
			distance := math.Sqrt(dx*dx + dy*dy)
			if distance < 0.1 { // Very close to center (0.1mm threshold)
				continue
			}
			drawPoint(ctx, point.x, point.y, style.PointSize, style.PointShape, style.PointStrokeColor.ToCanvasColor())
		}
	}
}

// drawSeriesFill draws the filled area for a series
func (c *Chart) drawSeriesFill(ctx *canvas.Context, points []struct{ x, y float64 }, fill FillStyle) {
	col := style.FillColor.ToCanvasColor()
	// Convert to RGBA to apply opacity
	if rgba, ok := col.(color.RGBA); ok {
		rgba.A = uint8(float64(rgba.A) * style.FillOpacity)
		ctx.SetFillColor(rgba)
	} else {
		ctx.SetFillColor(col)
	}
	ctx.Fill()
}

// drawSeriesLine draws the line for a series
func drawSeriesLine(ctx *canvas.Context, points []struct{ x, y float64 }, line LineStyle) {
	if len(points) == 0 {
		return
	}

	ctx.SetStrokeColor(line.Color.ToCanvasColor())
	ctx.SetStrokeWidth(line.Thickness)

	ctx.MoveTo(points[0].x, points[0].y)
	for i := 1; i < len(points); i++ {
		ctx.LineTo(points[i].x, points[i].y)
	}
	ctx.Close() // Connect back to start
	ctx.Stroke()
}

// drawPoint draws a point with the specified shape
func drawPoint(ctx *canvas.Context, x, y, size float64, shape PointShape, col color.Color) {
	ctx.SetFillColor(col)
	ctx.SetStrokeColor(col)
	ctx.SetStrokeWidth(1)

	switch shape {
	case PointShapeCircle:
		circle := canvas.Circle(size / 2)
		ctx.DrawPath(x, y, circle)
		ctx.Fill()
	case PointShapeSquare:
		rect := canvas.Rectangle(size, size)
		ctx.DrawPath(x-size/2, y-size/2, rect)
		ctx.Fill()
	case PointShapeTriangle:
		triangle := canvas.RegularPolygon(3, size/2, true)
		ctx.DrawPath(x, y, triangle)
		ctx.Fill()
	case PointShapeDiamond:
		diamond := canvas.RegularPolygon(4, size/2, true)
		ctx.DrawPath(x, y, diamond)
		ctx.Fill()
	}
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

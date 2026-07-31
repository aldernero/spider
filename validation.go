package spider

import (
	"fmt"
	"math"

	"github.com/tdewolff/canvas"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// validateColors reports unparseable color strings. Without this an unknown color
// name renders as transparent and a malformed hex value renders as black, so a
// typo silently produces a wrong chart instead of an error.
func (c *Chart) validateColors() error {
	check := func(field string, col Color) error {
		if err := col.Validate(); err != nil {
			return &ValidationError{Field: field, Message: err.Error()}
		}
		return nil
	}

	colors := map[string]Color{
		"options.background":                   c.Options.Background,
		"options.foreground":                   c.Options.Foreground,
		"options.title_style.color":            c.Options.TitleStyle.Color,
		"options.subtitle_style.color":         c.Options.SubtitleStyle.Color,
		"options.plot_options.outline_color":   c.Options.PlotOptions.OutlineColor,
		"options.axis_options.line_color":      c.Options.AxisOptions.LineColor,
		"options.axis_options.label_style":     c.Options.AxisOptions.LabelStyle.Color,
		"options.axis_options.tick_label":      c.Options.AxisOptions.TickLabelStyle.Color,
		"options.legend_options.outline_color": c.Options.LegendOptions.OutlineColor,
		"options.legend_options.style.color":   c.Options.LegendOptions.LegendStyle.Color,
	}
	for field, col := range colors {
		if err := check(field, col); err != nil {
			return err
		}
	}

	for i, col := range c.Options.Colors {
		if err := check(fmt.Sprintf("options.colors[%d]", i), col); err != nil {
			return err
		}
	}

	for _, series := range c.Data.Series {
		field := "series " + series.Name
		for name, col := range map[string]Color{
			"line_color":       series.Options.LineColor,
			"fill_color":       series.Options.FillColor,
			"point_color":      series.Options.PointStrokeColor,
			"point_fill_color": series.Options.PointFillColor,
		} {
			if err := check(field+"."+name, col); err != nil {
				return err
			}
		}
	}

	return nil
}

// chartFonts names each font face a chart needs, mapped to the style it comes from.
func (c *Chart) chartFonts() map[string]*Font {
	return map[string]*Font{
		"title":        &c.Options.TitleStyle,
		"subtitle":     &c.Options.SubtitleStyle,
		"axis_label":   &c.Options.AxisOptions.LabelStyle,
		"tick_label":   &c.Options.AxisOptions.TickLabelStyle,
		"legend_label": &c.Options.LegendOptions.LegendStyle,
	}
}

// ValidateChart validates a chart configuration
func (c *Chart) validate() error {
	// Validate all fonts load correctly
	c.fonts = make(map[string]*canvas.FontFace)
	for name, style := range c.chartFonts() {
		// Resolve the text color here rather than in the style: an unset color
		// would otherwise convert to transparent and render invisible text.
		resolved := *style
		resolved.Color = c.foreground(style.Color)
		face, err := resolved.loadFontFace(c.Options.DefaultFontName, c.Options.DefaultFontPath)
		if err != nil {
			return fmt.Errorf("failed to load %s font: %w", name, err)
		}
		if face == nil {
			return &ValidationError{
				Field:   "fonts",
				Message: fmt.Sprintf("failed to load %s font", name),
			}
		}
		c.fonts[name] = face
	}

	if err := c.validateColors(); err != nil {
		return err
	}

	// Validate axes count
	if len(c.Data.Axes) < 3 {
		return &ValidationError{
			Field:   "axes",
			Message: "at least 3 axes are required",
		}
	}
	if len(c.Data.Axes) > MaxAxes {
		return &ValidationError{
			Field:   "axes",
			Message: fmt.Sprintf("maximum %d axes allowed, got %d", MaxAxes, len(c.Data.Axes)),
		}
	}

	// Validate series count
	if len(c.Data.Series) > MaxSeries {
		return &ValidationError{
			Field:   "series",
			Message: fmt.Sprintf("maximum %d series allowed, got %d", MaxSeries, len(c.Data.Series)),
		}
	}

	// Collect axis names
	axisNames := make([]string, len(c.Data.Axes))
	axisNameMap := make(map[string]bool)
	for i, axis := range c.Data.Axes {
		if axis.Name == "" {
			return &ValidationError{
				Field:   "axes",
				Message: fmt.Sprintf("axis at index %d has no name", i),
			}
		}
		if axisNameMap[axis.Name] {
			return &ValidationError{
				Field:   "axes",
				Message: fmt.Sprintf("duplicate axis name: %s", axis.Name),
			}
		}
		axisNames[i] = axis.Name
		axisNameMap[axis.Name] = true
	}

	// Validate series data
	for i, series := range c.Data.Series {
		if series.Name == "" {
			return &ValidationError{
				Field:   "series",
				Message: fmt.Sprintf("series at index %d has no name", i),
			}
		}
		if err := series.ValidateData(axisNames); err != nil {
			return fmt.Errorf("series %s: %w", series.Name, err)
		}
	}

	// Validate chart options
	if c.Options.Width <= 0 {
		return &ValidationError{
			Field:   "options.width",
			Message: "width must be positive",
		}
	}
	if c.Options.Height <= 0 {
		return &ValidationError{
			Field:   "options.height",
			Message: "height must be positive",
		}
	}
	if c.Options.PlotOptions.Scale <= 0 || c.Options.PlotOptions.Scale > 1.0 {
		return &ValidationError{
			Field:   "options.plot_scale",
			Message: "plot_scale must be between 0 and 1",
		}
	}

	// The series palettes are indexed modulo their length, so an empty one is a
	// division by zero rather than a missing style.
	if len(c.Options.Colors) == 0 {
		return &ValidationError{
			Field:   "options.colors",
			Message: "at least one series color is required",
		}
	}
	if len(c.Options.PointMarkers) == 0 {
		return &ValidationError{
			Field:   "options.point_markers",
			Message: "at least one point marker is required",
		}
	}

	// Negative tick counts silently blank the axis rather than erroring
	if c.Options.AxisOptions.MajorTicks < 0 {
		return &ValidationError{
			Field:   "options.axis_options.major_ticks",
			Message: fmt.Sprintf("major_ticks must not be negative, got %d", c.Options.AxisOptions.MajorTicks),
		}
	}
	if c.Options.AxisOptions.MinorTicks < 0 {
		return &ValidationError{
			Field:   "options.axis_options.minor_ticks",
			Message: fmt.Sprintf("minor_ticks must not be negative, got %d", c.Options.AxisOptions.MinorTicks),
		}
	}

	// Validate legend options
	if c.Options.LegendOptions.MinWidth <= 0 {
		c.Options.LegendOptions.MinWidth = 0
	}
	if c.Options.LegendOptions.MaxWidth <= 0 {
		c.Options.LegendOptions.MaxWidth = math.Inf(1)
	}
	if c.Options.LegendOptions.MinHeight <= 0 {
		c.Options.LegendOptions.MinHeight = 0
	}
	if c.Options.LegendOptions.MaxHeight <= 0 {
		c.Options.LegendOptions.MaxHeight = math.Inf(1)
	}
	if c.Options.LegendOptions.MinWidth > c.Options.LegendOptions.MaxWidth {
		return &ValidationError{
			Field:   "options.legend_options.min_width",
			Message: "min_width must be less than max_width",
		}
	}
	if c.Options.LegendOptions.MinHeight > c.Options.LegendOptions.MaxHeight {
		return &ValidationError{
			Field:   "options.legend_options.min_height",
			Message: "min_height must be less than max_height",
		}
	}

	return nil
}

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

// ValidateChart validates a chart configuration
func (c *Chart) validate() error {
	// Validate all fonts load correctly
	c.fonts = make(map[string]*canvas.FontFace)
	face, err := c.Options.TitleStyle.loadFontFace(c.Options.DefaultFontName, c.Options.DefaultFontPath)
	if err != nil || face == nil {
		return fmt.Errorf("failed to load title style font: %w", err)
	}
	c.fonts["title"] = face
	face, err = c.Options.SubtitleStyle.loadFontFace(c.Options.DefaultFontName, c.Options.DefaultFontPath)
	if err != nil || face == nil {
		return fmt.Errorf("failed to load subtitle style font: %w", err)
	}
	c.fonts["subtitle"] = face
	face, err = c.Options.AxisOptions.LabelStyle.loadFontFace(c.Options.DefaultFontName, c.Options.DefaultFontPath)
	if err != nil || face == nil {
		return fmt.Errorf("failed to load axis label style font: %w", err)
	}
	c.fonts["axis_label"] = face
	face, err = c.Options.AxisOptions.TickLabelStyle.loadFontFace(c.Options.DefaultFontName, c.Options.DefaultFontPath)
	if err != nil || face == nil {
		return fmt.Errorf("failed to load tick label style font: %w", err)
	}
	c.fonts["tick_label"] = face
	face, err = c.Options.LegendOptions.LegendStyle.loadFontFace(c.Options.DefaultFontName, c.Options.DefaultFontPath)
	if err != nil || face == nil {
		return fmt.Errorf("failed to load legend label style font: %w", err)
	}
	c.fonts["legend_label"] = face
	if len(c.fonts) != 5 {
		return &ValidationError{
			Field:   "fonts",
			Message: fmt.Sprintf("expected 5 fonts, got %d", len(c.fonts)),
		}
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

	return nil
}

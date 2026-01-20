package spider

import (
	"fmt"
	"math"

	"github.com/tdewolff/canvas"
)

// ChartOptions represents options for the overall chart
type ChartOptions struct {
	Width            float64       `json:"width" yaml:"width"`                                       // Chart width
	Height           float64       `json:"height" yaml:"height"`                                     // Chart height
	Background       Color         `json:"background,omitempty" yaml:"background,omitempty"`         // Background color
	Foreground       Color         `json:"foreground,omitempty" yaml:"foreground,omitempty"`         // Foreground color
	Title            string        `json:"title,omitempty" yaml:"title,omitempty"`                   // Chart title
	TitleStyle       Font          `json:"title_style,omitempty" yaml:"title_style,omitempty"`       // Title font style
	TitleMargin      float64       `json:"title_margin" yaml:"title_margin"`                         // Title margin in millimeters
	Subtitle         string        `json:"subtitle,omitempty" yaml:"subtitle,omitempty"`             // Chart subtitle
	SubtitleStyle    Font          `json:"subtitle_style,omitempty" yaml:"subtitle_style,omitempty"` // Subtitle font style
	SubtitleMargin   float64       `json:"subtitle_margin" yaml:"subtitle_margin"`                   // Subtitle margin in millimeters
	PlotOptions      PlotOptions   `json:"plot_options" yaml:"plot_options"`                         // Plot options
	AxisOptions      AxisOptions   `json:"axis_options" yaml:"axis_options"`                         // Axis options
	SeriesOptions    SeriesOptions `json:"series_options" yaml:"series_options"`                     // Series options
	LegendOptions    LegendOptions `json:"legend_options" yaml:"legend_options"`                     // Legend options
	Colors           []Color       `json:"colors" yaml:"colors"`                                     // Colors for the series
	PointMarkers     []PointShape  `json:"point_markers" yaml:"point_markers"`                       // Point markers for the series
	PageMargin       float64       `json:"page_margin" yaml:"page_margin"`                           // Page margin in millimeters
	DefaultFontName  string        `json:"default_font_name" yaml:"default_font_name"`               // Default font name
	DefaultFontPath  string        `json:"default_font_path" yaml:"default_font_path"`               // Default font path
	ShowTitle        bool          `json:"show_title" yaml:"show_title"`                             // Whether to show the title
	ShowSubtitle     bool          `json:"show_subtitle" yaml:"show_subtitle"`                       // Whether to show the subtitle
	ShowLegend       bool          `json:"show_legend" yaml:"show_legend"`                           // Whether to show the legend
	ShowAxisNames    bool          `json:"show_axis_labels" yaml:"show_axis_labels"`                 // Whether to show the axis labels
	ShowTicks        bool          `json:"show_ticks" yaml:"show_ticks"`                             // Whether to show the ticks
	ShowTickLabels   bool          `json:"show_tick_labels" yaml:"show_tick_labels"`                 // Whether to show the tick labels
	ShowPointMarkers bool          `json:"show_point_markers" yaml:"show_point_markers"`             // Whether to show the point markers
}

// DefaultChartOptions returns default chart options

// DefaultChartOptions returns default chart options
func DefaultChartOptions() ChartOptions {
	return ChartOptions{
		Width:            DefaultChartWidth,
		Height:           DefaultChartHeight,
		Background:       Color("white"),
		Foreground:       Color("black"),
		TitleStyle:       DefaultTitleStyle(),
		SubtitleStyle:    DefaultSubtitleStyle(),
		PlotOptions:      DefaultPlotOptions(),
		AxisOptions:      DefaultAxisOptions(),
		SeriesOptions:    DefaultSeriesOptions(),
		LegendOptions:    DefaultLegendOptions(),
		Colors:           DefaultSeriesColors,
		PointMarkers:     DefaultPointMarkers,
		PageMargin:       DefaultPageMargin,
		ShowTitle:        true,
		ShowSubtitle:     true,
		ShowLegend:       true,
		ShowAxisNames:    true,
		ShowTicks:        true,
		ShowTickLabels:   true,
		ShowPointMarkers: true,
	}
}

func DefaultTitleStyle() Font {
	return Font{
		Size:  DefaultTitleFontSize,
		Color: Color("#000000"),
	}
}

func DefaultSubtitleStyle() Font {
	return Font{
		Size:  DefaultSubtitleFontSize,
		Color: Color("#000000"),
	}
}

type PlotOptions struct {
	Scale            float64     `json:"scale" yaml:"scale"`                         // Scale of the plot in millimeters
	OutlineThickness float64     `json:"outline_thickness" yaml:"outline_thickness"` // Outline thickness in millimeters
	OutlineColor     Color       `json:"outline_color" yaml:"outline_color"`         // Outline color
	ConnectType      ConnectType `json:"connect_type" yaml:"connect_type"`           // Connect type for the plot
	Margin           float64     `json:"margin" yaml:"margin"`                       // Margin of the plot in millimeters
	Padding          float64     `json:"padding" yaml:"padding"`                     // Padding of the plot in millimeters
}

func DefaultPlotOptions() PlotOptions {
	return PlotOptions{
		Scale:            DefaultPlotScale,
		OutlineThickness: DefaultPlotOutlineThickness,
		OutlineColor:     Color("#000000"),
		ConnectType:      DefaultConnectType,
		Margin:           DefaultPlotMargin,
		Padding:          DefaultPlotPadding,
	}
}

// ChartData represents the data in the chart
type ChartData struct {
	Axes   []Axis   `json:"axes" yaml:"axes"`     // Axes definitions
	Series []Series `json:"series" yaml:"series"` // Data series
}

// Chart represents a complete spider chart
type Chart struct {
	Options      ChartOptions                `json:"options" yaml:"options"` // Chart options
	Data         ChartData                   `json:"data" yaml:"data"`       // Chart data
	titleRect    canvas.Rect                 `json:"-" yaml:"-"`             // Rectangle for the title
	subtitleRect canvas.Rect                 `json:"-" yaml:"-"`             // Rectangle for the subtitle
	plotRect     canvas.Rect                 `json:"-" yaml:"-"`             // Rectangle for the plot
	pageRect     canvas.Rect                 `json:"-" yaml:"-"`             // Rectangle for the page
	legendRect   canvas.Rect                 `json:"-" yaml:"-"`             // Rectangle for the legend
	fonts        map[string]*canvas.FontFace `json:"-" yaml:"-"`             // Fonts
}

// NewChart creates a new chart with the given options and data
func NewChartWithDataAndOptions(data ChartData, options ChartOptions) *Chart {
	return &Chart{
		Options: options,
		Data:    data,
	}
}

func NewChartWithData(data ChartData) *Chart {
	return NewChartWithDataAndOptions(data, DefaultChartOptions())
}

func NewChart() *Chart {
	return NewChartWithDataAndOptions(ChartData{}, DefaultChartOptions())
}

// Radius returns the radius of the plot area in millimeters
func (c *Chart) Radius() float64 {
	canvasWidth := c.Width()
	canvasHeight := c.Height()
	return c.Options.PlotOptions.Scale*math.Min(canvasWidth, canvasHeight)/2 - c.Options.PlotOptions.Padding
}

func (c *Chart) AddAxis(name string) error {
	// check that chart axes are not nil
	if c.Data.Axes == nil {
		c.Data.Axes = make([]Axis, 0)
	}
	// check if axis already exists
	for _, axis := range c.Data.Axes {
		if axis.Name == name {
			return &ValidationError{
				Field:   "axes",
				Message: fmt.Sprintf("axis %s already exists", name),
			}
		}
	}
	if len(c.Data.Axes) >= MaxAxes {
		return &ValidationError{
			Field:   "axes",
			Message: fmt.Sprintf("maximum %d axes allowed, got %d", MaxAxes, len(c.Data.Axes)),
		}
	}
	c.Data.Axes = append(c.Data.Axes, Axis{Name: name})
	return nil
}

func (c *Chart) AddSeries(name string, data map[string]float64) error {
	// check that chart series are not nil
	if c.Data.Series == nil {
		c.Data.Series = make([]Series, 0)
	}
	// check if series already exists
	for _, series := range c.Data.Series {
		if series.Name == name {
			return &ValidationError{
				Field:   "series",
				Message: fmt.Sprintf("series %s already exists", name),
			}
		}
	}
	if len(c.Data.Series) >= MaxSeries {
		return &ValidationError{
			Field:   "series",
			Message: fmt.Sprintf("maximum %d series allowed, got %d", MaxSeries, len(c.Data.Series)),
		}
	}

	// Initialize data map if nil
	if data == nil {
		data = make(map[string]float64)
	}

	c.Data.Series = append(c.Data.Series, Series{
		Name:    name,
		Data:    data,
		Options: DefaultSeriesOptions(),
	})
	return nil
}

// CanvasWidth returns the canvas width in millimeters
func (c *Chart) Width() float64 {
	return c.Options.Width
}

// CanvasHeight returns the canvas height in millimeters
func (c *Chart) Height() float64 {
	return c.Options.Height
}

// Draw draws the chart to the given canvas context
func (c *Chart) Draw(ctx *canvas.Context) error {
	// Validate chart before drawing
	if err := c.validate(); err != nil {
		return err
	}

	c.calcRects()

	// Draw background
	c.drawBackground(ctx)
	c.drawTitle(ctx)
	c.drawSubtitle(ctx)
	c.drawPlotBackground(ctx)
	c.drawAxes(ctx)
	c.drawSeries(ctx)
	c.drawLegend(ctx)

	return nil
}

func (c *Chart) calcRects() {
	w := c.Width()
	h := c.Height()

	plotWidth := w * c.Options.PlotOptions.Scale
	plotHeight := h * c.Options.PlotOptions.Scale
	plotX := (w - plotWidth - c.Options.PageMargin) / 2
	plotY := (h - plotHeight - c.Options.PageMargin) / 2
	c.pageRect = canvas.Rect{
		X0: c.Options.PageMargin,
		Y0: c.Options.PageMargin,
		X1: w - c.Options.PageMargin,
		Y1: h - c.Options.PageMargin,
	}
	c.plotRect = canvas.Rect{
		X0: plotX,
		Y0: plotY,
		X1: plotX + plotWidth,
		Y1: plotY + plotHeight,
	}
	// legend rect
	subtitleBottom := c.plotRect.Y1 + c.Options.PlotOptions.Margin
	switch c.Options.LegendOptions.Placement {
	case LegendPlacementTop:
		targetHeight := c.Options.LegendOptions.LegendStyle.Size*mmPerPt*smidge + c.Options.LegendOptions.Padding
		if targetHeight < c.Options.LegendOptions.MinHeight {
			targetHeight = c.Options.LegendOptions.MinHeight
			c.plotRect.Y0 -= targetHeight
			c.plotRect.Y1 -= targetHeight
		}
		c.legendRect = canvas.Rect{ // above plot + plot margin
			X0: c.Options.PageMargin,
			Y0: c.plotRect.Y1 + c.Options.PlotOptions.Margin,
			X1: w - c.Options.PageMargin,
			Y1: c.plotRect.Y1 + c.Options.PlotOptions.Margin + targetHeight,
		}
		subtitleBottom = c.legendRect.Y1 + c.Options.SubtitleMargin
	case LegendPlacementBottom:
		// targetHeight := c.plotRect.Y0 - c.Options.PageMargin
		// if targetHeight < c.Options.LegendOptions.MinHeight {
		// 	targetHeight = c.Options.LegendOptions.MinHeight
		// 	c.plotRect.Y0 += targetHeight
		// 	c.plotRect.Y1 += targetHeight
		// }
		c.legendRect = canvas.Rect{
			X0: c.Options.PageMargin,
			Y0: c.Options.PageMargin,
			X1: w - c.Options.PageMargin,
			Y1: c.plotRect.Y0 - c.Options.PlotOptions.Margin,
		}
	case LegendPlacementLeft:
		targetWidth := c.plotRect.X0 - c.Options.PageMargin
		if targetWidth < c.Options.LegendOptions.MinWidth {
			targetWidth = c.Options.LegendOptions.MinWidth
			c.plotRect.X0 += targetWidth
			c.plotRect.X1 += targetWidth
		}
		c.legendRect = canvas.Rect{
			X0: c.Options.PageMargin,
			Y0: c.plotRect.Y0,
			X1: c.plotRect.X0 - c.Options.PlotOptions.Margin,
			Y1: c.plotRect.Y1,
		}
	case LegendPlacementRight:
		targetWidth := w - c.plotRect.X1 - c.Options.PageMargin
		if targetWidth < c.Options.LegendOptions.MinWidth {
			targetWidth = c.Options.LegendOptions.MinWidth
			c.plotRect.X1 -= targetWidth
			c.plotRect.X0 -= targetWidth
		}
		c.legendRect = canvas.Rect{
			X0: c.plotRect.X1 + c.Options.PlotOptions.Margin,
			Y0: c.plotRect.Y0,
			X1: w - c.Options.PageMargin,
			Y1: c.plotRect.Y1,
		}
	}
	// subtitle rect
	subtitleHeight := c.fonts["subtitle"].LineHeight() * smidge
	if c.Options.Subtitle == "" {
		subtitleHeight = 0
	}
	c.subtitleRect = canvas.Rect{
		X0: c.Options.PageMargin,
		Y0: subtitleBottom,
		X1: w - c.Options.PageMargin,
		Y1: subtitleBottom + subtitleHeight,
	}
	// title rect
	c.titleRect = canvas.Rect{
		X0: c.Options.PageMargin,
		Y0: c.subtitleRect.Y1 + c.Options.TitleMargin,
		X1: w - c.Options.PageMargin,
		Y1: c.pageRect.Y1,
	}
}

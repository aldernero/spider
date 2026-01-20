package spider

import "math"

const (
	// DefaultChartWidth is the default width for the chart in millimeters
	DefaultChartWidth = 200.0

	// DefaultChartHeight is the default height for the chart in millimeters
	DefaultChartHeight = 200.0

	// AutoscaleAxisPadding is the default padding for the axis max value
	AutoscaleAxisPaddingFactor = 1.15

	// DefaultLabelOffset is the default offset for the axis label in millimeters
	DefaultLabelOffset = 3.0

	// DefaultConnectType is the default connection type for the plot shape
	DefaultConnectType = ConnectTypeCircle

	// Tau is 2π, used for angle calculations
	Tau = 2 * math.Pi

	// MaxAxes is the maximum number of axes allowed in a chart
	MaxAxes = 50

	// MaxSeries is the maximum number of series allowed in a chart
	MaxSeries = 20

	// DefaultPageMargin is the default margin for the page in millimeters
	DefaultPageMargin = 3.0

	// DefaultPlotScale is the default plot size as a percentage of canvas
	DefaultPlotScale = 0.6

	// DefaultPlotMargin is the default margin for the plot in millimeters
	DefaultPlotMargin = 3.0

	// DefaultPlotPadding is the default padding for the plot in millimeters
	DefaultPlotPadding = 5.0

	// DefaultTitleMargin is the default margin for the title in millimeters
	DefaultTitleMargin = 3.0

	// DefaultSubtitleMargin is the default margin for the subtitle in millimeters
	DefaultSubtitleMargin = 3.0

	// DefaultMajorTickCount is the default number of major ticks
	DefaultMajorTickCount = 5

	// DefaultMinorTickCount is the default number of minor ticks per major tick
	DefaultMinorTickCount = 2

	// DefaultMajorTickSize is the default size of major ticks in millimeters
	DefaultMajorTickLength = 2.0

	// DefaultMinorTickSize is the default size of minor ticks in millimeters
	DefaultMinorTickLength = 1.0

	// DefaultPlotOutlineThickness is the default outline thickness for the plot in millimeters
	DefaultPlotOutlineThickness = 1.0

	// DefaultLineThickness is the default line thickness in millimeters
	DefaultAxisLineThickness = 0.75

	// DefaultMajorTickLineThickness is the default major tick line thickness in millimeters
	DefaultMajorTickLineThickness = 0.5

	// DefaultMinorTickLineThickness is the default minor tick line thickness in millimeters
	DefaultMinorTickLineThickness = 0.25

	// DefaultSeriesLineThickness is the default line thickness for series in millimeters
	DefaultSeriesLineThickness = 0.75

	// DefaultLegendLineLength is the default line length for legend in millimeters
	DefaultLegendLineLength = 7.0

	// DefaultLegendLineThickness is the default line thickness for legend in millimeters
	DefaultLegendLineThickness = 0.6

	// DefaultPointSize is the default point size in millimeters
	DefaultPointSize = 2.0

	// DefaultFontSize is the default font size in points
	DefaultFontSize = 12.0

	// DefaultTitleFontSize is the default title font size in points
	DefaultTitleFontSize = 18.0

	// DefaultSubtitleFontSize is the default subtitle font size in points
	DefaultSubtitleFontSize = 14.0

	// DefaultLegendFontSize is the default legend font size in points
	DefaultLegendFontSize = 10.0

	// DefaultAxisLabelFontSize is the default axis label font size in points
	DefaultAxisLabelFontSize = 10.0

	// DefaultTickLabelFontSize is the default tick label font size in points
	DefaultTickLabelFontSize = 8.0

	smidge  = 1.000000001
	mmPerPt = 0.3527777777777778
)

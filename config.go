package spider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// NewChartFromFile creates a chart from a configuration file
// It automatically detects the format based on the file extension (.json, .yaml, .yml)
func NewChartFromFile(filename string) (*Chart, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return NewChartFromJSON(data)
	case ".yaml", ".yml":
		return NewChartFromYAML(data)
	default:
		// Try JSON first, then YAML
		if chart, err := NewChartFromJSON(data); err == nil {
			return chart, nil
		}
		return NewChartFromYAML(data)
	}
}

// NewChartFromJSON creates a chart from JSON data
func NewChartFromJSON(data []byte) (*Chart, error) {
	// Start with defaults
	chart := Chart{
		Options: DefaultChartOptions(),
		Data:    ChartData{},
	}

	// Unmarshal config file, which will override defaults for specified fields
	if err := json.Unmarshal(data, &chart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Apply defaults to any nested fields that might still be zero values
	chart = *applyDefaults(&chart)

	// Validate
	if err := chart.validate(); err != nil {
		return nil, fmt.Errorf("chart validation failed: %w", err)
	}

	return &chart, nil
}

// NewChartFromYAML creates a chart from YAML data
func NewChartFromYAML(data []byte) (*Chart, error) {
	// Start with defaults
	chart := Chart{
		Options: DefaultChartOptions(),
		Data:    ChartData{},
	}

	// Unmarshal config file, which will override defaults for specified fields
	if err := yaml.Unmarshal(data, &chart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Apply defaults to any nested fields that might still be zero values
	chart = *applyDefaults(&chart)

	// Validate
	if err := chart.validate(); err != nil {
		return nil, fmt.Errorf("chart validation failed: %w", err)
	}

	return &chart, nil
}

// defaultConfigChart returns a chart to seed a generated config file with. It
// carries a minimal set of axes and a series because a chart with no axes fails
// validation, which would make the generated config unusable as written.
func defaultConfigChart() *Chart {
	chart := NewChart()
	data := make(map[string]float64, 3)
	for i, name := range []string{"axis1", "axis2", "axis3"} {
		if err := chart.AddAxis(name); err != nil {
			panic(err) // unreachable: names are distinct and within MaxAxes
		}
		data[name] = float64(i + 1)
	}
	if err := chart.AddSeries("series1", data); err != nil {
		panic(err) // unreachable: the chart is empty at this point
	}
	return chart
}

func GenerateDefaultConfigJSON(filename string) error {
	jsonData, err := json.MarshalIndent(defaultConfigChart(), "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return os.WriteFile(filename, jsonData, 0644)
}

func GenerateDefaultConfigYAML(filename string) error {
	yamlData, err := yaml.Marshal(defaultConfigChart())
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return os.WriteFile(filename, yamlData, 0644)
}

// applyDefaults applies default values to a chart where they are missing
func applyDefaults(chart *Chart) *Chart {
	// Apply chart option defaults
	if chart.Options.PlotOptions.Scale == 0 {
		chart.Options.PlotOptions.Scale = DefaultPlotScale
	}
	if chart.Options.PlotOptions.ConnectType == "" {
		chart.Options.PlotOptions.ConnectType = DefaultConnectType
	}
	if chart.Options.TitleStyle.Size == 0 {
		chart.Options.TitleStyle.Size = DefaultTitleFontSize
	}
	if chart.Options.SubtitleStyle.Size == 0 {
		chart.Options.SubtitleStyle.Size = DefaultSubtitleFontSize
	}
	if chart.Options.Background == "" {
		chart.Options.Background = Color("transparent")
	}
	// Foreground is the fallback for axis lines, outlines and text
	if chart.Options.Foreground == "" {
		chart.Options.Foreground = Color("black")
	}
	// These palettes are indexed modulo their length while drawing, so an empty
	// one (`colors: []`) would be a division by zero rather than a plain chart.
	if len(chart.Options.Colors) == 0 {
		chart.Options.Colors = DefaultSeriesColors
	}
	if len(chart.Options.PointMarkers) == 0 {
		chart.Options.PointMarkers = DefaultPointMarkers
	}

	// Apply axis defaults
	if chart.Options.AxisOptions.MajorTicks == 0 {
		chart.Options.AxisOptions.MajorTicks = DefaultMajorTickCount
	}
	if chart.Options.AxisOptions.MinorTicks == 0 {
		chart.Options.AxisOptions.MinorTicks = DefaultMinorTickCount
	}
	if chart.Options.AxisOptions.MajorTickLength == 0 {
		chart.Options.AxisOptions.MajorTickLength = DefaultMajorTickLength
	}
	if chart.Options.AxisOptions.MinorTickLength == 0 {
		chart.Options.AxisOptions.MinorTickLength = DefaultMinorTickLength
	}
	if chart.Options.AxisOptions.MajorTickLineThickness == 0 {
		chart.Options.AxisOptions.MajorTickLineThickness = DefaultMajorTickLineThickness
	}
	if chart.Options.AxisOptions.MinorTickLineThickness == 0 {
		chart.Options.AxisOptions.MinorTickLineThickness = DefaultMinorTickLineThickness
	}

	// Apply series style defaults. Colors and point shapes are deliberately left
	// empty here: drawing resolves those against the chart palette per series.
	defaults := DefaultSeriesOptions()
	for i := range chart.Data.Series {
		opts := &chart.Data.Series[i].Options
		if opts.LineThickness == 0 {
			opts.LineThickness = defaults.LineThickness
		}
		if opts.PointSize == 0 {
			opts.PointSize = defaults.PointSize
		}
		// A config that sets no point opacity means "opaque", not "invisible"
		if opts.PointFillOpacity == 0 {
			opts.PointFillOpacity = defaults.PointFillOpacity
		}
	}

	// Apply legend defaults
	if chart.Options.LegendOptions.LegendStyle.Size == 0 {
		chart.Options.LegendOptions.LegendStyle.Size = DefaultFontSize
	}
	if chart.Options.LegendOptions.Padding == 0 {
		chart.Options.LegendOptions.Padding = 2.0
	}
	if chart.Options.LegendOptions.OutlineThickness == 0 {
		chart.Options.LegendOptions.OutlineThickness = 0.5
	}
	if chart.Options.LegendOptions.Placement == "" {
		chart.Options.LegendOptions.Placement = LegendPlacementRight
	}

	return chart
}

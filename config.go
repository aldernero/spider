package spider

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// NewChartFromFile creates a chart from a configuration file
// It automatically detects the format based on the file extension (.json, .yaml, .yml)
func NewChartFromFile(filename string) (*Chart, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
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
	var chart Chart
	if err := json.Unmarshal(data, &chart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Apply defaults
	chart = *applyDefaults(&chart)

	// Validate
	if err := chart.validate(); err != nil {
		return nil, fmt.Errorf("chart validation failed: %w", err)
	}

	return &chart, nil
}

// NewChartFromYAML creates a chart from YAML data
func NewChartFromYAML(data []byte) (*Chart, error) {
	var chart Chart
	if err := yaml.Unmarshal(data, &chart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Apply defaults
	chart = *applyDefaults(&chart)

	// Validate
	if err := chart.validate(); err != nil {
		return nil, fmt.Errorf("chart validation failed: %w", err)
	}

	return &chart, nil
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
	if chart.Options.TitleStyle.Color == "" {
		chart.Options.TitleStyle.Color = Color("#000000")
	}
	if chart.Options.SubtitleStyle.Size == 0 {
		chart.Options.SubtitleStyle.Size = DefaultSubtitleFontSize
	}
	if chart.Options.SubtitleStyle.Color == "" {
		chart.Options.SubtitleStyle.Color = Color("#000000")
	}
	if chart.Options.Background == "" {
		chart.Options.Background = Color("transparent")
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

	// Apply series style defaults
	for i := range chart.Data.Series {
		if chart.Data.Series[i].Options.LineThickness == 0 {
			chart.Data.Series[i].Options.LineThickness = DefaultSeriesLineThickness
		}
	}

	// Apply legend defaults
	if chart.Options.LegendOptions.LegendStyle.Size == 0 {
		chart.Options.LegendOptions.LegendStyle.Size = DefaultFontSize
	}
	if chart.Options.LegendOptions.LegendStyle.Color == "" {
		chart.Options.LegendOptions.LegendStyle.Color = Color("#000000")
	}
	if chart.Options.LegendOptions.Padding == 0 {
		chart.Options.LegendOptions.Padding = 2.0
	}
	if chart.Options.LegendOptions.OutlineThickness == 0 {
		chart.Options.LegendOptions.OutlineThickness = 0.5
	}
	if chart.Options.LegendOptions.OutlineColor == "" {
		chart.Options.LegendOptions.OutlineColor = Color("#000000")
	}
	if chart.Options.LegendOptions.Placement == "" {
		chart.Options.LegendOptions.Placement = LegendPlacementRight
	}

	return chart
}

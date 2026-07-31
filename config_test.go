package spider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testYAML = `
options:
  width: 250
  height: 250
  title: "Test Chart"
  plot_options:
    connect_type: polygon
data:
  axes:
    - name: "alpha"
      max: 100
    - name: "beta"
    - name: "gamma"
  series:
    - name: "one"
      data:
        alpha: 50
        beta: 2
        gamma: 3
      options:
        line_color: "#3B82F6"
        fill_opacity: 0.7
`

const testJSON = `{
  "options": {"width": 250, "height": 250, "title": "Test Chart"},
  "data": {
    "axes": [{"name": "alpha", "max": 100}, {"name": "beta"}, {"name": "gamma"}],
    "series": [{"name": "one", "data": {"alpha": 50, "beta": 2, "gamma": 3}}]
  }
}`

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestNewChartFromYAML(t *testing.T) {
	chart, err := NewChartFromYAML([]byte(testYAML))
	if err != nil {
		t.Fatalf("NewChartFromYAML: %v", err)
	}

	if chart.Options.Title != "Test Chart" || chart.Options.Width != 250 {
		t.Errorf("options not applied: %+v", chart.Options)
	}
	if chart.Options.PlotOptions.ConnectType != ConnectTypePolygon {
		t.Errorf("ConnectType = %q, want polygon", chart.Options.PlotOptions.ConnectType)
	}
	if len(chart.Data.Axes) != 3 || len(chart.Data.Series) != 1 {
		t.Fatalf("got %d axes and %d series, want 3 and 1", len(chart.Data.Axes), len(chart.Data.Series))
	}
	if chart.Data.Axes[0].Max != 100 {
		t.Errorf("axis max = %v, want 100", chart.Data.Axes[0].Max)
	}
	if chart.Data.Series[0].Options.LineColor != "#3B82F6" {
		t.Errorf("series color = %q, want #3B82F6", chart.Data.Series[0].Options.LineColor)
	}
}

func TestNewChartFromJSON(t *testing.T) {
	chart, err := NewChartFromJSON([]byte(testJSON))
	if err != nil {
		t.Fatalf("NewChartFromJSON: %v", err)
	}
	if chart.Options.Title != "Test Chart" || len(chart.Data.Axes) != 3 {
		t.Errorf("config not applied: %+v", chart.Options)
	}
}

func TestNewChartFromFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		content string
	}{
		{"yaml extension", "chart.yaml", testYAML},
		{"yml extension", "chart.yml", testYAML},
		{"json extension", "chart.json", testJSON},
		// Without a recognized extension the format is sniffed: JSON, then YAML
		{"no extension with json", "chart", testJSON},
		{"no extension with yaml", "chart", testYAML},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart, err := NewChartFromFile(writeTemp(t, tt.file, tt.content))
			if err != nil {
				t.Fatalf("NewChartFromFile: %v", err)
			}
			if len(chart.Data.Axes) != 3 {
				t.Errorf("got %d axes, want 3", len(chart.Data.Axes))
			}
		})
	}
}

func TestNewChartFromFileErrors(t *testing.T) {
	if _, err := NewChartFromFile(filepath.Join(t.TempDir(), "nope.yaml")); err == nil {
		t.Error("a missing file should be an error")
	}
	if _, err := NewChartFromFile(writeTemp(t, "bad.yaml", "options: [this is not a mapping")); err == nil {
		t.Error("malformed YAML should be an error")
	}
	if _, err := NewChartFromFile(writeTemp(t, "short.yaml", "data:\n  axes:\n    - name: only-one\n")); err == nil {
		t.Error("a config that fails validation should be an error")
	}
}

// A config may set `colors: []`. The palette is indexed modulo its length while
// drawing, so an empty one used to be an integer divide by zero.
func TestConfigWithEmptyPaletteDoesNotPanic(t *testing.T) {
	yaml := strings.Replace(testYAML, "options:\n  width: 250",
		"options:\n  colors: []\n  point_markers: []\n  width: 250", 1)

	chart, err := NewChartFromYAML([]byte(yaml))
	if err != nil {
		t.Fatalf("NewChartFromYAML: %v", err)
	}
	if len(chart.Options.Colors) == 0 || len(chart.Options.PointMarkers) == 0 {
		t.Fatal("an empty palette should be restored to the defaults")
	}
	if err := chart.SaveSVG(filepath.Join(t.TempDir(), "out.svg")); err != nil {
		t.Fatalf("SaveSVG: %v", err)
	}
}

func TestApplyDefaults(t *testing.T) {
	chart := &Chart{}
	applyDefaults(chart)

	opts := chart.Options
	if opts.PlotOptions.Scale != DefaultPlotScale {
		t.Errorf("plot scale = %v, want %v", opts.PlotOptions.Scale, DefaultPlotScale)
	}
	if opts.PlotOptions.ConnectType != DefaultConnectType {
		t.Errorf("connect type = %q, want %q", opts.PlotOptions.ConnectType, DefaultConnectType)
	}
	if opts.TitleStyle.Size != DefaultTitleFontSize || opts.SubtitleStyle.Size != DefaultSubtitleFontSize {
		t.Errorf("font sizes not defaulted: %+v", opts)
	}
	if opts.AxisOptions.MajorTicks != DefaultMajorTickCount || opts.AxisOptions.MinorTicks != DefaultMinorTickCount {
		t.Errorf("tick counts not defaulted: %+v", opts.AxisOptions)
	}
	if opts.AxisOptions.MajorTickLength != DefaultMajorTickLength || opts.AxisOptions.MinorTickLength != DefaultMinorTickLength {
		t.Errorf("tick lengths not defaulted: %+v", opts.AxisOptions)
	}
	if len(opts.Colors) == 0 || len(opts.PointMarkers) == 0 {
		t.Error("palettes not defaulted")
	}
	if opts.LegendOptions.Placement == "" || opts.LegendOptions.LegendStyle.Size == 0 {
		t.Errorf("legend options not defaulted: %+v", opts.LegendOptions)
	}
}

// applyDefaults used to fill only LineThickness, so a series loaded from a
// config got PointFillOpacity 0 and rendered its point fills invisible.
func TestApplyDefaultsFillsSeriesOptions(t *testing.T) {
	chart := &Chart{Data: ChartData{Series: []Series{{Name: "s"}}}}
	applyDefaults(chart)

	defaults := DefaultSeriesOptions()
	got := chart.Data.Series[0].Options
	if got.LineThickness != defaults.LineThickness {
		t.Errorf("LineThickness = %v, want %v", got.LineThickness, defaults.LineThickness)
	}
	if got.PointSize != defaults.PointSize {
		t.Errorf("PointSize = %v, want %v", got.PointSize, defaults.PointSize)
	}
	if got.PointFillOpacity != defaults.PointFillOpacity {
		t.Errorf("PointFillOpacity = %v, want %v", got.PointFillOpacity, defaults.PointFillOpacity)
	}
}

func TestApplyDefaultsKeepsExplicitValues(t *testing.T) {
	chart := &Chart{}
	chart.Options.PlotOptions.Scale = 0.9
	chart.Options.AxisOptions.MajorTicks = 7
	chart.Options.Background = "red"

	applyDefaults(chart)

	if chart.Options.PlotOptions.Scale != 0.9 || chart.Options.AxisOptions.MajorTicks != 7 {
		t.Errorf("explicit values were overwritten: %+v", chart.Options)
	}
	if chart.Options.Background != "red" {
		t.Errorf("Background = %q, want red", chart.Options.Background)
	}
}

// A generated config is meant to be edited and loaded back, so it has to carry
// enough axes to pass validation.
func TestGenerateDefaultConfigRoundTrips(t *testing.T) {
	for _, tt := range []struct {
		name     string
		file     string
		generate func(string) error
	}{
		{"yaml", "config.yaml", GenerateDefaultConfigYAML},
		{"json", "config.json", GenerateDefaultConfigJSON},
	} {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), tt.file)
			if err := tt.generate(path); err != nil {
				t.Fatalf("generate: %v", err)
			}
			chart, err := NewChartFromFile(path)
			if err != nil {
				t.Fatalf("the generated config does not load back: %v", err)
			}
			if err := chart.SaveSVG(filepath.Join(t.TempDir(), "out.svg")); err != nil {
				t.Errorf("the generated config does not render: %v", err)
			}
		})
	}
}

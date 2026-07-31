package spider

import (
	"math"
	"testing"
)

func TestAxisGetMax(t *testing.T) {
	tests := []struct {
		name       string
		axis       Axis
		seriesData []map[string]float64
		want       float64
	}{
		{
			name: "explicit max wins over data",
			axis: Axis{Name: "a", Max: 50},
			seriesData: []map[string]float64{
				{"a": 999},
			},
			want: 50,
		},
		{
			name:       "explicit max with no data",
			axis:       Axis{Name: "a", Max: 50},
			seriesData: nil,
			want:       50,
		},
		{
			name: "auto scales the largest value",
			axis: Axis{Name: "a"},
			seriesData: []map[string]float64{
				{"a": 10},
				{"a": 40},
				{"a": 25},
			},
			want: 40 * AutoscaleAxisPaddingFactor,
		},
		{
			name: "auto ignores other axes",
			axis: Axis{Name: "a"},
			seriesData: []map[string]float64{
				{"a": 10, "b": 1000},
			},
			want: 10 * AutoscaleAxisPaddingFactor,
		},
		{
			name:       "no data falls back to one",
			axis:       Axis{Name: "a"},
			seriesData: nil,
			want:       1.0,
		},
		{
			name: "missing key falls back to one",
			axis: Axis{Name: "a"},
			seriesData: []map[string]float64{
				{"b": 100},
			},
			want: 1.0,
		},
		{
			name: "all zero values fall back to one",
			axis: Axis{Name: "a"},
			seriesData: []map[string]float64{
				{"a": 0},
			},
			want: 1.0,
		},
		{
			name: "only negative values fall back to one",
			axis: Axis{Name: "a"},
			seriesData: []map[string]float64{
				{"a": -50},
			},
			want: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.axis.GetMax(tt.seriesData); math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("GetMax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChartAxisMaxima(t *testing.T) {
	c := newTestChart(t)
	c.Data.Axes[1].Max = 500

	maxima := c.axisMaxima()

	if len(maxima) != len(c.Data.Axes) {
		t.Fatalf("got %d maxima, want %d", len(maxima), len(c.Data.Axes))
	}
	if got := maxima["B"]; got != 500 {
		t.Errorf("explicit max B = %v, want 500", got)
	}
	// A holds 1 and 4 across the two test series, so autoscale uses 4
	if want := 4 * AutoscaleAxisPaddingFactor; math.Abs(maxima["A"]-want) > 1e-9 {
		t.Errorf("auto max A = %v, want %v", maxima["A"], want)
	}
}

func TestDefaultAxisOptions(t *testing.T) {
	opts := DefaultAxisOptions()
	if opts.MajorTicks < 0 || opts.MinorTicks < 0 {
		t.Errorf("default tick counts must not be negative: %+v", opts)
	}
	if !opts.ShowAxis || !opts.ShowName || !opts.ShowTicks || !opts.ShowTickLabels {
		t.Errorf("default axis options should show everything, got %+v", opts)
	}
}

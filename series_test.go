package spider

import (
	"errors"
	"testing"
)

func TestSeriesGetDataValue(t *testing.T) {
	s := Series{Name: "s", Data: map[string]float64{"a": 1.5, "b": 0}}

	if got := s.GetDataValue("a"); got != 1.5 {
		t.Errorf("GetDataValue(a) = %v, want 1.5", got)
	}
	if got := s.GetDataValue("b"); got != 0 {
		t.Errorf("GetDataValue(b) = %v, want 0", got)
	}
	if got := s.GetDataValue("missing"); got != 0 {
		t.Errorf("GetDataValue(missing) = %v, want 0", got)
	}

	var empty Series
	if got := empty.GetDataValue("a"); got != 0 {
		t.Errorf("GetDataValue on nil data = %v, want 0", got)
	}
}

func TestSeriesValidateData(t *testing.T) {
	axisNames := []string{"a", "b"}
	tests := []struct {
		name    string
		data    map[string]float64
		wantErr bool
	}{
		{"exact match", map[string]float64{"a": 1, "b": 2}, false},
		{"missing key", map[string]float64{"a": 1}, true},
		{"extra key", map[string]float64{"a": 1, "b": 2, "c": 3}, true},
		{"empty data", map[string]float64{}, true},
		{"nil data", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Series{Name: "s", Data: tt.data}
			err := s.ValidateData(axisNames)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				var vErr *ValidationError
				if !errors.As(err, &vErr) {
					t.Errorf("error is %T, want *ValidationError", err)
				}
			}
		})
	}
}

func TestGetAllSeriesData(t *testing.T) {
	series := []Series{
		{Name: "s1", Data: map[string]float64{"a": 1}},
		{Name: "s2", Data: map[string]float64{"a": 2}},
	}

	got := getAllSeriesData(series)

	if len(got) != 2 {
		t.Fatalf("got %d data maps, want 2", len(got))
	}
	if got[0]["a"] != 1 || got[1]["a"] != 2 {
		t.Errorf("got %v, want the series data in order", got)
	}
	if len(getAllSeriesData(nil)) != 0 {
		t.Error("getAllSeriesData(nil) should be empty")
	}
}

func TestResolveSeriesOptionsFillsFromPalette(t *testing.T) {
	c := newTestChart(t)
	c.Options.Colors = []Color{"#111111", "#222222"}
	c.Options.PointMarkers = []PointShape{PointShapeSquare}

	opts := c.resolveSeriesOptions(0)

	// An unset fill color must resolve to the palette color. Leaving it empty
	// made it transparent, which combined with FillOpacity rendered as black.
	for name, got := range map[string]Color{
		"LineColor":        opts.LineColor,
		"FillColor":        opts.FillColor,
		"PointStrokeColor": opts.PointStrokeColor,
		"PointFillColor":   opts.PointFillColor,
	} {
		if got != "#111111" {
			t.Errorf("%s = %q, want the palette color #111111", name, got)
		}
	}
	if opts.PointShape != PointShapeSquare {
		t.Errorf("PointShape = %q, want square", opts.PointShape)
	}
	if opts.PointSize == 0 || opts.LineThickness == 0 || opts.PointLineThickness == 0 {
		t.Errorf("sizes should be defaulted, got %+v", opts)
	}

	// The palette cycles, and a second series takes the next color
	if got := c.resolveSeriesOptions(1).LineColor; got != "#222222" {
		t.Errorf("second series LineColor = %q, want #222222", got)
	}
}

func TestResolveSeriesOptionsKeepsExplicitValues(t *testing.T) {
	c := newTestChart(t)
	c.Data.Series[0].Options = SeriesOptions{
		LineColor:     "#ABCDEF",
		FillColor:     "transparent",
		PointShape:    PointShapeDiamond,
		PointSize:     9,
		LineThickness: 3,
	}

	opts := c.resolveSeriesOptions(0)

	if opts.LineColor != "#ABCDEF" || opts.FillColor != "transparent" {
		t.Errorf("explicit colors were overwritten: %+v", opts)
	}
	if opts.PointShape != PointShapeDiamond || opts.PointSize != 9 || opts.LineThickness != 3 {
		t.Errorf("explicit style values were overwritten: %+v", opts)
	}
}

func TestResolveSeriesOptionsDoesNotMutateChart(t *testing.T) {
	c := newTestChart(t)
	before := c.Data.Series[0].Options

	c.resolveSeriesOptions(0)

	if c.Data.Series[0].Options != before {
		t.Errorf("resolveSeriesOptions mutated the series: %+v, want %+v", c.Data.Series[0].Options, before)
	}
}

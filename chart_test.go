package spider

import (
	"errors"
	"strings"
	"testing"
)

// newTestChart returns a minimal chart that passes validation: three axes and
// two series covering them.
func newTestChart(t *testing.T) *Chart {
	t.Helper()
	c := NewChart()
	for _, name := range []string{"A", "B", "C"} {
		if err := c.AddAxis(name); err != nil {
			t.Fatalf("AddAxis(%q): %v", name, err)
		}
	}
	if err := c.AddSeries("s1", map[string]float64{"A": 1, "B": 2, "C": 3}); err != nil {
		t.Fatalf("AddSeries: %v", err)
	}
	if err := c.AddSeries("s2", map[string]float64{"A": 4, "B": 5, "C": 6}); err != nil {
		t.Fatalf("AddSeries: %v", err)
	}
	return c
}

func TestAddAxis(t *testing.T) {
	c := NewChart()

	if err := c.AddAxis("a"); err != nil {
		t.Fatalf("AddAxis: %v", err)
	}
	if len(c.Data.Axes) != 1 || c.Data.Axes[0].Name != "a" {
		t.Fatalf("axes = %+v, want one axis named a", c.Data.Axes)
	}

	err := c.AddAxis("a")
	if err == nil {
		t.Fatal("AddAxis with a duplicate name should fail")
	}
	var vErr *ValidationError
	if !errors.As(err, &vErr) {
		t.Errorf("error is %T, want *ValidationError", err)
	}
	if len(c.Data.Axes) != 1 {
		t.Errorf("a rejected axis should not be appended, got %d", len(c.Data.Axes))
	}
}

func TestAddAxisEnforcesLimit(t *testing.T) {
	c := NewChart()
	for i := 0; i < MaxAxes; i++ {
		if err := c.AddAxis(string(rune('a'+i%26)) + strings.Repeat("x", i/26+1)); err != nil {
			t.Fatalf("AddAxis %d: %v", i, err)
		}
	}
	if err := c.AddAxis("one-too-many"); err == nil {
		t.Errorf("AddAxis past MaxAxes (%d) should fail", MaxAxes)
	}
}

func TestAddSeries(t *testing.T) {
	c := NewChart()

	if err := c.AddSeries("s", map[string]float64{"a": 1}); err != nil {
		t.Fatalf("AddSeries: %v", err)
	}
	if err := c.AddSeries("s", nil); err == nil {
		t.Error("AddSeries with a duplicate name should fail")
	}
	if len(c.Data.Series) != 1 {
		t.Errorf("a rejected series should not be appended, got %d", len(c.Data.Series))
	}
}

func TestAddSeriesInitializesNilData(t *testing.T) {
	c := NewChart()
	if err := c.AddSeries("s", nil); err != nil {
		t.Fatalf("AddSeries: %v", err)
	}
	if c.Data.Series[0].Data == nil {
		t.Error("nil data should be initialized to an empty map, not left nil")
	}
	c.Data.Series[0].Data["a"] = 1 // must not panic
}

func TestAddSeriesEnforcesLimit(t *testing.T) {
	c := NewChart()
	for i := 0; i < MaxSeries; i++ {
		if err := c.AddSeries(strings.Repeat("s", i+1), nil); err != nil {
			t.Fatalf("AddSeries %d: %v", i, err)
		}
	}
	if err := c.AddSeries("one-too-many", nil); err == nil {
		t.Errorf("AddSeries past MaxSeries (%d) should fail", MaxSeries)
	}
}

func TestRadius(t *testing.T) {
	c := NewChart()
	c.Options.Width = 200
	c.Options.Height = 200
	c.Options.PlotOptions.Scale = 0.5
	c.Options.PlotOptions.Padding = 10

	// Before layout, the radius is derived from the canvas: 200*0.5/2 - 10
	if got, want := c.Radius(), 40.0; got != want {
		t.Errorf("Radius() = %v, want %v", got, want)
	}

	// Padding larger than the plot must not produce a negative radius
	c.Options.PlotOptions.Padding = 1000
	if got := c.Radius(); got < 0 {
		t.Errorf("Radius() = %v, want a non-negative value", got)
	}
}

// The plot circle is drawn at Radius() inside plotRect, so a radius wider than
// the rect would overflow the space reserved for it.
func TestRadiusFitsPlotRectOnNonSquareCanvas(t *testing.T) {
	c := newTestChart(t)
	c.Options.Width = 400
	c.Options.Height = 150
	if err := c.validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	c.calcRects()

	r := c.Radius()
	if w, h := c.plotRect.W(), c.plotRect.H(); 2*r > w+1e-9 || 2*r > h+1e-9 {
		t.Errorf("diameter %v exceeds plot rect %vx%v", 2*r, w, h)
	}
}

func TestCalcRectsPlacements(t *testing.T) {
	for _, placement := range []LegendPlacement{
		LegendPlacementTop, LegendPlacementBottom,
		LegendPlacementLeft, LegendPlacementRight, LegendPlacementNone,
	} {
		t.Run(string(placement), func(t *testing.T) {
			c := newTestChart(t)
			c.Options.LegendOptions.Placement = placement
			if err := c.validate(); err != nil {
				t.Fatalf("validate: %v", err)
			}
			c.calcRects()

			if c.plotRect.W() <= 0 || c.plotRect.H() <= 0 {
				t.Errorf("plot rect is empty: %+v", c.plotRect)
			}
			if placement == LegendPlacementNone && (c.legendRect.W() != 0 || c.legendRect.H() != 0) {
				t.Errorf("placement none should leave an empty legend rect, got %+v", c.legendRect)
			}
		})
	}
}

// calcRects reads the loaded font faces; reaching it without validate should
// degrade rather than dereference a nil map entry.
func TestCalcRectsWithoutValidateDoesNotPanic(t *testing.T) {
	c := newTestChart(t)
	c.Options.Subtitle = "sub"
	c.calcRects()
}

func TestDrawRequiresValidChart(t *testing.T) {
	c := NewChart() // no axes
	if err := c.SaveSVG(t.TempDir() + "/out.svg"); err == nil {
		t.Error("drawing a chart with no axes should fail validation")
	}
}

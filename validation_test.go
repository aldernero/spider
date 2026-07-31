package spider

import (
	"strings"
	"testing"
)

func TestValidateAcceptsMinimalChart(t *testing.T) {
	c := newTestChart(t)
	if err := c.validate(); err != nil {
		t.Fatalf("validate() = %v, want nil", err)
	}
	if len(c.fonts) != 5 {
		t.Errorf("loaded %d fonts, want 5", len(c.fonts))
	}
	for _, name := range []string{"title", "subtitle", "axis_label", "tick_label", "legend_label"} {
		if c.fonts[name] == nil {
			t.Errorf("font %q was not loaded", name)
		}
	}
}

func TestValidateRejections(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Chart)
		wantErr string
	}{
		{
			name:    "too few axes",
			mutate:  func(c *Chart) { c.Data.Axes = c.Data.Axes[:2] },
			wantErr: "at least 3 axes",
		},
		{
			name: "too many axes",
			mutate: func(c *Chart) {
				c.Data.Axes = make([]Axis, MaxAxes+1)
				for i := range c.Data.Axes {
					c.Data.Axes[i].Name = strings.Repeat("a", i+1)
				}
				c.Data.Series = nil
			},
			wantErr: "maximum",
		},
		{
			name: "too many series",
			mutate: func(c *Chart) {
				c.Data.Series = make([]Series, MaxSeries+1)
			},
			wantErr: "maximum",
		},
		{
			name:    "unnamed axis",
			mutate:  func(c *Chart) { c.Data.Axes[1].Name = "" },
			wantErr: "no name",
		},
		{
			name:    "duplicate axis name",
			mutate:  func(c *Chart) { c.Data.Axes[1].Name = "A" },
			wantErr: "duplicate",
		},
		{
			name:    "unnamed series",
			mutate:  func(c *Chart) { c.Data.Series[0].Name = "" },
			wantErr: "no name",
		},
		{
			name:    "series missing an axis",
			mutate:  func(c *Chart) { delete(c.Data.Series[0].Data, "B") },
			wantErr: "missing data for axis",
		},
		{
			name:    "series with an unknown axis",
			mutate:  func(c *Chart) { c.Data.Series[0].Data["Z"] = 1 },
			wantErr: "extra data key",
		},
		{
			name:    "zero width",
			mutate:  func(c *Chart) { c.Options.Width = 0 },
			wantErr: "width must be positive",
		},
		{
			name:    "negative height",
			mutate:  func(c *Chart) { c.Options.Height = -1 },
			wantErr: "height must be positive",
		},
		{
			name:    "zero plot scale",
			mutate:  func(c *Chart) { c.Options.PlotOptions.Scale = 0 },
			wantErr: "plot_scale",
		},
		{
			name:    "plot scale above one",
			mutate:  func(c *Chart) { c.Options.PlotOptions.Scale = 1.5 },
			wantErr: "plot_scale",
		},
		{
			// Drawing indexes the palette modulo its length, so an empty one
			// used to be an integer divide by zero.
			name:    "empty color palette",
			mutate:  func(c *Chart) { c.Options.Colors = []Color{} },
			wantErr: "at least one series color",
		},
		{
			name:    "empty point marker list",
			mutate:  func(c *Chart) { c.Options.PointMarkers = []PointShape{} },
			wantErr: "at least one point marker",
		},
		{
			name:    "negative major ticks",
			mutate:  func(c *Chart) { c.Options.AxisOptions.MajorTicks = -1 },
			wantErr: "major_ticks",
		},
		{
			name:    "negative minor ticks",
			mutate:  func(c *Chart) { c.Options.AxisOptions.MinorTicks = -1 },
			wantErr: "minor_ticks",
		},
		{
			name: "min width above max width",
			mutate: func(c *Chart) {
				c.Options.LegendOptions.MinWidth = 100
				c.Options.LegendOptions.MaxWidth = 10
			},
			wantErr: "min_width",
		},
		{
			name: "min height above max height",
			mutate: func(c *Chart) {
				c.Options.LegendOptions.MinHeight = 100
				c.Options.LegendOptions.MaxHeight = 10
			},
			wantErr: "min_height",
		},
		{
			name:    "invalid background color",
			mutate:  func(c *Chart) { c.Options.Background = "chartreusey" },
			wantErr: "unknown color",
		},
		{
			name:    "invalid palette color",
			mutate:  func(c *Chart) { c.Options.Colors = []Color{"#GGGGGG"} },
			wantErr: "invalid hex color",
		},
		{
			name:    "invalid series color",
			mutate:  func(c *Chart) { c.Data.Series[0].Options.LineColor = "definitely-not-a-color" },
			wantErr: "unknown color",
		},
		{
			name:    "invalid series fill color",
			mutate:  func(c *Chart) { c.Data.Series[0].Options.FillColor = "#12345" },
			wantErr: "invalid hex color",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChart(t)
			tt.mutate(c)

			err := c.validate()
			if err == nil {
				t.Fatalf("validate() = nil, want an error mentioning %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("validate() = %q, want it to mention %q", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNormalizesLegendBounds(t *testing.T) {
	c := newTestChart(t)
	c.Options.LegendOptions.MaxWidth = 0
	c.Options.LegendOptions.MaxHeight = 0
	c.Options.LegendOptions.MinWidth = -5

	if err := c.validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	// Unset maxima become unbounded rather than zero, which would collapse the legend
	if c.Options.LegendOptions.MaxWidth <= 0 || c.Options.LegendOptions.MaxHeight <= 0 {
		t.Errorf("unset maxima should become unbounded, got %+v", c.Options.LegendOptions)
	}
	if c.Options.LegendOptions.MinWidth != 0 {
		t.Errorf("negative min width should normalize to 0, got %v", c.Options.LegendOptions.MinWidth)
	}
}

func TestValidateReportsMissingFont(t *testing.T) {
	c := newTestChart(t)
	c.Options.TitleStyle.FontName = "definitely-not-an-installed-font-name"

	err := c.validate()
	if err == nil {
		t.Fatal("validate() = nil, want an error for an unloadable font")
	}
	// The old code wrapped a nil error here and rendered "%!w(<nil>)"
	if strings.Contains(err.Error(), "%!w") {
		t.Errorf("error formats a nil wrap: %q", err)
	}
}

func TestValidationErrorMessage(t *testing.T) {
	err := &ValidationError{Field: "axes", Message: "boom"}
	if got, want := err.Error(), "validation error in axes: boom"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

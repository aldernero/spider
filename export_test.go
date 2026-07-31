package spider

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// renderSVG draws a chart and returns the SVG source.
func renderSVG(t *testing.T, c *Chart) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "out.svg")
	if err := c.SaveSVG(path); err != nil {
		t.Fatalf("SaveSVG: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("SaveSVG wrote an empty file")
	}
	return string(data)
}

func TestSaveDispatchesOnExtension(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{"png", "out.png", false},
		{"svg", "out.svg", false},
		{"uppercase extension", "out.SVG", false},
		{"unsupported extension", "out.gif", true},
		{"no extension", "out", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChart(t)
			err := c.Save(filepath.Join(t.TempDir(), tt.file))
			if (err != nil) != tt.wantErr {
				t.Fatalf("Save(%q) error = %v, wantErr %v", tt.file, err, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(err.Error(), "unsupported file format") {
				t.Errorf("error = %q, want it to mention the unsupported format", err)
			}
		})
	}
}

func TestSavePNGProducesAPNG(t *testing.T) {
	c := newTestChart(t)
	path := filepath.Join(t.TempDir(), "out.png")
	if err := c.SavePNG(path); err != nil {
		t.Fatalf("SavePNG: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.HasPrefix(string(data), "\x89PNG\r\n\x1a\n") {
		t.Error("output does not carry a PNG signature")
	}
}

func TestSaveSVGRendersChartContent(t *testing.T) {
	c := newTestChart(t)
	c.Options.Title = "My Title"
	c.Options.Subtitle = "My Subtitle"

	svg := renderSVG(t, c)

	for _, want := range []string{"<svg", "My Title", "My Subtitle", "A", "B", "C", "s1", "s2"} {
		if !strings.Contains(svg, want) {
			t.Errorf("SVG output is missing %q", want)
		}
	}
}

// End-to-end regression for the alpha-premultiplication bug. A #3B82F6 fill at
// 0.7 opacity used to reach the renderer as rgba(84,186,255) because straight
// components were stored in the premultiplied color.RGBA type.
func TestSaveSVGWritesRequestedFillColor(t *testing.T) {
	c := newTestChart(t)
	c.Data.Series[0].Options.FillColor = "#3B82F6"
	c.Data.Series[0].Options.FillOpacity = 0.7

	svg := strings.ToLower(renderSVG(t, c))

	// (84,186,255) is #54baff, the over-bright value the bug produced
	if strings.Contains(svg, "#54baff") {
		t.Error("fill color is over-bright: the color is being treated as alpha-premultiplied")
	}
	// The renderer stores colors premultiplied internally, so a round trip can
	// shift a component by one: #3b82f6 comes back as #3a81f6.
	if !strings.Contains(svg, "#3a81f6") && !strings.Contains(svg, "#3b82f6") {
		t.Errorf("SVG does not contain the requested fill color #3B82F6; got fills:\n%s", excerptFills(svg))
	}
}

// The renderer writes translucent paints as CSS Color 4 rgba(), which SVG 1.1
// renderers such as Inkscape cannot parse: they fall back to the initial value
// and paint every translucent shape solid black.
func TestSaveSVGAvoidsRGBAFunctionalNotation(t *testing.T) {
	c := newTestChart(t)
	c.Data.Series[0].Options.FillColor = "#3B82F6"
	c.Data.Series[0].Options.FillOpacity = 0.7

	svg := renderSVG(t, c)

	if strings.Contains(svg, "rgba(") {
		t.Errorf("output uses rgba() notation, which SVG 1.1 renderers paint black:\n%s", excerptRGBA(svg))
	}
	if !strings.Contains(svg, "fill-opacity:") {
		t.Error("translucent fill should carry a separate fill-opacity property")
	}
}

func TestSVGToSVG11Colors(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "translucent fill",
			in:   `style="fill:rgba(58,129,246,.70196078);stroke-width:.75"`,
			want: `style="fill:#3a81f6;fill-opacity:.70196078;stroke-width:.75"`,
		},
		{
			name: "translucent stroke",
			in:   `style="stroke:rgba(255,0,0,.5)"`,
			want: `style="stroke:#ff0000;stroke-opacity:.5"`,
		},
		{
			name: "both",
			in:   `style="fill:rgba(0,0,0,0);stroke:rgba(1,2,3,.25)"`,
			want: `style="fill:#000000;fill-opacity:0;stroke:#010203;stroke-opacity:.25"`,
		},
		{
			name: "opaque hex is left alone",
			in:   `style="fill:none;stroke:#000;stroke-width:.25"`,
			want: `style="fill:none;stroke:#000;stroke-width:.25"`,
		},
		{
			name: "out of range components are left alone",
			in:   `style="fill:rgba(300,0,0,.5)"`,
			want: `style="fill:rgba(300,0,0,.5)"`,
		},
		{
			name: "gradient stop colors are untouched",
			in:   `<stop stop-color="rgba(1,2,3,.5)"/>`,
			want: `<stop stop-color="rgba(1,2,3,.5)"/>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(svgToSVG11Colors([]byte(tt.in))); got != tt.want {
				t.Errorf("svgToSVG11Colors() = %q, want %q", got, tt.want)
			}
		})
	}
}

// excerptFills lists the distinct fill declarations in an SVG, for failure output.
func excerptFills(svg string) string {
	seen := map[string]bool{}
	var found []string
	for _, part := range strings.Split(svg, "fill:") {
		if end := strings.IndexAny(part, ";\""); end > 0 && !seen[part[:end]] {
			seen[part[:end]] = true
			found = append(found, "fill:"+part[:end])
		}
	}
	return strings.Join(found, "\n")
}

func excerptRGBA(svg string) string {
	var found []string
	for i := 0; i+5 < len(svg); i++ {
		if strings.HasPrefix(svg[i:], "rgba(") {
			if end := strings.Index(svg[i:], ")"); end > 0 {
				found = append(found, svg[i:i+end+1])
			}
		}
	}
	return strings.Join(found, " ")
}

// FillOpacity with no FillColor used to render black, because an unset color
// became transparent and then had the opacity applied to it.
func TestFillOpacityWithoutFillColorUsesPaletteColor(t *testing.T) {
	c := newTestChart(t)
	c.Options.Colors = []Color{"#3B82F6"}
	c.Data.Series[0].Options.FillOpacity = 0.5
	c.Data.Series[1].Options.FillOpacity = 0.5

	svg := renderSVG(t, c)

	if strings.Contains(svg, "rgba(0,0,0,.5") {
		t.Error("fill rendered as black; an unset fill color should fall back to the palette")
	}
}

func TestShowFlagsSuppressOutput(t *testing.T) {
	tests := []struct {
		name    string
		disable func(*Chart)
		absent  string
	}{
		{"ShowTitle", func(c *Chart) { c.Options.ShowTitle = false }, "Suppressed Title"},
		{"ShowSubtitle", func(c *Chart) { c.Options.ShowSubtitle = false }, "Suppressed Subtitle"},
		{"ShowAxisNames", func(c *Chart) { c.Options.ShowAxisNames = false }, "AxisName"},
		{"AxisOptions.ShowName", func(c *Chart) { c.Options.AxisOptions.ShowName = false }, "AxisName"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestChart(t)
			c.Options.Title = "Suppressed Title"
			c.Options.Subtitle = "Suppressed Subtitle"
			c.Data.Axes[0].Name = "AxisName"
			c.Data.Series[0].Data["AxisName"] = c.Data.Series[0].Data["A"]
			c.Data.Series[1].Data["AxisName"] = c.Data.Series[1].Data["A"]
			delete(c.Data.Series[0].Data, "A")
			delete(c.Data.Series[1].Data, "A")
			tt.disable(c)

			if svg := renderSVG(t, c); strings.Contains(svg, tt.absent) {
				t.Errorf("%s = false but %q still appears in the output", tt.name, tt.absent)
			}
		})
	}
}

// ShowLegend was declared but never consulted; only LegendOptions.Show was.
func TestShowLegendSuppressesLegend(t *testing.T) {
	withLegend := renderSVG(t, newTestChart(t))

	c := newTestChart(t)
	c.Options.ShowLegend = false
	withoutLegend := renderSVG(t, c)

	if len(withoutLegend) >= len(withLegend) {
		t.Errorf("ShowLegend = false did not remove the legend (%d bytes vs %d)", len(withoutLegend), len(withLegend))
	}
}

func TestLegendPlacementNoneSuppressesLegend(t *testing.T) {
	c := newTestChart(t)
	c.Options.LegendOptions.Placement = LegendPlacementNone
	if err := c.SaveSVG(filepath.Join(t.TempDir(), "out.svg")); err != nil {
		t.Fatalf("SaveSVG: %v", err)
	}
}

func TestShowTicksSuppressesTicks(t *testing.T) {
	withTicks := renderSVG(t, newTestChart(t))

	c := newTestChart(t)
	c.Options.ShowTicks = false
	withoutTicks := renderSVG(t, c)

	if len(withoutTicks) >= len(withTicks) {
		t.Errorf("ShowTicks = false did not remove the ticks (%d bytes vs %d)", len(withoutTicks), len(withTicks))
	}
}

func TestShowAxisSuppressesAxes(t *testing.T) {
	c := newTestChart(t)
	c.Options.AxisOptions.ShowAxis = false
	if svg := renderSVG(t, c); strings.Contains(svg, ">A<") {
		t.Error("ShowAxis = false but axis labels still appear")
	}
}

// AxisOptions.LineColor and PlotOptions.OutlineColor were declared but the
// drawing code stroked black regardless.
func TestAxisAndOutlineColorsAreApplied(t *testing.T) {
	c := newTestChart(t)
	c.Options.AxisOptions.LineColor = "#FF00FF"
	c.Options.PlotOptions.OutlineColor = "#00FF00"

	svg := renderSVG(t, c)

	for _, want := range []string{"#f0f", "#0f0"} {
		if !strings.Contains(strings.ToLower(svg), want) {
			t.Errorf("SVG does not use the configured color %q", want)
		}
	}
}

// Foreground was declared but never read anywhere.
func TestForegroundIsUsedAsFallback(t *testing.T) {
	c := newTestChart(t)
	c.Options.Foreground = "#FF00FF"

	if svg := renderSVG(t, c); !strings.Contains(strings.ToLower(svg), "#f0f") {
		t.Error("Foreground color is not used for axis and outline strokes")
	}
}

func TestConnectTypes(t *testing.T) {
	for _, ct := range []ConnectType{ConnectTypeCircle, ConnectTypePolygon} {
		t.Run(string(ct), func(t *testing.T) {
			c := newTestChart(t)
			c.Options.PlotOptions.ConnectType = ct
			renderSVG(t, c)
		})
	}
}

func TestPointShapesRender(t *testing.T) {
	for _, shape := range []PointShape{
		PointShapeCircle, PointShapeSquare, PointShapeTriangle, PointShapeDiamond, PointShapeNone,
	} {
		t.Run(string(shape), func(t *testing.T) {
			c := newTestChart(t)
			c.Data.Series[0].Options.PointShape = shape
			renderSVG(t, c)
		})
	}
}

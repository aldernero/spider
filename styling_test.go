package spider

import (
	"image/color"
	"testing"
)

// straight returns the non-premultiplied components of a color, which is the
// form colors are written in and so the form worth asserting on.
func straight(c color.Color) [4]uint8 {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	return [4]uint8{n.R, n.G, n.B, n.A}
}

func TestColorToCanvasColor(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  [4]uint8
	}{
		{"six digit hex", "#3B82F6", [4]uint8{0x3B, 0x82, 0xF6, 0xFF}},
		{"six digit hex lowercase", "#3b82f6", [4]uint8{0x3B, 0x82, 0xF6, 0xFF}},
		{"eight digit hex opaque", "#3B82F6FF", [4]uint8{0x3B, 0x82, 0xF6, 0xFF}},
		{"eight digit hex zero alpha keeps rgb", "#3B82F600", [4]uint8{0x3B, 0x82, 0xF6, 0}},
		{"three digit shorthand", "#F00", [4]uint8{0xFF, 0, 0, 0xFF}},
		{"four digit shorthand", "#F00F", [4]uint8{0xFF, 0, 0, 0xFF}},
		{"shorthand doubles each digit", "#3B82", [4]uint8{0x33, 0xBB, 0x88, 0x22}},
		{"named color", "red", [4]uint8{0xFF, 0, 0, 0xFF}},
		{"named color mixed case", "ReD", [4]uint8{0xFF, 0, 0, 0xFF}},
		{"named color with spaces", "  red  ", [4]uint8{0xFF, 0, 0, 0xFF}},
		{"transparent keyword", "transparent", [4]uint8{0, 0, 0, 0}},
		{"empty means unset", "", [4]uint8{0, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := straight(tt.color.ToCanvasColor()); got != tt.want {
				t.Errorf("ToCanvasColor() = %v, want %v", got, tt.want)
			}
			if err := tt.color.Validate(); err != nil {
				t.Errorf("Validate() = %v, want nil", err)
			}
		})
	}
}

// A malformed color used to be swallowed: bad hex rendered as black and an
// unknown name rendered as transparent, so a typo silently produced a wrong
// chart instead of an error.
func TestColorValidateRejectsMalformed(t *testing.T) {
	for _, c := range []Color{"#ZZZZZZ", "#", "#1", "#12", "#12345", "#3B82F6F", "#123456789", "notacolor", "0000FF"} {
		t.Run(string(c), func(t *testing.T) {
			if err := c.Validate(); err == nil {
				t.Errorf("Validate(%q) = nil, want an error", c)
			}
		})
	}
}

// Regression for the alpha-premultiplication bug: color.RGBA is premultiplied by
// definition, so storing straight components with a partial alpha made canvas
// read the color back over-bright and clipped. #3B82F6 at 0.7 opacity rendered
// as rgba(84,186,255) instead of rgba(59,130,246).
func TestToCanvasColorWithOpacityIsNotPremultiplied(t *testing.T) {
	got := Color("#3B82F6").ToCanvasColorWithOpacity(0.7)

	if _, ok := got.(color.NRGBA); !ok {
		t.Fatalf("got %T, want color.NRGBA (color.RGBA would be interpreted as premultiplied)", got)
	}

	// Straight components survive the round trip through the color model
	if got, want := straight(got), [4]uint8{0x3B, 0x82, 0xF6, 179}; got != want {
		t.Errorf("as NRGBA = %v, want %v", got, want)
	}

	// ...and the premultiplied form a renderer sees is correctly scaled down,
	// rather than left brighter than the alpha allows.
	r, g, b, a := got.RGBA()
	if r > a || g > a || b > a {
		t.Errorf("RGBA() = (%d,%d,%d,%d): components must not exceed alpha", r, g, b, a)
	}
	if wantR := uint32(0x3B * 179 / 255 * 257); absDiff(r, wantR) > 300 {
		t.Errorf("premultiplied red = %d, want about %d", r, wantR)
	}
}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestToCanvasColorWithOpacityClampsAndScales(t *testing.T) {
	tests := []struct {
		name    string
		color   Color
		opacity float64
		wantA   uint8
	}{
		{"fully opaque", "#FF0000", 1.0, 255},
		{"half", "#FF0000", 0.5, 128},
		{"zero", "#FF0000", 0, 0},
		{"negative clamps to zero", "#FF0000", -5, 0},
		{"above one clamps to one", "#FF0000", 42, 255},
		// Opacity scales the color's own alpha rather than replacing it
		{"scales existing alpha", "#FF000080", 0.5, 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := color.NRGBAModel.Convert(tt.color.ToCanvasColorWithOpacity(tt.opacity)).(color.NRGBA)
			if got.A != tt.wantA {
				t.Errorf("alpha = %d, want %d", got.A, tt.wantA)
			}
		})
	}
}

func TestParseUint8(t *testing.T) {
	tests := []struct {
		in      string
		want    uint8
		wantErr bool
	}{
		{"00", 0, false},
		{"ff", 255, false},
		{"FF", 255, false},
		{"3b", 0x3B, false},
		{"A0", 0xA0, false},
		{"zz", 0, true},
		{"0z", 0, true},
		{"g0", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := parseUint8(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseUint8(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("parseUint8(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestDefaultPalettesAreValid(t *testing.T) {
	if len(DefaultSeriesColors) == 0 {
		t.Fatal("DefaultSeriesColors is empty; drawing indexes it modulo its length")
	}
	if len(DefaultPointMarkers) == 0 {
		t.Fatal("DefaultPointMarkers is empty; drawing indexes it modulo its length")
	}
	for _, c := range DefaultSeriesColors {
		if err := c.Validate(); err != nil {
			t.Errorf("default color %q is invalid: %v", c, err)
		}
	}
}

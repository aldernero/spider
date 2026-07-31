package spider

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/tdewolff/canvas"
)

// Color represents a color that can be specified as hex, named color, or RGBA
// It will be converted to color.Color when used
type Color string

// DefaultColors is a list of default colors for the series
var DefaultSeriesColors = []Color{
	Color("#677ad1"),
	Color("#6fac5d"),
	Color("#b94663"),
	Color("#9750a1"),
	Color("#bc7d39"),
}

var DefaultPointMarkers = []PointShape{
	PointShapeCircle,
	PointShapeSquare,
	PointShapeTriangle,
	PointShapeDiamond,
}

// ToCanvasColor converts a Color string to a color.Color.
// Supports hex colors (#RGB, #RGBA, #RRGGBB, #RRGGBBAA) and named colors.
// Colors that cannot be parsed convert to transparent; use Validate to detect them.
func (c Color) ToCanvasColor() color.Color {
	col, err := c.parse()
	if err != nil {
		return canvas.Transparent
	}
	return col
}

// ToCanvasColorWithOpacity converts a Color to a color.Color, scaling its alpha
// by opacity (clamped to [0, 1]).
//
// The result is a non-premultiplied color.NRGBA. Using color.RGBA here would be a
// bug: that type is alpha-premultiplied by definition, so straight (r, g, b)
// components paired with a partial alpha are read back as over-bright and clipped.
func (c Color) ToCanvasColorWithOpacity(opacity float64) color.Color {
	opacity = math.Min(1, math.Max(0, opacity))
	col := color.NRGBAModel.Convert(c.ToCanvasColor()).(color.NRGBA)
	col.A = uint8(math.Round(float64(col.A) * opacity))
	return col
}

// Validate reports whether the color can be parsed. An empty color is valid and
// means "unset"; callers substitute their own default.
func (c Color) Validate() error {
	_, err := c.parse()
	return err
}

// parse converts the color to a color.Color, reporting malformed hex values and
// unknown color names rather than silently substituting a default.
func (c Color) parse() (color.Color, error) {
	s := strings.TrimSpace(string(c))
	if s == "" || strings.EqualFold(s, "transparent") {
		return canvas.Transparent, nil
	}
	if s[0] == '#' {
		return parseHexColor(s)
	}
	if col, ok := namedColors[strings.ToLower(s)]; ok {
		return col, nil
	}
	return nil, fmt.Errorf("unknown color %q: expected a named color or hex value (#RGB, #RGBA, #RRGGBB, #RRGGBBAA)", s)
}

// Color utility functions

// parseHexColor parses a hex color string (#RGB, #RGBA, #RRGGBB or #RRGGBBAA)
func parseHexColor(hex string) (color.Color, error) {
	digits := hex[1:]
	// Expand shorthand notation, where each digit is doubled: #abc -> #aabbcc
	switch len(digits) {
	case 3, 4:
		expanded := make([]byte, 0, 2*len(digits))
		for i := 0; i < len(digits); i++ {
			expanded = append(expanded, digits[i], digits[i])
		}
		digits = string(expanded)
	case 6, 8:
	default:
		return nil, fmt.Errorf("invalid hex color %q: expected 3, 4, 6 or 8 digits, got %d", hex, len(digits))
	}

	// Alpha defaults to fully opaque when not supplied
	components := [4]uint8{0, 0, 0, 255}
	for i := 0; i < len(digits)/2; i++ {
		val, err := parseUint8(digits[2*i : 2*i+2])
		if err != nil {
			return nil, fmt.Errorf("invalid hex color %q: %w", hex, err)
		}
		components[i] = val
	}
	// NRGBA, not RGBA: hex components are straight, not alpha-premultiplied
	return color.NRGBA{R: components[0], G: components[1], B: components[2], A: components[3]}, nil
}

// parseUint8 parses a two-character hex string to uint8
func parseUint8(s string) (uint8, error) {
	var val uint8
	for i := 0; i < 2; i++ {
		var digit uint8
		switch {
		case s[i] >= '0' && s[i] <= '9':
			digit = s[i] - '0'
		case s[i] >= 'a' && s[i] <= 'f':
			digit = s[i] - 'a' + 10
		case s[i] >= 'A' && s[i] <= 'F':
			digit = s[i] - 'A' + 10
		default:
			return 0, fmt.Errorf("invalid hex digit %q", s[i])
		}
		val = val*16 + digit
	}
	return val, nil
}

// namedColors maps CSS color names to their values. It is built once at package
// init: it used to be constructed on every color conversion, of which there are
// many per chart.
//
// from https://github.com/tdewolff/canvas/blob/master/colors_defs.go
var namedColors = map[string]color.Color{
	"aliceblue":            canvas.Aliceblue,
	"antiquewhite":         canvas.Antiquewhite,
	"aqua":                 canvas.Aqua,
	"aquamarine":           canvas.Aquamarine,
	"azure":                canvas.Azure,
	"beige":                canvas.Beige,
	"bisque":               canvas.Bisque,
	"black":                canvas.Black,
	"blanchedalmond":       canvas.Blanchedalmond,
	"blue":                 canvas.Blue,
	"blueviolet":           canvas.Blueviolet,
	"brown":                canvas.Brown,
	"burlywood":            canvas.Burlywood,
	"cadetblue":            canvas.Cadetblue,
	"chartreuse":           canvas.Chartreuse,
	"chocolate":            canvas.Chocolate,
	"coral":                canvas.Coral,
	"cornflowerblue":       canvas.Cornflowerblue,
	"cornsilk":             canvas.Cornsilk,
	"crimson":              canvas.Crimson,
	"cyan":                 canvas.Cyan,
	"darkblue":             canvas.Darkblue,
	"darkcyan":             canvas.Darkcyan,
	"darkgoldenrod":        canvas.Darkgoldenrod,
	"darkgray":             canvas.Darkgray,
	"darkgreen":            canvas.Darkgreen,
	"darkgrey":             canvas.Darkgrey,
	"darkkhaki":            canvas.Darkkhaki,
	"darkmagenta":          canvas.Darkmagenta,
	"darkolivegreen":       canvas.Darkolivegreen,
	"darkorange":           canvas.Darkorange,
	"darkorchid":           canvas.Darkorchid,
	"darkred":              canvas.Darkred,
	"darksalmon":           canvas.Darksalmon,
	"darkseagreen":         canvas.Darkseagreen,
	"darkslateblue":        canvas.Darkslateblue,
	"darkslategray":        canvas.Darkslategray,
	"darkslategrey":        canvas.Darkslategrey,
	"darkturquoise":        canvas.Darkturquoise,
	"darkviolet":           canvas.Darkviolet,
	"deeppink":             canvas.Deeppink,
	"deepskyblue":          canvas.Deepskyblue,
	"dimgray":              canvas.Dimgray,
	"dimgrey":              canvas.Dimgrey,
	"dodgerblue":           canvas.Dodgerblue,
	"firebrick":            canvas.Firebrick,
	"floralwhite":          canvas.Floralwhite,
	"forestgreen":          canvas.Forestgreen,
	"fuchsia":              canvas.Fuchsia,
	"gainsboro":            canvas.Gainsboro,
	"ghostwhite":           canvas.Ghostwhite,
	"gold":                 canvas.Gold,
	"goldenrod":            canvas.Goldenrod,
	"gray":                 canvas.Gray,
	"green":                canvas.Green,
	"greenyellow":          canvas.Greenyellow,
	"grey":                 canvas.Grey,
	"honeydew":             canvas.Honeydew,
	"hotpink":              canvas.Hotpink,
	"indianred":            canvas.Indianred,
	"indigo":               canvas.Indigo,
	"ivory":                canvas.Ivory,
	"khaki":                canvas.Khaki,
	"lavender":             canvas.Lavender,
	"lavenderblush":        canvas.Lavenderblush,
	"lawngreen":            canvas.Lawngreen,
	"lemonchiffon":         canvas.Lemonchiffon,
	"lightblue":            canvas.Lightblue,
	"lightcoral":           canvas.Lightcoral,
	"lightcyan":            canvas.Lightcyan,
	"lightgoldenrodyellow": canvas.Lightgoldenrodyellow,
	"lightgray":            canvas.Lightgray,
	"lightgreen":           canvas.Lightgreen,
	"lightgrey":            canvas.Lightgrey,
	"lightpink":            canvas.Lightpink,
	"lightsalmon":          canvas.Lightsalmon,
	"lightseagreen":        canvas.Lightseagreen,
	"lightskyblue":         canvas.Lightskyblue,
	"lightslategray":       canvas.Lightslategray,
	"lightslategrey":       canvas.Lightslategrey,
	"lightsteelblue":       canvas.Lightsteelblue,
	"lightyellow":          canvas.Lightyellow,
	"lime":                 canvas.Lime,
	"limegreen":            canvas.Limegreen,
	"linen":                canvas.Linen,
	"magenta":              canvas.Magenta,
	"maroon":               canvas.Maroon,
	"mediumaquamarine":     canvas.Mediumaquamarine,
	"mediumblue":           canvas.Mediumblue,
	"mediumorchid":         canvas.Mediumorchid,
	"mediumpurple":         canvas.Mediumpurple,
	"mediumseagreen":       canvas.Mediumseagreen,
	"mediumslateblue":      canvas.Mediumslateblue,
	"mediumspringgreen":    canvas.Mediumspringgreen,
	"mediumturquoise":      canvas.Mediumturquoise,
	"mediumvioletred":      canvas.Mediumvioletred,
	"midnightblue":         canvas.Midnightblue,
	"mintcream":            canvas.Mintcream,
	"mistyrose":            canvas.Mistyrose,
	"moccasin":             canvas.Moccasin,
	"navajowhite":          canvas.Navajowhite,
	"navy":                 canvas.Navy,
	"oldlace":              canvas.Oldlace,
	"olive":                canvas.Olive,
	"olivedrab":            canvas.Olivedrab,
	"orange":               canvas.Orange,
	"orchid":               canvas.Orchid,
	"palegoldenrod":        canvas.Palegoldenrod,
	"palegreen":            canvas.Palegreen,
	"paleturquoise":        canvas.Paleturquoise,
	"palevioletred":        canvas.Palevioletred,
	"papayawhip":           canvas.Papayawhip,
	"peachpuff":            canvas.Peachpuff,
	"peru":                 canvas.Peru,
	"pink":                 canvas.Pink,
	"plum":                 canvas.Plum,
	"powderblue":           canvas.Powderblue,
	"purple":               canvas.Purple,
	"red":                  canvas.Red,
	"rosybrown":            canvas.Rosybrown,
	"royalblue":            canvas.Royalblue,
	"saddlebrown":          canvas.Saddlebrown,
	"salmon":               canvas.Salmon,
	"sandybrown":           canvas.Sandybrown,
	"seagreen":             canvas.Seagreen,
	"seashell":             canvas.Seashell,
	"sienna":               canvas.Sienna,
	"silver":               canvas.Silver,
	"skyblue":              canvas.Skyblue,
	"slateblue":            canvas.Slateblue,
	"slategray":            canvas.Slategray,
	"slategrey":            canvas.Slategrey,
	"snow":                 canvas.Snow,
	"springgreen":          canvas.Springgreen,
	"steelblue":            canvas.Steelblue,
	"tan":                  canvas.Tan,
	"teal":                 canvas.Teal,
	"thistle":              canvas.Thistle,
	"tomato":               canvas.Tomato,
	"turquoise":            canvas.Turquoise,
	"violet":               canvas.Violet,
	"wheat":                canvas.Wheat,
	"white":                canvas.White,
	"whitesmoke":           canvas.Whitesmoke,
	"yellow":               canvas.Yellow,
	"yellowgreen":          canvas.Yellowgreen,
	"transparent":          canvas.Transparent,
}

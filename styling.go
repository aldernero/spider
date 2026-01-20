package spider

import (
	"image/color"
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

// ToCanvasColor converts a Color string to a color.Color
// Supports hex colors (#RRGGBB, #RRGGBBAA) and named colors
func (c Color) ToCanvasColor() color.Color {
	if c == "" || strings.ToLower(string(c)) == "transparent" {
		return canvas.Transparent
	}
	// Try to parse as hex color first
	if len(c) > 0 && c[0] == '#' {
		return parseHexColor(string(c))
	}
	// Try named colors
	if col := parseNamedColor(string(c)); col != nil {
		return *col
	}
	// Default to transparent if parsing fails
	return canvas.Transparent
}

func (c Color) ToCanvasColorWithOpacity(opacity float64) color.Color {
	col := c.ToCanvasColor()
	r, g, b, _ := col.RGBA()
	// RGBA() returns uint32 values in range 0-65535, convert to uint8 (0-255)
	// Divide by 257 for accurate conversion (65535/255 = 257)
	// Clamp opacity to valid range [0, 1]
	if opacity < 0 {
		opacity = 0
	} else if opacity > 1 {
		opacity = 1
	}
	return color.RGBA{
		R: uint8(r / 257),
		G: uint8(g / 257),
		B: uint8(b / 257),
		A: uint8(opacity * 255),
	}
}

// Color utility functions

// parseHexColor parses a hex color string (#RRGGBB or #RRGGBBAA)
func parseHexColor(hex string) color.Color {
	if len(hex) == 7 { // #RRGGBB
		var r, g, b uint8
		if _, err := parseUint8(hex[1:3], &r); err != nil {
			return canvas.Black
		}
		if _, err := parseUint8(hex[3:5], &g); err != nil {
			return canvas.Black
		}
		if _, err := parseUint8(hex[5:7], &b); err != nil {
			return canvas.Black
		}
		return color.RGBA{R: r, G: g, B: b, A: 255}
	} else if len(hex) == 9 { // #RRGGBBAA
		var r, g, b, a uint8
		if _, err := parseUint8(hex[1:3], &r); err != nil {
			return canvas.Black
		}
		if _, err := parseUint8(hex[3:5], &g); err != nil {
			return canvas.Black
		}
		if _, err := parseUint8(hex[5:7], &b); err != nil {
			return canvas.Black
		}
		if _, err := parseUint8(hex[7:9], &a); err != nil {
			return canvas.Black
		}
		return color.RGBA{R: r, G: g, B: b, A: a}
	}
	return canvas.Black
}

// parseUint8 parses a two-character hex string to uint8
func parseUint8(s string, result *uint8) (int, error) {
	var val uint8
	for i := 0; i < 2; i++ {
		var digit uint8
		if s[i] >= '0' && s[i] <= '9' {
			digit = uint8(s[i] - '0')
		} else if s[i] >= 'a' && s[i] <= 'f' {
			digit = uint8(s[i] - 'a' + 10)
		} else if s[i] >= 'A' && s[i] <= 'F' {
			digit = uint8(s[i] - 'A' + 10)
		} else {
			return 0, nil
		}
		val = val*16 + digit
	}
	*result = val
	return 2, nil
}

// parseNamedColor parses common named colors (case-insensitive)
func parseNamedColor(name string) *color.Color {
	// Convert to lowercase for case-insensitive matching
	nameLower := strings.ToLower(name)

	// from https://github.com/tdewolff/canvas/blob/master/colors_defs.go
	colors := map[string]color.Color{
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
	if col, ok := colors[nameLower]; ok {
		return &col
	}
	return nil
}

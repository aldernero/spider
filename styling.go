package spider

import (
	"image/color"
	"strings"

	"github.com/tdewolff/canvas"
)

// Color represents a color that can be specified as hex, named color, or RGBA
// It will be converted to color.Color when used
type Color string

// ToCanvasColor converts a Color string to a color.Color
// Supports hex colors (#RRGGBB, #RRGGBBAA), named colors, and "rgba(r,g,b,a)" format
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

	colors := map[string]color.Color{
		"black":       canvas.Black,
		"white":       canvas.White,
		"red":         canvas.Red,
		"green":       canvas.Green,
		"blue":        canvas.Blue,
		"yellow":      canvas.Yellow,
		"cyan":        canvas.Cyan,
		"magenta":     canvas.Magenta,
		"gray":        canvas.Gray,
		"grey":        canvas.Gray,
		"orange":      color.RGBA{R: 255, G: 165, B: 0, A: 255},
		"purple":      color.RGBA{R: 128, G: 0, B: 128, A: 255},
		"pink":        color.RGBA{R: 255, G: 192, B: 203, A: 255},
		"brown":       color.RGBA{R: 165, G: 42, B: 42, A: 255},
		"transparent": canvas.Transparent,
	}
	if col, ok := colors[nameLower]; ok {
		return &col
	}
	return nil
}

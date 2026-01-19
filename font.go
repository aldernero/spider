package spider

import (
	"fmt"

	"codeberg.org/go-fonts/liberation/liberationsansregular"
	"github.com/tdewolff/canvas"
)

type Font struct {
	FontName string  `json:"name" yaml:"name"`   // Font name
	FontPath string  `json:"path" yaml:"path"`   // Font path
	Size     float64 `json:"size" yaml:"size"`   // Font size in points
	Color    Color   `json:"color" yaml:"color"` // Font color
}

func (f *Font) loadFontFace() (*canvas.FontFace, error) {
	if f.FontPath == "" && f.FontName == "" {
		// Try system fonts first (more reliable, avoids CFF errors)
		// Try common system fonts in order of preference
		systemFonts := []string{"Liberation Sans", "DejaVu Sans", "Helvetica", "Arial"}
		for _, fontName := range systemFonts {
			fontFamily := canvas.NewFontFamily(fontName)
			if err := fontFamily.LoadSystemFont(fontName, canvas.FontRegular); err == nil {
				face := fontFamily.Face(f.Size, f.Color.ToCanvasColor(), canvas.FontRegular, canvas.FontNormal)
				if face != nil {
					return face, nil
				}
			}
		}

		// Last resort: use default system font
		defaultFamily := canvas.NewFontFamily("default")
		if err := defaultFamily.LoadSystemFont("sans-serif", canvas.FontRegular); err == nil {
			face := defaultFamily.Face(f.Size, f.Color.ToCanvasColor(), canvas.FontRegular, canvas.FontNormal)
			if face != nil {
				return face, nil
			}
		}

		// Fallback to embedded font as last resort
		embeddedFamily := canvas.NewFontFamily("embedded")
		err := embeddedFamily.LoadFont(liberationsansregular.TTF, 0, canvas.FontRegular)
		if err != nil {
			return nil, fmt.Errorf("failed to load embedded font: %w", err)
		}
		face := embeddedFamily.Face(f.Size, f.Color.ToCanvasColor(), canvas.FontRegular, canvas.FontNormal)
		if face == nil {
			return nil, fmt.Errorf("failed to load embedded font")
		}
		return face, nil
	}
	family := canvas.NewFontFamily("user-specified")
	if f.FontPath != "" {
		if err := family.LoadFontFile(f.FontPath, canvas.FontRegular); err != nil {
			return nil, fmt.Errorf("failed to load font from file: %w", err)
		}
	} else {
		if err := family.LoadSystemFont(f.FontName, canvas.FontRegular); err != nil {
			return nil, fmt.Errorf("failed to load system font: %w", err)
		}
	}
	face := family.Face(f.Size, f.Color.ToCanvasColor(), canvas.FontRegular, canvas.FontNormal)
	return face, nil
}

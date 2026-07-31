package spider

import (
	"fmt"
	"sync"

	"codeberg.org/go-fonts/liberation/liberationsansregular"
	"github.com/tdewolff/canvas"
)

type Font struct {
	FontName string  `json:"name" yaml:"name"`   // Font name
	FontPath string  `json:"path" yaml:"path"`   // Font path
	Size     float64 `json:"size" yaml:"size"`   // Font size in points
	Color    Color   `json:"color" yaml:"color"` // Font color
}

// fontFamilies caches parsed font families by (name, path).
//
// Parsing a font file dominates chart rendering: it accounts for roughly two
// thirds of the time and four fifths of the allocations of a single Draw, and
// Draw re-validates (and so re-resolves every font) on each call. Deriving faces
// from an already-parsed family instead is about four orders of magnitude cheaper.
var (
	fontFamiliesMu sync.RWMutex
	fontFamilies   = map[string]*canvas.FontFamily{}
)

// loadFontFamily returns a parsed font family, loading it on first use.
// The returned family is shared; canvas.FontFamily.Face only reads it, so
// deriving faces from a shared family concurrently is safe.
func loadFontFamily(fontName, fontPath string) (*canvas.FontFamily, error) {
	key := fontName + "\x00" + fontPath

	fontFamiliesMu.RLock()
	family, ok := fontFamilies[key]
	fontFamiliesMu.RUnlock()
	if ok {
		return family, nil
	}

	family, err := newFontFamily(fontName, fontPath)
	if err != nil {
		return nil, err
	}

	fontFamiliesMu.Lock()
	defer fontFamiliesMu.Unlock()
	// Another goroutine may have raced us here; either family is equivalent, so
	// keep whichever landed first and let ours be collected.
	if existing, ok := fontFamilies[key]; ok {
		return existing, nil
	}
	fontFamilies[key] = family
	return family, nil
}

// newFontFamily loads a font family from a path or system font name. A font the
// caller explicitly asked for is an error when it cannot be loaded, rather than a
// silent substitution; only when no font is requested does it fall back through
// the well-known system fonts to the embedded one.
func newFontFamily(fontName, fontPath string) (*canvas.FontFamily, error) {
	if fontPath != "" {
		family := canvas.NewFontFamily(fontPath)
		if err := family.LoadFontFile(fontPath, canvas.FontRegular); err != nil {
			return nil, fmt.Errorf("failed to load font from %s: %w", fontPath, err)
		}
		return family, nil
	}

	if fontName != "" {
		family := canvas.NewFontFamily(fontName)
		if err := family.LoadSystemFont(fontName, canvas.FontRegular); err != nil {
			return nil, fmt.Errorf("failed to load system font %q: %w", fontName, err)
		}
		return family, nil
	}

	for _, name := range []string{"Liberation Sans", "DejaVu Sans", "sans-serif"} {
		family := canvas.NewFontFamily(name)
		if err := family.LoadSystemFont(name, canvas.FontRegular); err == nil {
			return family, nil
		}
	}

	// No usable system font, so fall back to the font compiled into the binary.
	family := canvas.NewFontFamily("embedded")
	if err := family.LoadFont(liberationsansregular.TTF, 0, canvas.FontRegular); err != nil {
		return nil, fmt.Errorf("failed to load embedded font: %w", err)
	}
	return family, nil
}

// loadFontFace resolves the font face for this style, falling back to the
// chart-level default font when the style does not name one of its own.
func (f *Font) loadFontFace(defaultFontName, defaultFontPath string) (*canvas.FontFace, error) {
	name, path := f.FontName, f.FontPath
	if name == "" && path == "" {
		name, path = defaultFontName, defaultFontPath
	}

	family, err := loadFontFamily(name, path)
	if err != nil {
		return nil, err
	}

	size := f.Size
	if size <= 0 {
		size = DefaultFontSize
	}
	face := family.Face(size, f.Color.ToCanvasColor(), canvas.FontRegular, canvas.FontNormal)
	if face == nil {
		return nil, fmt.Errorf("failed to create %gpt font face", size)
	}
	return face, nil
}

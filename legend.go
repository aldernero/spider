package spider

// LegendConfig represents configuration for the legend
type LegendOptions struct {
	Show             bool            `json:"show" yaml:"show"`                           // Whether to show the legend
	Placement        LegendPlacement `json:"placement" yaml:"placement"`                 // Legend placement
	MinWidth         float64         `json:"min_width" yaml:"min_width"`                 // Minimum width in millimeters
	MaxWidth         float64         `json:"max_width" yaml:"max_width"`                 // Maximum width in millimeters
	MinHeight        float64         `json:"min_height" yaml:"min_height"`               // Minimum height in millimeters
	MaxHeight        float64         `json:"max_height" yaml:"max_height"`               // Maximum height in millimeters
	LineLength       float64         `json:"line_length" yaml:"line_length"`             // Line length in millimeters
	LineThickness    float64         `json:"line_thickness" yaml:"line_thickness"`       // Line thickness in millimeters
	OutlineThickness float64         `json:"outline_thickness" yaml:"outline_thickness"` // Outline thickness in millimeters
	OutlineColor     Color           `json:"outline_color" yaml:"outline_color"`         // Outline color
	LegendStyle      Font            `json:"style" yaml:"style"`                         // Legend text style
	Padding          float64         `json:"padding,omitempty" yaml:"padding,omitempty"` // Padding around legend items
	ShowOutline      bool            `json:"show_outline" yaml:"show_outline"`           // Whether to show the outline
}

// DefaultLegendOptions returns default legend options
func DefaultLegendOptions() LegendOptions {
	return LegendOptions{
		Show:             true,
		Placement:        LegendPlacementBottom,
		LineLength:       DefaultLegendLineLength,
		LineThickness:    DefaultLegendLineThickness,
		OutlineThickness: 0.5,
		OutlineColor:     Color("#000000"),
		LegendStyle: Font{
			Size:  DefaultLegendFontSize,
			Color: Color("#000000"),
		},
		Padding:     2.0,
		ShowOutline: false,
	}
}

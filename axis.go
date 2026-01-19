package spider

func DefaultAxisLabelStyle() Font {
	return Font{
		Size:  DefaultAxisLabelFontSize,
		Color: Color("#000000"),
	}
}

func DefaultTickLabelStyle() Font {
	return Font{
		Size:  DefaultTickLabelFontSize,
		Color: Color("#000000"),
	}
}

type AxisOptions struct {
	LineThickness          float64 `json:"line_thickness" yaml:"line_thickness"`                       // Thickness of the axis line in millimeters
	LineColor              Color   `json:"line_color" yaml:"line_color"`                               // Color of the axis line
	LabelOffset            float64 `json:"label_offset" yaml:"label_offset"`                           // Offset for the axis label in millimeters
	LabelStyle             Font    `json:"label_style" yaml:"label_style"`                             // Style for the axis label
	TickLabelStyle         Font    `json:"tick_label_style" yaml:"tick_label_style"`                   // Style for the tick labels
	MajorTicks             int     `json:"major_ticks" yaml:"major_ticks"`                             // Number of major ticks
	MinorTicks             int     `json:"minor_ticks" yaml:"minor_ticks"`                             // Number of minor ticks
	MajorTickLength        float64 `json:"major_tick_length" yaml:"major_tick_length"`                 // Length of the major ticks in millimeters
	MinorTickLength        float64 `json:"minor_tick_length" yaml:"minor_tick_length"`                 // Length of the minor ticks in millimeters
	MajorTickLineThickness float64 `json:"major_tick_line_thickness" yaml:"major_tick_line_thickness"` // Thickness of the major tick lines in millimeters
	MinorTickLineThickness float64 `json:"minor_tick_line_thickness" yaml:"minor_tick_line_thickness"` // Thickness of the minor tick lines in millimeters
	ShowName               bool    `json:"show_name" yaml:"show_name"`                                 // Whether to show the axis name
	ShowAxis               bool    `json:"show_axis" yaml:"show_axis"`                                 // Whether to show the axis
	ShowTicks              bool    `json:"show_ticks" yaml:"show_ticks"`                               // Whether to show the ticks
	ShowTickLabels         bool    `json:"show_tick_labels" yaml:"show_tick_labels"`                   // Whether to show the tick labels
}

// DefaultAxisOptions returns a default axis options

func DefaultAxisOptions() AxisOptions {
	return AxisOptions{
		LineThickness:          DefaultAxisLineThickness,
		LineColor:              Color("#000000"),
		LabelOffset:            DefaultLabelOffset,
		LabelStyle:             DefaultAxisLabelStyle(),
		TickLabelStyle:         DefaultTickLabelStyle(),
		MajorTicks:             DefaultMajorTickCount,
		MinorTicks:             DefaultMinorTickCount,
		MajorTickLength:        DefaultMajorTickLength,
		MinorTickLength:        DefaultMinorTickLength,
		MajorTickLineThickness: DefaultMajorTickLineThickness,
		MinorTickLineThickness: DefaultMinorTickLineThickness,
		ShowName:               true,
		ShowAxis:               true,
		ShowTicks:              true,
		ShowTickLabels:         true,
	}
}

// Axis represents a single axis in the spider chart
type Axis struct {
	Name string  `json:"name" yaml:"name"`                   // Axis name
	Max  float64 `json:"max,omitempty" yaml:"max,omitempty"` // Maximum value (zero means auto-calculate)
}

// GetMax returns the maximum value for the axis, calculating it from series data if needed
func (a *Axis) GetMax(seriesData []map[string]float64) float64 {
	if a.Max > 0 {
		return a.Max
	}
	// Auto-calculate max from series data
	max := 0.0
	for _, data := range seriesData {
		if val, ok := data[a.Name]; ok && val > max {
			max = val
		}
	}
	// Add autoscale padding
	if max > 0 {
		return max * AutoscaleAxisPaddingFactor
	}
	return 1.0 // Default to 1.0 if no data
}

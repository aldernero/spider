package spider

import "slices"

type SeriesOptions struct {
	LineThickness      float64    `json:"line_thickness" yaml:"line_thickness"`             // Thickness of the line in millimeters
	LineColor          Color      `json:"line_color" yaml:"line_color"`                     // Color of the line
	FillOpacity        float64    `json:"fill_opacity" yaml:"fill_opacity"`                 // Opacity of the fill
	FillColor          Color      `json:"fill_color" yaml:"fill_color"`                     // Color of the fill
	PointSize          float64    `json:"point_size" yaml:"point_size"`                     // Size of the point in millimeters
	PointLineThickness float64    `json:"point_line_thickness" yaml:"point_line_thickness"` // Thickness of the point line in millimeters
	PointStrokeColor   Color      `json:"point_color" yaml:"point_color"`                   // Color of the point
	PointFillColor     Color      `json:"point_fill_color" yaml:"point_fill_color"`         // Color of the point fill
	PointFillOpacity   float64    `json:"point_fill_opacity" yaml:"point_fill_opacity"`     // Opacity of the point fill
	PointShape         PointShape `json:"point_shape" yaml:"point_shape"`                   // Shape of the point
}

// DefaultSeriesOptions returns a default series options
func DefaultSeriesOptions() SeriesOptions {
	return SeriesOptions{
		LineThickness:      DefaultSeriesLineThickness,
		LineColor:          "",
		FillOpacity:        0.0,
		FillColor:          "",
		PointSize:          DefaultPointSize,
		PointLineThickness: 0.0,
		PointStrokeColor:   "",
		PointFillColor:     "",
		PointFillOpacity:   0.0,
		PointShape:         "",
	}
}

// Series represents a data series in the spider chart
type Series struct {
	Name    string             `json:"name" yaml:"name"`       // Series name
	Data    map[string]float64 `json:"data" yaml:"data"`       // Data values keyed by axis name
	Options SeriesOptions      `json:"options" yaml:"options"` // Series options
}

// GetDataValue returns the data value for a given axis name, or 0 if not found
func (s *Series) GetDataValue(axisName string) float64 {
	if val, ok := s.Data[axisName]; ok {
		return val
	}
	return 0.0
}

// ValidateData checks if the series data matches the provided axis names
func (s *Series) ValidateData(axisNames []string) error {
	// Check that all axes have corresponding data
	for _, axisName := range axisNames {
		if _, ok := s.Data[axisName]; !ok {
			return &ValidationError{
				Field:   "series.data",
				Message: "missing data for axis: " + axisName,
			}
		}
	}
	// Check that there are no extra data keys
	for key := range s.Data {
		if !slices.Contains(axisNames, key) {
			return &ValidationError{
				Field:   "series.data",
				Message: "extra data key not matching any axis: " + key,
			}
		}
	}
	return nil
}

// getAllSeriesData extracts all data maps from series for max calculation
func getAllSeriesData(series []Series) []map[string]float64 {
	result := make([]map[string]float64, len(series))
	for i, s := range series {
		result[i] = s.Data
	}
	return result
}

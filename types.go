package spider

// ScaleType represents the type of scale for an axis
type ScaleType string

const (
	// ScaleTypeLinear represents a linear scale
	ScaleTypeLinear ScaleType = "linear"

	// ScaleTypeLog10 represents a base-10 logarithmic scale
	ScaleTypeLog10 ScaleType = "log10"

	// ScaleTypeLog2 represents a base-2 logarithmic scale
	ScaleTypeLog2 ScaleType = "log2"
)

// ConnectType represents how the plot shape connects points
type ConnectType string

const (
	// ConnectTypeCircle represents a circular plot shape
	ConnectTypeCircle ConnectType = "circle"

	// ConnectTypePolygon represents a regular polygon plot shape
	ConnectTypePolygon ConnectType = "polygon"
)

// PointShape represents the shape of data points
type PointShape string

const (
	// PointShapeCircle represents a circular point
	PointShapeCircle PointShape = "circle"

	// PointShapeSquare represents a square point
	PointShapeSquare PointShape = "square"

	// PointShapeTriangle represents a triangular point
	PointShapeTriangle PointShape = "triangle"

	// PointShapeDiamond represents a diamond-shaped point
	PointShapeDiamond PointShape = "diamond"

	// PointShapeNone represents no point (hidden)
	PointShapeNone PointShape = "none"
)

// LegendPlacement represents where the legend should be placed
type LegendPlacement string

const (
	// LegendPlacementTop places the legend at the top
	LegendPlacementTop LegendPlacement = "top"

	// LegendPlacementBottom places the legend at the bottom
	LegendPlacementBottom LegendPlacement = "bottom"

	// LegendPlacementLeft places the legend on the left
	LegendPlacementLeft LegendPlacement = "left"

	// LegendPlacementRight places the legend on the right
	LegendPlacementRight LegendPlacement = "right"

	// LegendPlacementNone disables the legend
	LegendPlacementNone LegendPlacement = "none"
)

func (s ScaleType) String() string {
	return string(s)
}

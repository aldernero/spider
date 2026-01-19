# spider

A Go library and CLI tool for generating spider plots (radar charts) with **independent axis scales**. Unlike traditional radar charts where all axes share the same scale, `spider` allows each axis to have its own independent scale, making it ideal for comparing metrics with vastly different ranges.

## Features

- **Independent Axis Scales**: Each axis can have its own min/max values, perfect for comparing metrics with different ranges
- **Multiple Scale Types**: Support for linear, base-10 logarithmic, and base-2 logarithmic scales
- **Flexible Configuration**: Create charts programmatically or from JSON/YAML configuration files
- **Rich Styling Options**: Customize line colors, fill opacity, point shapes, and more
- **Automatic Scaling**: Auto-calculate axis maximums from series data when not specified
- **Tick Configuration**: Configurable major and minor ticks with labels
- **Legend Support**: Customizable legend with multiple placement options
- **Multiple Export Formats**: Export charts as PNG or SVG
- **CLI Tool**: Simple command-line interface for generating charts from config files

## Installation

```bash
go get github.com/aldernero/spider
```

## Quick Start

### Using the Library

```go
package main

import (
    "log"
    "github.com/aldernero/spider"
    "github.com/tdewolff/canvas"
    "github.com/tdewolff/canvas/renderers"
)

func main() {
    // Create a chart
    chart := spider.NewChart(spider.ChartOptions{
        Width:       500,
        Height:      500,
        PlotScale:   0.7,
        ConnectType: spider.ConnectTypePolygon,
        UnitType:    spider.UnitTypePixels,
        Title:       "Performance Metrics",
    }, spider.ChartData{
        Series: []spider.Series{
            {
                Name: "System A",
                Data: map[string]float64{
                    "throughput": 1000000,
                    "latency":    50,
                    "cost":      5000,
                },
            },
        },
        Axes: []spider.Axis{
            {Name: "throughput", Max: floatPtr(2000000)},
            {Name: "latency", Max: floatPtr(100)},
            {Name: "cost", Max: floatPtr(10000)},
        },
    })

    // Draw and save
    c := canvas.New(chart.CanvasWidth(), chart.CanvasHeight())
    ctx := canvas.NewContext(c)
    chart.Draw(ctx)
    renderers.Write("chart.png", c, chart.Resolution())
}

func floatPtr(f float64) *float64 {
    return &f
}
```

### Using Configuration Files

Create a `chart.yaml` file:

```yaml
options:
  width: 500
  height: 500
  plot_scale: 0.7
  connect_type: polygon
  unit_type: pixels
  title: "Performance Comparison"

data:
  series:
    - name: "System A"
      data:
        throughput: 1000000
        latency: 50
        cost: 5000
      style:
        line:
          color: "#FF0000"
          thickness: 1.0
        fill:
          color: "#FF0000"
          opacity: 0.2
        point:
          shape: circle
          size: 3.0
          show: true

  axes:
    - name: "throughput"
      max: 2000000
      options:
        scale: linear
        show_name: true
        show_axis: true
    - name: "latency"
      max: 100
      options:
        scale: linear
    - name: "cost"
      max: 10000
      options:
        scale: linear
```

Then use the CLI tool:

```bash
go run ./cmd/spider-cli -config chart.yaml -output chart.png
```

Or build and use the CLI:

```bash
go build -o spider-cli ./cmd/spider-cli
./spider-cli -config chart.yaml -output chart.png
```

## Configuration File Format

The library supports both JSON and YAML configuration files. The structure matches the `Chart` type:

- **options**: Chart-level settings (width, height, title, plot scale, etc.)
- **data.series**: Array of data series with styling options
- **data.axes**: Array of axis definitions with scale and tick configurations

### Scale Types

- `linear`: Standard linear scale (default)
- `log10`: Base-10 logarithmic scale
- `log2`: Base-2 logarithmic scale

### Point Shapes

- `circle`: Circular points (default)
- `square`: Square points
- `triangle`: Triangular points
- `diamond`: Diamond-shaped points
- `none`: Hide points

### Legend Placement

- `top`, `bottom`, `left`, `right`
- `top-left`, `top-right`, `bottom-left`, `bottom-right`

## API Overview

### Core Types

- `Chart`: Main chart type containing options and data
- `Axis`: Represents a single axis with scale and tick configuration
- `Series`: Data series with styling options
- `ChartOptions`: Chart-level configuration (size, title, legend, etc.)

### Key Functions

- `NewChart(options, data)`: Create a chart programmatically
- `NewChartFromFile(filename)`: Load chart from JSON/YAML file
- `Save(chart, filename)`: Save chart to PNG or SVG (auto-detects format)
- `SavePNG(chart, filename)`: Save as PNG
- `SaveSVG(chart, filename)`: Save as SVG

### Auto-Max Calculation

If an axis doesn't specify a `max` value, it will be automatically calculated from the series data with 10% padding. This makes it easy to create charts without manually setting all axis ranges.

## Examples

See `cmd/main.go` for a complete programmatic example, or check the `cmd/spider-cli` directory for the CLI tool implementation.

## Limitations

- Maximum 50 axes per chart
- Maximum 20 series per chart
- Log scales require positive maximum values

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

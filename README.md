# spider

A Go library and CLI tool for generating spider plots (radar charts) with **independent axis scales**. Unlike traditional radar charts where all axes share the same scale, `spider` allows each axis to have its own independent scale, making it ideal for comparing metrics with vastly different ranges.

## Features

- **Independent Axis Scales**: Each axis can have its own min/max values, perfect for comparing metrics with different ranges and visualizing tradeoffs
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
  chart := spider.NewChart()
  // add axes
  chart.AddAxis("axis1")
  chart.AddAxis("axis2")
  chart.AddAxis("axis3")
  chart.AddAxis("axis4")
  chart.AddAxis("axis5")
  // add series with datapoints
  chart.AddSeries("series1", map[string]float64{
    "axis1": 1000,
    "axis2": 2.0,
    "axis3": 3.0,
    "axis4": 1000000,
    "axis5": 5.0,
  })
  chart.AddSeries("series2", map[string]float64{
    "axis1": 1500,
    "axis2": 1.0,
    "axis3": 2.5,
    "axis4": 2100000,
    "axis5": 12.0,
  })
  // customize
	chart.Options.Subtitle = "Subtitle"
	chart.Options.Title = "Title"
  // save chart
	if err := chart.Save("output.png"); err != nil {
		log.Fatalf("Failed to save chart: %v", err)
	}
}

```
The code produces the following spider chart
<img width="756" height="756" alt="output" src="https://github.com/user-attachments/assets/d2bcefc3-4d31-448d-a3f1-3b1faf054155" />


### Using Configuration Files

Create a `chart.yaml` file:

```yaml
options:
  connect_type: polygon
  title: "Performance Comparison"

data:
  axes:
    - name: "throughput"
    - name: "latency"
      max: 100
    - name: "cost"

  series:
    - name: "System A"
      data:
        throughput: 1000000
        latency: 50
        cost: 5000
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

See the `examples` folder for more details.

### Point Shapes

- `circle`: Circular points (default)
- `square`: Square points
- `triangle`: Triangular points
- `diamond`: Diamond-shaped points
- `none`: Hide points

### Legend Placement

- `top`, `bottom`, `left`, `right`

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

Some output from the `examples/` folder:

<img width="512" height="512" alt="output" src="https://github.com/user-attachments/assets/dc42b4a9-1362-4c5c-ae4c-f1f88d70ad37" />
<img width="512" height="512" alt="output" src="https://github.com/user-attachments/assets/7dde1c32-4bba-46fd-a2a8-bdf010cb64ca" />
<img width="512" height="512" alt="output" src="https://github.com/user-attachments/assets/5a40596c-3be8-43a1-8e7c-888e52e19dbb" />


## Limitations

- Maximum 50 axes per chart
- Maximum 20 series per chart

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

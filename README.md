# spider

A Go library and CLI tool for generating spider plots (radar charts) with **independent axis scales**. Unlike traditional radar charts where all axes share the same scale, `spider` allows each axis to have its own independent scale, making it ideal for comparing metrics with vastly different ranges.

## Features

- **Independent Axis Scales**: Each axis can have its own maximum value, perfect for comparing metrics with different ranges and visualizing tradeoffs
- **Flexible Configuration**: Create charts programmatically or from JSON/YAML configuration files
- **Rich Styling Options**: Customize line colors, fill opacity, point shapes, and more
- **Automatic Scaling**: Auto-calculate axis maximums from series data when not specified
- **Tick Configuration**: Configurable major and minor ticks with labels
- **Legend Support**: Customizable legend with multiple placement options
- **Multiple Export Formats**: Export charts as PNG or SVG
- **CLI Tool**: Simple command-line interface for generating charts from config files

## Installation

As a library:

```bash
go get github.com/aldernero/spider
```

The CLI, from source:

```bash
go install github.com/aldernero/spider/cmd/spider-cli@latest
```

Or download a prebuilt binary for Linux, macOS or Windows (amd64/arm64) from the
[releases page](https://github.com/aldernero/spider/releases).

## Quick Start

### Using the Library

```go
package main

import (
	"log"

	"github.com/aldernero/spider"
)

func main() {
	// Create a chart
	chart := spider.NewChart()

	// Add axes
	for _, name := range []string{"axis1", "axis2", "axis3", "axis4", "axis5"} {
		if err := chart.AddAxis(name); err != nil {
			log.Fatalf("Failed to add axis %s: %v", name, err)
		}
	}

	// Add series with datapoints
	if err := chart.AddSeries("series1", map[string]float64{
		"axis1": 1000,
		"axis2": 2.0,
		"axis3": 3.0,
		"axis4": 1000000,
		"axis5": 5.0,
	}); err != nil {
		log.Fatalf("Failed to add series: %v", err)
	}
	if err := chart.AddSeries("series2", map[string]float64{
		"axis1": 1500,
		"axis2": 1.0,
		"axis3": 2.5,
		"axis4": 2100000,
		"axis5": 12.0,
	}); err != nil {
		log.Fatalf("Failed to add series: %v", err)
	}

	// Customize
	chart.Options.Title = "Title"
	chart.Options.Subtitle = "Subtitle"

	// Save chart
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
  title: "Performance Comparison"
  plot_options:
    connect_type: polygon

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

- `top`, `bottom`, `left`, `right`, `none`

### Colors

Colors are written as CSS color names (`red`, `cornflowerblue`, case-insensitive),<img width="756" height="756" alt="output" src="https://github.com/user-attachments/assets/8061e447-c867-463d-94b3-db1cb8deace3" />

hex values (`#RGB`, `#RGBA`, `#RRGGBB`, `#RRGGBBAA`), or `transparent`. An
unrecognized color is a validation error rather than a silently wrong chart.

Where a series does not set its own color, it takes the next one from
`options.colors`, cycling when there are more series than colors. `options.foreground`
sets the fallback for axis lines, the plot outline, and text.

## API Overview

### Core Types

- `Chart`: Main chart type containing options and data
- `Axis`: Represents a single axis with scale and tick configuration
- `Series`: Data series with styling options
- `ChartOptions`: Chart-level configuration (size, title, legend, etc.)

### Key Functions

- `NewChart()`: Create an empty chart with default options
- `NewChartWithData(data)`: Create a chart from data, with default options
- `NewChartWithDataAndOptions(data, options)`: Create a fully specified chart
- `NewChartFromFile(filename)`: Load a chart from a JSON/YAML file
- `(*Chart).AddAxis(name)` / `(*Chart).AddSeries(name, data)`: Build a chart programmatically
- `(*Chart).Save(filename)`: Save to PNG or SVG (auto-detects from the extension)
- `(*Chart).SavePNG(filename)` / `(*Chart).SaveSVG(filename)`: Save in a specific format

### Auto-Max Calculation

If an axis doesn't specify a `max` value, it is calculated from the series data with
15% headroom (`AutoscaleAxisPaddingFactor`). This makes it easy to create charts
without manually setting all axis ranges. An axis with no data defaults to a max of 1.

## Examples

See `examples/from-code/main.go` for a complete programmatic example, `examples/devops`
and `examples/ssd-compare` for config-driven ones, or the `cmd/spider-cli` directory for
the CLI tool implementation.

Some output from the `examples/` folder:

<img width="945" height="945" alt="output" src="https://github.com/user-attachments/assets/f99f2832-9947-41d7-9272-d7b2b7d6bab7" />

<img width="756" height="756" alt="output" src="https://github.com/user-attachments/assets/40000910-f251-4a9e-9775-fc1f03ec6114" />

<img width="756" height="756" alt="output" src="https://github.com/user-attachments/assets/ce73ac4b-aa31-4e3d-a506-5a741376f452" />



## Limitations

- Minimum 3, maximum 50 axes per chart
- Maximum 20 series per chart
- Every series must supply a value for every axis
- Axis scales are linear; the `ScaleType` values are declared but not yet applied
- A `Chart` is not safe for concurrent use, though rendering separate charts
  from multiple goroutines is fine

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

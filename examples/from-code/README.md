# Programmatic Chart Creation Example

This example demonstrates how to create a spider chart completely programmatically using Go code, without relying on configuration files.

## Overview

This example creates a spider chart with:
- **7 axes** (A, B, C, D, E, F, G) with maximum values of 10⁰, 10¹, 10², 10³, 10⁴, 10⁵, 10⁶ respectively
- **5 series** (Series 1 through Series 5) with normally distributed random values for each axis

## Key Features

1. **Programmatic Chart Creation**: Shows how to build a chart entirely from code using the `spider` package API
2. **Dynamic Data Generation**: Uses normally distributed random numbers to generate realistic-looking data
3. **Independent Axis Scales**: Each axis has its own maximum value (powers of 10), demonstrating how spider charts handle different scales

## Code Structure

The example demonstrates:

- Creating a new chart with `spider.NewChart()`
- Setting chart options (width, height, title, subtitle, etc.)
- Adding axes programmatically with `chart.AddAxis()`
- Setting axis maximum values
- Adding series with data using `chart.AddSeries()`
- Generating normally distributed random values
- Saving the chart to both SVG and PNG formats

## Running the Example

To generate the chart:

```bash
# From the repository root
cd examples/from-code
go run main.go
```

This will generate:
- `output.svg` - Vector format
- `output.png` - Raster format

## Data Generation

The example uses a Box-Muller transform to generate normally distributed random numbers:
- Mean: `max / 2` (halfway between 0 and max)
- Standard deviation: `max / 6` (ensures most values stay within [0, max])
- Values are clamped to ensure they stay within the valid range [0, max]

This creates realistic-looking data that demonstrates how the chart handles different scales across axes.

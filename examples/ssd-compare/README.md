# SSD Performance Comparison Example

This example demonstrates how to use spider charts to compare SSDs from different vendors across multiple performance metrics with independent scales. It uses a minimal config file, relying mostly on defaults.

## Use Case

When selecting an SSD, you need to evaluate multiple factors simultaneously:
- **Cost**: Purchase price per unit
- **4k Read IOPS**: Random read performance (important for database workloads)
- **4k Write IOPS**: Random write performance (important for logging and caching)
- **1M Read MBps**: Sequential read throughput (important for large file operations)
- **1M Write MBps**: Sequential write throughput (important for backups and data transfer)

Traditional radar charts require all axes to share the same scale, making it impossible to meaningfully compare metrics with vastly different ranges (e.g., cost in dollars vs. IOPS in thousands). Spider charts with independent scales solve this problem.

## Chart Configuration

This chart compares three SSD vendors:

- **Vendor A (Budget)**: Low-cost option with moderate performance
  - Cost: $120
  - 4k Read IOPS: 25,000
  - 4k Write IOPS: 20,000
  - 1M Read MBps: 500
  - 1M Write MBps: 450

- **Vendor B (Balanced)**: Mid-range option offering good value
  - Cost: $250
  - 4k Read IOPS: 60,000
  - 4k Write IOPS: 55,000
  - 1M Read MBps: 1,800
  - 1M Write MBps: 1,600

- **Vendor C (Performance)**: High-end option with excellent performance
  - Cost: $450
  - 4k Read IOPS: 95,000
  - 4k Write IOPS: 85,000
  - 1M Read MBps: 3,200
  - 1M Write MBps: 2,800

## Generating the Chart

To generate the chart from the configuration file:

```bash
# From the repository root
go run ./cmd/spider-cli -config examples/ssd-compare/config.yaml -output examples/ssd-compare/output.png
```

Or if you've built the CLI tool:

```bash
./spider-cli -config examples/ssd-compare/config.yaml -output examples/ssd-compare/output.png
```

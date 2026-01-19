# SSD Performance Comparison Example

This example demonstrates how to use spider charts to compare SSDs from different vendors across multiple performance metrics with independent scales.

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

## Axis Descriptions

1. **Cost** (0-500 USD): The purchase price of the SSD. Lower is better for budget-conscious buyers.

2. **4k Read IOPS** (0-100,000): Random read operations per second using 4KB blocks. Higher values indicate better performance for database and application workloads.

3. **4k Write IOPS** (0-90,000): Random write operations per second using 4KB blocks. Important for write-heavy workloads like logging and caching.

4. **1M Read MBps** (0-3,500): Sequential read throughput in megabytes per second using 1MB blocks. Higher values mean faster large file reads.

5. **1M Write MBps** (0-3,000): Sequential write throughput in megabytes per second using 1MB blocks. Important for backup and data transfer operations.

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

## Key Insights

This chart makes it easy to see:

- **Price-to-performance ratio**: Vendor B offers the best balance, with roughly 2x the performance of Vendor A for about 2x the cost, while Vendor C provides diminishing returns.

- **Workload suitability**: 
  - For budget-constrained projects: Vendor A
  - For general-purpose workloads: Vendor B
  - For performance-critical applications: Vendor C

- **Independent scales**: Each metric is displayed on its own scale, allowing meaningful comparison across different units (dollars, IOPS, MBps).

## Customization

You can modify `config.yaml` to:
- Add more vendors or remove existing ones
- Adjust axis maximums to focus on specific ranges
- Change colors and styling for each series
- Modify chart size, title, or legend placement
- Switch between linear and logarithmic scales for specific axes

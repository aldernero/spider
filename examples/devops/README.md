# DevOps Software Version Comparison Example

This example demonstrates how spider charts can be used to compare different versions of software across multiple cloud deployment metrics, highlighting typical tradeoffs in software development.

## Use Case

When evaluating software versions for cloud deployment, you need to consider multiple factors:
- **Performance**: Requests per second the system can handle
- **Reliability**: Uptime percentage and system stability
- **Cloud Cost**: Monthly infrastructure costs
- **Scalability**: Maximum concurrent users the system can support
- **Security Score**: Security posture and compliance rating
- **Developer Experience**: Deployment time in minutes (lower is better)

Different versions often make different tradeoffs. This chart helps visualize those tradeoffs to make informed decisions.

## Chart Configuration

This chart compares two software versions:

- **Version 1.0**: Stable, cost-effective baseline
  - Performance: 3,500 req/s
  - Reliability: 95% uptime
  - Cloud Cost: $1,200/month
  - Scalability: 25,000 concurrent users
  - Security Score: 85/100
  - Developer Experience: 45 minutes deployment time

- **Version 2.0**: Enhanced version with improved performance
  - Performance: 8,500 req/s (2.4x improvement)
  - Reliability: 99.5% uptime (improved)
  - Cloud Cost: $3,800/month (3.2x increase)
  - Scalability: 85,000 concurrent users (3.4x improvement)
  - Security Score: 95/100 (improved)
  - Developer Experience: 15 minutes deployment time (3x faster)

## Axis Descriptions

1. **Performance** (0-10,000 req/s): Requests per second the system can handle. Higher values indicate better throughput and responsiveness.

2. **Reliability** (0-100%): System uptime percentage. Higher values mean more stable and available systems.

3. **Cloud Cost** (0-5,000 USD/month): Monthly infrastructure costs including compute, storage, and networking. Lower is better for cost optimization.

4. **Scalability** (0-100,000 users): Maximum number of concurrent users the system can support. Higher values indicate better ability to handle growth.

5. **Security Score** (0-100): Security posture rating based on compliance, vulnerability management, and best practices. Higher values indicate better security.

6. **Developer Experience** (0-60 minutes): Time required to deploy the application. Lower values indicate faster deployment cycles and better developer productivity.

## Generating the Chart

To generate the chart from the configuration file:

```bash
# From the repository root
go run ./cmd/spider-cli -config examples/devops/config.yaml -output examples/devops/output.png
```

Or if you've built the CLI tool:

```bash
./spider-cli -config examples/devops/config.yaml -output examples/devops/output.png
```

## Key Insights

This chart clearly illustrates the tradeoffs between versions:

- **Performance vs Cost**: Version 2.0 provides significantly better performance but at 3.2x the cost. The decision depends on whether the performance gains justify the additional expense.

- **Reliability Improvement**: Version 2.0 shows improved reliability (99.5% vs 95%), which may be critical for production systems.

- **Scalability**: Version 2.0 can handle 3.4x more concurrent users, making it better suited for growing applications.

- **Developer Experience**: Version 2.0 has much faster deployment times (15 vs 45 minutes), improving development velocity.

- **Security**: Both versions have good security scores, but Version 2.0 shows improvement, which may be important for compliance requirements.

## Decision Making

Use this chart to help answer questions like:
- Is the performance improvement worth the cost increase?
- Does the improved reliability justify the migration effort?
- Will the scalability gains support projected growth?
- How important is faster deployment for your team?

The independent scales allow you to compare metrics with different units and ranges, making it easier to see the full picture of tradeoffs.

## Customization

You can modify `config.yaml` to:
- Add more versions for comparison
- Adjust axis maximums to focus on specific ranges
- Change colors and styling for each version
- Modify chart size, title, or legend placement
- Use logarithmic scales for metrics that span orders of magnitude (e.g., scalability)

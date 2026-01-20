# DevOps Software Version Comparison Example

This example demonstrates how spider charts can be used to compare different versions of software across multiple cloud deployment metrics, highlighting typical tradeoffs in software development. It shows
extra customization in the config file compared to ssd-compare example.

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

## Customization

You can modify `config.yaml` to:
- Add more versions for comparison
- Adjust axis maximums to focus on specific ranges
- Change colors and styling for each version
- Modify chart size, title, or legend placement
# Driftlock Claude Code Plugin

Real-time compression-based anomaly detection for Claude Code sessions.

## Features

- **Code Change Detection**: Analyzes code changes for anomalies using compression-based algorithms
- **AI Interaction Monitoring**: Tracks Claude responses and tool usage patterns
- **Real-time Alerts**: Immediate notification of detected anomalies
- **Configurable Thresholds**: Customize detection sensitivity per project
- **Cost-Effective**: Intelligent batching reduces API calls
- **Privacy-First**: Local processing with optional cloud analysis

## Quick Start

1. **Install the Plugin**
   ```bash
   # In your Claude Code workspace
   claude plugin install driftlock
   ```

2. **Configure API Keys**
   ```bash
   export DRIFTLOCK_API_KEY="your-api-key"
   export DRIFTLOCK_TENANT_ID="your-tenant-id"
   ```

3. **Start Monitoring**
   ```bash
   # Monitor current session
   /monitor

   # Analyze recent changes
   /detect changes

   # Check status
   /status
   ```

## Configuration

Create `.driftlock.yaml` in your project root:

```yaml
detection:
  model: haiku-4.5          # AI model for analysis
  threshold: 0.7           # Anomaly confidence threshold
  batch_size: 50           # Events to batch

files:
  include:
    - "**/*.go"
    - "**/*.js"
    - "**/*.ts"
    - "**/*.py"
  exclude:
    - "**/*_test.go"
    - "node_modules/**"

notifications:
  real_time: true
  severity_threshold: 0.8
```

## Commands

### `/detect [scope]`
Detect anomalies in code changes.
- `scope`: `changes` (default), `session`, `file:path`

### `/monitor [--threshold=X] [--batch-size=N]`
Start real-time monitoring of the session.
- `threshold`: Anomaly detection threshold (0.0-1.0)
- `batch-size`: Number of events to batch

### `/status [--format=table|json]`
Show current monitoring status and recent anomalies.

### `/config [key] [value]`
View or update configuration.
- Examples: `/config model sonnet`, `/config threshold 0.8`

## Integration Examples

### GitHub Actions
```yaml
- name: Driftlock Analysis
  uses: driftlock/action@v1
  with:
    api-key: ${{ secrets.DRIFTLOCK_API_KEY }}
    analyze-changes: true
    model: haiku-4.5
    fail-on-severe: true
```

### VS Code
The plugin integrates with VS Code's Driftlock extension for:
- Real-time anomaly indicators
- Code annotations
- Detailed reports in sidebar

## Pricing

- **Free Tier**: 10,000 events/month, 7-day retention
- **Radar**: $5/month + AI costs, 500K events
- **Lock**: $15/month + AI costs, 5M events
- **Orbit**: $50/month + AI costs, unlimited

AI costs are passed through with a 15% margin.

## Privacy & Security

- All baseline models stored locally
- Only anomalies sent to cloud (optional)
- End-to-end encryption for API communications
- PII automatically filtered
- Audit trail maintained locally

## Support

- Documentation: https://docs.driftlock.ai
- Issues: https://github.com/driftlock/claude-code-plugin/issues
- Community: https://discord.gg/driftlock

## License

MIT License - see LICENSE file for details.
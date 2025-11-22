# CLI Reference

The Driftlock CLI (`driftlock`) is the fundamental tool for local anomaly detection and pipeline integration. It brings the power of the Universal Radar to your terminal.

## Installation

```bash
curl -sL https://driftlock.net/install.sh | bash
```

## Commands

### `driftlock scan`

Scans a stream of data (stdin) for anomalies.

```bash
# Pipe logs directly into driftlock
tail -f /var/log/nginx/access.log | driftlock scan --format=nginx

# Scan a JSON file
cat transactions.json | driftlock scan --window=100
```

**Flags:**
*   `--window, -w`: Sliding window size (default: 50).
*   `--format, -f`: Input format (json, nginx, csv, syslog).
*   `--threshold, -t`: Sensitivity threshold (0.0 - 1.0). Default: 0.5.

### `driftlock auth`

Authenticates the CLI with your Shannon Labs account.

```bash
driftlock auth login
```

### `driftlock config`

Manage local configuration and entropy models.

```bash
driftlock config set-tier radar
```

## Exit Codes

*   `0`: Nominal. No anomalies detected.
*   `1`: Anomaly Detected. (Useful for CI/CD pipelines).
*   `2`: System Error.

---

*© 2025 Shannon Labs. See Bad. Stop Bad.*

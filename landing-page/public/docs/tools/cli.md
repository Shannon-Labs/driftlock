# Driftlock CLI

The Driftlock Command Line Interface (CLI) allows you to test, manage, and debug your anomaly detection streams directly from your terminal.

## Installation

Install via npm:

```bash
npm install -g @driftlock/cli
```

Or run via npx:

```bash
npx @driftlock/cli <command>
```

## Authentication

Authenticate with your API key.

```bash
driftlock login
# Enter your API key when prompted
```

Or set the `DRIFTLOCK_API_KEY` environment variable.

## Commands

### `detect`

Send a single event or a file of events for anomaly detection.

```bash
# Single event
driftlock detect --stream my-stream --json '{"value": 100}'

# From file (newline delimited JSON)
driftlock detect --stream logs --file ./logs.jsonl
```

### `stream`

Stream events from stdin for real-time detection.

```bash
# Pipe logs to Driftlock
tail -f /var/log/nginx/access.log | driftlock stream --format nginx
```

### `history`

View recent anomalies for a stream.

```bash
driftlock history --stream payment-service --limit 10
```

### `config`

Manage CLI configuration.

```bash
driftlock config set output json
```

## Usage Examples

### Testing a New Stream

1. Create a file `test-events.json` with normal data:
   ```json
   {"latency": 100}
   {"latency": 102}
   {"latency": 98}
   ```

2. Feed it to Driftlock:
   ```bash
   driftlock detect --stream test-stream --file test-events.json
   ```

3. Send an anomaly:
   ```bash
   driftlock detect --stream test-stream --json '{"latency": 5000}'
   ```

4. Verify detection:
   ```
   [ANOMALY DETECTED]
   Confidence: 0.99
   Why: Significant deviation from recent pattern
   ```

### CI/CD Integration

Use the CLI in your CI pipeline to verify that recent changes haven't introduced performance anomalies.

```yaml
# .github/workflows/test.yml
- name: Check for Performance Regressions
  run: |
    npm install -g @driftlock/cli
    ./run-load-tests.sh > metrics.jsonl
    driftlock detect --stream ci-performance --file metrics.jsonl --fail-on-anomaly
```

## Reference

Run `driftlock --help` for a full list of commands and options.

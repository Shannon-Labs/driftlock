# Driftlock Playground

Run a self-serve UI to test `/v1/detect` with JSON/NDJSON input. Paste data (or load the bundled sample) and the playground will auto-tune detector settings, call the API, and surface anomalies for you - no parameter fiddling required.

## Configuration

The playground connects to the Driftlock API. Configure the API URL via environment variables:

### Local Development
```bash
cd playground
npm install
cp .env.example .env
# Edit .env and set: VITE_API_BASE_URL=http://localhost:8080
npm run dev
```

### Production/Deployed API
```bash
# Edit .env and set: VITE_API_BASE_URL=https://api.driftlock.net
# Or set it when building:
VITE_API_BASE_URL=https://api.driftlock.net npm run build
```

## Development
```bash
cd playground
npm install
cp .env.example .env # set VITE_API_BASE_URL
npm run dev
```

## Build
```bash
npm run build
npm run preview
```

## Plug-and-Play Flow

- **Auto sample:** on first load we fetch `public/samples/payments.ndjson` (derived from `test-data/financial-demo.json`) so the experience works even before you paste data. Swap the default by replacing that file.
- **Auto-run:** every new paste/upload or sample selection parses the payload, derives parameters, and re-runs detection after a short debounce. Advanced overrides also re-run automatically.
- **Validation:** invalid JSON/NDJSON is surfaced immediately, before we ever call the API.

## Automatic Parameter Derivation

- Baseline size defaults to `min(400, max(50, floor(events * 0.2)))`, additionally clamped so it never exceeds the available events.
- Window and hop are pinned to `1`, matching the current server defaults.
- Algorithm defaults to `zstd`. Overrides are optional and tucked behind an “Advanced settings” accordion.
- We only add query parameters to the API (and the Curl snippet) when they differ from the server defaults, so the shared commands stay concise.

## API Health Check

The playground automatically checks API connectivity on load, shows the status badge in the header, and re-checks every 30 seconds before issuing new detection requests.



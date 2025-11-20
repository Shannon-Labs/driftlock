# Postman Collection

Explore the Driftlock API interactively using our official Postman Collection.

## Quick Start

1. **[Run in Postman](https://god.gw.postman.com/run-collection/YOUR_COLLECTION_ID)** (Click to import)
2. Or download the [collection JSON file](/downloads/driftlock.postman_collection.json)

## Setup

After importing the collection:

1. Select the **"Driftlock API"** environment (or create a new one).
2. Set the `api_key` variable to your Driftlock API key.
3. Set the `base_url` to `https://driftlock-api-o6kjgrsowq-uc.a.run.app`.

## Included Requests

### Authentication
- **Verify Key**: Check if your API key is valid and see permissions.

### Detection
- **Detect Anomaly (Single)**: Send a single event for analysis.
- **Detect Anomaly (Batch)**: Send multiple events in one request.
- **Detect Anomaly (Stream)**: Simulate a stream of events.

### Management
- **List Anomalies**: Get a history of detected anomalies.
- **Get Anomaly Details**: View full diagnostics for a specific anomaly.
- **Get Usage**: Check your current API usage and limits.

## Example Workflow

1. **Authenticate**: Run the "Verify Key" request to ensure you're connected.
2. **Establish Baseline**: Run the "Detect Anomaly (Batch)" request with 10-20 "normal" events to help the system learn the baseline for a new stream ID.
3. **Test Anomaly**: Modify one event in the "Detect Anomaly (Single)" request to be significantly different (e.g., change a value by 10x).
4. **Verify**: Check the response to see `detected: true` and a high `confidence` score.

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `base_url` | API Endpoint | `https://driftlock-api-o6kjgrsowq-uc.a.run.app` |
| `api_key` | Your Secret Key | `dl_sk_...` |
| `stream_id` | Default Stream ID | `test-stream` |

## Scripts

The collection includes pre-request and test scripts to:
- Automatically generate timestamps
- Validate JSON schemas
- Visualize anomaly scores in the "Visualize" tab

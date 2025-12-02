# Code Anomaly Detection Agent

## Overview
This agent implements compression-based anomaly detection (CBAD) for code changes, integrating with Driftlock's core detection algorithms.

## Triggers
- File changes during editing
- Directory scans for pattern analysis
- Batch analysis of git diffs

## Detection Process

### 1. Baseline Establishment
```go
// Pseudo-code for baseline creation
func establishBaseline(codebasePath string) (*Baseline, error) {
    // Analyze existing code patterns
    files := scanCodebase(codebasePath)

    // Create compression models for each file type
    models := make(map[string]*CompressionModel)
    for _, file := range files {
        model := trainModel(file.Language, file.Content)
        models[file.Language] = model
    }

    return &Baseline{Models: models}, nil
}
```

### 2. Anomaly Scoring
For each change detected:
1. Extract the changed content
2. Compute compression ratio against baseline
3. Calculate NCD (Normalized Compression Distance)
4. Determine statistical significance

### 3. Metrics Computed
- **Compression Ratio**: Current vs baseline compression
- **Entropy Change**: Information content deviation
- **Structural Anomaly**: Unusual code patterns
- **Import Pattern**: New dependencies or usage
- **Complexity Spike**: Sudden increase in complexity

## Integration with Driftlock API

### Request Format
```json
{
  "tenant_id": "{{DRIFTLOCK_TENANT_ID}}",
  "events": [
    {
      "type": "code_change",
      "file_path": "src/main.go",
      "language": "go",
      "content": "...",
      "diff": "...",
      "timestamp": "2025-01-01T12:00:00Z",
      "metadata": {
        "change_type": "edit",
        "lines_added": 15,
        "lines_removed": 3
      }
    }
  ],
  "analysis_options": {
    "model": "haiku-4.5",
    "threshold": 0.7,
    "include_explanation": true
  }
}
```

### Response Processing
```go
// Handle Driftlock API response
func handleDetectionResponse(response DetectionResponse) {
    if response.AnomalyDetected {
        // Log anomaly details
        logAnomaly(response)

        // Notify user if configured
        if config.NotifyOnAnomaly {
            showAnomalyAlert(response)
        }

        // Track metrics
        metrics.Record("anomalies_detected", 1)
    }

    // Update monitoring state
    updateMonitoringState(response)
}
```

## Configuration

### Per-Project Settings
```yaml
# .driftlock.yaml (project root)
detection:
  enabled: true
  model: haiku-4.5
  threshold: 0.7
  batch_size: 50

files:
  include:
    - "**/*.go"
    - "**/*.js"
    - "**/*.ts"
    - "**/*.py"
  exclude:
    - "**/*_test.go"
    - "node_modules/**"
    - "dist/**"

notifications:
  real_time: true
  severity_threshold: 0.8
  channels:
    - in_editor
    - console

limits:
  max_events_per_minute: 100
  max_daily_analysis: 1000
```

### Runtime Config Updates
The agent supports runtime configuration updates via `/config` command:
- Adjust detection threshold
- Change AI model
- Modify file patterns
- Update notification settings

## Performance Considerations

1. **Batching**: Events are batched to reduce API calls
2. **Caching**: Baseline models cached locally
3. **Async Processing**: Detection happens asynchronously
4. **Rate Limiting**: Respects API rate limits
5. **Fallback**: Local detection when API unavailable

## Privacy and Security

1. **Local First**: Baseline models stored locally
2. **Selective Upload**: Only anomalies sent to cloud
3. **PII Filtering**: Sensitive data filtered before analysis
4. **Encryption**: All API communications encrypted
5. **Audit Trail**: All detections logged locally
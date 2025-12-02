# AI Interaction Monitor Agent

## Overview
This agent monitors Claude Code interactions to detect anomalies in AI behavior, response patterns, and tool usage that might indicate issues or optimization opportunities.

## Monitoring Scope

### 1. Response Patterns
- **Response Length**: Unusually long/short responses
- **Token Usage**: Spike in token consumption
- **Response Time**: Latency anomalies
- **Repetition**: Duplicate or near-duplicate content
- **Error Rate**: Increase in failures or retries

### 2. Tool Usage Patterns
- **Tool Call Frequency**: Excessive tool usage
- **Tool Sequence**: Unusual tool combination patterns
- **Parameter Anomalies**: Unexpected parameter values
- **Failure Patterns**: Repeated tool failures

### 3. Code Generation Quality
- **Syntax Errors**: Increase in generated code errors
- **Security Issues**: Detection of potential security vulnerabilities
- **Performance Issues**: Generated code with performance anti-patterns
- **Best Practice Deviation**: Departure from coding standards

## Implementation Details

### Event Capture
```go
type InteractionEvent struct {
    SessionID    string    `json:"session_id"`
    Timestamp    time.Time `json:"timestamp"`
    EventType    string    `json:"event_type"` // request, response, tool_call
    Content      string    `json:"content,omitempty"`
    TokensUsed   int       `json:"tokens_used,omitempty"`
    ResponseTime int       `json:"response_time,omitempty"`
    Tools        []ToolUse `json:"tools,omitempty"`
    ErrorCode    string    `json:"error_code,omitempty"`
}

type ToolUse struct {
    Name      string                 `json:"name"`
    Arguments map[string]interface{} `json:"arguments"`
    Success   bool                   `json:"success"`
    Duration  int                    `json:"duration_ms"`
}
```

### Anomaly Detection Logic
```go
func (m *InteractionMonitor) detectAnomalies(event InteractionEvent) (*AnomalyResult, error) {
    // 1. Check for response anomalies
    if event.EventType == "response" {
        if err := m.checkResponseAnomalies(event); err != nil {
            return &AnomalyResult{
                Type: "response_anomaly",
                Severity: calculateSeverity(event),
                Details: err.Error(),
            }, nil
        }
    }

    // 2. Analyze tool usage patterns
    if len(event.Tools) > 0 {
        if anomaly := m.analyzeToolUsage(event.Tools); anomaly != nil {
            return anomaly, nil
        }
    }

    // 3. Check for session-level anomalies
    if anomaly := m.checkSessionAnomalies(event.SessionID, event); anomaly != nil {
        return anomaly, nil
    }

    return nil, nil
}
```

## Specific Anomaly Types

### 1. Token Consumption Spike
```go
func (m *InteractionMonitor) checkTokenSpike(event InteractionEvent) bool {
    recent := m.getRecentTokens(event.SessionID, 5*time.Minute)
    average := calculateAverage(recent)

    // Alert if current request is 3x average
    return event.TokensUsed > average*3
}
```

### 2. Repetitive Pattern Detection
```go
func (m *InteractionMonitor) detectRepetition(content string) bool {
    // Hash content for comparison
    hash := hashContent(content)

    // Check recent responses for similarity
    for _, recent := range m.recentResponses {
        similarity := calculateJaccard(hash, recent.hash)
        if similarity > 0.9 { // 90% similarity threshold
            return true
        }
    }

    return false
}
```

### 3. Tool Call Loop Detection
```go
func (m *InteractionMonitor) detectToolLoop(tools []ToolUse) bool {
    // Build tool sequence
    sequence := make([]string, len(tools))
    for i, tool := range tools {
        sequence[i] = tool.Name
    }

    // Check for repeating patterns
    return hasRepeatingPattern(sequence, 3) // 3+ repetitions
}
```

## Integration with Driftlock

### API Integration
```json
{
  "tenant_id": "{{DRIFTLOCK_TENANT_ID}}",
  "events": [
    {
      "type": "ai_interaction",
      "session_id": "sess_123",
      "interaction_type": "response",
      "metrics": {
        "response_time_ms": 2500,
        "tokens_used": 1250,
        "tools_called": 3,
        "error_count": 0
      },
      "anomaly_score": 0.82,
      "content_hash": "abc123",
      "timestamp": "2025-01-01T12:00:00Z"
    }
  ]
}
```

### Local Fallback Detection
When API is unavailable, maintain local detection:
```go
func (m *InteractionMonitor) localDetection(event InteractionEvent) {
    // Maintain rolling statistics
    m.updateStats(event)

    // Check against local thresholds
    if m.stats.ResponseTime.Avg*2 < event.ResponseTime {
        m.notifyLocalAnomaly("response_time_spike")
    }

    if m.stats.TokensPerRequest.Avg*3 < event.TokensUsed {
        m.notifyLocalAnomaly("token_consumption_spike")
    }
}
```

## Configuration

### Monitoring Profiles
```yaml
# .driftlock-monitoring.yaml
monitoring:
  profile: "default" # default, strict, relaxed

thresholds:
  response_time_ms: 5000      # Alert if > 5s
  tokens_per_request: 10000   # Alert if > 10k tokens
  tool_calls_per_request: 10  # Alert if > 10 tools
  error_rate_percent: 10      # Alert if > 10% errors

patterns:
  detect_repetition: true
  detect_loops: true
  detect_security_issues: true

batching:
  max_events: 50
  flush_interval: 30s

local_fallback:
  enabled: true
  stats_window: 100  # Keep last 100 events
```

## Privacy Considerations

1. **Content Hashing**: Only store hashes, not actual content
2. **Local Stats**: Metrics computed locally, only anomalies sent
3. **User Consent**: Clear opt-in for monitoring
4. **Data Minimization**: Send minimum necessary data
5. **Retention**: Limited local retention period

## Notifications

### In-Editor Alerts
```go
type AnomalyAlert struct {
    Type        string    `json:"type"`
    Severity    float64   `json:"severity"`
    Message     string    `json:"message"`
    Suggestion  string    `json:"suggestion"`
    Timestamp   time.Time `json:"timestamp"`
    SessionID   string    `json:"session_id"`
}

// Example alerts
{
  "type": "response_time_spike",
  "severity": 0.75,
  "message": "Response time is 3x normal average",
  "suggestion": "Consider breaking down complex requests"
}
```

### Dashboard Integration
Metrics exposed for dashboard:
- Response time distribution
- Token usage trends
- Tool usage patterns
- Anomaly frequency
- Error rates

## Performance Impact

1. **Minimal Overhead**: < 1ms per event
2. **Async Processing**: Non-blocking
3. **Memory Efficient**: Fixed-size buffers
4. **Smart Sampling**: Sample high-frequency events
5. **Background Sync**: Batch API calls
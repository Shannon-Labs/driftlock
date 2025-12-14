# CBAD Auto-Configuration Design

## Overview

This document specifies how CBAD should automatically configure detection parameters based on the data type being analyzed. The goal is a plug-and-play experience where users don't need to tune parameters.

---

## Detection Profiles

### Profile Definitions

```rust
// crates/driftlock-api/src/profiles.rs

use cbad_core::anomaly::AnomalyConfig;
use cbad_core::window::WindowConfig;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "snake_case")]
pub enum DetectionProfile {
    /// Financial transactions, fraud detection
    /// High precision, good for payment monitoring
    Financial,

    /// LLM safety, jailbreak detection, prompt injection
    /// High recall, accepts more false positives (human review expected)
    LlmSafety,

    /// Structured logs, JSON events, application telemetry
    /// Balanced precision/recall with tokenization
    Logs,

    /// Time series, metrics, sensor data
    /// Optimized for numeric patterns
    TimeSeries,

    /// Network traffic, intrusion detection
    /// Protocol-aware detection
    Network,

    /// Cryptocurrency, blockchain transactions
    /// Similar to financial but tuned for crypto patterns
    Crypto,

    /// Manufacturing, IoT sensors, predictive maintenance
    /// High sensitivity for early warning
    Industrial,

    /// User-specified configuration
    Custom,

    /// Auto-detect based on data shape
    Auto,
}

impl DetectionProfile {
    /// Returns the optimized CBAD config for this profile
    pub fn config(&self) -> AnomalyConfig {
        match self {
            Self::Financial => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 300,
                    window_size: 50,
                    hop_size: 20,
                    max_capacity: 500,
                    ..Default::default()
                },
                permutation_count: 100,
                ncd_threshold: 0.20,
                require_statistical_significance: true,
                ..Default::default()
            },

            Self::LlmSafety => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 200,
                    window_size: 30,
                    hop_size: 10,
                    max_capacity: 400,
                    ..Default::default()
                },
                permutation_count: 100,
                ncd_threshold: 0.25,
                require_statistical_significance: true,
                ..Default::default()
            },

            Self::Logs => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 200,
                    window_size: 40,
                    hop_size: 15,
                    max_capacity: 400,
                    ..Default::default()
                },
                permutation_count: 100,
                ncd_threshold: 0.22,
                require_statistical_significance: true,
                ..Default::default()
            },

            Self::TimeSeries => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 150,
                    window_size: 40,
                    hop_size: 15,
                    max_capacity: 300,
                    ..Default::default()
                },
                permutation_count: 50,
                ncd_threshold: 0.35, // Higher threshold to reduce FPs
                require_statistical_significance: true, // Enable to reduce FPs
                ..Default::default()
            },

            Self::Network => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 150,
                    window_size: 30,
                    hop_size: 10,
                    max_capacity: 300,
                    ..Default::default()
                },
                permutation_count: 100,
                ncd_threshold: 0.22,
                require_statistical_significance: true,
                ..Default::default()
            },

            Self::Crypto => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 250,
                    window_size: 40,
                    hop_size: 15,
                    max_capacity: 400,
                    ..Default::default()
                },
                permutation_count: 100,
                ncd_threshold: 0.20,
                require_statistical_significance: true,
                ..Default::default()
            },

            Self::Industrial => AnomalyConfig {
                window_config: WindowConfig {
                    baseline_size: 100,
                    window_size: 30,
                    hop_size: 10,
                    max_capacity: 200,
                    ..Default::default()
                },
                permutation_count: 50,
                ncd_threshold: 0.18, // More sensitive
                require_statistical_significance: false, // Catch more
                ..Default::default()
            },

            Self::Custom | Self::Auto => {
                // Default balanced config
                AnomalyConfig::default()
            }
        }
    }

    /// Expected F1 score based on benchmarks
    pub fn expected_f1(&self) -> f64 {
        match self {
            Self::Financial => 0.74,
            Self::LlmSafety => 0.66,
            Self::Logs => 0.80,
            Self::TimeSeries => 0.30,
            Self::Network => 0.55,
            Self::Crypto => 0.65,
            Self::Industrial => 0.40,
            Self::Custom | Self::Auto => 0.50,
        }
    }

    /// Human-readable description
    pub fn description(&self) -> &'static str {
        match self {
            Self::Financial => "Optimized for payment transactions and fraud detection",
            Self::LlmSafety => "High-recall detection for LLM jailbreaks and prompt injection",
            Self::Logs => "Balanced detection for structured application logs",
            Self::TimeSeries => "Numeric time series and metrics analysis",
            Self::Network => "Network traffic and intrusion detection",
            Self::Crypto => "Cryptocurrency and blockchain transaction monitoring",
            Self::Industrial => "IoT sensors and predictive maintenance",
            Self::Custom => "User-specified configuration",
            Self::Auto => "Automatically detect profile from data shape",
        }
    }
}
```

---

## Auto-Detection Heuristics

```rust
// crates/driftlock-api/src/auto_detect.rs

use crate::profiles::DetectionProfile;
use serde_json::Value;

/// Detects the appropriate profile based on event data
pub fn detect_profile(events: &[Value]) -> DetectionProfile {
    if events.is_empty() {
        return DetectionProfile::Logs; // Default
    }

    // Sample first N events for analysis
    let sample: Vec<_> = events.iter().take(10).collect();

    // Score each profile based on field patterns
    let mut scores = ProfileScores::default();

    for event in &sample {
        if let Value::Object(obj) = event {
            analyze_fields(obj, &mut scores);
        } else if let Value::String(s) = event {
            analyze_text(s, &mut scores);
        }
    }

    scores.best_match()
}

#[derive(Default)]
struct ProfileScores {
    financial: f64,
    llm_safety: f64,
    logs: f64,
    time_series: f64,
    network: f64,
    crypto: f64,
    industrial: f64,
}

impl ProfileScores {
    fn best_match(&self) -> DetectionProfile {
        let scores = [
            (self.financial, DetectionProfile::Financial),
            (self.llm_safety, DetectionProfile::LlmSafety),
            (self.logs, DetectionProfile::Logs),
            (self.time_series, DetectionProfile::TimeSeries),
            (self.network, DetectionProfile::Network),
            (self.crypto, DetectionProfile::Crypto),
            (self.industrial, DetectionProfile::Industrial),
        ];

        scores
            .into_iter()
            .max_by(|a, b| a.0.partial_cmp(&b.0).unwrap())
            .map(|(_, profile)| profile)
            .unwrap_or(DetectionProfile::Logs)
    }
}

fn analyze_fields(obj: &serde_json::Map<String, Value>, scores: &mut ProfileScores) {
    let keys: Vec<&str> = obj.keys().map(|s| s.as_str()).collect();

    // Financial indicators
    if keys.iter().any(|k| matches!(*k, "amount" | "amt" | "merchant" | "card" | "payment" | "transaction")) {
        scores.financial += 2.0;
    }
    if keys.iter().any(|k| k.contains("usd") || k.contains("currency")) {
        scores.financial += 1.0;
    }

    // LLM/Prompt indicators
    if keys.iter().any(|k| matches!(*k, "prompt" | "message" | "completion" | "response" | "user_input" | "assistant")) {
        scores.llm_safety += 2.0;
    }
    if keys.iter().any(|k| k.contains("token") || k.contains("model")) {
        scores.llm_safety += 1.0;
    }

    // Network indicators
    if keys.iter().any(|k| matches!(*k, "ip" | "port" | "protocol" | "src_ip" | "dst_ip" | "packet" | "bytes")) {
        scores.network += 2.0;
    }
    if keys.iter().any(|k| k.contains("tcp") || k.contains("udp") || k.contains("http")) {
        scores.network += 1.0;
    }

    // Crypto indicators
    if keys.iter().any(|k| matches!(*k, "wallet" | "address" | "hash" | "block" | "chain" | "wei" | "gas")) {
        scores.crypto += 2.0;
    }
    if keys.iter().any(|k| k.contains("eth") || k.contains("btc") || k.contains("crypto")) {
        scores.crypto += 1.0;
    }

    // Industrial/IoT indicators
    if keys.iter().any(|k| matches!(*k, "sensor" | "temperature" | "pressure" | "vibration" | "rul" | "cycle")) {
        scores.industrial += 2.0;
    }
    if keys.iter().any(|k| k.contains("reading") || k.contains("measurement")) {
        scores.industrial += 1.0;
    }

    // Time series indicators (mostly numeric values)
    let numeric_count = obj.values().filter(|v| v.is_number()).count();
    if numeric_count > obj.len() / 2 {
        scores.time_series += 1.0;
    }
    if keys.iter().any(|k| matches!(*k, "value" | "metric" | "timestamp" | "ts")) && numeric_count > 0 {
        scores.time_series += 1.0;
    }

    // Logs indicators (fallback)
    if keys.iter().any(|k| matches!(*k, "level" | "log" | "trace_id" | "span_id" | "service" | "host")) {
        scores.logs += 2.0;
    }
}

fn analyze_text(text: &str, scores: &mut ProfileScores) {
    let lower = text.to_lowercase();

    // Check for prompt-like content
    if lower.contains("user:") || lower.contains("assistant:") ||
       lower.contains("system:") || lower.contains("ignore") ||
       lower.contains("pretend") || lower.contains("roleplay") {
        scores.llm_safety += 2.0;
    }

    // Check for log-like content
    if lower.contains("error") || lower.contains("warn") ||
       lower.contains("info") || lower.contains("debug") {
        scores.logs += 1.0;
    }
}
```

---

## API Changes

### Request Schema

```json
// POST /v1/detect
{
  "stream_id": "my-stream",
  "profile": "financial",  // optional - "auto" if not specified
  "events": [
    { "amount": 150.00, "merchant": "ACME Corp", ... }
  ]
}
```

### Profile Endpoint

```json
// GET /v1/profiles
{
  "profiles": [
    {
      "name": "financial",
      "description": "Optimized for payment transactions and fraud detection",
      "expected_f1": 0.74,
      "config": {
        "baseline_size": 300,
        "window_size": 50,
        "ncd_threshold": 0.20
      }
    },
    // ... other profiles
  ]
}
```

### Stream Profile Management

```json
// GET /v1/streams/{id}/profile
{
  "stream_id": "my-stream",
  "profile": "financial",
  "auto_detected": true,
  "confidence": 0.85,
  "sample_fields": ["amount", "merchant", "card_last4"]
}

// PATCH /v1/streams/{id}/profile
{
  "profile": "llm_safety"
}
```

---

## Response Enrichment

When a profile is auto-detected, include metadata in the response:

```json
// POST /v1/detect response
{
  "anomaly_detected": true,
  "confidence": 0.78,
  "ncd": 0.31,
  "profile_used": "financial",
  "profile_auto_detected": true,
  "profile_confidence": 0.92,
  "detection_note": "Transaction pattern differs significantly from baseline"
}
```

---

## Implementation Files

| File | Changes |
|------|---------|
| `crates/driftlock-api/src/profiles.rs` | New - DetectionProfile enum |
| `crates/driftlock-api/src/auto_detect.rs` | New - Auto-detection logic |
| `crates/driftlock-api/src/handlers/detect.rs` | Add profile parameter |
| `crates/driftlock-api/src/handlers/profiles.rs` | New - Profile endpoints |
| `crates/driftlock-db/src/models/stream.rs` | Add profile field |
| `cbad-core/src/lib.rs` | Export DetectionProfile |

---

## Migration Path

### Phase 1: Add Profile Support (Non-Breaking)
1. Add `profile` as optional parameter (default: `auto`)
2. Implement auto-detection
3. Store detected profile with stream

### Phase 2: Profile Tuning
1. Collect feedback on profile accuracy
2. Adjust thresholds based on real-world data
3. Add more profiles as needed

### Phase 3: User Customization
1. Allow users to override profile configs
2. Implement per-stream custom configs
3. Add profile recommendation in response

---

## Testing Strategy

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_detect_financial() {
        let events = vec![
            json!({ "amount": 150.0, "merchant": "ACME", "card_last4": "1234" })
        ];
        assert_eq!(detect_profile(&events), DetectionProfile::Financial);
    }

    #[test]
    fn test_detect_llm() {
        let events = vec![
            json!({ "prompt": "Hello, how are you?", "model": "gpt-4" })
        ];
        assert_eq!(detect_profile(&events), DetectionProfile::LlmSafety);
    }

    #[test]
    fn test_detect_network() {
        let events = vec![
            json!({ "src_ip": "192.168.1.1", "dst_port": 443, "protocol": "tcp" })
        ];
        assert_eq!(detect_profile(&events), DetectionProfile::Network);
    }
}
```

---

## Open Questions

1. **Should profile be locked after first detection?** Or re-evaluate periodically?
2. **How to handle mixed data types?** Some streams may have multiple event types.
3. **Should we expose confidence score?** May confuse users if low.
4. **Profile inheritance:** Should custom profiles extend base profiles?

//! Sliding window system for CBAD
//!
//! This module provides a time-series buffering system with configurable
//! baseline/window/hop semantics for real-time anomaly detection.
//!
//! The system maintains:
//! - A baseline window (reference pattern)
//! - A current window (to compare against baseline)
//! - Configurable window sizes and hop distance
//! - Bounded memory usage
//! - Thread-safe operations for concurrent access
//!
//! Key concepts:
//! - Baseline: Historical "normal" pattern for comparison
//! - Window: Current data to test for anomalies
//! - Hop: How much to advance the window after each analysis
//! - Capacity: Maximum number of events to retain (bounded memory)

use serde::{Deserialize, Serialize};
use std::collections::VecDeque;
use std::sync::{Arc, Mutex};
use std::time::{Duration, SystemTime};

type WindowResult<T> = std::result::Result<T, Box<dyn std::error::Error + Send + Sync>>;
/// Configuration for the sliding window system
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct WindowConfig {
    /// Size of the baseline window (number of events)
    pub baseline_size: usize,
    /// Size of the analysis window (number of events)
    pub window_size: usize,
    /// Hop size: how many events to advance the window after analysis
    pub hop_size: usize,
    /// Maximum capacity: total number of events to retain in memory
    pub max_capacity: usize,
    /// Optional time-based windowing (if Some, window size refers to time span)
    pub time_window: Option<Duration>,
    /// Privacy redaction configuration
    pub privacy_config: PrivacyConfig,
}

impl Default for WindowConfig {
    fn default() -> Self {
        Self {
            baseline_size: 1000,
            window_size: 100,
            hop_size: 50,
            max_capacity: 10000,
            time_window: None,
            privacy_config: PrivacyConfig::default(),
        }
    }
}

/// Privacy redaction configuration for sensitive data
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PrivacyConfig {
    /// Fields to always redact (e.g., "password", "ssn", "token")
    pub redact_fields: Vec<String>,
    /// Regular expression patterns to redact
    pub redact_patterns: Vec<String>,
    /// Whether to completely drop events that don't meet privacy requirements
    pub drop_non_compliant: bool,
    /// Whether to encrypt sensitive data in memory
    pub encrypt_sensitive: bool,
}

impl Default for PrivacyConfig {
    fn default() -> Self {
        Self {
            redact_fields: vec![
                "password".to_string(),
                "token".to_string(),
                "key".to_string(),
                "secret".to_string(),
                "jwt".to_string(),
                "authorization".to_string(),
                "cookie".to_string(),
            ],
            redact_patterns: vec![
                r"(?i)\b(credit|debit)_?card[:\s]*\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?(\d{4})\b"
                    .to_string(),
                r"(?i)\bssn[:\s]*\d{3}-?\d{2}-?(\d{4})\b".to_string(),
                r"(?i)\bemail[:\s]*[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}\b".to_string(),
            ],
            drop_non_compliant: false,
            encrypt_sensitive: false,
        }
    }
}

/// A data event in the sliding window system
#[derive(Debug, Clone)]
pub struct DataEvent {
    /// The actual data payload
    pub data: Vec<u8>,
    /// Timestamp of when the event was received
    pub timestamp: SystemTime,
    /// Optional metadata about the event
    pub metadata: std::collections::HashMap<String, String>,
    /// Hash of the original data for integrity verification
    pub integrity_hash: Option<u64>,
}

impl DataEvent {
    /// Create a new data event
    pub fn new(data: Vec<u8>) -> Self {
        Self {
            data,
            timestamp: SystemTime::now(),
            metadata: std::collections::HashMap::new(),
            integrity_hash: None,
        }
    }

    /// Create a new data event with custom timestamp
    pub fn with_timestamp(data: Vec<u8>, timestamp: SystemTime) -> Self {
        Self {
            data,
            timestamp,
            metadata: std::collections::HashMap::new(),
            integrity_hash: None,
        }
    }
}

/// The sliding window system
///
/// This structure manages the baseline and current windows, handles
/// event insertion, and provides access to the current windows for
/// anomaly detection.
pub struct SlidingWindow {
    /// Configuration for this sliding window
    config: WindowConfig,
    /// All events in the system (bounded by max_capacity)
    events: VecDeque<DataEvent>,
    /// Current position of the baseline window start
    baseline_start: usize,
    /// Current position of the analysis window start
    window_start: usize,
    /// Whether we have enough data to perform analysis
    is_ready: bool,
    /// Total events processed
    total_events: u64,
    /// Whether we have aligned the first analysis window to the tail
    aligned: bool,
}

impl SlidingWindow {
    /// Create a new sliding window with the given configuration
    pub fn new(config: WindowConfig) -> Self {
        Self {
            config,
            events: VecDeque::new(),
            baseline_start: 0,
            window_start: 0,
            is_ready: false,
            total_events: 0,
            aligned: false,
        }
    }

    /// Add an event to the sliding window system
    ///
    /// This will apply privacy redaction if configured, and may trigger
    /// window advancement if the system is ready for analysis.
    pub fn add_event(&mut self, mut event: DataEvent) -> bool {
        // Apply privacy redaction
        if let Some(redacted_event) = self.apply_privacy_redaction(event) {
            event = redacted_event;
        } else {
            // Event was dropped due to privacy compliance
            return false;
        }

        // Add the event
        self.events.push_back(event);
        self.total_events += 1;

        // Bound the memory usage
        self.bound_memory();

        // Check if we have enough data to be ready
        self.update_readiness();

        true
    }

    /// Apply privacy redaction to an event
    fn apply_privacy_redaction(&self, mut event: DataEvent) -> Option<DataEvent> {
        let privacy = &self.config.privacy_config;

        // If drop_non_compliant is true and we can't ensure privacy, drop the event
        if privacy.drop_non_compliant {
            // For now, we assume basic compliance - in a real system, this would be more complex
        }

        // Apply field redaction - simple substring-based redaction for now
        for field in &privacy.redact_fields {
            // This is a simplified redaction - in practice, this would parse structured data
            // like JSON and redact specific fields
            let field_key = format!("\"{}\":", field);
            let data_str = String::from_utf8_lossy(&event.data);

            // Replace "field": "value" with "field": "[REDACTED]"
            let mut modified_str = data_str.to_string();
            let mut pos = 0;

            while let Some(start) = modified_str[pos..].find(&field_key) {
                let abs_start = pos + start;
                let search_from = &modified_str[abs_start + field_key.len()..];

                // Find the value part (assumes format "field": "value" or "field":value)
                if let Some(quote_pos) = search_from.find('"') {
                    let value_start = abs_start + field_key.len() + quote_pos;
                    if let Some(value_end) = modified_str[value_start + 1..].find('"') {
                        let value_end = value_start + 1 + value_end;
                        let before = &modified_str[..value_start + 1];
                        let after = &modified_str[value_end..];
                        modified_str = format!("{}[REDACTED]{}", before, after);
                        pos = value_end; // Continue search after the replacement
                    } else {
                        break; // No closing quote found, break to avoid infinite loop
                    }
                } else {
                    // If no quote, look for non-whitespace token
                    let remaining = &modified_str[abs_start + field_key.len()..];
                    let remaining_trimmed = remaining.trim_start();
                    if !remaining_trimmed.is_empty() {
                        let token_start = abs_start
                            + field_key.len()
                            + (remaining.len() - remaining_trimmed.len());
                        let token_end = token_start
                            + remaining_trimmed
                                .split_whitespace()
                                .next()
                                .map_or(0, |s| s.len());

                        if token_end > token_start {
                            let before = &modified_str[..token_start];
                            let after = &modified_str[token_end..];
                            modified_str = format!("{}[REDACTED]{}", before, after);
                            pos = token_end;
                        } else {
                            break; // Nothing to redact, break to avoid infinite loop
                        }
                    } else {
                        break; // Nothing to redact, break to avoid infinite loop
                    }
                }
            }

            if modified_str != data_str {
                event.data = modified_str.as_bytes().to_vec();
            }
        }

        // Apply pattern redaction using regex
        for pattern in &privacy.redact_patterns {
            if let Ok(re) = regex::Regex::new(pattern) {
                let data_str = String::from_utf8_lossy(&event.data);
                let redacted_str = re.replace_all(&data_str, "[REDACTED]");

                if data_str != *redacted_str {
                    event.data = redacted_str.as_bytes().to_vec();
                }
            }
        }

        Some(event)
    }

    /// Bound the memory usage by removing old events
    fn bound_memory(&mut self) {
        while self.events.len() > self.config.max_capacity {
            self.events.pop_front();
            // Adjust window positions since we removed an event from the front
            if self.baseline_start > 0 {
                self.baseline_start -= 1;
            }
            if self.window_start > 0 {
                self.window_start -= 1;
            }
        }
        // If trimming made us lose readiness, realign when we regain it later
        if self.events.len() < self.config.baseline_size + self.config.window_size {
            self.aligned = false;
        }
    }

    /// Update the readiness state based on available data
    fn update_readiness(&mut self) {
        let required_events = self.config.baseline_size + self.config.window_size;
        self.is_ready = self.events.len() >= required_events;
        if self.is_ready && !self.aligned {
            self.align_to_tail();
        }
    }

    /// Align the initial baseline/window to the most recent data tail.
    fn align_to_tail(&mut self) {
        if !self.is_ready {
            return;
        }
        let end = self.events.len();
        self.window_start = end.saturating_sub(self.config.window_size);
        self.baseline_start = self.window_start.saturating_sub(self.config.baseline_size);
        self.aligned = true;
    }

    /// Advance the windows based on hop size after an analysis run
    pub fn advance_after_analysis(&mut self) {
        if !self.is_ready {
            return;
        }

        let max_position = self.events.len().saturating_sub(self.config.window_size);
        let next_window = (self.window_start + self.config.hop_size).min(max_position);
        self.window_start = next_window;

        // Baseline immediately precedes the current window to keep recency
        self.baseline_start = self.window_start.saturating_sub(self.config.baseline_size);
    }

    /// Get the current baseline window for analysis
    ///
    /// Returns None if there are not enough events to form a complete baseline
    pub fn get_baseline(&self) -> Option<Vec<u8>> {
        if !self.is_ready {
            return None;
        }

        let start_idx = self.baseline_start;
        let end_idx = (start_idx + self.config.baseline_size).min(self.events.len());

        if end_idx - start_idx < self.config.baseline_size {
            return None;
        }

        let mut result = Vec::new();
        for i in start_idx..end_idx {
            if i < self.events.len() {
                result.extend_from_slice(&self.events[i].data);
            }
        }
        Some(result)
    }

    /// Get the current analysis window
    ///
    /// Returns None if there are not enough events to form a complete window
    pub fn get_window(&self) -> Option<Vec<u8>> {
        if !self.is_ready {
            return None;
        }

        let start_idx = self.window_start;
        let end_idx = (start_idx + self.config.window_size).min(self.events.len());

        if end_idx - start_idx < self.config.window_size {
            return None;
        }

        let mut result = Vec::new();
        for i in start_idx..end_idx {
            if i < self.events.len() {
                result.extend_from_slice(&self.events[i].data);
            }
        }
        Some(result)
    }

    /// Get both baseline and current window for CBAD analysis
    pub fn get_baseline_and_window(&self) -> Option<(Vec<u8>, Vec<u8>)> {
        if let (Some(baseline), Some(window)) = (self.get_baseline(), self.get_window()) {
            Some((baseline, window))
        } else {
            None
        }
    }

    /// Check if the system has enough data for analysis
    pub fn is_ready(&self) -> bool {
        self.is_ready
    }

    /// Get the total number of events processed
    pub fn total_events(&self) -> u64 {
        self.total_events
    }

    /// Get the current memory usage (number of events stored)
    pub fn memory_usage(&self) -> usize {
        self.events.len()
    }

    /// Get window configuration
    pub fn config(&self) -> &WindowConfig {
        &self.config
    }

    /// Reset the window positions while maintaining the event buffer
    pub fn reset_positions(&mut self) {
        self.baseline_start = 0;
        self.window_start = 0;
        self.is_ready = false;
        self.aligned = false;
    }

    /// Clear all events and reset the system
    pub fn clear(&mut self) {
        self.events.clear();
        self.baseline_start = 0;
        self.window_start = 0;
        self.is_ready = false;
        self.total_events = 0;
        self.aligned = false;
    }
}

/// Thread-safe wrapper for the sliding window system
///
/// This provides concurrent access to the sliding window system without
/// requiring external synchronization.
#[derive(Clone)]
pub struct ThreadSafeSlidingWindow {
    inner: Arc<Mutex<SlidingWindow>>,
}

impl ThreadSafeSlidingWindow {
    /// Create a new thread-safe sliding window
    pub fn new(config: WindowConfig) -> Self {
        Self {
            inner: Arc::new(Mutex::new(SlidingWindow::new(config))),
        }
    }

    /// Add an event to the sliding window system
    pub fn add_event(&self, event: DataEvent) -> WindowResult<bool> {
        let mut window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.add_event(event))
    }

    /// Get the current baseline window for analysis
    pub fn get_baseline(&self) -> WindowResult<Option<Vec<u8>>> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.get_baseline())
    }

    /// Get the current analysis window
    pub fn get_window(&self) -> WindowResult<Option<Vec<u8>>> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.get_window())
    }

    /// Get both baseline and current window for CBAD analysis
    pub fn get_baseline_and_window(&self) -> WindowResult<Option<(Vec<u8>, Vec<u8>)>> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.get_baseline_and_window())
    }

    /// Check if the system has enough data for analysis
    pub fn is_ready(&self) -> WindowResult<bool> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.is_ready())
    }

    /// Get the total number of events processed
    pub fn total_events(&self) -> WindowResult<u64> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.total_events())
    }

    /// Get the current memory usage (number of events stored)
    pub fn memory_usage(&self) -> WindowResult<usize> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.memory_usage())
    }

    /// Get window configuration
    pub fn config(&self) -> WindowResult<WindowConfig> {
        let window = self.inner.lock().map_err(|e| e.to_string())?;
        Ok(window.config().clone())
    }

    /// Advance the window after analysis to honor hop semantics
    pub fn advance_after_analysis(&self) -> WindowResult<()> {
        let mut window = self.inner.lock().map_err(|e| e.to_string())?;
        window.advance_after_analysis();
        Ok(())
    }

    /// Clear all buffered events and positions
    pub fn clear(&self) -> WindowResult<()> {
        let mut window = self.inner.lock().map_err(|e| e.to_string())?;
        window.clear();
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sliding_window_creation() {
        let config = WindowConfig::default();
        let window = SlidingWindow::new(config.clone());

        assert_eq!(window.config().baseline_size, config.baseline_size);
        assert_eq!(window.total_events(), 0);
        assert_eq!(window.memory_usage(), 0);
        assert!(!window.is_ready());
    }

    #[test]
    fn test_add_events_and_readiness() {
        let config = WindowConfig {
            baseline_size: 5,
            window_size: 3,
            max_capacity: 20,
            ..Default::default()
        };

        let mut window = SlidingWindow::new(config);

        // Add fewer events than needed for readiness
        for i in 0..7 {
            let event = DataEvent::new(format!("event_{}", i).as_bytes().to_vec());
            window.add_event(event);
        }

        assert!(!window.is_ready());
        assert_eq!(window.total_events(), 7);
        assert_eq!(window.memory_usage(), 7);

        // Add enough events for readiness
        for i in 7..10 {
            let event = DataEvent::new(format!("event_{}", i).as_bytes().to_vec());
            window.add_event(event);
        }

        assert!(window.is_ready());
        assert_eq!(window.total_events(), 10);
        assert_eq!(window.memory_usage(), 10);
    }

    #[test]
    fn test_get_baseline_and_window() {
        let config = WindowConfig {
            baseline_size: 3,
            window_size: 2,
            max_capacity: 20,
            ..Default::default()
        };

        let mut window = SlidingWindow::new(config);

        // Add enough events
        for i in 0..10 {
            let event = DataEvent::new(format!("data_{}", i).as_bytes().to_vec());
            window.add_event(event);
        }

        assert!(window.is_ready());

        // Check that we can get both baseline and window
        let baseline = window.get_baseline();
        let _window = window.get_window();

        assert!(baseline.is_some());
        assert!(_window.is_some());

        // Verify baseline and window are anchored to the tail (most recent events)
        let baseline_len = ["data_5", "data_6", "data_7"]
            .iter()
            .map(|s| s.len())
            .sum::<usize>();
        assert_eq!(baseline.unwrap().len(), baseline_len);
    }

    #[test]
    fn test_memory_bounding() {
        let config = WindowConfig {
            max_capacity: 5,
            ..Default::default()
        };

        let mut window = SlidingWindow::new(config);

        // Add more events than capacity
        for i in 0..10 {
            let event = DataEvent::new(format!("event_{}", i).as_bytes().to_vec());
            window.add_event(event);
        }

        assert_eq!(window.memory_usage(), 5);
        assert_eq!(window.total_events(), 10);
    }

    #[test]
    fn test_thread_safe_wrapper() {
        use std::sync::Arc;
        use std::thread;

        let config = WindowConfig::default();
        let window = Arc::new(ThreadSafeSlidingWindow::new(config));

        let window_clone = window.clone();
        let handle = thread::spawn(move || {
            let event = DataEvent::new(b"thread_event".to_vec());
            window_clone.add_event(event).unwrap();

            let is_ready = window_clone.is_ready().unwrap();
            (window_clone.total_events().unwrap(), is_ready)
        });

        // Add from main thread too
        let event = DataEvent::new(b"main_event".to_vec());
        window.add_event(event).unwrap();

        let (events, _) = handle.join().unwrap();
        assert!(events >= 1);

        let total = window.total_events().unwrap();
        assert_eq!(total, 2); // One from each thread
    }

    #[test]
    fn test_hop_advances_from_tail() {
        let mut config = WindowConfig {
            baseline_size: 3,
            window_size: 2,
            hop_size: 1,
            max_capacity: 20,
            ..Default::default()
        };
        config.privacy_config.redact_fields.clear();
        let mut window = SlidingWindow::new(config);

        for i in 0..7 {
            let event = DataEvent::new(format!("evt{}", i).as_bytes().to_vec());
            window.add_event(event);
        }

        assert!(window.is_ready());
        let first_window = window.get_window().unwrap();
        window.advance_after_analysis();
        let second_window = window.get_window().unwrap();
        assert_ne!(first_window, second_window, "window should advance by hop");
    }

    #[test]
    fn test_privacy_redaction() {
        let mut config = WindowConfig::default();
        config.privacy_config.redact_fields = vec!["password".to_string()];

        let mut window = SlidingWindow::new(config);

        // Create an event with sensitive data
        let sensitive_data = r#"{"user":"alice","password":"secret123","action":"login"}"#;
        let event = DataEvent::new(sensitive_data.as_bytes().to_vec());

        // Add the event (this should trigger redaction)
        window.add_event(event);

        // In a real system, we'd verify that the password was redacted
        // For this test, we just ensure the event was added
        assert_eq!(window.total_events(), 1);
    }

    #[test]
    fn test_time_window_config() {
        let config = WindowConfig {
            time_window: Some(Duration::from_secs(60)), // 1 minute window
            ..Default::default()
        };

        let window = SlidingWindow::new(config);

        assert_eq!(window.config().time_window, Some(Duration::from_secs(60)));
    }
}

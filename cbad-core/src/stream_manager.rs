use crate::anomaly::AnomalyConfig;
use crate::api::{CbadDetector, CbadDetectorBuilder, DetectionProfile, DetectionRecord};
use crate::error::{CbadError, Result};
use crate::storage::{AnomalyFilter, AnomalyStore, CorrelatedAnomalies, InMemoryAnomalyStore};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::{Duration, SystemTime};

/// Manage multiple independent detectors and route events by stream id.
pub struct StreamManager {
    streams: Mutex<HashMap<String, CbadDetector>>,
    store: Arc<dyn AnomalyStore>,
}

impl Default for StreamManager {
    fn default() -> Self {
        Self::new()
    }
}

impl StreamManager {
    pub fn new() -> Self {
        Self::with_store(Arc::new(InMemoryAnomalyStore::default()))
    }

    pub fn with_store(store: Arc<dyn AnomalyStore>) -> Self {
        Self {
            streams: Mutex::new(HashMap::new()),
            store,
        }
    }

    /// Create a new named stream with its own detector/baseline.
    pub fn create_stream(
        &self,
        name: impl Into<String>,
        profile: DetectionProfile,
        config: Option<AnomalyConfig>,
    ) -> Result<()> {
        let name = name.into();
        let detector = CbadDetectorBuilder::new()
            .with_profile(profile)
            .with_config(config.unwrap_or_default())
            .with_store(self.store.clone())
            .named(&name)
            .build()?;

        let mut guard = self
            .streams
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;
        guard.insert(name, detector);
        Ok(())
    }

    /// Route an event to a stream; returns an anomaly if detected.
    pub fn ingest(&self, stream: &str, event: Vec<u8>) -> Result<Option<DetectionRecord>> {
        let guard = self
            .streams
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;
        let detector = guard
            .get(stream)
            .ok_or_else(|| CbadError::InvalidConfig(format!("unknown stream {}", stream)))?;
        detector.ingest_event(event)
    }

    /// Retrieve correlated anomalies across streams within the given window.
    pub fn correlate_anomalies(&self, window: Duration) -> Result<Vec<CorrelatedAnomalies>> {
        let since = SystemTime::now()
            .checked_sub(window)
            .ok_or_else(|| CbadError::InvalidConfig("invalid correlation window".into()))?;

        let results = futures::executor::block_on(self.store.query(AnomalyFilter {
            since: Some(since),
            until: None,
            stream: None,
            only_anomalies: true,
        }))?;

        Ok(vec![CorrelatedAnomalies {
            window,
            anomalies: results,
        }])
    }

    /// List all stream names currently managed.
    pub fn list_streams(&self) -> Result<Vec<String>> {
        let guard = self
            .streams
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;
        Ok(guard.keys().cloned().collect())
    }

    /// Check if a stream exists.
    pub fn has_stream(&self, name: &str) -> Result<bool> {
        let guard = self
            .streams
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;
        Ok(guard.contains_key(name))
    }

    /// Remove a stream from management.
    pub fn remove_stream(&self, name: &str) -> Result<bool> {
        let mut guard = self
            .streams
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;
        Ok(guard.remove(name).is_some())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_stream_manager_creation() {
        let manager = StreamManager::new();
        let streams = manager.list_streams().expect("list_streams");
        assert!(streams.is_empty());
    }

    #[test]
    fn test_create_stream() {
        let manager = StreamManager::new();

        manager
            .create_stream("api-logs", DetectionProfile::Balanced, None)
            .expect("create_stream");

        assert!(manager.has_stream("api-logs").expect("has_stream"));
        assert!(!manager.has_stream("other").expect("has_stream"));

        let streams = manager.list_streams().expect("list_streams");
        assert_eq!(streams.len(), 1);
        assert!(streams.contains(&"api-logs".to_string()));
    }

    #[test]
    fn test_create_multiple_streams() {
        let manager = StreamManager::new();

        manager
            .create_stream("api-logs", DetectionProfile::Sensitive, None)
            .expect("create stream 1");
        manager
            .create_stream("db-metrics", DetectionProfile::Strict, None)
            .expect("create stream 2");
        manager
            .create_stream("user-events", DetectionProfile::Balanced, None)
            .expect("create stream 3");

        let streams = manager.list_streams().expect("list_streams");
        assert_eq!(streams.len(), 3);
    }

    #[test]
    fn test_ingest_to_stream() {
        let manager = StreamManager::new();

        manager
            .create_stream("test-stream", DetectionProfile::Balanced, None)
            .expect("create_stream");

        // Ingest some events
        for i in 0..10 {
            let result = manager.ingest("test-stream", format!("event {}", i).into_bytes());
            assert!(result.is_ok());
        }
    }

    #[test]
    fn test_ingest_to_unknown_stream_fails() {
        let manager = StreamManager::new();

        let result = manager.ingest("unknown-stream", b"event".to_vec());
        assert!(result.is_err());

        let err = result.unwrap_err();
        match err {
            CbadError::InvalidConfig(msg) => {
                assert!(msg.contains("unknown stream"));
            }
            _ => panic!("expected InvalidConfig error"),
        }
    }

    #[test]
    fn test_remove_stream() {
        let manager = StreamManager::new();

        manager
            .create_stream("temp-stream", DetectionProfile::Balanced, None)
            .expect("create_stream");

        assert!(manager.has_stream("temp-stream").expect("has_stream"));

        let removed = manager.remove_stream("temp-stream").expect("remove_stream");
        assert!(removed);

        assert!(!manager.has_stream("temp-stream").expect("has_stream"));

        // Removing again should return false
        let removed_again = manager.remove_stream("temp-stream").expect("remove_stream");
        assert!(!removed_again);
    }

    #[test]
    fn test_correlate_anomalies_empty() {
        let manager = StreamManager::new();

        manager
            .create_stream("stream1", DetectionProfile::Balanced, None)
            .expect("create_stream");

        let correlated = manager
            .correlate_anomalies(Duration::from_secs(300))
            .expect("correlate_anomalies");

        // Should return one CorrelatedAnomalies entry (even if empty)
        assert_eq!(correlated.len(), 1);
        assert!(correlated[0].anomalies.is_empty());
    }

    #[test]
    fn test_stream_manager_with_custom_store() {
        let store = Arc::new(InMemoryAnomalyStore::default());
        let manager = StreamManager::with_store(store.clone());

        manager
            .create_stream("test", DetectionProfile::Balanced, None)
            .expect("create_stream");

        // Both manager and test can share the store
        assert!(manager.has_stream("test").expect("has_stream"));
    }

    #[test]
    fn test_create_stream_with_custom_config() {
        let manager = StreamManager::new();

        let custom_config = AnomalyConfig {
            ncd_threshold: 0.5,
            p_value_threshold: 0.1,
            permutation_count: 50,
            ..Default::default()
        };

        manager
            .create_stream("custom", DetectionProfile::Balanced, Some(custom_config))
            .expect("create_stream");

        assert!(manager.has_stream("custom").expect("has_stream"));
    }

    #[test]
    fn test_stream_manager_default() {
        let manager = StreamManager::default();
        let streams = manager.list_streams().expect("list_streams");
        assert!(streams.is_empty());
    }
}

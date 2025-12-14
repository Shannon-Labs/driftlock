//! # Anomaly Storage & Persistence
//!
//! This module provides storage backends for persisting anomalies and detector state.
//!
//! ## Available Backends
//!
//! - [`InMemoryAnomalyStore`] - Fast, ephemeral storage for development and testing
//!
//! ## Planned Backends (Future Work)
//!
//! The following storage backends are planned for future releases:
//!
//! ### SQLite Storage
//! ```ignore
//! // Future API (not yet implemented)
//! use cbad_core::storage::SqliteAnomalyStore;
//!
//! let store = SqliteAnomalyStore::open("anomalies.db").await?;
//! let detector = CbadDetector::builder()
//!     .with_store(Arc::new(store))
//!     .build()?;
//! ```
//!
//! SQLite backend will provide:
//! - Single-file persistence with ACID guarantees
//! - Automatic schema migrations
//! - Full-text search on anomaly explanations
//! - Efficient time-range queries with indexes
//!
//! ### PostgreSQL Storage
//! ```ignore
//! // Future API (not yet implemented)
//! use cbad_core::storage::PostgresAnomalyStore;
//!
//! let store = PostgresAnomalyStore::connect("postgres://localhost/cbad").await?;
//! let detector = CbadDetector::builder()
//!     .with_store(Arc::new(store))
//!     .build()?;
//! ```
//!
//! PostgreSQL backend will provide:
//! - Production-grade persistence for distributed deployments
//! - Partitioned tables for time-series data
//! - Stream-based filtering with efficient indexes
//! - Integration with existing PostgreSQL infrastructure
//!
//! ## Implementing Custom Backends
//!
//! Implement the [`AnomalyStore`] trait for custom storage backends:
//!
//! ```ignore
//! use cbad_core::storage::{AnomalyStore, AnomalyFilter, AnomalyId, StoredAnomaly, Feedback};
//! use cbad_core::anomaly::AnomalyResult;
//! use async_trait::async_trait;
//!
//! struct MyCustomStore { /* ... */ }
//!
//! #[async_trait]
//! impl AnomalyStore for MyCustomStore {
//!     async fn store(&self, anomaly: &AnomalyResult, stream: Option<&str>) -> Result<AnomalyId> {
//!         // Store anomaly and return unique ID
//!     }
//!
//!     async fn query(&self, filter: AnomalyFilter) -> Result<Vec<StoredAnomaly>> {
//!         // Query anomalies matching filter
//!     }
//!
//!     async fn get_feedback(&self, id: AnomalyId) -> Result<Option<Feedback>> {
//!         // Retrieve feedback for anomaly
//!     }
//!
//!     async fn set_feedback(&self, id: AnomalyId, feedback: Feedback) -> Result<()> {
//!         // Store feedback (false positive / confirmed)
//!     }
//! }
//! ```

use crate::anomaly::{AnomalyConfig, AnomalyDetector, AnomalyResult};
use crate::error::{CbadError, Result};
use crate::window::WindowState;
use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs::File;
use std::io::BufReader;
use std::path::Path;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Arc, Mutex};
use std::time::{Duration, SystemTime};

/// Stable identifier for stored anomalies.
pub type AnomalyId = u64;

/// Persisted anomaly plus metadata for querying.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StoredAnomaly {
    pub id: AnomalyId,
    pub result: AnomalyResult,
    pub observed_at: SystemTime,
    pub stream: Option<String>,
    pub feedback: Option<Feedback>,
}

/// Simple feedback signal to power auto-tuning.
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq, Eq)]
pub enum Feedback {
    Confirmed,
    FalsePositive,
}

/// Filter for querying persisted anomalies.
#[derive(Debug, Clone, Default)]
pub struct AnomalyFilter {
    pub stream: Option<String>,
    pub since: Option<SystemTime>,
    pub until: Option<SystemTime>,
    pub only_anomalies: bool,
}

/// Trait describing persistence backends (SQLite/Postgres/in-memory).
#[async_trait]
pub trait AnomalyStore: Send + Sync {
    async fn store(&self, anomaly: &AnomalyResult, stream: Option<&str>) -> Result<AnomalyId>;
    async fn query(&self, filter: AnomalyFilter) -> Result<Vec<StoredAnomaly>>;
    async fn get_feedback(&self, id: AnomalyId) -> Result<Option<Feedback>>;

    /// Optional hook to record feedback (default is no-op).
    async fn set_feedback(&self, _id: AnomalyId, _feedback: Feedback) -> Result<()> {
        Ok(())
    }
}

/// Minimal in-memory store for development/testing.
#[derive(Default, Clone)]
pub struct InMemoryAnomalyStore {
    next_id: Arc<AtomicU64>,
    inner: Arc<Mutex<HashMap<AnomalyId, StoredAnomaly>>>,
}

#[async_trait]
impl AnomalyStore for InMemoryAnomalyStore {
    async fn store(&self, anomaly: &AnomalyResult, stream: Option<&str>) -> Result<AnomalyId> {
        let id = self.next_id.fetch_add(1, Ordering::Relaxed) + 1;
        let stored = StoredAnomaly {
            id,
            result: anomaly.clone(),
            observed_at: SystemTime::now(),
            stream: stream.map(|s| s.to_string()),
            feedback: None,
        };

        let mut guard = self
            .inner
            .lock()
            .map_err(|e| CbadError::StorageError(e.to_string()))?;
        guard.insert(id, stored);
        Ok(id)
    }

    async fn query(&self, filter: AnomalyFilter) -> Result<Vec<StoredAnomaly>> {
        let guard = self
            .inner
            .lock()
            .map_err(|e| CbadError::StorageError(e.to_string()))?;
        let mut out = Vec::new();
        for (_id, anomaly) in guard.iter() {
            if let Some(stream) = &filter.stream {
                if anomaly.stream.as_deref() != Some(stream) {
                    continue;
                }
            }
            if let Some(since) = filter.since {
                if anomaly.observed_at < since {
                    continue;
                }
            }
            if let Some(until) = filter.until {
                if anomaly.observed_at > until {
                    continue;
                }
            }
            if filter.only_anomalies && !anomaly.result.is_anomaly {
                continue;
            }
            out.push(anomaly.clone());
        }
        Ok(out)
    }

    async fn get_feedback(&self, id: AnomalyId) -> Result<Option<Feedback>> {
        let guard = self
            .inner
            .lock()
            .map_err(|e| CbadError::StorageError(e.to_string()))?;
        Ok(guard.get(&id).and_then(|a| a.feedback))
    }

    async fn set_feedback(&self, id: AnomalyId, feedback: Feedback) -> Result<()> {
        let mut guard = self
            .inner
            .lock()
            .map_err(|e| CbadError::StorageError(e.to_string()))?;
        if let Some(existing) = guard.get_mut(&id) {
            existing.feedback = Some(feedback);
        }
        Ok(())
    }
}

/// Serializable detector state for persistence and hot restart.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct DetectorState {
    pub config: AnomalyConfig,
    pub window: WindowState,
}

/// Window of anomalies that appear correlated across streams.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CorrelatedAnomalies {
    pub window: Duration,
    pub anomalies: Vec<StoredAnomaly>,
}

impl DetectorState {
    pub fn to_reader(&self, writer: &mut impl std::io::Write) -> Result<()> {
        serde_json::to_writer(writer, self).map_err(|e| CbadError::Serialization(e.to_string()))
    }

    pub fn from_reader(reader: impl std::io::Read) -> Result<Self> {
        serde_json::from_reader(reader).map_err(|e| CbadError::Serialization(e.to_string()))
    }
}

impl AnomalyDetector {
    /// Persist the detector state (config + window buffers) to memory.
    pub fn save_state(&self) -> Result<DetectorState> {
        let window = self
            .window()
            .snapshot_state()
            .map_err(|e| CbadError::StorageError(e.to_string()))?;

        Ok(DetectorState {
            config: self.config().clone(),
            window,
        })
    }

    /// Restore a detector from a persisted state blob.
    pub fn from_state(state: DetectorState) -> Result<Self> {
        let detector =
            AnomalyDetector::new(state.config.clone()).map_err(CbadError::Compression)?;
        detector
            .window()
            .restore_state(state.window)
            .map_err(|e| CbadError::StorageError(e.to_string()))?;
        Ok(detector)
    }

    /// Persist state to a file on disk.
    pub fn save_state_to_path<P: AsRef<Path>>(&self, path: P) -> Result<()> {
        let state = self.save_state()?;
        let mut file = File::create(path)?;
        state.to_reader(&mut file)
    }

    /// Load state from a file path.
    pub fn load_state_from_path<P: AsRef<Path>>(path: P) -> Result<Self> {
        let file = File::open(path)?;
        let reader = BufReader::new(file);
        let state = DetectorState::from_reader(reader)?;
        Self::from_state(state)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::metrics::AnomalyMetrics;

    #[test]
    fn test_in_memory_store_store_and_query() {
        let store = InMemoryAnomalyStore::default();
        let mut metrics = AnomalyMetrics::new();
        metrics.is_anomaly = true;
        metrics.ncd = 0.5;
        metrics.p_value = 0.01;

        let result = crate::anomaly::AnomalyResult::new(metrics, 0.05);

        // Store an anomaly
        let id = futures::executor::block_on(store.store(&result, Some("test-stream")))
            .expect("store should succeed");
        assert!(id > 0);

        // Query it back
        let queried = futures::executor::block_on(store.query(AnomalyFilter {
            stream: Some("test-stream".to_string()),
            only_anomalies: true,
            ..Default::default()
        }))
        .expect("query should succeed");

        assert_eq!(queried.len(), 1);
        assert_eq!(queried[0].id, id);
        assert_eq!(queried[0].stream.as_deref(), Some("test-stream"));
    }

    #[test]
    fn test_in_memory_store_feedback() {
        let store = InMemoryAnomalyStore::default();
        let mut metrics = AnomalyMetrics::new();
        metrics.is_anomaly = true;

        let result = crate::anomaly::AnomalyResult::new(metrics, 0.05);

        let id = futures::executor::block_on(store.store(&result, None)).expect("store");

        // Initially no feedback
        let feedback = futures::executor::block_on(store.get_feedback(id)).expect("get_feedback");
        assert!(feedback.is_none());

        // Set feedback
        futures::executor::block_on(store.set_feedback(id, Feedback::FalsePositive))
            .expect("set_feedback");

        // Verify feedback
        let feedback = futures::executor::block_on(store.get_feedback(id)).expect("get_feedback");
        assert_eq!(feedback, Some(Feedback::FalsePositive));
    }

    #[test]
    fn test_in_memory_store_filter_by_time() {
        let store = InMemoryAnomalyStore::default();
        let mut metrics = AnomalyMetrics::new();
        metrics.is_anomaly = true;

        let result = crate::anomaly::AnomalyResult::new(metrics, 0.05);

        // Store anomaly
        futures::executor::block_on(store.store(&result, None)).expect("store");

        // Query with time filter (should find it since it was just stored)
        let since = SystemTime::now() - Duration::from_secs(60);
        let queried = futures::executor::block_on(store.query(AnomalyFilter {
            since: Some(since),
            only_anomalies: true,
            ..Default::default()
        }))
        .expect("query");

        assert_eq!(queried.len(), 1);

        // Query with future time filter (should not find it)
        let future_time = SystemTime::now() + Duration::from_secs(60);
        let queried = futures::executor::block_on(store.query(AnomalyFilter {
            since: Some(future_time),
            ..Default::default()
        }))
        .expect("query");

        assert_eq!(queried.len(), 0);
    }

    #[test]
    fn test_in_memory_store_only_anomalies_filter() {
        let store = InMemoryAnomalyStore::default();

        // Store an anomaly
        let mut metrics = AnomalyMetrics::new();
        metrics.is_anomaly = true;
        let anomaly_result = crate::anomaly::AnomalyResult::new(metrics, 0.05);
        futures::executor::block_on(store.store(&anomaly_result, None)).expect("store");

        // Store a non-anomaly
        let mut metrics2 = AnomalyMetrics::new();
        metrics2.is_anomaly = false;
        let non_anomaly_result = crate::anomaly::AnomalyResult::new(metrics2, 0.05);
        futures::executor::block_on(store.store(&non_anomaly_result, None)).expect("store");

        // Query only anomalies
        let anomalies_only = futures::executor::block_on(store.query(AnomalyFilter {
            only_anomalies: true,
            ..Default::default()
        }))
        .expect("query");
        assert_eq!(anomalies_only.len(), 1);

        // Query all
        let all =
            futures::executor::block_on(store.query(AnomalyFilter::default())).expect("query");
        assert_eq!(all.len(), 2);
    }

    #[test]
    fn test_detector_state_serialization() {
        let config = AnomalyConfig::default();
        let state = DetectorState {
            config,
            window: WindowState {
                events: vec![],
                baseline_start: 0,
                window_start: 0,
                total_events: 100,
                aligned: true,
            },
        };

        // Serialize to bytes
        let mut buf = Vec::new();
        state.to_reader(&mut buf).expect("serialize");

        // Deserialize
        let restored = DetectorState::from_reader(buf.as_slice()).expect("deserialize");
        assert_eq!(restored.window.total_events, 100);
        assert!(restored.window.aligned);
    }

    #[test]
    fn test_stored_anomaly_structure() {
        let stored = StoredAnomaly {
            id: 42,
            result: crate::anomaly::AnomalyResult::new(AnomalyMetrics::new(), 0.05),
            observed_at: SystemTime::now(),
            stream: Some("test".to_string()),
            feedback: Some(Feedback::Confirmed),
        };

        assert_eq!(stored.id, 42);
        assert_eq!(stored.stream.as_deref(), Some("test"));
        assert_eq!(stored.feedback, Some(Feedback::Confirmed));
    }

    #[test]
    fn test_feedback_equality() {
        assert_eq!(Feedback::Confirmed, Feedback::Confirmed);
        assert_eq!(Feedback::FalsePositive, Feedback::FalsePositive);
        assert_ne!(Feedback::Confirmed, Feedback::FalsePositive);
    }

    #[test]
    fn test_anomaly_filter_default() {
        let filter = AnomalyFilter::default();
        assert!(filter.stream.is_none());
        assert!(filter.since.is_none());
        assert!(filter.until.is_none());
        assert!(!filter.only_anomalies);
    }

    #[test]
    fn test_correlated_anomalies_structure() {
        let correlated = CorrelatedAnomalies {
            window: Duration::from_secs(300),
            anomalies: vec![],
        };

        assert_eq!(correlated.window.as_secs(), 300);
        assert!(correlated.anomalies.is_empty());
    }
}

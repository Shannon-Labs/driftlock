use crate::anomaly::{AnomalyConfig, AnomalyDetector, AnomalyResult};
use crate::compression::CompressionAlgorithm;
use crate::error::{CbadError, Result};
use crate::storage::{AnomalyId, AnomalyStore, Feedback, InMemoryAnomalyStore};
use crate::tokenizer::{TokenizerConfig, TokenizerStats};
use crate::window::WindowConfig;
use csv::ReaderBuilder;
use futures::executor::block_on;
use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader};
use std::path::Path;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::{Arc, Mutex};
use std::time::{Duration, Instant};

#[cfg(feature = "runtime-tokio")]
use async_stream::try_stream;
#[cfg(feature = "runtime-tokio")]
use futures::Stream;
#[cfg(feature = "runtime-tokio")]
use tokio::io::{AsyncBufRead, AsyncBufReadExt};

/// Preset detection profiles for fast onboarding.
#[derive(Debug, Clone)]
pub enum DetectionProfile {
    /// Low false positives, may miss subtle anomalies
    Strict,
    /// Balanced detection (default)
    Balanced,
    /// High sensitivity, more false positives acceptable
    Sensitive,
    /// Custom thresholds
    Custom(CustomProfile),
}

/// Customizable profile for power users.
#[derive(Debug, Clone)]
pub struct CustomProfile {
    pub ncd_threshold: f64,
    pub p_value_threshold: f64,
    pub permutation_count: usize,
    pub tokenizer_config: Option<TokenizerConfig>,
    pub window_config: Option<WindowConfig>,
    pub compression_algorithm: Option<CompressionAlgorithm>,
}

impl Default for CustomProfile {
    fn default() -> Self {
        Self {
            ncd_threshold: 0.3,
            p_value_threshold: 0.05,
            permutation_count: 1000,
            tokenizer_config: None,
            window_config: None,
            compression_algorithm: None,
        }
    }
}

impl DetectionProfile {
    fn apply(&self, base: AnomalyConfig) -> AnomalyConfig {
        match self {
            DetectionProfile::Balanced => base,
            DetectionProfile::Strict => AnomalyConfig {
                ncd_threshold: 0.35,
                p_value_threshold: 0.01,
                permutation_count: base.permutation_count.max(1500),
                require_statistical_significance: true,
                ..base
            },
            DetectionProfile::Sensitive => AnomalyConfig {
                ncd_threshold: 0.22,
                p_value_threshold: 0.1,
                permutation_count: base.permutation_count.min(200),
                require_statistical_significance: false,
                ..base
            },
            DetectionProfile::Custom(custom) => AnomalyConfig {
                ncd_threshold: custom.ncd_threshold,
                p_value_threshold: custom.p_value_threshold,
                permutation_count: custom.permutation_count,
                tokenizer_config: custom.tokenizer_config.or(base.tokenizer_config),
                window_config: custom.window_config.clone().unwrap_or(base.window_config),
                compression_algorithm: custom
                    .compression_algorithm
                    .unwrap_or(base.compression_algorithm),
                ..base
            },
        }
    }
}

/// Baseline configuration for batch ingestion.
#[derive(Debug, Clone)]
pub enum BaselineStrategy {
    /// Use the first N events as baseline before scoring.
    UseFirstN(usize),
    /// Provide an explicit baseline payload.
    ExplicitBaseline(Vec<Vec<u8>>),
    /// Assume baseline already populated (streaming only).
    WindowOnly,
}

/// CSV ingestion knobs.
#[derive(Debug, Clone)]
pub struct CsvConfig {
    pub has_headers: bool,
    pub delimiter: u8,
    pub column: Option<usize>,
}

impl Default for CsvConfig {
    fn default() -> Self {
        Self {
            has_headers: true,
            delimiter: b',',
            column: None,
        }
    }
}

/// Parquet ingestion configuration.
///
/// **Requires the `parquet` feature flag.**
#[derive(Debug, Clone, Default)]
pub struct ParquetConfig {
    /// Specific column to extract for anomaly detection.
    /// If None, the entire row is serialized as JSON.
    pub column: Option<String>,
    /// Batch size for reading rows (default: 1024).
    pub batch_size: usize,
}

impl ParquetConfig {
    /// Create a new ParquetConfig targeting a specific column.
    pub fn with_column(column: impl Into<String>) -> Self {
        Self {
            column: Some(column.into()),
            batch_size: 1024,
        }
    }
}

/// Backpressure controls for async ingestion.
#[derive(Debug, Clone)]
pub struct BackpressureConfig {
    pub max_in_flight: usize,
}

impl Default for BackpressureConfig {
    fn default() -> Self {
        Self {
            max_in_flight: 1024,
        }
    }
}

/// Progress callback data.
#[derive(Debug, Clone)]
pub struct ProgressUpdate {
    pub processed: u64,
    pub anomalies: u64,
    pub dataset_len: Option<u64>,
}

/// Optional callback for long-running batch jobs.
pub type ProgressCallback = Arc<dyn Fn(ProgressUpdate) + Send + Sync>;

/// Anomaly result tagged with a durable identifier for feedback.
#[derive(Debug, Clone)]
pub struct DetectionRecord {
    pub id: AnomalyId,
    pub result: AnomalyResult,
}

/// Production-ready detector wrapper with ingestion helpers.
pub struct CbadDetector {
    detector: Mutex<AnomalyDetector>,
    store: Option<Arc<dyn AnomalyStore>>,
    auto_tuner: AutoTuner,
    observability: Observability,
    progress: Option<ProgressCallback>,
    id_gen: AtomicU64,
    detections: Mutex<HashMap<AnomalyId, AnomalyResult>>,
    #[cfg(feature = "runtime-tokio")]
    backpressure: Option<tokio::sync::Semaphore>,
    name: String,
}

impl CbadDetector {
    pub fn builder() -> CbadDetectorBuilder {
        CbadDetectorBuilder::new()
    }

    /// Ingest a single event and optionally get an anomaly detection.
    pub fn ingest_event(&self, data: Vec<u8>) -> Result<Option<DetectionRecord>> {
        let detector = self
            .detector
            .lock()
            .map_err(|e| CbadError::InvalidConfig(format!("detector poisoned: {}", e)))?;

        detector
            .add_data(data)
            .map_err(|e| CbadError::StorageError(e.to_string()))?;

        self.observability.on_event();
        self.emit_progress(None);

        if !detector
            .is_ready()
            .map_err(|e| CbadError::StorageError(e.to_string()))?
        {
            return Ok(None);
        }

        let start = Instant::now();
        let detection = detector
            .detect_anomaly()
            .map_err(|e| CbadError::StorageError(e.to_string()))?;
        let Some(result) = detection else {
            return Ok(None);
        };

        let id = self.id_gen.fetch_add(1, Ordering::Relaxed) + 1;
        self.observability
            .on_detection(result.is_anomaly, start.elapsed());
        self.persist_detection(id, &result);

        if result.is_anomaly {
            self.maybe_store(&result)?;
        }

        Ok(Some(DetectionRecord { id, result }))
    }

    /// Batch process an in-memory dataset.
    pub fn analyze_batch<I>(&self, dataset: I) -> Result<Vec<DetectionRecord>>
    where
        I: IntoIterator,
        I::Item: AsRef<[u8]>,
    {
        let mut out = Vec::new();
        for entry in dataset {
            if let Some(record) = self.ingest_event(entry.as_ref().to_vec())? {
                out.push(record);
            }
        }
        Ok(out)
    }

    /// Analyze newline-delimited JSON/NDJSON datasets from disk with streaming.
    pub fn analyze_file<P: AsRef<Path>>(&self, path: P) -> Result<Vec<DetectionRecord>> {
        let path = path.as_ref();
        let ext = path
            .extension()
            .and_then(|e| e.to_str())
            .unwrap_or("")
            .to_lowercase();

        if ext == "csv" {
            return self.analyze_csv(path, CsvConfig::default());
        }
        if ext == "jsonl" || ext == "ndjson" {
            return self.analyze_json_lines(path);
        }
        if ext == "json" {
            // Treat as array of events.
            let file = File::open(path)?;
            let reader = BufReader::new(file);
            let value: serde_json::Value = serde_json::from_reader(reader)
                .map_err(|e| CbadError::Serialization(e.to_string()))?;
            let Some(arr) = value.as_array() else {
                return Err(CbadError::UnsupportedFormat(
                    "expected JSON array for dataset ingestion".into(),
                ));
            };

            let mut out = Vec::new();
            for item in arr {
                let bytes = serde_json::to_vec(item)
                    .map_err(|e| CbadError::Serialization(e.to_string()))?;
                if let Some(record) = self.ingest_event(bytes)? {
                    out.push(record);
                }
            }
            return Ok(out);
        }

        Err(CbadError::UnsupportedFormat(format!(
            "extension '{}' not supported (use csv/jsonl/ndjson/json)",
            ext
        )))
    }

    /// Analyze CSV data using a configurable column.
    pub fn analyze_csv<P: AsRef<Path>>(
        &self,
        path: P,
        config: CsvConfig,
    ) -> Result<Vec<DetectionRecord>> {
        let mut reader = ReaderBuilder::new()
            .has_headers(config.has_headers)
            .delimiter(config.delimiter)
            .from_path(path)?;

        let mut out = Vec::new();
        for record in reader.records() {
            let record = record.map_err(|e| CbadError::Serialization(e.to_string()))?;
            let bytes = if let Some(idx) = config.column {
                record
                    .get(idx)
                    .ok_or_else(|| {
                        CbadError::InvalidConfig(format!("column {} out of bounds", idx))
                    })?
                    .as_bytes()
                    .to_vec()
            } else {
                record.as_slice().as_bytes().to_vec()
            };

            if let Some(detected) = self.ingest_event(bytes)? {
                out.push(detected);
            }
        }
        Ok(out)
    }

    fn analyze_json_lines<P: AsRef<Path>>(&self, path: P) -> Result<Vec<DetectionRecord>> {
        let file = File::open(path)?;
        let reader = BufReader::new(file);
        let mut out = Vec::new();
        for line in reader.lines() {
            let line = line?;
            if line.trim().is_empty() {
                continue;
            }
            if let Some(record) = self.ingest_event(line.into_bytes())? {
                out.push(record);
            }
        }
        Ok(out)
    }

    /// Analyze Parquet files for anomaly detection.
    ///
    /// **Requires the `parquet` feature flag.**
    ///
    /// This method reads Parquet files and processes each row for anomaly detection.
    /// The specified column (or all columns serialized as JSON) is used as the event payload.
    ///
    /// # Example (when parquet feature is enabled)
    /// ```ignore
    /// let detector = CbadDetector::builder().build()?;
    /// let config = ParquetConfig {
    ///     column: Some("message".to_string()),
    ///     batch_size: 1024,
    /// };
    /// let anomalies = detector.analyze_parquet("logs.parquet", config)?;
    /// ```
    #[cfg(feature = "parquet")]
    pub fn analyze_parquet<P: AsRef<Path>>(
        &self,
        path: P,
        config: ParquetConfig,
    ) -> Result<Vec<DetectionRecord>> {
        use parquet::file::reader::{FileReader, SerializedFileReader};
        use parquet::record::Row;

        let file = File::open(path)?;
        let reader = SerializedFileReader::new(file)
            .map_err(|e| CbadError::Serialization(format!("parquet error: {}", e)))?;

        let mut out = Vec::new();
        let mut iter = reader
            .get_row_iter(None)
            .map_err(|e| CbadError::Serialization(format!("parquet iter error: {}", e)))?;

        while let Some(row_result) = iter.next() {
            let row = row_result
                .map_err(|e| CbadError::Serialization(format!("parquet row error: {}", e)))?;

            let bytes = if let Some(ref col) = config.column {
                // Extract specific column
                Self::extract_parquet_column(&row, col)?
            } else {
                // Serialize entire row as JSON
                format!("{:?}", row).into_bytes()
            };

            if let Some(record) = self.ingest_event(bytes)? {
                out.push(record);
            }
        }

        Ok(out)
    }

    #[cfg(feature = "parquet")]
    fn extract_parquet_column(row: &parquet::record::Row, column: &str) -> Result<Vec<u8>> {
        // Simple extraction - find column by name
        for (name, field) in row.get_column_iter() {
            if name == column {
                return Ok(format!("{:?}", field).into_bytes());
            }
        }
        Err(CbadError::InvalidConfig(format!(
            "column '{}' not found in parquet file",
            column
        )))
    }

    /// Stub for Parquet analysis when feature is not enabled.
    #[cfg(not(feature = "parquet"))]
    pub fn analyze_parquet<P: AsRef<Path>>(
        &self,
        _path: P,
        _config: ParquetConfig,
    ) -> Result<Vec<DetectionRecord>> {
        Err(CbadError::UnsupportedFormat(
            "parquet support requires the 'parquet' feature flag: cargo build --features parquet"
                .into(),
        ))
    }

    /// Asynchronous streaming ingestion (Tokio runtime).
    #[cfg(feature = "runtime-tokio")]
    pub fn stream_analyze<'a, R>(
        &'a self,
        mut reader: R,
    ) -> impl Stream<Item = Result<DetectionRecord>> + 'a
    where
        R: AsyncBufRead + Unpin + Send + 'a,
    {
        try_stream! {
            loop {
                let mut buf = String::new();
                let read = reader.read_line(&mut buf).await?;
                if read == 0 {
                    break;
                }

                if let Some(record) = self.ingest_event(buf.into_bytes())? {
                    yield record;
                }
            }
        }
    }

    /// Persist detector state to disk.
    pub fn save_state<P: AsRef<Path>>(&self, path: P) -> Result<()> {
        let detector = self
            .detector
            .lock()
            .map_err(|e| CbadError::InvalidConfig(format!("detector poisoned: {}", e)))?;
        detector.save_state_to_path(path)
    }

    /// Load a detector from persisted state.
    pub fn load_state(path: impl AsRef<Path>) -> Result<Self> {
        let detector = AnomalyDetector::load_state_from_path(path)?;
        Ok(CbadDetector {
            detector: Mutex::new(detector),
            store: None,
            auto_tuner: AutoTuner::default(),
            observability: Observability::default(),
            progress: None,
            id_gen: AtomicU64::new(0),
            detections: Mutex::new(HashMap::new()),
            #[cfg(feature = "runtime-tokio")]
            backpressure: None,
            name: "restored".into(),
        })
    }

    /// Report current metrics for observability/Prometheus export.
    pub fn metrics(&self) -> DetectorMetricsSnapshot {
        let detector = self.detector.lock();
        let tokenizer_stats = detector.as_ref().ok().and_then(|d| d.tokenizer_stats());
        let config = detector
            .as_ref()
            .ok()
            .map(|d| d.config().clone())
            .unwrap_or_default();

        self.observability.snapshot(tokenizer_stats, &config)
    }

    /// Mark a detection as false positive and auto-tune thresholds.
    pub fn mark_false_positive(&self, id: AnomalyId) -> Result<()> {
        let record = {
            let map = self
                .detections
                .lock()
                .map_err(|e| CbadError::StorageError(e.to_string()))?;
            map.get(&id).cloned()
        };

        if record.is_none() {
            return Err(CbadError::InvalidConfig(format!(
                "anomaly id {} not found",
                id
            )));
        }

        self.auto_tuner.adjust(false, &self.detector)?;

        if let Some(store) = &self.store {
            block_on(store.set_feedback(id, Feedback::FalsePositive))?;
        }
        Ok(())
    }

    /// Confirm a detection and auto-tune toward higher sensitivity.
    pub fn confirm_anomaly(&self, id: AnomalyId) -> Result<()> {
        let record = {
            let map = self
                .detections
                .lock()
                .map_err(|e| CbadError::StorageError(e.to_string()))?;
            map.get(&id).cloned()
        };

        if record.is_none() {
            return Err(CbadError::InvalidConfig(format!(
                "anomaly id {} not found",
                id
            )));
        }

        self.auto_tuner.adjust(true, &self.detector)?;

        if let Some(store) = &self.store {
            block_on(store.set_feedback(id, Feedback::Confirmed))?;
        }
        Ok(())
    }

    #[cfg(feature = "runtime-tokio")]
    pub async fn add_data_with_backpressure(
        &self,
        data: Vec<u8>,
        config: BackpressureConfig,
    ) -> Result<Option<DetectionRecord>> {
        let semaphore = if let Some(sema) = &self.backpressure {
            sema.clone()
        } else {
            tokio::sync::Semaphore::new(config.max_in_flight)
        };

        let permit = semaphore
            .acquire()
            .await
            .map_err(|_| CbadError::ResourceExhausted {
                resource: "backpressure channel",
                limit: config.max_in_flight,
            })?;
        let result = self.ingest_event(data)?;
        drop(permit);
        Ok(result)
    }

    fn persist_detection(&self, id: AnomalyId, result: &AnomalyResult) {
        if let Ok(mut map) = self.detections.lock() {
            map.insert(id, result.clone());
        }
    }

    fn maybe_store(&self, result: &AnomalyResult) -> Result<Option<AnomalyId>> {
        if let Some(store) = &self.store {
            let id = block_on(store.store(result, Some(&self.name)))?;
            Ok(Some(id))
        } else {
            Ok(None)
        }
    }

    fn emit_progress(&self, dataset_len: Option<u64>) {
        if let Some(cb) = &self.progress {
            cb(ProgressUpdate {
                processed: self.observability.events_processed.load(Ordering::Relaxed),
                anomalies: self
                    .observability
                    .anomalies_detected
                    .load(Ordering::Relaxed),
                dataset_len,
            });
        }
    }
}

/// Builder for ergonomic detector setup.
pub struct CbadDetectorBuilder {
    config: AnomalyConfig,
    profile: DetectionProfile,
    baseline_strategy: Option<BaselineStrategy>,
    progress: Option<ProgressCallback>,
    store: Option<Arc<dyn AnomalyStore>>,
    name: String,
    #[cfg(feature = "runtime-tokio")]
    backpressure: Option<BackpressureConfig>,
}

impl CbadDetectorBuilder {
    pub fn new() -> Self {
        Self {
            config: AnomalyConfig::default(),
            profile: DetectionProfile::Balanced,
            baseline_strategy: None,
            progress: None,
            store: None,
            name: "cbad-detector".into(),
            #[cfg(feature = "runtime-tokio")]
            backpressure: None,
        }
    }

    pub fn with_profile(mut self, profile: DetectionProfile) -> Self {
        self.profile = profile;
        self
    }

    pub fn with_config(mut self, config: AnomalyConfig) -> Self {
        self.config = config;
        self
    }

    pub fn with_baseline_strategy(mut self, strategy: BaselineStrategy) -> Self {
        self.baseline_strategy = Some(strategy);
        self
    }

    pub fn on_progress<F>(mut self, callback: F) -> Self
    where
        F: Fn(ProgressUpdate) + Send + Sync + 'static,
    {
        self.progress = Some(Arc::new(callback));
        self
    }

    pub fn with_store(mut self, store: Arc<dyn AnomalyStore>) -> Self {
        self.store = Some(store);
        self
    }

    pub fn named(mut self, name: impl Into<String>) -> Self {
        self.name = name.into();
        self
    }

    #[cfg(feature = "runtime-tokio")]
    pub fn with_backpressure(mut self, config: BackpressureConfig) -> Self {
        self.backpressure = Some(config);
        self
    }

    pub fn build(self) -> Result<CbadDetector> {
        let mut config = self.profile.apply(self.config);
        if let Some(BaselineStrategy::UseFirstN(n)) = &self.baseline_strategy {
            config.window_config.baseline_size = *n;
        }

        let detector = AnomalyDetector::new(config.clone()).map_err(CbadError::Compression)?;

        // Preload baseline if provided explicitly.
        if let Some(BaselineStrategy::ExplicitBaseline(baseline)) = &self.baseline_strategy {
            for entry in baseline {
                detector
                    .add_data(entry.clone())
                    .map_err(|e| CbadError::StorageError(e.to_string()))?;
            }
        }

        Ok(CbadDetector {
            detector: Mutex::new(detector),
            store: self.store,
            auto_tuner: AutoTuner::default(),
            observability: Observability::default(),
            progress: self.progress,
            id_gen: AtomicU64::new(0),
            detections: Mutex::new(HashMap::new()),
            #[cfg(feature = "runtime-tokio")]
            backpressure: self
                .backpressure
                .map(|cfg| tokio::sync::Semaphore::new(cfg.max_in_flight)),
            name: self.name,
        })
    }
}

#[derive(Debug, Default)]
struct Observability {
    events_processed: AtomicU64,
    anomalies_detected: AtomicU64,
    detection_cycles: AtomicU64,
    latency_nanos: AtomicU64,
}

impl Observability {
    fn on_event(&self) {
        self.events_processed.fetch_add(1, Ordering::Relaxed);
    }

    fn on_detection(&self, is_anomaly: bool, latency: Duration) {
        self.detection_cycles.fetch_add(1, Ordering::Relaxed);
        if is_anomaly {
            self.anomalies_detected.fetch_add(1, Ordering::Relaxed);
        }
        self.latency_nanos
            .fetch_add(latency.as_nanos() as u64, Ordering::Relaxed);
    }

    fn snapshot(
        &self,
        tokenizer_stats: Option<TokenizerStats>,
        config: &AnomalyConfig,
    ) -> DetectorMetricsSnapshot {
        let detections = self.detection_cycles.load(Ordering::Relaxed);
        let latency_total = self.latency_nanos.load(Ordering::Relaxed);
        let avg_latency = if detections > 0 {
            latency_total as f64 / detections as f64 / 1_000_000_000.0
        } else {
            0.0
        };

        DetectorMetricsSnapshot {
            events_processed: self.events_processed.load(Ordering::Relaxed),
            anomalies_detected: self.anomalies_detected.load(Ordering::Relaxed),
            detection_cycles: detections,
            avg_detection_latency_secs: avg_latency,
            baseline_size: config.window_config.baseline_size as u64,
            window_size: config.window_config.window_size as u64,
            tokenizer_stats,
        }
    }
}

impl Default for CbadDetectorBuilder {
    fn default() -> Self {
        Self::new()
    }
}

/// Snapshot of metrics for Prometheus/telemetry export.
#[derive(Debug, Clone)]
pub struct DetectorMetricsSnapshot {
    pub events_processed: u64,
    pub anomalies_detected: u64,
    pub detection_cycles: u64,
    pub avg_detection_latency_secs: f64,
    pub baseline_size: u64,
    pub window_size: u64,
    pub tokenizer_stats: Option<TokenizerStats>,
}

impl DetectorMetricsSnapshot {
    /// Export metrics in Prometheus text format.
    ///
    /// Returns a string suitable for serving at `/metrics` endpoints.
    /// Example output:
    /// ```text
    /// # HELP cbad_events_processed_total Total number of events processed
    /// # TYPE cbad_events_processed_total counter
    /// cbad_events_processed_total 12345
    /// ```
    pub fn to_prometheus_text(&self) -> String {
        self.to_prometheus_text_with_labels(&[])
    }

    /// Export metrics in Prometheus text format with custom labels.
    ///
    /// # Arguments
    /// * `labels` - Slice of (key, value) label pairs to add to all metrics
    pub fn to_prometheus_text_with_labels(&self, labels: &[(&str, &str)]) -> String {
        let mut out = String::with_capacity(2048);
        let label_str = Self::format_labels(labels);

        // Events processed counter
        out.push_str("# HELP cbad_events_processed_total Total number of events ingested\n");
        out.push_str("# TYPE cbad_events_processed_total counter\n");
        out.push_str(&format!(
            "cbad_events_processed_total{} {}\n",
            label_str, self.events_processed
        ));

        // Anomalies detected counter
        out.push_str("# HELP cbad_anomalies_detected_total Total number of anomalies detected\n");
        out.push_str("# TYPE cbad_anomalies_detected_total counter\n");
        out.push_str(&format!(
            "cbad_anomalies_detected_total{} {}\n",
            label_str, self.anomalies_detected
        ));

        // Detection cycles counter
        out.push_str("# HELP cbad_detection_cycles_total Total number of detection cycles run\n");
        out.push_str("# TYPE cbad_detection_cycles_total counter\n");
        out.push_str(&format!(
            "cbad_detection_cycles_total{} {}\n",
            label_str, self.detection_cycles
        ));

        // Average detection latency gauge
        out.push_str(
            "# HELP cbad_detection_latency_seconds Average detection latency in seconds\n",
        );
        out.push_str("# TYPE cbad_detection_latency_seconds gauge\n");
        out.push_str(&format!(
            "cbad_detection_latency_seconds{} {:.9}\n",
            label_str, self.avg_detection_latency_secs
        ));

        // Baseline size gauge
        out.push_str("# HELP cbad_baseline_size Current baseline window size\n");
        out.push_str("# TYPE cbad_baseline_size gauge\n");
        out.push_str(&format!(
            "cbad_baseline_size{} {}\n",
            label_str, self.baseline_size
        ));

        // Window size gauge
        out.push_str("# HELP cbad_window_size Current analysis window size\n");
        out.push_str("# TYPE cbad_window_size gauge\n");
        out.push_str(&format!(
            "cbad_window_size{} {}\n",
            label_str, self.window_size
        ));

        // Tokenizer stats (if available)
        if let Some(ref stats) = self.tokenizer_stats {
            out.push_str(
                "# HELP cbad_tokenizer_replacements_total Tokenizer replacements by type\n",
            );
            out.push_str("# TYPE cbad_tokenizer_replacements_total counter\n");

            let types = [
                ("uuid", stats.uuid_count),
                ("hash", stats.hash_count),
                ("jwt", stats.jwt_count),
                ("base64", stats.base64_count),
                ("ip", stats.ip_count),
                ("url", stats.url_count),
                ("domain", stats.domain_count),
                ("email", stats.email_count),
                ("timestamp", stats.timestamp_count),
                ("numeric", stats.numeric_count),
                ("cloud_id", stats.cloud_id_count),
                ("json", stats.json_canonicalized_count),
            ];

            for (token_type, count) in types {
                let type_labels = if labels.is_empty() {
                    format!("{{type=\"{}\"}}", token_type)
                } else {
                    let base = Self::format_labels(labels);
                    // Insert type label before closing brace
                    format!("{}type=\"{}\"}}", &base[..base.len() - 1], token_type)
                };
                out.push_str(&format!(
                    "cbad_tokenizer_replacements_total{} {}\n",
                    type_labels, count
                ));
            }

            out.push_str("# HELP cbad_tokenizer_bytes_saved_total Bytes saved by tokenization\n");
            out.push_str("# TYPE cbad_tokenizer_bytes_saved_total counter\n");
            out.push_str(&format!(
                "cbad_tokenizer_bytes_saved_total{} {}\n",
                label_str, stats.bytes_saved
            ));
        }

        out
    }

    fn format_labels(labels: &[(&str, &str)]) -> String {
        if labels.is_empty() {
            String::new()
        } else {
            let inner: Vec<String> = labels
                .iter()
                .map(|(k, v)| format!("{}=\"{}\"", k, Self::escape_label_value(v)))
                .collect();
            format!("{{{}}}", inner.join(","))
        }
    }

    fn escape_label_value(value: &str) -> String {
        value
            .replace('\\', "\\\\")
            .replace('"', "\\\"")
            .replace('\n', "\\n")
    }
}

#[derive(Default)]
struct AutoTuneState {
    false_positive: u64,
    confirmed: u64,
}

#[derive(Default)]
struct AutoTuner {
    state: Mutex<AutoTuneState>,
}

impl AutoTuner {
    fn adjust(&self, confirmed: bool, detector: &Mutex<AnomalyDetector>) -> Result<()> {
        let mut guard = self
            .state
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;

        if confirmed {
            guard.confirmed += 1;
        } else {
            guard.false_positive += 1;
        }

        let mut detector = detector
            .lock()
            .map_err(|e| CbadError::InvalidConfig(e.to_string()))?;

        let mut config = detector.config().clone();
        let sensitivity_delta = if confirmed { -0.01 } else { 0.02 };
        config.ncd_threshold = (config.ncd_threshold + sensitivity_delta).clamp(0.05, 1.0);

        // Adjust p-value the opposite direction for FP vs TP signals.
        if confirmed {
            config.p_value_threshold = (config.p_value_threshold + 0.01).min(0.2);
        } else {
            config.p_value_threshold = (config.p_value_threshold * 0.9).max(0.005);
        }

        detector.update_config(config);
        Ok(())
    }
}

impl Default for CbadDetector {
    fn default() -> Self {
        Self::builder()
            .with_store(Arc::new(InMemoryAnomalyStore::default()))
            .build()
            .expect("default detector should be constructible")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_prometheus_text_format() {
        let snapshot = DetectorMetricsSnapshot {
            events_processed: 12345,
            anomalies_detected: 42,
            detection_cycles: 100,
            avg_detection_latency_secs: 0.007,
            baseline_size: 1000,
            window_size: 100,
            tokenizer_stats: None,
        };

        let text = snapshot.to_prometheus_text();

        assert!(text.contains("cbad_events_processed_total 12345"));
        assert!(text.contains("cbad_anomalies_detected_total 42"));
        assert!(text.contains("cbad_detection_cycles_total 100"));
        assert!(text.contains("cbad_baseline_size 1000"));
        assert!(text.contains("cbad_window_size 100"));
        assert!(text.contains("# TYPE cbad_events_processed_total counter"));
        assert!(text.contains("# TYPE cbad_detection_latency_seconds gauge"));
    }

    #[test]
    fn test_prometheus_text_with_labels() {
        let snapshot = DetectorMetricsSnapshot {
            events_processed: 100,
            anomalies_detected: 5,
            detection_cycles: 10,
            avg_detection_latency_secs: 0.001,
            baseline_size: 500,
            window_size: 50,
            tokenizer_stats: None,
        };

        let text =
            snapshot.to_prometheus_text_with_labels(&[("stream", "api-logs"), ("env", "prod")]);

        assert!(text.contains(r#"cbad_events_processed_total{stream="api-logs",env="prod"} 100"#));
        assert!(text.contains(r#"cbad_anomalies_detected_total{stream="api-logs",env="prod"} 5"#));
    }

    #[test]
    fn test_prometheus_text_with_tokenizer_stats() {
        let snapshot = DetectorMetricsSnapshot {
            events_processed: 1000,
            anomalies_detected: 10,
            detection_cycles: 50,
            avg_detection_latency_secs: 0.005,
            baseline_size: 500,
            window_size: 100,
            tokenizer_stats: Some(TokenizerStats {
                jwt_count: 5,
                uuid_count: 100,
                hash_count: 50,
                base64_count: 20,
                ip_count: 10,
                url_count: 8,
                domain_count: 6,
                email_count: 4,
                timestamp_count: 3,
                numeric_count: 2,
                cloud_id_count: 1,
                json_canonicalized_count: 200,
                bytes_saved: 5000,
            }),
        };

        let text = snapshot.to_prometheus_text();

        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="uuid"} 100"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="hash"} 50"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="jwt"} 5"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="base64"} 20"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="ip"} 10"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="url"} 8"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="domain"} 6"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="email"} 4"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="timestamp"} 3"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="numeric"} 2"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="cloud_id"} 1"#));
        assert!(text.contains(r#"cbad_tokenizer_replacements_total{type="json"} 200"#));
        assert!(text.contains("cbad_tokenizer_bytes_saved_total 5000"));
    }

    #[test]
    fn test_label_escaping() {
        let snapshot = DetectorMetricsSnapshot {
            events_processed: 1,
            anomalies_detected: 0,
            detection_cycles: 0,
            avg_detection_latency_secs: 0.0,
            baseline_size: 100,
            window_size: 10,
            tokenizer_stats: None,
        };

        // Test escaping of special characters
        let text = snapshot.to_prometheus_text_with_labels(&[("path", "/api/v1/\"test\"")]);
        assert!(text.contains(r#"path="/api/v1/\"test\"""#));
    }

    #[test]
    fn test_detector_builder_profiles() {
        // Test strict profile
        let detector = CbadDetector::builder()
            .with_profile(DetectionProfile::Strict)
            .build()
            .expect("strict profile should build");
        let metrics = detector.metrics();
        assert_eq!(metrics.baseline_size, 1000); // Default baseline

        // Test sensitive profile
        let detector = CbadDetector::builder()
            .with_profile(DetectionProfile::Sensitive)
            .build()
            .expect("sensitive profile should build");
        assert!(detector.metrics().baseline_size > 0);

        // Test custom profile
        let custom = CustomProfile {
            ncd_threshold: 0.5,
            p_value_threshold: 0.1,
            permutation_count: 100,
            ..Default::default()
        };
        let detector = CbadDetector::builder()
            .with_profile(DetectionProfile::Custom(custom))
            .build()
            .expect("custom profile should build");
        assert!(detector.metrics().baseline_size > 0);
    }

    #[test]
    fn test_baseline_strategy_first_n() {
        let detector = CbadDetector::builder()
            .with_baseline_strategy(BaselineStrategy::UseFirstN(50))
            .build()
            .expect("baseline strategy should work");

        // The baseline size should be set to 50
        let config_baseline = detector.metrics().baseline_size;
        assert_eq!(config_baseline, 50);
    }

    #[test]
    fn test_baseline_strategy_explicit() {
        let baseline_data: Vec<Vec<u8>> = (0..10)
            .map(|i| format!("baseline event {}", i).into_bytes())
            .collect();

        let detector = CbadDetector::builder()
            .with_baseline_strategy(BaselineStrategy::ExplicitBaseline(baseline_data))
            .build()
            .expect("explicit baseline should work");

        // The explicit baseline is added directly to the inner detector's window,
        // so we verify by checking that the detector was constructed successfully.
        // The events_processed metric only counts events added via ingest_event().
        let metrics = detector.metrics();
        assert_eq!(metrics.baseline_size, 1000); // Default baseline size
    }

    #[test]
    fn test_progress_callback() {
        use std::sync::atomic::{AtomicU64, Ordering};

        let progress_count = Arc::new(AtomicU64::new(0));
        let progress_count_clone = progress_count.clone();

        let detector = CbadDetector::builder()
            .on_progress(move |_update| {
                progress_count_clone.fetch_add(1, Ordering::Relaxed);
            })
            .build()
            .expect("detector with progress callback should build");

        // Add some events
        for i in 0..5 {
            let _ = detector.ingest_event(format!("event {}", i).into_bytes());
        }

        // Progress callback should have been called
        assert!(progress_count.load(Ordering::Relaxed) >= 5);
    }

    #[test]
    fn test_detector_with_store() {
        let store = Arc::new(InMemoryAnomalyStore::default());
        let detector = CbadDetector::builder()
            .with_store(store.clone())
            .named("test-stream")
            .build()
            .expect("detector with store should build");

        // Verify it works
        let metrics = detector.metrics();
        assert_eq!(metrics.events_processed, 0);
    }

    #[test]
    fn test_csv_config_defaults() {
        let config = CsvConfig::default();
        assert!(config.has_headers);
        assert_eq!(config.delimiter, b',');
        assert!(config.column.is_none());
    }

    #[test]
    fn test_backpressure_config_defaults() {
        let config = BackpressureConfig::default();
        assert_eq!(config.max_in_flight, 1024);
    }

    #[test]
    fn test_parquet_config_defaults() {
        let config = ParquetConfig::default();
        assert!(config.column.is_none());
        assert_eq!(config.batch_size, 0); // Default impl sets to 0
    }

    #[test]
    fn test_parquet_config_with_column() {
        let config = ParquetConfig::with_column("message");
        assert_eq!(config.column.as_deref(), Some("message"));
        assert_eq!(config.batch_size, 1024);
    }

    #[test]
    #[cfg(not(feature = "parquet"))]
    fn test_parquet_stub_returns_error() {
        let detector = CbadDetector::builder().build().expect("build");
        let result = detector.analyze_parquet("test.parquet", ParquetConfig::default());
        assert!(result.is_err());

        let err = result.unwrap_err();
        match err {
            CbadError::UnsupportedFormat(msg) => {
                assert!(msg.contains("parquet"));
                assert!(msg.contains("feature"));
            }
            _ => panic!("expected UnsupportedFormat error"),
        }
    }
}

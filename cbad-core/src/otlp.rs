//! OpenTelemetry OTLP ingestion glue.
//! Feature-gated behind `otlp` to avoid pulling grpc deps by default.

use crate::anomaly::AnomalyConfig;
use crate::api::{CbadDetector, DetectionRecord};
use crate::error::{CbadError, Result};
use async_trait::async_trait;
use opentelemetry_otlp::proto::collector::logs::v1::{
    logs_service_server::{LogsService, LogsServiceServer},
    ExportLogsServiceRequest, ExportLogsServiceResponse,
};
use opentelemetry_otlp::proto::collector::metrics::v1::{
    metrics_service_server::{MetricsService, MetricsServiceServer},
    ExportMetricsServiceRequest, ExportMetricsServiceResponse,
};
use opentelemetry_otlp::proto::collector::trace::v1::{
    trace_service_server::{TraceService, TraceServiceServer},
    ExportTraceServiceRequest, ExportTraceServiceResponse,
};
use opentelemetry_otlp::proto::logs::v1::LogRecord;
use opentelemetry_otlp::proto::metrics::v1::Metric;
use opentelemetry_otlp::proto::trace::v1::Span;
use std::net::SocketAddr;
use std::sync::Arc;
use tonic::{Request, Response, Status};

/// Process OTLP log records with CBAD.
#[derive(Clone)]
pub struct LogProcessor {
    detector: Arc<CbadDetector>,
}

impl LogProcessor {
    pub fn new(detector: Arc<CbadDetector>) -> Self {
        Self { detector }
    }

    pub fn process(&self, record: &LogRecord) -> Result<Option<DetectionRecord>> {
        let serialized = serialize_log_record(record);
        self.detector.ingest_event(serialized)
    }
}

/// Process OTLP metrics payloads.
#[derive(Clone)]
pub struct MetricProcessor {
    detector: Arc<CbadDetector>,
}

impl MetricProcessor {
    pub fn new(detector: Arc<CbadDetector>) -> Self {
        Self { detector }
    }

    pub fn process(&self, metric: &Metric) -> Result<Option<DetectionRecord>> {
        let serialized = format!("{:?}", metric).into_bytes();
        self.detector.ingest_event(serialized)
    }
}

/// Process OTLP trace spans.
#[derive(Clone)]
pub struct TraceProcessor {
    detector: Arc<CbadDetector>,
}

impl TraceProcessor {
    pub fn new(detector: Arc<CbadDetector>) -> Self {
        Self { detector }
    }

    pub fn process(&self, span: &Span) -> Result<Option<DetectionRecord>> {
        let serialized = format!("{:?}", span).into_bytes();
        self.detector.ingest_event(serialized)
    }
}

/// Standalone OTLP receiver suitable for embedding in pipelines.
pub struct CbadOtlpReceiver {
    addr: SocketAddr,
    detector: Arc<CbadDetector>,
    on_anomaly: Option<Arc<dyn Fn(DetectionRecord) + Send + Sync>>,
}

impl CbadOtlpReceiver {
    pub fn bind(addr: impl AsRef<str>) -> Result<Self> {
        let addr: SocketAddr = addr
            .as_ref()
            .parse()
            .map_err(|e| CbadError::InvalidConfig(format!("invalid bind addr: {}", e)))?;

        Ok(Self {
            addr,
            detector: Arc::new(CbadDetector::builder().build()?),
            on_anomaly: None,
        })
    }

    pub fn with_detector_config(mut self, config: AnomalyConfig) -> Result<Self> {
        self.detector = Arc::new(CbadDetector::builder().with_config(config).build()?);
        Ok(self)
    }

    pub fn on_anomaly<F>(mut self, callback: F) -> Self
    where
        F: Fn(DetectionRecord) + Send + Sync + 'static,
    {
        self.on_anomaly = Some(Arc::new(callback));
        self
    }

    pub async fn start(self) -> Result<()> {
        let state = Arc::new(self);

        tonic::transport::Server::builder()
            .add_service(LogsServiceServer::new(LogService {
                state: state.clone(),
            }))
            .add_service(MetricsServiceServer::new(MetricService {
                state: state.clone(),
            }))
            .add_service(TraceServiceServer::new(TraceServiceImpl { state }))
            .serve(state.addr)
            .await
            .map_err(|e| CbadError::InvalidConfig(format!("otlp server error: {}", e)))
    }
}

#[derive(Clone)]
struct LogService {
    state: Arc<CbadOtlpReceiver>,
}

#[derive(Clone)]
struct MetricService {
    state: Arc<CbadOtlpReceiver>,
}

#[derive(Clone)]
struct TraceServiceImpl {
    state: Arc<CbadOtlpReceiver>,
}

#[async_trait]
impl LogsService for LogService {
    async fn export(
        &self,
        request: Request<ExportLogsServiceRequest>,
    ) -> Result<Response<ExportLogsServiceResponse>, Status> {
        let payload = request.into_inner();
        for resource in payload.resource_logs {
            for scope in resource.scope_logs {
                for log in scope.log_records {
                    if let Ok(Some(record)) =
                        self.state.detector.ingest_event(serialize_log_record(&log))
                    {
                        if let Some(cb) = &self.state.on_anomaly {
                            cb(record);
                        }
                    }
                }
            }
        }

        Ok(Response::new(ExportLogsServiceResponse {
            partial_success: None,
        }))
    }
}

#[async_trait]
impl MetricsService for MetricService {
    async fn export(
        &self,
        request: Request<ExportMetricsServiceRequest>,
    ) -> Result<Response<ExportMetricsServiceResponse>, Status> {
        let payload = request.into_inner();
        for resource in payload.resource_metrics {
            for scope in resource.scope_metrics {
                for metric in scope.metrics {
                    if let Ok(Some(record)) = self
                        .state
                        .detector
                        .ingest_event(format!("{:?}", metric).into_bytes())
                    {
                        if let Some(cb) = &self.state.on_anomaly {
                            cb(record);
                        }
                    }
                }
            }
        }

        Ok(Response::new(ExportMetricsServiceResponse {
            partial_success: None,
        }))
    }
}

#[async_trait]
impl TraceService for TraceServiceImpl {
    async fn export(
        &self,
        request: Request<ExportTraceServiceRequest>,
    ) -> Result<Response<ExportTraceServiceResponse>, Status> {
        let payload = request.into_inner();
        for resource in payload.resource_spans {
            for scope in resource.scope_spans {
                for span in scope.spans {
                    if let Ok(Some(record)) = self
                        .state
                        .detector
                        .ingest_event(format!("{:?}", span).into_bytes())
                    {
                        if let Some(cb) = &self.state.on_anomaly {
                            cb(record);
                        }
                    }
                }
            }
        }

        Ok(Response::new(ExportTraceServiceResponse {
            partial_success: None,
        }))
    }
}

fn serialize_log_record(record: &LogRecord) -> Vec<u8> {
    let mut fields = Vec::new();
    if let Some(body) = &record.body {
        fields.push(format!("body={:?}", body));
    }
    for attr in &record.attributes {
        fields.push(format!("{}={:?}", attr.key, attr.value));
    }
    fields.join(" ").into_bytes()
}

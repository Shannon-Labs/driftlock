Driftlock CBAD Collector Processor

This is a skeleton OpenTelemetry Collector processor named `driftlock_cbad`.
It will route logs/metrics through CBAD metrics and produce glass-box explanations.

Usage (custom build required)
- Use `opentelemetry-collector-builder` to create a custom distribution that includes this processor.
- Example manifest (manifest.yaml) is included as a starting point.

Status
- Processor factories and configs are sketched; `cbad.ComputeMetrics` integration is TODO.


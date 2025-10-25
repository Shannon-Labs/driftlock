package cbad

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/your-org/driftlock/api-server/internal/models"
	"github.com/your-org/driftlock/api-server/internal/storage"
	"github.com/your-org/driftlock/api-server/internal/stream"
	"github.com/your-org/driftlock/collector-processor/driftlockcbad"
)

// Detector integrates CBAD core with the API server
type Detector struct {
	storage  *storage.Storage
	streamer *stream.Streamer
	config   *models.DetectionConfig
}

// NewDetector creates a new CBAD detector
func NewDetector(storage *storage.Storage, streamer *stream.Streamer) (*Detector, error) {
	// Validate CBAD library
	if err := driftlockcbad.ValidateLibrary(); err != nil {
		return nil, fmt.Errorf("CBAD library validation failed: %w", err)
	}

	// Load active configuration
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := storage.GetActiveConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	log.Printf("CBAD detector initialized with config: NCD=%.2f, p-value=%.3f", config.NCDThreshold, config.PValueThreshold)

	return &Detector{
		storage:  storage,
		streamer: streamer,
		config:   config,
	}, nil
}

// ProcessTelemetry analyzes telemetry data for anomalies
func (d *Detector) ProcessTelemetry(ctx context.Context, baseline []byte, window []byte, streamType models.StreamType) error {
	start := time.Now()

	// Compute CBAD metrics using Rust FFI
	metrics, err := driftlockcbad.ComputeMetrics(baseline, window, 42, 1000)
	if err != nil {
		return fmt.Errorf("CBAD computation failed: %w", err)
	}

	// Check if anomaly meets threshold
	if !metrics.IsAnomaly || metrics.PValue > d.config.PValueThreshold {
		// No anomaly detected
		return nil
	}

	// Create anomaly record
	anomalyCreate := &models.AnomalyCreate{
		Timestamp:           time.Now(),
		StreamType:          streamType,
		NCDScore:            metrics.NCD,
		PValue:              metrics.PValue,
		GlassBoxExplanation: metrics.GetAnomalyExplanation(),
		CompressionBaseline: metrics.BaselineCompressionRatio,
		CompressionWindow:   metrics.WindowCompressionRatio,
		CompressionCombined: (metrics.BaselineCompressionRatio + metrics.WindowCompressionRatio) / 2,
		ConfidenceLevel:     metrics.ConfidenceLevel,
		Metadata: map[string]interface{}{
			"baseline_entropy": metrics.BaselineEntropy,
			"window_entropy":   metrics.WindowEntropy,
			"processing_time_ms": time.Since(start).Milliseconds(),
		},
		Tags: []string{
			string(streamType),
			fmt.Sprintf("severity:%s", getSeverity(metrics)),
		},
	}

	// Store in database
	anomaly, err := d.storage.CreateAnomaly(ctx, anomalyCreate)
	if err != nil {
		return fmt.Errorf("failed to store anomaly: %w", err)
	}

	// Broadcast to SSE clients
	if d.streamer != nil {
		d.streamer.BroadcastAnomaly(anomaly)
	}

	log.Printf("Anomaly detected: ID=%s, stream=%s, NCD=%.3f, p=%.4f", anomaly.ID, streamType, metrics.NCD, metrics.PValue)

	return nil
}

// ReloadConfig reloads the detection configuration
func (d *Detector) ReloadConfig(ctx context.Context) error {
	config, err := d.storage.GetActiveConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	d.config = config
	log.Printf("CBAD config reloaded: NCD=%.2f, p-value=%.3f", config.NCDThreshold, config.PValueThreshold)

	return nil
}

// GetConfig returns the current detection configuration
func (d *Detector) GetConfig() *models.DetectionConfig {
	return d.config
}

// Helper function to determine severity
func getSeverity(m *driftlockcbad.Metrics) string {
	if m.PValue < 0.001 && m.NCD > 0.5 {
		return "critical"
	}
	if m.PValue < 0.01 && m.NCD > 0.3 {
		return "high"
	}
	if m.PValue < 0.05 {
		return "medium"
	}
	return "low"
}

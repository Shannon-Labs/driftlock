package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/Hmbown/driftlock/api-server/internal/models"
)

// setupTestDB creates a temporary PostgreSQL container for testing
func setupTestDB(t *testing.T) (*Storage, func()) {
	ctx := context.Background()

	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	// Get the host and port
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Create connection string
	connString := fmt.Sprintf("postgres://testuser:password@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Create storage instance
	storage, err := NewPostgres(connString)
	require.NoError(t, err)

	// Run migrations
	err = runMigrations(storage.DB())
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		storage.Close()
		container.Terminate(ctx)
	}

	return storage, cleanup
}

// runMigrations executes the database migrations
func runMigrations(db *sql.DB) error {
	// Create anomalies table
	_, err := db.Exec(`
		CREATE TABLE anomalies (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			timestamp TIMESTAMP NOT NULL,
			stream_type VARCHAR(50) NOT NULL,
			ncd_score DECIMAL(10, 8) NOT NULL,
			p_value DECIMAL(10, 8) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			glass_box_explanation TEXT NOT NULL,
			detailed_explanation TEXT,
			compression_baseline DECIMAL(10, 4) NOT NULL,
			compression_window DECIMAL(10, 4) NOT NULL,
			compression_combined DECIMAL(10, 4) NOT NULL,
			compression_ratio_change DECIMAL(10, 4) NOT NULL,
			baseline_entropy DECIMAL(10, 4),
			window_entropy DECIMAL(10, 4),
			entropy_change DECIMAL(10, 4),
			confidence_level DECIMAL(10, 8) NOT NULL,
			is_statistically_significant BOOLEAN NOT NULL,
			baseline_data JSONB,
			window_data JSONB,
			metadata JSONB,
			tags TEXT[],
			acknowledged_by VARCHAR(255),
			acknowledged_at TIMESTAMP,
			dismissed_by VARCHAR(255),
			dismissed_at TIMESTAMP,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create anomalies table: %w", err)
	}

	// Create detection_config table
	_, err = db.Exec(`
		CREATE TABLE detection_config (
			id SERIAL PRIMARY KEY,
			ncd_threshold DECIMAL(10, 8) NOT NULL DEFAULT 0.2,
			p_value_threshold DECIMAL(10, 8) NOT NULL DEFAULT 0.05,
			baseline_size INTEGER NOT NULL DEFAULT 100,
			window_size INTEGER NOT NULL DEFAULT 100,
			hop_size INTEGER NOT NULL DEFAULT 50,
			stream_overrides JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_by VARCHAR(255),
			notes TEXT,
			is_active BOOLEAN NOT NULL DEFAULT true
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create detection_config table: %w", err)
	}

	// Create performance_metrics table
	_, err = db.Exec(`
		CREATE TABLE performance_metrics (
			id BIGSERIAL PRIMARY KEY,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			metric_type VARCHAR(100) NOT NULL,
			endpoint VARCHAR(255),
			duration_ms DECIMAL(10, 4) NOT NULL,
			success BOOLEAN NOT NULL,
			error_message TEXT,
			metadata JSONB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create performance_metrics table: %w", err)
	}

	// Create audit_log table (for compliance)
	_, err = db.Exec(`
		CREATE TABLE audit_log (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id VARCHAR(255),
			action VARCHAR(100) NOT NULL,
			resource_type VARCHAR(50) NOT NULL,
			resource_id VARCHAR(255) NOT NULL,
			old_values JSONB,
			new_values JSONB
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create audit_log table: %w", err)
	}

	// Insert default active config
	_, err = db.Exec(`
		INSERT INTO detection_config (ncd_threshold, p_value_threshold, baseline_size, window_size, hop_size)
		VALUES (0.2, 0.05, 100, 100, 50);
	`)
	if err != nil {
		return fmt.Errorf("failed to insert default config: %w", err)
	}

	return nil
}

// createTestAnomaly creates a test anomaly for use in tests
func createTestAnomaly() *models.AnomalyCreate {
	now := time.Now()
	baselineData := map[string]interface{}{"key1": "value1", "key2": 123}
	windowData := map[string]interface{}{"key1": "value1_modified", "key2": 124}
	metadata := map[string]interface{}{"source": "test", "version": "1.0"}

	return &models.AnomalyCreate{
		Timestamp:           now,
		StreamType:          models.StreamTypeLogs,
		NCDScore:            0.25,
		PValue:              0.01,
		GlassBoxExplanation: "Compression ratio increased significantly due to reduced entropy",
		CompressionBaseline: 0.6,
		CompressionWindow:   0.3,
		CompressionCombined: 0.45,
		ConfidenceLevel:     0.95,
		BaselineData:        baselineData,
		WindowData:          windowData,
		Metadata:            metadata,
		Tags:                []string{"test", "integration"},
	}
}

func TestCreateAnomaly_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	anomalyCreate := createTestAnomaly()
	
	anomaly, err := storage.CreateAnomaly(context.Background(), anomalyCreate)
	assert.NoError(t, err)
	assert.NotNil(t, anomaly)
	assert.Equal(t, anomalyCreate.StreamType, anomaly.StreamType)
	assert.Equal(t, anomalyCreate.NCDScore, anomaly.NCDScore)
	assert.Equal(t, anomalyCreate.PValue, anomaly.PValue)
	assert.Equal(t, anomalyCreate.GlassBoxExplanation, anomaly.GlassBoxExplanation)
	assert.Equal(t, anomalyCreate.CompressionBaseline, anomaly.CompressionBaseline)
	assert.Equal(t, anomalyCreate.CompressionWindow, anomaly.CompressionWindow)
	assert.Equal(t, anomalyCreate.CompressionCombined, anomaly.CompressionCombined)
	assert.Equal(t, anomalyCreate.ConfidenceLevel, anomaly.ConfidenceLevel)
	assert.Equal(t, anomalyCreate.BaselineData, anomaly.BaselineData)
	assert.Equal(t, anomalyCreate.WindowData, anomaly.WindowData)
	assert.Equal(t, anomalyCreate.Metadata, anomaly.Metadata)
	assert.Equal(t, anomalyCreate.Tags, anomaly.Tags)
	assert.Equal(t, models.StatusPending, anomaly.Status)
	assert.True(t, anomaly.IsStatisticallySignificant)
	assert.Equal(t, 50.0, anomaly.CompressionRatioChange) // Calculated as ((0.3-0.6)/0.6)*100
	
	// Check if the ID was generated
	assert.NotEqual(t, uuid.Nil, anomaly.ID)
	
	// Check timestamps
	assert.False(t, anomaly.CreatedAt.IsZero())
	assert.False(t, anomaly.UpdatedAt.IsZero())
}

func TestGetAnomaly_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// First create an anomaly
	anomalyCreate := createTestAnomaly()
	createdAnomaly, err := storage.CreateAnomaly(context.Background(), anomalyCreate)
	assert.NoError(t, err)
	assert.NotNil(t, createdAnomaly)
	
	// Then retrieve it
	retrievedAnomaly, err := storage.GetAnomaly(context.Background(), createdAnomaly.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedAnomaly)
	
	// Verify all fields match
	assert.Equal(t, createdAnomaly.ID, retrievedAnomaly.ID)
	assert.Equal(t, createdAnomaly.Timestamp, retrievedAnomaly.Timestamp)
	assert.Equal(t, createdAnomaly.StreamType, retrievedAnomaly.StreamType)
	assert.Equal(t, createdAnomaly.NCDScore, retrievedAnomaly.NCDScore)
	assert.Equal(t, createdAnomaly.PValue, retrievedAnomaly.PValue)
	assert.Equal(t, createdAnomaly.GlassBoxExplanation, retrievedAnomaly.GlassBoxExplanation)
	assert.Equal(t, createdAnomaly.CompressionBaseline, retrievedAnomaly.CompressionBaseline)
	assert.Equal(t, createdAnomaly.CompressionWindow, retrievedAnomaly.CompressionWindow)
	assert.Equal(t, createdAnomaly.CompressionCombined, retrievedAnomaly.CompressionCombined)
	assert.Equal(t, createdAnomaly.ConfidenceLevel, retrievedAnomaly.ConfidenceLevel)
	assert.Equal(t, createdAnomaly.BaselineData, retrievedAnomaly.BaselineData)
	assert.Equal(t, createdAnomaly.WindowData, retrievedAnomaly.WindowData)
	assert.Equal(t, createdAnomaly.Metadata, retrievedAnomaly.Metadata)
	assert.Equal(t, createdAnomaly.Tags, retrievedAnomaly.Tags)
	assert.Equal(t, createdAnomaly.IsStatisticallySignificant, retrievedAnomaly.IsStatisticallySignificant)
	assert.Equal(t, createdAnomaly.CompressionRatioChange, retrievedAnomaly.CompressionRatioChange)
	assert.Equal(t, models.StatusPending, retrievedAnomaly.Status)
}

func TestGetAnomaly_NotFound_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Try to get an anomaly that doesn't exist
	id := uuid.New()
	anomaly, err := storage.GetAnomaly(context.Background(), id)
	
	assert.Error(t, err)
	assert.Nil(t, anomaly)
	assert.Contains(t, err.Error(), "anomaly not found")
}

func TestListAnomalies_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple anomalies
	anomaly1 := createTestAnomaly()
	anomaly1.StreamType = models.StreamTypeMetrics
	anomaly1.Timestamp = time.Now().Add(-1 * time.Hour)
	created1, err := storage.CreateAnomaly(context.Background(), anomaly1)
	assert.NoError(t, err)
	
	anomaly2 := createTestAnomaly()
	anomaly2.StreamType = models.StreamTypeLogs
	anomaly2.Timestamp = time.Now()
	created2, err := storage.CreateAnomaly(context.Background(), anomaly2)
	assert.NoError(t, err)

	// Test listing all anomalies
	filter := &models.AnomalyFilter{}
	result, err := storage.ListAnomalies(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Anomalies, 2)
	
	// Results should be ordered by timestamp DESC (newest first)
	assert.Equal(t, created2.ID, result.Anomalies[0].ID)
	assert.Equal(t, created1.ID, result.Anomalies[1].ID)

	// Test filtering by stream type
	filter.StreamType = &anomaly1.StreamType
	result, err = storage.ListAnomalies(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Anomalies, 1)
	assert.Equal(t, created1.ID, result.Anomalies[0].ID)

	// Reset filter for next test
	filter.StreamType = nil
	
	// Test filtering by status
	status := models.StatusPending
	filter.Status = &status
	result, err = storage.ListAnomalies(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Anomalies, 2)

	// Test pagination
	filter.Limit = 1
	filter.Offset = 0
	result, err = storage.ListAnomalies(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Anomalies, 1)
	assert.Equal(t, created2.ID, result.Anomalies[0].ID)
	assert.True(t, result.HasMore)

	// Test second page
	filter.Offset = 1
	result, err = storage.ListAnomalies(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Anomalies, 1)
	assert.Equal(t, created1.ID, result.Anomalies[0].ID)
	assert.False(t, result.HasMore)
}

func TestUpdateAnomalyStatus_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create an anomaly
	anomalyCreate := createTestAnomaly()
	createdAnomaly, err := storage.CreateAnomaly(context.Background(), anomalyCreate)
	assert.NoError(t, err)
	assert.NotNil(t, createdAnomaly)

	// Update the status
	update := &models.AnomalyUpdate{
		Status: models.StatusAcknowledged,
		Notes:  stringPtr("Test note for acknowledged status"),
	}
	
	err = storage.UpdateAnomalyStatus(context.Background(), createdAnomaly.ID, update, "test-user")
	assert.NoError(t, err)

	// Retrieve the updated anomaly to verify
	retrievedAnomaly, err := storage.GetAnomaly(context.Background(), createdAnomaly.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.StatusAcknowledged, retrievedAnomaly.Status)
	assert.Equal(t, "Test note for acknowledged status", *retrievedAnomaly.Notes)
	assert.Equal(t, "test-user", retrievedAnomaly.AcknowledgedBy)
	assert.NotNil(t, retrievedAnomaly.AcknowledgedAt)
}

func TestUpdateAnomalyStatus_NotFound_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Try to update status for an anomaly that doesn't exist
	id := uuid.New()
	update := &models.AnomalyUpdate{
		Status: models.StatusAcknowledged,
	}
	
	err := storage.UpdateAnomalyStatus(context.Background(), id, update, "test-user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "anomaly not found")
}

func TestGetActiveConfig_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the default active config
	config, err := storage.GetActiveConfig(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, config)
	
	// Verify default values
	assert.Equal(t, 0.2, config.NCDThreshold)
	assert.Equal(t, 0.05, config.PValueThreshold)
	assert.Equal(t, 100, config.BaselineSize)
	assert.Equal(t, 100, config.WindowSize)
	assert.Equal(t, 50, config.HopSize)
	assert.True(t, config.IsActive)
	assert.False(t, config.CreatedAt.IsZero())
	assert.False(t, config.UpdatedAt.IsZero())
}

func TestUpdateConfig_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create an update
	newThreshold := 0.3
	newPValue := 0.01
	newBaselineSize := 200
	newWindowSize := 200
	newHopSize := 100
	newNotes := "Updated for testing"
	
	update := &models.DetectionConfigUpdate{
		NCDThreshold:    &newThreshold,
		PValueThreshold: &newPValue,
		BaselineSize:    &newBaselineSize,
		WindowSize:      &newWindowSize,
		HopSize:         &newHopSize,
		Notes:           &newNotes,
	}

	config, err := storage.UpdateConfig(context.Background(), update)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	
	// Verify updated values
	assert.Equal(t, newThreshold, config.NCDThreshold)
	assert.Equal(t, newPValue, config.PValueThreshold)
	assert.Equal(t, newBaselineSize, config.BaselineSize)
	assert.Equal(t, newWindowSize, config.WindowSize)
	assert.Equal(t, newHopSize, config.HopSize)
	assert.Equal(t, newNotes, *config.Notes)
}

func TestRecordPerformanceMetric_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	metric := &models.PerformanceMetric{
		MetricType: "api_request",
		Endpoint:   stringPtr("/api/anomalies"),
		DurationMs: 123.45,
		Success:    true,
		Metadata:   map[string]interface{}{"method": "GET", "user_id": "test-user"},
	}

	err := storage.RecordPerformanceMetric(context.Background(), metric)
	assert.NoError(t, err)

	// Verify the metric was recorded by querying the database directly
	var id int64
	var endpoint *string
	var duration float64
	var success bool
	var metadataJSON []byte
	
	err = storage.db.QueryRow(`
		SELECT id, endpoint, duration_ms, success, metadata 
		FROM performance_metrics 
		WHERE metric_type = $1
		ORDER BY timestamp DESC LIMIT 1
	`, metric.MetricType).Scan(&id, &endpoint, &duration, &success, &metadataJSON)
	
	assert.NoError(t, err)
	assert.Equal(t, *metric.Endpoint, *endpoint)
	assert.Equal(t, metric.DurationMs, duration)
	assert.Equal(t, metric.Success, success)
	
	// Verify metadata was stored as JSON
	var storedMetadata map[string]interface{}
	err = json.Unmarshal(metadataJSON, &storedMetadata)
	assert.NoError(t, err)
	assert.Equal(t, metric.Metadata, storedMetadata)
}

func TestConnectionPooling_Integration(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple anomalies concurrently to test connection pooling
	errChan := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			anomalyCreate := createTestAnomaly()
			_, err := storage.CreateAnomaly(context.Background(), anomalyCreate)
			errChan <- err
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		err := <-errChan
		assert.NoError(t, err)
	}

	// Verify all anomalies were created
	filter := &models.AnomalyFilter{}
	result, err := storage.ListAnomalies(context.Background(), filter)
	assert.NoError(t, err)
	assert.Equal(t, 10, result.Total)
}

func stringPtr(s string) *string {
	return &s
}
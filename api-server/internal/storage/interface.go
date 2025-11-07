package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
)

// AnomalyStorage defines the interface for anomaly storage operations
type AnomalyStorage interface {
	CreateAnomaly(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error)
	GetAnomaly(ctx context.Context, id uuid.UUID) (*models.Anomaly, error)
	ListAnomalies(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error)
	UpdateAnomalyStatus(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error
}

// ConfigStorage defines the interface for configuration storage operations
type ConfigStorage interface {
	GetActiveConfig(ctx context.Context) (*models.DetectionConfig, error)
	UpdateConfig(ctx context.Context, update *models.DetectionConfigUpdate) (*models.DetectionConfig, error)
}

// MetricsStorage defines the interface for metrics storage operations
type MetricsStorage interface {
	RecordPerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error
}

// Storage defines the full interface for all storage operations
type StorageInterface interface {
	AnomalyStorage
	ConfigStorage
	MetricsStorage
	Close() error
	Ping(ctx context.Context) error
	DB() *sql.DB
}

// Ensure Storage implements StorageInterface
var _ StorageInterface = (*Storage)(nil)

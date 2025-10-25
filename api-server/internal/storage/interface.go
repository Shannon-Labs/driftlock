package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/your-org/driftlock/api-server/internal/models"
)

// AnomalyStorage defines the interface for anomaly storage operations
type AnomalyStorage interface {
	CreateAnomaly(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error)
	GetAnomaly(ctx context.Context, id uuid.UUID) (*models.Anomaly, error)
	ListAnomalies(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error)
	UpdateAnomalyStatus(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error
}

// Ensure Storage implements AnomalyStorage
var _ AnomalyStorage = (*Storage)(nil)

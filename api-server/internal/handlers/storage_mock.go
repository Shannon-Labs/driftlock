package handlers

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/driftlock/api-server/internal/models"
)

// MockStorage is a mock implementation of the storage interface
type MockStorage struct {
	mock.Mock
}

// AnomalyStorage methods
func (m *MockStorage) CreateAnomaly(ctx context.Context, create *models.AnomalyCreate) (*models.Anomaly, error) {
	args := m.Called(ctx, create)
	return args.Get(0).(*models.Anomaly), args.Error(1)
}

func (m *MockStorage) GetAnomaly(ctx context.Context, id uuid.UUID) (*models.Anomaly, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Anomaly), args.Error(1)
}

func (m *MockStorage) ListAnomalies(ctx context.Context, filter *models.AnomalyFilter) (*models.AnomalyListResponse, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.AnomalyListResponse), args.Error(1)
}

func (m *MockStorage) UpdateAnomalyStatus(ctx context.Context, id uuid.UUID, update *models.AnomalyUpdate, username string) error {
	args := m.Called(ctx, id, update, username)
	return args.Error(0)
}

// Config methods
func (m *MockStorage) GetActiveConfig(ctx context.Context) (*models.DetectionConfig, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.DetectionConfig), args.Error(1)
}

func (m *MockStorage) UpdateConfig(ctx context.Context, update *models.DetectionConfigUpdate) (*models.DetectionConfig, error) {
	args := m.Called(ctx, update)
	return args.Get(0).(*models.DetectionConfig), args.Error(1)
}

// Performance metrics methods
func (m *MockStorage) RecordPerformanceMetric(ctx context.Context, metric *models.PerformanceMetric) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

// Storage interface methods that might be needed (for completeness)
func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) DB() *sql.DB {
	args := m.Called()
	return args.Get(0).(*sql.DB)
}
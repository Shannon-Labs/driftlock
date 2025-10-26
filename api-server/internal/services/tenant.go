package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/your-org/driftlock/api-server/internal/models"
)

// TenantService handles tenant operations
type TenantService struct {
	db *sql.DB
}

// NewTenantService creates a new tenant service
func NewTenantService(db *sql.DB) *TenantService {
	return &TenantService{
		db: db,
	}
}

// CreateTenant creates a new tenant
func (s *TenantService) CreateTenant(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	// Validate tenant data
	if tenant.Name == "" {
		return nil, fmt.Errorf("tenant name is required")
	}
	if tenant.Domain == "" {
		return nil, fmt.Errorf("tenant domain is required")
	}

	// Set default values
	if tenant.Status == "" {
		tenant.Status = models.TenantStatusTrial
	}
	if tenant.Plan == "" {
		tenant.Plan = models.TenantPlanTrial
	}

	// Set timestamps
	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	// Initialize usage
	tenant.Usage = models.TenantUsage{
		LastReset: now,
	}

	// Set default quotas based on plan
	tenant.Quotas = s.getDefaultQuotasForPlan(tenant.Plan)

	// Insert tenant into database
	query := `
		INSERT INTO tenants (id, name, domain, email, status, plan, created_at, updated_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	err := s.db.QueryRowContext(
		ctx, query,
		tenant.ID, tenant.Name, tenant.Domain, tenant.Email, tenant.Status, tenant.Plan,
		tenant.CreatedAt, tenant.UpdatedAt, tenant.ExpiresAt,
	).Scan()
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// Create tenant schema if using schema isolation
	if tenant.Status == models.TenantStatusActive {
		if err := s.createTenantSchema(ctx, tenant.ID); err != nil {
			return nil, fmt.Errorf("failed to create tenant schema: %w", err)
		}
	}

	return tenant, nil
}

// GetTenant retrieves a tenant by ID
func (s *TenantService) GetTenant(ctx context.Context, tenantID string) (*models.Tenant, error) {
	tenant := &models.Tenant{}

	query := `
		SELECT id, name, domain, email, status, plan, created_at, updated_at, expires_at
		FROM tenants
		WHERE id = $1
	`

	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(
		&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Email, &tenant.Status, &tenant.Plan,
		&tenant.CreatedAt, &tenant.UpdatedAt, &tenant.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Get tenant usage
	usage, err := s.getTenantUsage(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant usage: %w", err)
	}
	tenant.Usage = *usage

	// Get tenant quotas
	quotas, err := s.getTenantQuotas(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant quotas: %w", err)
	}
	tenant.Quotas = *quotas

	return tenant, nil
}

// UpdateTenant updates an existing tenant
func (s *TenantService) UpdateTenant(ctx context.Context, tenant *models.Tenant) (*models.Tenant, error) {
	// Validate tenant data
	if tenant.ID == "" {
		return nil, fmt.Errorf("tenant ID is required")
	}

	// Set updated timestamp
	tenant.UpdatedAt = time.Now()

	query := `
		UPDATE tenants
		SET name = $2, domain = $3, email = $4, status = $5, plan = $6, updated_at = $7, expires_at = $8
		WHERE id = $1
	`

	_, err := s.db.ExecContext(
		ctx, query,
		tenant.ID, tenant.Name, tenant.Domain, tenant.Email, tenant.Status, tenant.Plan,
		tenant.UpdatedAt, tenant.ExpiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return s.GetTenant(ctx, tenant.ID)
}

// DeleteTenant deletes a tenant
func (s *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	// Check if tenant exists (just verify it exists without using the data)
	_, err := s.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("tenant not found: %w", err)
	}

	// Drop tenant schema if using schema isolation
	if err := s.dropTenantSchema(ctx, tenantID); err != nil {
		return fmt.Errorf("failed to drop tenant schema: %w", err)
	}

	// Delete tenant from database
	query := `DELETE FROM tenants WHERE id = $1`
	_, err = s.db.ExecContext(ctx, query, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	return nil
}

// ListTenants retrieves all tenants
func (s *TenantService) ListTenants(ctx context.Context) ([]*models.Tenant, error) {
	query := `
		SELECT id, name, domain, email, status, plan, created_at, updated_at, expires_at
		FROM tenants
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer rows.Close()

	var tenants []*models.Tenant
	for rows.Next() {
		tenant := &models.Tenant{}
		err := rows.Scan(
			&tenant.ID, &tenant.Name, &tenant.Domain, &tenant.Email, &tenant.Status, &tenant.Plan,
			&tenant.CreatedAt, &tenant.UpdatedAt, &tenant.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant: %w", err)
		}

		// Get tenant usage
		usage, err := s.getTenantUsage(ctx, tenant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tenant usage: %w", err)
		}
		tenant.Usage = *usage

		// Get tenant quotas
		quotas, err := s.getTenantQuotas(ctx, tenant.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tenant quotas: %w", err)
		}
		tenant.Quotas = *quotas

		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

// UpdateTenantUsage updates resource usage for a tenant
func (s *TenantService) UpdateTenantUsage(ctx context.Context, tenantID string, usageType string, amount int) error {
	// Get current usage
	usage, err := s.getTenantUsage(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant usage: %w", err)
	}

	// Reset daily counters if needed
	now := time.Now()
	if now.Sub(usage.LastReset) >= 24*time.Hour {
		usage.AnomaliesDetected = 0
		usage.EventsProcessed = 0
		usage.APIRequestsToday = 0
		usage.LastReset = now
	}

	// Update the specific usage counter
	switch usageType {
	case "anomalies":
		usage.AnomaliesDetected += amount
	case "events":
		usage.EventsProcessed += amount
	case "api_requests":
		usage.APIRequestsToday += amount
	case "storage":
		usage.StorageUsedGB += amount
	default:
		return fmt.Errorf("unknown usage type: %s", usageType)
	}

	// Update usage in database
	query := `
		UPDATE tenant_usage
		SET anomalies_detected = $2, events_processed = $3, storage_used_gb = $4, 
		    api_requests_today = $5, last_reset = $6
		WHERE tenant_id = $1
	`

	_, err = s.db.ExecContext(
		ctx, query,
		tenantID, usage.AnomaliesDetected, usage.EventsProcessed,
		usage.StorageUsedGB, usage.APIRequestsToday, usage.LastReset,
	)
	if err != nil {
		return fmt.Errorf("failed to update tenant usage: %w", err)
	}

	return nil
}

// CheckTenantQuota checks if a tenant has exceeded their quota
func (s *TenantService) CheckTenantQuota(ctx context.Context, tenantID string, resourceType string) (bool, error) {
	tenant, err := s.GetTenant(ctx, tenantID)
	if err != nil {
		return false, fmt.Errorf("failed to get tenant: %w", err)
	}

	switch resourceType {
	case "anomalies":
		return tenant.Usage.AnomaliesDetected >= tenant.Quotas.MaxAnomaliesPerDay, nil
	case "events":
		return tenant.Usage.EventsProcessed >= tenant.Quotas.MaxEventsPerDay, nil
	case "storage":
		return tenant.Usage.StorageUsedGB >= tenant.Quotas.MaxStorageGB, nil
	case "api_requests":
		// Check requests per minute
		// This is a simplified check - in production, you'd want to track this more precisely
		return tenant.Usage.APIRequestsToday >= tenant.Quotas.MaxAPIRequestsPerMin, nil
	default:
		return false, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

// Helper functions

func (s *TenantService) getDefaultQuotasForPlan(plan string) models.TenantQuotas {
	switch plan {
	case models.TenantPlanTrial:
		return models.TenantQuotas{
			MaxAnomaliesPerDay:   100,
			MaxEventsPerDay:      1000,
			MaxStorageGB:         10,
			MaxAPIRequestsPerMin: 60,
		}
	case models.TenantPlanStarter:
		return models.TenantQuotas{
			MaxAnomaliesPerDay:   500,
			MaxEventsPerDay:      5000,
			MaxStorageGB:         50,
			MaxAPIRequestsPerMin: 300,
		}
	case models.TenantPlanPro:
		return models.TenantQuotas{
			MaxAnomaliesPerDay:   2000,
			MaxEventsPerDay:      20000,
			MaxStorageGB:         200,
			MaxAPIRequestsPerMin: 1200,
		}
	case models.TenantPlanEnterprise:
		return models.TenantQuotas{
			MaxAnomaliesPerDay:   10000,
			MaxEventsPerDay:      100000,
			MaxStorageGB:         1000,
			MaxAPIRequestsPerMin: 6000,
		}
	default:
		// Return trial quotas as default
		return s.getDefaultQuotasForPlan(models.TenantPlanTrial)
	}
}

func (s *TenantService) getTenantUsage(ctx context.Context, tenantID string) (*models.TenantUsage, error) {
	usage := &models.TenantUsage{}

	query := `
		SELECT anomalies_detected, events_processed, storage_used_gb, api_requests_today, last_reset
		FROM tenant_usage
		WHERE tenant_id = $1
	`

	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(
		&usage.AnomaliesDetected, &usage.EventsProcessed,
		&usage.StorageUsedGB, &usage.APIRequestsToday, &usage.LastReset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant usage: %w", err)
	}

	return usage, nil
}

func (s *TenantService) getTenantQuotas(ctx context.Context, tenantID string) (*models.TenantQuotas, error) {
	quotas := &models.TenantQuotas{}

	query := `
		SELECT max_anomalies_per_day, max_events_per_day, max_storage_gb, max_api_requests_per_min
		FROM tenant_quotas
		WHERE tenant_id = $1
	`

	err := s.db.QueryRowContext(ctx, query, tenantID).Scan(
		&quotas.MaxAnomaliesPerDay, &quotas.MaxEventsPerDay,
		&quotas.MaxStorageGB, &quotas.MaxAPIRequestsPerMin,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant quotas: %w", err)
	}

	return quotas, nil
}

func (s *TenantService) createTenantSchema(ctx context.Context, tenantID string) error {
	// This is a simplified implementation
	// In a real production environment, you would:
	// 1. Create a new database schema
	// 2. Create all tables in the new schema
	// 3. Set up proper permissions

	schemaName := fmt.Sprintf("tenant_%s", tenantID)
	query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create tenant schema: %w", err)
	}

	return nil
}

func (s *TenantService) dropTenantSchema(ctx context.Context, tenantID string) error {
	// This is a simplified implementation
	// In a real production environment, you would:
	// 1. Drop all tables in the schema
	// 2. Drop the schema itself

	schemaName := fmt.Sprintf("tenant_%s", tenantID)
	query := fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schemaName)

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop tenant schema: %w", err)
	}

	return nil
}

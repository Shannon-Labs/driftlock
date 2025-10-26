package models

import (
	"time"
)

// Tenant represents a customer organization
type Tenant struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Domain    string     `json:"domain" db:"domain"`
	Email     string     `json:"email" db:"email"`   // Primary contact email
	Status    string     `json:"status" db:"status"` // "active", "suspended", "trial"
	Plan      string     `json:"plan" db:"plan"`     // "trial", "starter", "pro", "enterprise"
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`

	// Resource usage tracking
	Usage  TenantUsage  `json:"usage"`
	Quotas TenantQuotas `json:"quotas"`
}

// TenantUsage tracks resource consumption for a tenant
type TenantUsage struct {
	AnomaliesDetected int       `json:"anomalies_detected"`
	EventsProcessed   int       `json:"events_processed"`
	StorageUsedGB     int       `json:"storage_used_gb"`
	APIRequestsToday  int       `json:"api_requests_today"`
	LastReset         time.Time `json:"last_reset"`
}

// TenantQuotas represents resource limits for a tenant
type TenantQuotas struct {
	MaxAnomaliesPerDay   int `json:"max_anomalies_per_day"`
	MaxEventsPerDay      int `json:"max_events_per_day"`
	MaxStorageGB         int `json:"max_storage_gb"`
	MaxAPIRequestsPerMin int `json:"max_api_requests_per_min"`
}

// TenantStatus constants
const (
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
	TenantStatusTrial     = "trial"
)

// TenantPlan constants
const (
	TenantPlanTrial      = "trial"
	TenantPlanStarter    = "starter"
	TenantPlanPro        = "pro"
	TenantPlanEnterprise = "enterprise"
)

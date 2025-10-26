package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name"`
	Password  string         `json:"-" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Tenant struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	OwnerID     uint           `json:"-" gorm:"not null"`
	Owner       User           `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Settings    TenantSettings `json:"settings" gorm:"foreignKey:TenantID"`
	Anomalies   []Anomaly      `json:"anomalies" gorm:"foreignKey:TenantID"`
}

type TenantSettings struct {
	ID                    uint   `json:"id" gorm:"primaryKey"`
	TenantID              uint   `json:"tenant_id" gorm:"uniqueIndex;not null"`
	LogAnomalyEnabled     bool   `json:"log_anomaly_enabled" gorm:"default:true"`
	MetricAnomalyEnabled  bool   `json:"metric_anomaly_enabled" gorm:"default:true"`
	AnomalyThresholdMs    int    `json:"anomaly_threshold_ms" gorm:"default:1000"`
	AlertEmails           string `json:"alert_emails"`
	AlertWebhookURL       string `json:"alert_webhook_url"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type Anomaly struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenant_id" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null"` // "log", "metric", etc.
	Severity    string    `json:"severity"`             // "low", "medium", "high", "critical"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Source      string    `json:"source"` // where the anomaly was detected
	DetectedAt  time.Time `json:"detected_at"`
	Tags        string    `json:"tags"`         // JSON string of tags
	RawData     string    `json:"raw_data"`     // original data that triggered the anomaly
	Resolved    bool      `json:"resolved"`     // whether the anomaly has been acknowledged/resolved
	ResolvedAt  *time.Time `json:"resolved_at"` // when the anomaly was resolved
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
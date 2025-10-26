package config

import (
	"os"
)

type AnalyticsConfig struct {
	GA4MeasurementID string
	GA4APIKey        string
	CloudflareAPIKey string
	CloudflareAccountID string
	EnableTracking   bool
	EnableAuditLogs  bool
}

func LoadAnalyticsConfig() *AnalyticsConfig {
	return &AnalyticsConfig{
		GA4MeasurementID: os.Getenv("GA4_MEASUREMENT_ID"),
		GA4APIKey:        os.Getenv("GA4_API_KEY"),
		CloudflareAPIKey: os.Getenv("CLOUDFLARE_API_KEY"),
		CloudflareAccountID: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		EnableTracking:   os.Getenv("ENABLE_ANALYTICS") == "true",
		EnableAuditLogs:  os.Getenv("ENABLE_AUDIT_LOGS") == "true",
	}
}
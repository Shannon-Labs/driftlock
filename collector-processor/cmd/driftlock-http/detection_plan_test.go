package main

import "testing"

// These tests ensure we always require statistical significance by default so
// production callers don't get noisy anomaly spam.

func TestBuildDetectionSettingsRequiresSignificance(t *testing.T) {
	cfg := config{}
	settings := streamSettings{}
	settings.applyDefaults()
	stream := streamRecord{}

	plan := buildDetectionSettings(cfg, stream, settings, nil)
	if !plan.RequireStatisticalSignificance {
		t.Fatalf("expected RequireStatisticalSignificance to default to true for production detection plan")
	}
}

func TestBuildDemoDetectionSettingsRequiresSignificance(t *testing.T) {
	cfg := config{}
	plan := buildDemoDetectionSettings(cfg, nil)
	if !plan.RequireStatisticalSignificance {
		t.Fatalf("expected RequireStatisticalSignificance to default to true for demo detection plan")
	}
}

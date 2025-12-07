package main

import (
	"testing"
)

func TestDetectionPlanCapsPermutationAndSizes(t *testing.T) {
	cfg := config{
		DefaultAlgo: "zstd",
		Seed:        42,
		MaxEvents:   1000,
	}
	stream := streamRecord{Seed: 0, Compressor: ""}
	settings := streamSettings{
		BaselineSize:     5000,
		WindowSize:       2000,
		HopSize:          1500,
		NCDThreshold:     0.3,
		PValueThreshold:  0.05,
		PermutationCount: 5000,
		Compressor:       "zstd",
	}

	plan := buildDetectionSettings(cfg, stream, settings, nil)

	if plan.PermutationCount > 300 {
		t.Fatalf("permutation count should be capped, got %d", plan.PermutationCount)
	}
	if plan.BaselineSize > cfg.MaxEvents {
		t.Fatalf("baseline size should respect MaxEvents cap: %d > %d", plan.BaselineSize, cfg.MaxEvents)
	}
	if plan.WindowSize > cfg.MaxEvents {
		t.Fatalf("window size should respect MaxEvents cap: %d > %d", plan.WindowSize, cfg.MaxEvents)
	}
	if plan.HopSize > plan.WindowSize {
		t.Fatalf("hop size should not exceed window size: hop=%d window=%d", plan.HopSize, plan.WindowSize)
	}
}

func TestValidateConfigOverride(t *testing.T) {
	cfg := config{
		DefaultBaseline: 400,
		DefaultWindow:   50,
		MaxEvents:       1000,
	}

	// Valid override
	valid := &configOverride{
		BaselineSize:     intPtr(100),
		WindowSize:       intPtr(50),
		HopSize:          intPtr(25),
		PermutationCount: intPtr(200),
	}
	if err := validateConfigOverride(cfg, valid); err != nil {
		t.Fatalf("valid override rejected: %v", err)
	}

	// Baseline too small
	tooSmallBaseline := &configOverride{BaselineSize: intPtr(5)}
	if err := validateConfigOverride(cfg, tooSmallBaseline); err == nil {
		t.Fatal("expected error for too-small baseline")
	}

	// Window exceeds baseline
	windowTooBig := &configOverride{
		BaselineSize: intPtr(20),
		WindowSize:   intPtr(30),
	}
	if err := validateConfigOverride(cfg, windowTooBig); err == nil {
		t.Fatal("expected error for window > baseline")
	}

	// Hop exceeds window
	hopTooBig := &configOverride{
		BaselineSize: intPtr(100),
		WindowSize:   intPtr(50),
		HopSize:      intPtr(60),
	}
	if err := validateConfigOverride(cfg, hopTooBig); err == nil {
		t.Fatal("expected error for hop > window")
	}

	// Permutation too large
	permTooLarge := &configOverride{PermutationCount: intPtr(200000)}
	if err := validateConfigOverride(cfg, permTooLarge); err == nil {
		t.Fatal("expected error for large permutation_count")
	}

	// Capacity overflow: baseline+3*window exceeds maxEvents
	capacityOverflow := &configOverride{
		BaselineSize: intPtr(900),
		WindowSize:   intPtr(400),
	}
	if err := validateConfigOverride(cfg, capacityOverflow); err == nil {
		t.Fatal("expected error for capacity overflow")
	}
}

func intPtr(v int) *int { return &v }

package main

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

func TestComputeAutoTuneAdjustment_HighFalsePositiveRate(t *testing.T) {
	stats := FeedbackStats{
		TotalFeedback:     30,
		FalsePositives:    9,
		Confirmed:         5,
		FalsePositiveRate: 0.30,
		AvgFPNCD:          0.29,
	}

	currentNCD := 0.30
	currentPValue := 0.05

	newNCD, newPValue, reason, shouldAdjust := computeAutoTuneAdjustment(stats, currentNCD, currentPValue, defaultAutoTuneConfig)
	if !shouldAdjust {
		t.Fatalf("expected adjustment for high false positive rate")
	}
	if reason != "high_false_positive_rate" {
		t.Fatalf("unexpected reason: %s", reason)
	}
	if newNCD <= currentNCD {
		t.Fatalf("expected NCD threshold to increase, got %f", newNCD)
	}
	if newPValue >= currentPValue {
		t.Fatalf("expected P-value threshold to decrease, got %f", newPValue)
	}
	expectedNCD := 0.3375
	expectedPValue := 0.0375
	if math.Abs(newNCD-expectedNCD) > 1e-4 {
		t.Fatalf("unexpected NCD value: %f", newNCD)
	}
	if math.Abs(newPValue-expectedPValue) > 1e-4 {
		t.Fatalf("unexpected P-value: %f", newPValue)
	}
	if newNCD < 0.1 || newNCD > 0.8 {
		t.Fatalf("NCD threshold out of bounds: %f", newNCD)
	}
	if newPValue < 0.001 || newPValue > 0.2 {
		t.Fatalf("P-value threshold out of bounds: %f", newPValue)
	}
}

func TestComputeAutoTuneAdjustment_InsufficientFeedback(t *testing.T) {
	stats := FeedbackStats{
		TotalFeedback:     5,
		FalsePositiveRate: 0.20,
	}

	currentNCD := 0.30
	currentPValue := 0.05

	newNCD, newPValue, reason, shouldAdjust := computeAutoTuneAdjustment(stats, currentNCD, currentPValue, defaultAutoTuneConfig)
	if shouldAdjust {
		t.Fatalf("did not expect adjustment with insufficient feedback")
	}
	if reason != "insufficient_feedback" {
		t.Fatalf("unexpected reason: %s", reason)
	}
	if newNCD != currentNCD || newPValue != currentPValue {
		t.Fatalf("thresholds should remain unchanged when not adjusting")
	}
}

func TestComputeAutoTuneAdjustment_LowDetectionRate(t *testing.T) {
	// When FP rate is very low (< half of target) and there are confirmed anomalies,
	// the system should increase sensitivity by lowering the NCD threshold
	stats := FeedbackStats{
		TotalFeedback:     30,
		FalsePositives:    0,  // 0% FP rate - well below 2.5% threshold (half of 5% target)
		Confirmed:         10, // Has confirmed anomalies, proving real anomalies exist
		FalsePositiveRate: 0.0,
	}

	currentNCD := 0.30
	currentPValue := 0.05

	newNCD, _, reason, shouldAdjust := computeAutoTuneAdjustment(
		stats, currentNCD, currentPValue, defaultAutoTuneConfig)

	if !shouldAdjust {
		t.Fatalf("expected adjustment for low detection rate (FP rate 0%% < 2.5%% threshold)")
	}
	if reason != "low_detection_rate" {
		t.Fatalf("unexpected reason: %s, expected low_detection_rate", reason)
	}
	if newNCD >= currentNCD {
		t.Fatalf("expected NCD threshold to decrease for more sensitivity, got %f >= %f", newNCD, currentNCD)
	}
}

func TestComputeAdaptiveWindowSizes_MemoryBounded(t *testing.T) {
	chars := StreamCharacteristics{
		AvgEventsPerHour:   5000,
		AvgEventSizeBytes:  5 * 1024 * 1024, // large events to trigger memory bound
		AvgBaselineEntropy: 7.0,
		PatternDiversity:   0.8,
	}

	baseline, window := computeAdaptiveWindowSizes(chars, defaultAdaptiveConfig)

	if baseline != 100 {
		t.Fatalf("expected baseline to be memory-clamped to 100, got %d", baseline)
	}
	if window != 12 {
		t.Fatalf("expected window to scale with baseline to 12, got %d", window)
	}
	if window > baseline {
		t.Fatalf("window should not exceed baseline, got window=%d baseline=%d", window, baseline)
	}
}

func TestApplyProfileDefaultsAndCustom(t *testing.T) {
	t.Run("strict profile uses defaults", func(t *testing.T) {
		plan := detectionPlan{}
		applyProfile(&plan, ProfileStrict, nil, nil)

		if plan.NCDThreshold != profileDefaults[ProfileStrict].NCDThreshold {
			t.Fatalf("strict profile NCD mismatch: %f", plan.NCDThreshold)
		}
		if plan.PValueThreshold != profileDefaults[ProfileStrict].PValueThreshold {
			t.Fatalf("strict profile p-value mismatch: %f", plan.PValueThreshold)
		}
		if plan.BaselineSize != profileDefaults[ProfileStrict].BaselineSize {
			t.Fatalf("strict profile baseline mismatch: %d", plan.BaselineSize)
		}
		if plan.WindowSize != profileDefaults[ProfileStrict].WindowSize {
			t.Fatalf("strict profile window mismatch: %d", plan.WindowSize)
		}
	})

	t.Run("custom profile uses tuned thresholds", func(t *testing.T) {
		tunedNCD := 0.55
		tunedPValue := 0.02
		plan := detectionPlan{}
		applyProfile(&plan, ProfileCustom, &tunedNCD, &tunedPValue)

		if plan.NCDThreshold != tunedNCD {
			t.Fatalf("expected tuned NCD threshold %f, got %f", tunedNCD, plan.NCDThreshold)
		}
		if plan.PValueThreshold != tunedPValue {
			t.Fatalf("expected tuned p-value %f, got %f", tunedPValue, plan.PValueThreshold)
		}
		if plan.BaselineSize != profileDefaults[ProfileCustom].BaselineSize {
			t.Fatalf("custom profile baseline mismatch: %d", plan.BaselineSize)
		}
		if plan.WindowSize != profileDefaults[ProfileCustom].WindowSize {
			t.Fatalf("custom profile window mismatch: %d", plan.WindowSize)
		}
	})
}

func TestShouldThrottleAutoTune(t *testing.T) {
	now := time.Now()
	cfg := defaultAutoTuneConfig

	if shouldThrottleAutoTune(nil, cfg, now) {
		t.Fatalf("nil last tune should not throttle")
	}

	lastRecent := now.Add(-cfg.CooldownPeriod / 2)
	if !shouldThrottleAutoTune(&lastRecent, cfg, now) {
		t.Fatalf("recent tune should throttle")
	}

	lastOld := now.Add(-cfg.CooldownPeriod * 2)
	if shouldThrottleAutoTune(&lastOld, cfg, now) {
		t.Fatalf("old tune should not throttle")
	}
}

func TestDeriveStreamCharacteristics(t *testing.T) {
	events := []json.RawMessage{
		json.RawMessage(`{"msg":"a","size":1}`),
		json.RawMessage(`{"msg":"b","size":2}`),
		json.RawMessage(`{"msg":"c","size":3}`),
	}

	chars := deriveStreamCharacteristics(events)
	if chars.AvgEventSizeBytes == 0 {
		t.Fatalf("expected AvgEventSizeBytes to be computed")
	}
	if chars.PatternDiversity < 0 {
		t.Fatalf("expected PatternDiversity to be non-negative")
	}
	if chars.AvgBaselineEntropy <= 0 {
		t.Fatalf("expected entropy to be > 0")
	}
}

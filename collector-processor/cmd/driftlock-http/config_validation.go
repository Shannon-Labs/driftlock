package main

import "fmt"

// validateConfigOverride ensures override values cannot exhaust resources.
func validateConfigOverride(cfg config, override *configOverride) error {
	if override == nil {
		return nil
	}

	maxEvents := cfg.MaxEvents
	if maxEvents <= 0 {
		maxEvents = 10000
	}

	if override.BaselineSize != nil {
		if *override.BaselineSize < 10 || *override.BaselineSize > 10000 {
			return fmt.Errorf("baseline_size override must be between 10 and 10000")
		}
	}
	if override.WindowSize != nil {
		if *override.WindowSize < 1 || *override.WindowSize > 1000 {
			return fmt.Errorf("window_size override must be between 1 and 1000")
		}
		if override.BaselineSize != nil && *override.WindowSize > *override.BaselineSize {
			return fmt.Errorf("window_size cannot exceed baseline_size")
		}
	}
	if override.HopSize != nil && override.WindowSize != nil {
		if *override.HopSize <= 0 || *override.HopSize > *override.WindowSize {
			return fmt.Errorf("hop_size must be >0 and <= window_size")
		}
	}

	if override.PermutationCount != nil {
		if *override.PermutationCount < 1 || *override.PermutationCount > 10000 {
			return fmt.Errorf("permutation_count must be between 1 and 10000")
		}
	}

	// Capacity guardrail based on overrides or defaults
	baseline := max(cfg.DefaultBaseline, 10)
	window := max(cfg.DefaultWindow, 1)
	if override.BaselineSize != nil {
		baseline = *override.BaselineSize
	}
	if override.WindowSize != nil {
		window = *override.WindowSize
	}
	required := baseline + 3*window
	if required > maxEvents {
		return fmt.Errorf("baseline+3*window exceeds allowed capacity (%d > %d)", required, maxEvents)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

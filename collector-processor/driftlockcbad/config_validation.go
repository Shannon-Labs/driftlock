package driftlockcbad

import (
	"fmt"
)

// ValidateConfig checks if the detector configuration is valid
func ValidateConfig(cfg DetectorConfig) error {
	// p-value threshold should be between 0 and 1
	if cfg.PValueThreshold <= 0 || cfg.PValueThreshold > 1 {
		return fmt.Errorf("p_value_threshold must be between 0 and 1, got %f", cfg.PValueThreshold)
	}

	// NCD threshold should be positive and reasonable
	if cfg.NCDThreshold <= 0 || cfg.NCDThreshold > 2 {
		return fmt.Errorf("ncd_threshold must be between 0 and 2, got %f", cfg.NCDThreshold)
	}

	// Baseline size should be reasonable
	if cfg.BaselineSize < 10 || cfg.BaselineSize > 50000 {
		return fmt.Errorf("baseline_size must be between 10 and 50000, got %d", cfg.BaselineSize)
	}

	// Window size should not exceed baseline size (to avoid NCD distortion)
	if cfg.WindowSize > cfg.BaselineSize {
		return fmt.Errorf("window_size (%d) should not exceed baseline_size (%d) to avoid NCD distortion",
			cfg.WindowSize, cfg.BaselineSize)
	}

	// Hop size should be reasonable
	if cfg.HopSize <= 0 || cfg.HopSize > cfg.WindowSize {
		return fmt.Errorf("hop_size (%d) should be >0 and <= window_size (%d)", cfg.HopSize, cfg.WindowSize)
	}

	// Permutation count affects performance significantly
	if cfg.PermutationCount < 10 || cfg.PermutationCount > 1000 {
		return fmt.Errorf("permutation_count (%d) should be between 10 and 1000 for performance",
			cfg.PermutationCount)
	}

	// Compression algorithm should be supported
	validAlgos := map[string]bool{
		"zstd": true, "zlib": true, "gzip": true, "lz4": true,
	}
	if !validAlgos[cfg.CompressionAlgorithm] {
		return fmt.Errorf("unsupported compression algorithm: %s", cfg.CompressionAlgorithm)
	}

	// MaxCapacity should be sufficient for baseline + multiple windows
	minCapacity := cfg.BaselineSize + 3*cfg.WindowSize
	if cfg.MaxCapacity < minCapacity {
		return fmt.Errorf("max_capacity (%d) should be at least baseline_size + 3*window_size (%d)",
			cfg.MaxCapacity, minCapacity)
	}

	// Ensure drop/entropy/composite thresholds are non-negative
	if cfg.CompressionRatioDropThreshold < 0 || cfg.EntropyChangeThreshold < 0 || cfg.CompositeThreshold < 0 {
		return fmt.Errorf("thresholds must be non-negative")
	}

	return nil
}

// DefaultProductionConfig returns a production-ready configuration
func DefaultProductionConfig() DetectorConfig {
	return DetectorConfig{
		BaselineSize:         1000,   // 1000 events for stable baseline
		WindowSize:           200,    // 200 events for analysis window
		HopSize:              100,    // Advance by 100 events each time
		MaxCapacity:          2000,   // Enough for baseline + 5 windows
		PValueThreshold:      0.01,   // More strict for production
		NCDThreshold:         0.4,    // Adjusted based on testing
		PermutationCount:     100,    // Balanced for performance
		Seed:                 42,     // Fixed seed for reproducibility
		CompressionAlgorithm: "zstd", // Best compression ratio
	}
}

// DefaultDemoConfig returns a configuration optimized for demo
func DefaultDemoConfig() DetectorConfig {
	return DetectorConfig{
		BaselineSize:         30,     // Small for demo
		WindowSize:           10,     // Small window
		HopSize:              5,      // Small hop
		MaxCapacity:          100,    // Small capacity
		PValueThreshold:      0.05,   // Less strict for demo
		NCDThreshold:         0.3,    // Default threshold
		PermutationCount:     50,     // Fast for demo
		Seed:                 42,     // Fixed seed
		CompressionAlgorithm: "zstd", // Fast compression
	}
}

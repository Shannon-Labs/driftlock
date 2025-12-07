package driftlockcbad

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		config    DetectorConfig
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid config",
			config: DetectorConfig{
				BaselineSize:         100,
				WindowSize:           50,
				HopSize:              25,
				MaxCapacity:          300,
				PValueThreshold:      0.05,
				NCDThreshold:         0.3,
				PermutationCount:     100,
				CompressionAlgorithm: "zstd",
			},
			expectErr: false,
		},
		{
			name: "invalid p-value too low",
			config: DetectorConfig{
				PValueThreshold: -0.1,
			},
			expectErr: true,
			errMsg:    "p_value_threshold",
		},
		{
			name: "invalid p-value too high",
			config: DetectorConfig{
				PValueThreshold: 1.5,
			},
			expectErr: true,
			errMsg:    "p_value_threshold",
		},
		{
			name: "window size larger than baseline",
			config: DetectorConfig{
				BaselineSize: 100,
				WindowSize:   200,
			},
			expectErr: true,
			errMsg:    "window_size",
		},
		{
			name: "permutation count too high",
			config: DetectorConfig{
				PermutationCount: 2000,
			},
			expectErr: true,
			errMsg:    "permutation_count",
		},
		{
			name: "unsupported compression algorithm",
			config: DetectorConfig{
				CompressionAlgorithm: "unknown",
			},
			expectErr: true,
			errMsg:    "unsupported compression algorithm",
		},
		{
			name: "max capacity too small",
			config: DetectorConfig{
				BaselineSize: 1000,
				WindowSize:   500,
				MaxCapacity:  1000,
			},
			expectErr: true,
			errMsg:    "max_capacity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}
				if tt.errMsg != "" && err.Error() == "" {
					t.Fatalf("expected error containing %q, got empty error", tt.errMsg)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDefaultProductionConfig(t *testing.T) {
	config := DefaultProductionConfig()

	if err := ValidateConfig(config); err != nil {
		t.Fatalf("DefaultProductionConfig failed validation: %v", err)
	}

	// Check specific values
	if config.PermutationCount != 100 {
		t.Errorf("Expected permutation count 100, got %d", config.PermutationCount)
	}
	if config.BaselineSize != 1000 {
		t.Errorf("Expected baseline size 1000, got %d", config.BaselineSize)
	}
}

func TestDefaultDemoConfig(t *testing.T) {
	config := DefaultDemoConfig()

	if err := ValidateConfig(config); err != nil {
		t.Fatalf("DefaultDemoConfig failed validation: %v", err)
	}

	// Check specific values
	if config.PermutationCount != 50 {
		t.Errorf("Expected permutation count 50, got %d", config.PermutationCount)
	}
	if config.BaselineSize != 30 {
		t.Errorf("Expected baseline size 30, got %d", config.BaselineSize)
	}
}
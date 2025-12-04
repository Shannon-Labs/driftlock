package plans

import "testing"

func TestNormalizePlan(t *testing.T) {
	tests := []struct {
		input      string
		wantPlan   string
		wantNormed bool
	}{
		// Canonical plans should not be normalized
		{Pulse, Pulse, false},
		{Radar, Radar, false},
		{Tensor, Tensor, false},
		{Orbit, Orbit, false},

		// Legacy free tier names
		{"trial", Pulse, true},
		{"pilot", Pulse, true},
		{"starter", Pulse, true},

		// Legacy basic tier names
		{"basic", Radar, true},
		{"signal", Radar, true},

		// Legacy pro tier names
		{"pro", Tensor, true},
		{"lock", Tensor, true},
		{"transistor", Tensor, true},
		{"sentinel", Tensor, true},
		{"growth", Tensor, true},

		// Legacy enterprise tier names
		{"enterprise", Orbit, true},

		// Unknown plans default to Pulse
		{"unknown", Pulse, true},
		{"", Pulse, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotPlan, gotNormed := NormalizePlan(tt.input)
			if gotPlan != tt.wantPlan {
				t.Errorf("NormalizePlan(%q) plan = %q, want %q", tt.input, gotPlan, tt.wantPlan)
			}
			if gotNormed != tt.wantNormed {
				t.Errorf("NormalizePlan(%q) normalized = %v, want %v", tt.input, gotNormed, tt.wantNormed)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		plan string
		want bool
	}{
		{Pulse, true},
		{Radar, true},
		{Tensor, true},
		{Orbit, true},
		{"trial", false},
		{"pro", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.plan, func(t *testing.T) {
			if got := IsValid(tt.plan); got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.plan, got, tt.want)
			}
		})
	}
}

func TestGetLimit(t *testing.T) {
	tests := []struct {
		plan string
		want int64
	}{
		{Pulse, 10_000},
		{Radar, 500_000},
		{Tensor, 5_000_000},
		{Orbit, 1_000_000_000},
		{"unknown", 10_000}, // Defaults to Pulse
	}

	for _, tt := range tests {
		t.Run(tt.plan, func(t *testing.T) {
			if got := GetLimit(tt.plan); got != tt.want {
				t.Errorf("GetLimit(%q) = %d, want %d", tt.plan, got, tt.want)
			}
		})
	}
}

func TestIsPaid(t *testing.T) {
	tests := []struct {
		plan string
		want bool
	}{
		{Pulse, false},
		{Radar, true},
		{Tensor, true},
		{Orbit, true},
	}

	for _, tt := range tests {
		t.Run(tt.plan, func(t *testing.T) {
			if got := IsPaid(tt.plan); got != tt.want {
				t.Errorf("IsPaid(%q) = %v, want %v", tt.plan, got, tt.want)
			}
		})
	}
}

func TestIsEnterprise(t *testing.T) {
	tests := []struct {
		plan string
		want bool
	}{
		{Pulse, false},
		{Radar, false},
		{Tensor, false},
		{Orbit, true},
	}

	for _, tt := range tests {
		t.Run(tt.plan, func(t *testing.T) {
			if got := IsEnterprise(tt.plan); got != tt.want {
				t.Errorf("IsEnterprise(%q) = %v, want %v", tt.plan, got, tt.want)
			}
		})
	}
}

package main

// Detection profiles provide pre-configured sensitivity presets for anomaly detection.
// Users can select a profile per-stream to balance between sensitivity and false positive rate.

// DetectionProfile represents a sensitivity preset
type DetectionProfile string

const (
	ProfileSensitive DetectionProfile = "sensitive"
	ProfileBalanced  DetectionProfile = "balanced"
	ProfileStrict    DetectionProfile = "strict"
	ProfileCustom    DetectionProfile = "custom"
)

// ProfileDefaults defines threshold defaults for each profile
type ProfileDefaults struct {
	NCDThreshold     float64
	PValueThreshold  float64
	BaselineSize     int
	WindowSize       int
	HopSize          int
	PermutationCount int
	Description      string
}

// profileDefaults maps each profile to its default settings
var profileDefaults = map[DetectionProfile]ProfileDefaults{
	ProfileSensitive: {
		NCDThreshold:     0.20, // Lower = more sensitive
		PValueThreshold:  0.10, // Higher = more sensitive
		BaselineSize:     200,  // Smaller baseline = faster warm-up
		WindowSize:       30,
		HopSize:          10,
		PermutationCount: 500,
		Description:      "Lower thresholds, more anomalies reported. Best for critical systems requiring high detection rates.",
	},
	ProfileBalanced: {
		NCDThreshold:     0.30, // Current production default
		PValueThreshold:  0.05,
		BaselineSize:     400,
		WindowSize:       50,
		HopSize:          10,
		PermutationCount: 1000,
		Description:      "Default settings balancing detection rate and false positive rate.",
	},
	ProfileStrict: {
		NCDThreshold:     0.45, // Higher = less sensitive
		PValueThreshold:  0.01, // Lower = more strict statistical test
		BaselineSize:     800,  // Larger baseline for stability
		WindowSize:       100,
		HopSize:          20,
		PermutationCount: 1000,
		Description:      "Higher thresholds, only high-confidence anomalies. Best for noisy data or when false positives are costly.",
	},
	ProfileCustom: {
		// Custom profile uses tuned values from database
		// These are fallback defaults if tuned values are NULL
		NCDThreshold:     0.30,
		PValueThreshold:  0.05,
		BaselineSize:     400,
		WindowSize:       50,
		HopSize:          10,
		PermutationCount: 1000,
		Description:      "User-customized or auto-tuned thresholds.",
	},
}

// GetProfileDefaults returns the defaults for a given profile
func GetProfileDefaults(profile DetectionProfile) ProfileDefaults {
	if defaults, ok := profileDefaults[profile]; ok {
		return defaults
	}
	return profileDefaults[ProfileBalanced]
}

// ValidProfiles returns all valid profile names
func ValidProfiles() []string {
	return []string{
		string(ProfileSensitive),
		string(ProfileBalanced),
		string(ProfileStrict),
		string(ProfileCustom),
	}
}

// IsValidProfile checks if a profile name is valid
func IsValidProfile(profile string) bool {
	switch DetectionProfile(profile) {
	case ProfileSensitive, ProfileBalanced, ProfileStrict, ProfileCustom:
		return true
	}
	return false
}

// applyProfile applies profile settings to a detection plan.
// For custom profile, uses tuned values if available; otherwise falls back to defaults.
func applyProfile(plan *detectionPlan, profile DetectionProfile, tunedNCD, tunedPValue *float64) {
	defaults := GetProfileDefaults(profile)

	if profile == ProfileCustom {
		// For custom profile, use tuned values if available
		if tunedNCD != nil {
			plan.NCDThreshold = *tunedNCD
		} else {
			plan.NCDThreshold = defaults.NCDThreshold
		}
		if tunedPValue != nil {
			plan.PValueThreshold = *tunedPValue
		} else {
			plan.PValueThreshold = defaults.PValueThreshold
		}
	} else {
		// For preset profiles, use profile defaults
		plan.NCDThreshold = defaults.NCDThreshold
		plan.PValueThreshold = defaults.PValueThreshold
	}

	// Apply size defaults only if not already set by stream config
	if plan.BaselineSize == 0 {
		plan.BaselineSize = defaults.BaselineSize
	}
	if plan.WindowSize == 0 {
		plan.WindowSize = defaults.WindowSize
	}
	if plan.HopSize == 0 {
		plan.HopSize = defaults.HopSize
	}
	if plan.PermutationCount == 0 {
		plan.PermutationCount = defaults.PermutationCount
	}
}

// ProfileSummary returns a description of a profile for API responses
type ProfileSummary struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	NCDThreshold float64 `json:"ncd_threshold"`
	PValueThreshold float64 `json:"pvalue_threshold"`
	BaselineSize int     `json:"baseline_size"`
	WindowSize   int     `json:"window_size"`
}

// GetProfileSummaries returns summaries of all profiles for API documentation
func GetProfileSummaries() map[string]ProfileSummary {
	summaries := make(map[string]ProfileSummary)
	for profile, defaults := range profileDefaults {
		summaries[string(profile)] = ProfileSummary{
			Name:            string(profile),
			Description:     defaults.Description,
			NCDThreshold:    defaults.NCDThreshold,
			PValueThreshold: defaults.PValueThreshold,
			BaselineSize:    defaults.BaselineSize,
			WindowSize:      defaults.WindowSize,
		}
	}
	return summaries
}

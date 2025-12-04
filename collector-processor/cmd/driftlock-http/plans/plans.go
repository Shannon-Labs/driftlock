// Package plans defines the canonical plan tier names and their configurations.
// This is the single source of truth for plan names across the Driftlock codebase.
package plans

// Plan names - these are the canonical plan tier identifiers
const (
	// Pulse is the free tier: $0/mo, 10,000 events/mo, basic detection
	Pulse = "pulse"

	// Radar is the basic paid tier: $15/mo, 500,000 events/mo, email alerts
	Radar = "radar"

	// Tensor is the pro tier: $100/mo, 5,000,000 events/mo, DORA/NIS2 compliance
	Tensor = "tensor"

	// Orbit is the enterprise tier: custom pricing, unlimited events, dedicated support
	Orbit = "orbit"
)

// Legacy plan name mappings - used for backward compatibility during migration
// These are deprecated and will be removed in a future version
var LegacyNames = map[string]string{
	// Old free tier names
	"trial":   Pulse,
	"pilot":   Pulse,
	"starter": Pulse,

	// Old basic tier names
	"basic":  Radar,
	"signal": Radar,

	// Old pro tier names
	"pro":        Tensor,
	"lock":       Tensor,
	"transistor": Tensor,
	"sentinel":   Tensor,
	"growth":     Tensor,

	// Old enterprise tier names
	"enterprise": Orbit,
}

// PlanLimits defines monthly event limits for each plan tier
var PlanLimits = map[string]int64{
	Pulse:  10_000,
	Radar:  500_000,
	Tensor: 5_000_000,
	Orbit:  1_000_000_000, // Effectively unlimited
}

// ValidPlans contains all valid plan names
var ValidPlans = map[string]bool{
	Pulse:  true,
	Radar:  true,
	Tensor: true,
	Orbit:  true,
}

// NormalizePlan converts a plan name to its canonical form.
// Returns the canonical plan name and whether normalization was needed.
// If the plan is already canonical, returns the same name.
// If the plan is a legacy name, returns the mapped canonical name.
// If the plan is unknown, returns Pulse as the default.
func NormalizePlan(plan string) (canonical string, normalized bool) {
	// Already canonical
	if ValidPlans[plan] {
		return plan, false
	}

	// Check legacy mapping
	if canonical, ok := LegacyNames[plan]; ok {
		return canonical, true
	}

	// Unknown plan - default to free tier
	return Pulse, true
}

// IsValid returns true if the plan name is a valid canonical plan
func IsValid(plan string) bool {
	return ValidPlans[plan]
}

// GetLimit returns the monthly event limit for a plan.
// Returns the Pulse limit for unknown plans.
func GetLimit(plan string) int64 {
	if limit, ok := PlanLimits[plan]; ok {
		return limit
	}
	return PlanLimits[Pulse]
}

// IsPaid returns true if the plan is a paid tier (not Pulse)
func IsPaid(plan string) bool {
	return plan == Radar || plan == Tensor || plan == Orbit
}

// IsEnterprise returns true if the plan is the enterprise tier
func IsEnterprise(plan string) bool {
	return plan == Orbit
}

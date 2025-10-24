package driftlockcbad

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines the configuration for the driftlock CBAD processor
type Config struct {
	// Embed the component.Config interface to satisfy the component configuration contract
	component.Config `mapstructure:",squash"`

	WindowSize  int     `mapstructure:"window_size"`
	HopSize     int     `mapstructure:"hop_size"`
	Threshold   float64 `mapstructure:"threshold"`
	Determinism bool    `mapstructure:"determinism"`
}

package driftlockcbad

import (
    "go.opentelemetry.io/collector/config"
)

type Config struct {
    config.ProcessorSettings `mapstructure:",squash"`

    WindowSize  int     `mapstructure:"window_size"`
    HopSize     int     `mapstructure:"hop_size"`
    Threshold   float64 `mapstructure:"threshold"`
    Determinism bool    `mapstructure:"determinism"`
}


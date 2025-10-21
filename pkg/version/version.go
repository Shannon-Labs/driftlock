package version

import "os"

var v string

// Version returns a version string from env or a default.
func Version() string {
    if v != "" {
        return v
    }
    if ev := os.Getenv("DRIFTLOCK_VERSION"); ev != "" {
        return ev
    }
    return "dev"
}


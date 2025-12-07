package driftlockcbad

import "testing"

func TestRustCBADLibraryAvailable(t *testing.T) {
	if !IsAvailable() {
		t.Fatalf("CBAD library not available: %v", AvailabilityError())
	}
	t.Log("CBAD library successfully loaded")
}

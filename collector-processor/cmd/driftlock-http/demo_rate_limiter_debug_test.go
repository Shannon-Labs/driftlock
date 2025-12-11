package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDemoRateLimiterLeavesStaleIPs(t *testing.T) {
	limiter := newDemoRateLimiterWithWindow(50 * time.Millisecond)
	defer limiter.Close()
	limiter.limit = 2

	for i := 0; i < 20; i++ {
		ip := fmt.Sprintf("10.0.0.%d", i)
		if !limiter.allow(ip) {
			t.Fatalf("unexpected throttle for %s", ip)
		}
	}

	// Allow window to elapse so entries should be considered stale.
	time.Sleep(120 * time.Millisecond)

	limiter.mu.Lock()
	registrySize := len(limiter.requests)
	limiter.mu.Unlock()

	if registrySize >= 10 {
		t.Fatalf("stale demo limiter entries were not cleaned (size=%d)", registrySize)
	}
}

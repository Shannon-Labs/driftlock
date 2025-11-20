package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTenantRateLimiter_TenantLimit(t *testing.T) {
	rl := newTenantRateLimiter(1)
	tc := tenantContext{
		Tenant: tenantRecord{ID: uuid.New(), RateLimitRPS: 1},
		Key:    apiKeyRecord{},
	}

	first := rl.Allow(tc)
	if !first.Allowed {
		t.Fatalf("expected first request to be allowed")
	}

	second := rl.Allow(tc)
	if second.Allowed {
		t.Fatalf("expected second request to be rejected")
	}
	if second.RetryAfter <= 0 {
		t.Fatalf("expected RetryAfter to be > 0, got %v", second.RetryAfter)
	}
}

func TestTenantRateLimiter_KeyLimit(t *testing.T) {
	rl := newTenantRateLimiter(100)
	tc := tenantContext{
		Tenant: tenantRecord{ID: uuid.New(), RateLimitRPS: 100},
		Key:    apiKeyRecord{ID: uuid.New(), KeyRateLimit: 1},
	}

	first := rl.Allow(tc)
	if !first.Allowed {
		t.Fatalf("expected first request to be allowed")
	}

	second := rl.Allow(tc)
	if second.Allowed {
		t.Fatalf("expected second request to be rate limited")
	}
	if second.RetryAfter <= 0 {
		t.Fatalf("expected RetryAfter to be set")
	}
}

func TestDecorateRateLimitHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	decision := rateLimitDecision{
		Allowed:         false,
		RetryAfter:      1500 * time.Millisecond,
		TenantLimit:     10,
		TenantRemaining: 5,
		KeyLimit:        4,
		KeyRemaining:    2,
	}

	decorateRateLimitHeaders(w, decision)

	if got := w.Header().Get("X-RateLimit-Limit"); got != "10" {
		t.Fatalf("expected X-RateLimit-Limit=10, got %s", got)
	}
	if got := w.Header().Get("X-RateLimit-Remaining"); got != "2" {
		t.Fatalf("expected X-RateLimit-Remaining=2, got %s", got)
	}
	if got := w.Header().Get("Retry-After"); got == "" {
		t.Fatalf("expected Retry-After to be set")
	}
	if got := w.Header().Get("X-RateLimit-Reset"); got == "" {
		t.Fatalf("expected X-RateLimit-Reset to be set")
	}
}

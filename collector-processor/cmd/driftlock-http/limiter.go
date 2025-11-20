package main

import (
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
)

type tenantRateLimiter struct {
	defaultLimit int
	mu           sync.Mutex
	tenant       map[uuid.UUID]*tokenBucket
	keys         map[uuid.UUID]*tokenBucket
}

type rateLimitDecision struct {
	Allowed         bool
	RetryAfter      time.Duration
	TenantLimit     int
	TenantRemaining int
	KeyLimit        int
	KeyRemaining    int
}

func newTenantRateLimiter(defaultLimit int) *tenantRateLimiter {
	if defaultLimit <= 0 {
		defaultLimit = 60
	}
	return &tenantRateLimiter{
		defaultLimit: defaultLimit,
		tenant:       make(map[uuid.UUID]*tokenBucket),
		keys:         make(map[uuid.UUID]*tokenBucket),
	}
}

func (rl *tenantRateLimiter) Allow(tc tenantContext) rateLimitDecision {
	decision := rateLimitDecision{}
	if tc.Tenant.ID == uuid.Nil {
		decision.Allowed = true
		return decision
	}
	now := time.Now()
	tenantLimit := tc.Tenant.RateLimitRPS
	if tenantLimit <= 0 {
		tenantLimit = rl.defaultLimit
	}
	decision.TenantLimit = tenantLimit

	tenantBucket := rl.getBucket(rl.tenant, tc.Tenant.ID, tenantLimit)
	allowedTenant, retryTenant, tenantTokens := tenantBucket.allow(now)
	decision.TenantRemaining = tokensToRemaining(tenantTokens)
	if !allowedTenant {
		decision.RetryAfter = retryTenant
		return decision
	}

	if tc.Key.KeyRateLimit > 0 {
		keyBucket := rl.getBucket(rl.keys, tc.Key.ID, tc.Key.KeyRateLimit)
		allowedKey, retryKey, keyTokens := keyBucket.allow(now)
		decision.KeyLimit = tc.Key.KeyRateLimit
		decision.KeyRemaining = tokensToRemaining(keyTokens)
		decision.TenantRemaining = tokensToRemaining(tenantBucket.remaining(now))
		if !allowedKey {
			tenantBucket.refund()
			decision.TenantRemaining = tokensToRemaining(tenantBucket.remaining(now))
			decision.RetryAfter = retryKey
			return decision
		}
	} else {
		decision.KeyRemaining = decision.TenantRemaining
	}

	decision.Allowed = true
	return decision
}

func (rl *tenantRateLimiter) getBucket(registry map[uuid.UUID]*tokenBucket, id uuid.UUID, limit int) *tokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	bucket, ok := registry[id]
	if !ok {
		bucket = newTokenBucket(limit)
		registry[id] = bucket
		return bucket
	}
	bucket.updateLimit(limit)
	return bucket
}

type tokenBucket struct {
	mu        sync.Mutex
	rate      float64
	capacity  float64
	tokens    float64
	lastFill  time.Time
	allowance float64
}

func newTokenBucket(limit int) *tokenBucket {
	if limit <= 0 {
		limit = 1
	}
	now := time.Now()
	return &tokenBucket{
		rate:      float64(limit),
		capacity:  float64(limit),
		tokens:    float64(limit),
		lastFill:  now,
		allowance: 1.0,
	}
}

func (b *tokenBucket) updateLimit(limit int) {
	if limit <= 0 {
		limit = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill(time.Now())
	b.rate = float64(limit)
	b.capacity = float64(limit)
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
}

func (b *tokenBucket) allow(now time.Time) (bool, time.Duration, float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill(now)
	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true, 0, b.tokens
	}
	if b.rate == 0 {
		return false, time.Second, b.tokens
	}
	needed := 1.0 - b.tokens
	wait := time.Duration((needed / b.rate) * float64(time.Second))
	return false, wait, b.tokens
}

func (b *tokenBucket) refund() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens += 1.0
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
}

func (b *tokenBucket) remaining(now time.Time) float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill(now)
	return b.tokens
}

func (b *tokenBucket) refill(now time.Time) {
	if b.lastFill.IsZero() {
		b.lastFill = now
		return
	}
	elapsed := now.Sub(b.lastFill).Seconds()
	if elapsed <= 0 {
		return
	}
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.lastFill = now
}

func tokensToRemaining(v float64) int {
	if v <= 0 {
		return 0
	}
	return int(math.Floor(v))
}

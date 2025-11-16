package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ctxTenantKey struct{}

var tenantContextKey ctxTenantKey

func withAuth(s *store, limiter *tenantRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(r.Header.Get("X-Api-Key"))
		if raw == "" {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("missing api key"))
			return
		}
		keyID, err := parseAPIKey(raw)
		if err != nil {
			writeError(w, r, http.StatusUnauthorized, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		rec, err := s.resolveAPIKey(ctx, keyID)
		if err != nil || !verifyAPIKey(rec.KeyHash, raw) {
			writeError(w, r, http.StatusUnauthorized, fmt.Errorf("invalid api key"))
			return
		}
		tenant, ok := s.tenantByID(rec.TenantID)
		if !ok {
			// cache miss, reload once
			if err := s.loadCache(ctx); err != nil {
				writeError(w, r, http.StatusInternalServerError, fmt.Errorf("load config: %w", err))
				return
			}
			tenant, ok = s.tenantByID(rec.TenantID)
			if !ok {
				writeError(w, r, http.StatusForbidden, fmt.Errorf("tenant not found"))
				return
			}
		}
		tc := tenantContext{Tenant: tenant, Key: *rec}
		if limiter != nil {
			if allowed, retry := limiter.Allow(tc); !allowed {
				if retry > 0 {
					w.Header().Set("Retry-After", fmt.Sprintf("%.0f", math.Ceil(retry.Seconds())))
				}
				payload := apiErrorPayload{Error: apiError{
					Code:              "rate_limit_exceeded",
					Message:           "per-tenant rate limit exceeded",
					RequestID:         requestIDFrom(r.Context()),
					RetryAfterSeconds: int(math.Ceil(retry.Seconds())),
				}}
				writeJSON(w, r, http.StatusTooManyRequests, payload)
				return
			}
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), tenantContextKey, tc)))
	})
}

func parseAPIKey(raw string) (uuid.UUID, error) {
	if !strings.HasPrefix(raw, "dlk_") {
		return uuid.Nil, errors.New("invalid api key format")
	}
	rest := strings.TrimPrefix(raw, "dlk_")
	parts := strings.SplitN(rest, ".", 2)
	if len(parts) != 2 {
		return uuid.Nil, errors.New("invalid api key format")
	}
	id, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, err
	}
	if len(parts[1]) < 8 {
		return uuid.Nil, errors.New("invalid api key secret")
	}
	return id, nil
}

func tenantFromContext(ctx context.Context) (tenantContext, bool) {
	v := ctx.Value(tenantContextKey)
	if v == nil {
		return tenantContext{}, false
	}
	tc, ok := v.(tenantContext)
	return tc, ok
}

package ctxutil

import "context"

type ctxKey string

const (
    ctxOrgIDKey     ctxKey = "driftlock_org_id"
    ctxEventTypeKey ctxKey = "driftlock_event_type"
)

// WithEventContext adds organization ID and event type to context
func WithEventContext(ctx context.Context, orgID string, eventType string) context.Context {
    if orgID != "" {
        ctx = context.WithValue(ctx, ctxOrgIDKey, orgID)
    }
    if eventType != "" {
        ctx = context.WithValue(ctx, ctxEventTypeKey, eventType)
    }
    return ctx
}

// GetOrganizationID returns org ID from context if present
func GetOrganizationID(ctx context.Context) string {
    if v, ok := ctx.Value(ctxOrgIDKey).(string); ok {
        return v
    }
    return ""
}

// GetEventType returns event type from context if present
func GetEventType(ctx context.Context) string {
    if v, ok := ctx.Value(ctxEventTypeKey).(string); ok {
        return v
    }
    return ""
}


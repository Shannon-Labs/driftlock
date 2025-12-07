package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BillingAccessError represents an access denial due to billing status
type BillingAccessError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *BillingAccessError) Error() string {
	return e.Message
}

// checkBillingAccess verifies the tenant's billing status allows API access.
// Returns nil if access is allowed, or a BillingAccessError if blocked.
//
// Access policy:
//   - free, trialing, active: Always allowed
//   - grace_period: Allowed with warning header (soft enforcement)
//   - past_due: Blocked - payment failed and needs resolution
//   - canceled, expired: Blocked - subscription ended
func checkBillingAccess(ctx context.Context, store *store, tenantID uuid.UUID) (*BillingStatus, error) {
	bs, err := store.getBillingStatus(ctx, tenantID)
	if err != nil {
		// Fail open for availability - don't block on billing lookup errors
		// Log the error but allow the request to proceed
		Logger().Warn("Failed to check billing status, allowing request",
			zap.String("tenant_id", tenantID.String()),
			zap.Error(err))
		return nil, nil
	}

	switch bs.Status {
	case "free", "trialing", "active":
		// Full access
		return bs, nil

	case "grace_period":
		// Soft enforcement: allow but return status for warning headers
		return bs, nil

	case "past_due":
		return bs, &BillingAccessError{
			Status:  bs.Status,
			Message: "Payment failed. Please update your payment method to continue using the API.",
			Code:    "payment_required",
		}

	case "canceled", "expired":
		return bs, &BillingAccessError{
			Status:  bs.Status,
			Message: "Your subscription has ended. Please reactivate to continue using the API.",
			Code:    "subscription_inactive",
		}

	default:
		// Unknown status - fail open
		Logger().Warn("Unknown billing status, allowing request",
			zap.String("tenant_id", tenantID.String()),
			zap.String("status", bs.Status))
		return bs, nil
	}
}

// writeBillingError writes a billing access denial response
func writeBillingError(w http.ResponseWriter, r *http.Request, billingErr *BillingAccessError) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Billing-Status", billingErr.Status)

	// Use 402 Payment Required for billing issues
	w.WriteHeader(http.StatusPaymentRequired)

	resp := struct {
		Error   string `json:"error"`
		Code    string `json:"code"`
		Status  string `json:"billing_status"`
		HelpURL string `json:"help_url,omitempty"`
	}{
		Error:   billingErr.Message,
		Code:    billingErr.Code,
		Status:  billingErr.Status,
		HelpURL: "https://driftlock.io/billing",
	}

	json.NewEncoder(w).Encode(resp)
}

// decorateBillingHeaders adds billing status headers for warning/info purposes
func decorateBillingHeaders(w http.ResponseWriter, bs *BillingStatus) {
	if bs == nil {
		return
	}

	w.Header().Set("X-Billing-Status", bs.Status)
	w.Header().Set("X-Billing-Plan", bs.Plan)

	// Add grace period warning
	if bs.Status == "grace_period" && bs.GracePeriodEndsAt != nil {
		w.Header().Set("X-Grace-Period-Ends", bs.GracePeriodEndsAt.Format("2006-01-02"))
	}

	// Add trial info
	if bs.Status == "trialing" && bs.TrialDaysRemaining != nil {
		w.Header().Set("X-Trial-Days-Remaining", strconv.Itoa(*bs.TrialDaysRemaining))
	}
}

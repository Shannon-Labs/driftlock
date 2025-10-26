package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/Hmbown/driftlock/api-server/internal/models"
	"github.com/Hmbown/driftlock/api-server/internal/services"
)

// BillingService orchestrates billing operations between internal models and Stripe
type BillingService struct {
	stripeSvc   *StripeService
	tenantSvc   *services.TenantService
}

// NewBillingService creates a new billing service instance
func NewBillingService(stripeSvc *StripeService, tenantSvc *services.TenantService) *BillingService {
	return &BillingService{
		stripeSvc:   stripeSvc,
		tenantSvc:   tenantSvc,
	}
}

// SubscribeTenant subscribes a tenant to a plan
func (b *BillingService) SubscribeTenant(ctx context.Context, tenantID, planID string) error {
	// Get tenant details
	tenant, err := b.tenantSvc.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	// Create customer in Stripe if not exists (assuming customer ID is stored in tenant's metadata)
	// In a real implementation, you'd store the Stripe customer ID in the tenant model
	// For now, we'll create a new customer
	customer, err := b.stripeSvc.CreateCustomer(ctx, tenant)
	if err != nil {
		return fmt.Errorf("failed to create customer in Stripe: %w", err)
	}

	// Create subscription in Stripe
	_, err = b.stripeSvc.CreateSubscription(ctx, customer.ID, planID)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Update tenant status to active/active plan
	updatedTenant := &models.Tenant{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Domain:    tenant.Domain,
		Email:     tenant.Email,
		Status:    models.TenantStatusActive, // Change to active since subscribed
		Plan:      getTenantPlanFromStripePlan(planID), // Convert Stripe plan to our internal plan
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: time.Now(),
		ExpiresAt: tenant.ExpiresAt,
		Usage:     tenant.Usage,
		Quotas:    tenant.Quotas,
	}

	_, err = b.tenantSvc.UpdateTenant(ctx, updatedTenant)
	if err != nil {
		return fmt.Errorf("failed to update tenant after subscription: %w", err)
	}

	return nil
}

// CreateCheckoutSession creates a Stripe checkout session for a tenant
func (b *BillingService) CreateCheckoutSession(ctx context.Context, tenantID, planID string) (*stripe.CheckoutSession, error) {
	// Get tenant details
	tenant, err := b.tenantSvc.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Create checkout session
	session, err := b.stripeSvc.CreateCheckoutSession(ctx, tenant, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return session, nil
}

// CancelSubscription cancels a tenant's subscription
func (b *BillingService) CancelSubscription(ctx context.Context, tenantID string) error {
	// In a real implementation, you would look up the Stripe subscription ID
	// associated with this tenant and cancel it
	// For now, we'll just update the tenant status
	
	tenant, err := b.tenantSvc.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	updatedTenant := &models.Tenant{
		ID:        tenant.ID,
		Name:      tenant.Name,
		Domain:    tenant.Domain,
		Email:     tenant.Email,
		Status:    models.TenantStatusSuspended, // Set to suspended
		Plan:      tenant.Plan,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: time.Now(),
		ExpiresAt: tenant.ExpiresAt,
		Usage:     tenant.Usage,
		Quotas:    tenant.Quotas,
	}

	_, err = b.tenantSvc.UpdateTenant(ctx, updatedTenant)
	if err != nil {
		return fmt.Errorf("failed to update tenant after subscription cancellation: %w", err)
	}

	return nil
}

// ProcessUsage updates tenant usage and creates usage records in Stripe
func (b *BillingService) ProcessUsage(ctx context.Context, tenantID string, eventType string, value int) error {
	// Update tenant usage in the database
	err := b.tenantSvc.UpdateTenantUsage(ctx, tenantID, eventType, value)
	if err != nil {
		return fmt.Errorf("failed to update tenant usage: %w", err)
	}

	// In a real implementation, you would create usage records in Stripe here
	// This would involve:
	// 1. Looking up the subscription and subscription item ID for this tenant
	// 2. Creating usage records for metered billing
	
	return nil
}

// getTenantPlanFromStripePlan converts Stripe plan ID to internal tenant plan
func getTenantPlanFromStripePlan(stripePlanID string) string {
	switch stripePlanID {
	case "price_trial":
		return models.TenantPlanTrial
	case "price_starter":
		return models.TenantPlanStarter
	case "price_pro":
		return models.TenantPlanPro
	case "price_enterprise":
		return models.TenantPlanEnterprise
	default:
		return models.TenantPlanTrial // Default to trial
	}
}
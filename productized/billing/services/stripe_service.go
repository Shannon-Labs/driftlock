package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	stripe "github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/client"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/sub"
	"gorm.io/gorm"

	"driftlock/productized/api/database"
	"driftlock/productized/billing/models"
	"driftlock/productized/billing/config"
)

type StripeService struct {
	SC *client.API
	Config *config.StripeConfig
}

func NewStripeService(stripeConfig *config.StripeConfig) *StripeService {
	sc := &client.API{}
	sc.Init(stripeConfig.SecretKey, nil)

	return &StripeService{
		SC:     sc,
		Config: stripeConfig,
	}
}

func (s *StripeService) CreateCustomer(email, name string, userID uint) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		// Store our internal user ID in metadata
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	cust, err := customer.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe customer: %v", err)
	}

	return cust, nil
}

func (s *StripeService) CreateCheckoutSession(customerID, planID string, userID uint) (*stripe.CheckoutSession, error) {
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL:         stripe.String(s.Config.SuccessURL),
		CancelURL:          stripe.String(s.Config.CancelURL),
		AllowPromotionCodes: stripe.Bool(true),
		// Store our internal user ID in metadata
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	session, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %v", err)
	}

	return session, nil
}

func (s *StripeService) CreateCustomerPortalSession(customerID string) (*stripe.BillingPortalSession, error) {
	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(s.Config.SuccessURL), // URL to return customer to
	}

	portalSession, err := s.SC.BillingPortalSession.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing portal session: %v", err)
	}

	return portalSession, nil
}

func (s *StripeService) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{}
	sub, err := sub.Get(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %v", err)
	}

	return sub, nil
}

func (s *StripeService) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	sub, err := sub.Update(subscriptionID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription: %v", err)
	}

	// Update our database
	db := database.DB
	var userSub models.UserSubscription
	if err := db.Where("subscription_id = ?", subscriptionID).First(&userSub).Error; err != nil {
		return nil, fmt.Errorf("failed to find subscription in database: %v", err)
	}

	userSub.Status = "canceled"
	if err := db.Save(&userSub).Error; err != nil {
		return nil, fmt.Errorf("failed to update subscription status in database: %v", err)
	}

	return sub, nil
}

func (s *StripeService) GetCurrentPlan(userID uint) (*models.SubscriptionPlan, error) {
	db := database.DB
	
	var userSub models.UserSubscription
	if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&userSub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User has no active subscription, return free plan
			return &models.SubscriptionPlan{
				ID:          "free",
				Name:        "Free Plan",
				Description: "Basic plan with limited features",
				Price:       0,
				Currency:    s.Config.Currency,
				Interval:    "month",
			}, nil
		}
		return nil, fmt.Errorf("failed to get user subscription: %v", err)
	}

	var plan models.SubscriptionPlan
	if err := db.Where("id = ?", userSub.PlanID).First(&plan).Error; err != nil {
		return nil, fmt.Errorf("failed to get plan details: %v", err)
	}

	return &plan, nil
}

func (s *StripeService) RecordUsage(userID uint, featureName string, quantity int64) error {
	db := database.DB
	
	// Get the user's subscription
	var userSub models.UserSubscription
	if err := db.Where("user_id = ? AND status = ?", userID, "active").First(&userSub).Error; err != nil {
		return fmt.Errorf("user has no active subscription: %v", err)
	}

	// Check if the feature is available in the plan
	var feature models.Feature
	if err := db.Joins("JOIN subscription_plans ON features.plan_id = subscription_plans.id").
		Where("subscription_plans.id = ? AND features.name = ?", userSub.PlanID, featureName).First(&feature).Error; err != nil {
		return fmt.Errorf("feature %s not available in plan: %v", featureName, err)
	}

	// Check quota if it's not unlimited (-1)
	if feature.Quota != -1 {
		// Calculate usage for the current period
		startPeriod := userSub.CurrentPeriodStart
		endPeriod := userSub.CurrentPeriodEnd
		
		var totalUsage int64
		if err := db.Model(&models.UsageRecord{}).
			Where("user_id = ? AND feature_name = ? AND timestamp BETWEEN ? AND ?",
				userID, featureName, startPeriod, endPeriod).
			Count(&totalUsage).Error; err != nil {
			return fmt.Errorf("failed to calculate usage: %v", err)
		}

		if totalUsage + quantity > int64(feature.Quota) {
			return fmt.Errorf("usage limit exceeded for feature %s: %d/%d", featureName, totalUsage + quantity, feature.Quota)
		}
	}

	// Record the usage
	usageRecord := models.UsageRecord{
		UserID:         userID,
		SubscriptionID: userSub.SubscriptionID,
		FeatureName:    featureName,
		Quantity:       quantity,
		Timestamp:      time.Now(),
	}

	if err := db.Create(&usageRecord).Error; err != nil {
		return fmt.Errorf("failed to record usage: %v", err)
	}

	// If this is a Stripe metered usage, report to Stripe as well
	if s.Config.UsageEnabled {
		if err := s.reportUsageToStripe(userSub.SubscriptionID, featureName, quantity); err != nil {
			log.Printf("Failed to report usage to Stripe: %v", err)
			// Don't fail the operation if Stripe reporting fails
		}
	}

	return nil
}

func (s *StripeService) reportUsageToStripe(subscriptionID, featureName string, quantity int64) error {
	// This would be implemented with Stripe's usage record API
	// For now, we'll add a placeholder
	log.Printf("Would report usage to Stripe: sub=%s, feature=%s, quantity=%d", subscriptionID, featureName, quantity)
	return nil
}

func (s *StripeService) GetCustomerByID(customerID string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{}
	cust, err := customer.Get(customerID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %v", err)
	}

	return cust, nil
}
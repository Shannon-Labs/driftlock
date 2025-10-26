package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"driftlock/productized/api/models"
	"driftlock/productized/api/database"
	"driftlock/productized/billing/config"
	"driftlock/productized/billing/services"
	billingModels "driftlock/productized/billing/models"
)

var stripeService *services.StripeService

// Initialize the Stripe service with configuration
func InitBillingService(stripeConfig *config.StripeConfig) {
	stripeService = services.NewStripeService(stripeConfig)
}

func CreateCheckoutSession(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var req struct {
		PlanID string `json:"plan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create Stripe customer
	var userSub billingModels.UserSubscription
	err := database.DB.Where("user_id = ?", currentUser.ID).First(&userSub).Error
	var customerID string

	if err != nil {
		// Create new Stripe customer
		customer, err := stripeService.CreateCustomer(currentUser.Email, currentUser.Name, currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
			return
		}
		customerID = customer.ID
	} else {
		customerID = userSub.CustomerID
	}

	// Create checkout session
	session, err := stripeService.CreateCheckoutSession(customerID, req.PlanID, currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.ID,
		"url":        session.URL,
	})
}

func GetCustomerPortal(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var userSub billingModels.UserSubscription
	err := database.DB.Where("user_id = ?", currentUser.ID).First(&userSub).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User has no subscription"})
		return
	}

	portalSession, err := stripeService.CreateCustomerPortalSession(userSub.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer portal session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": portalSession.URL,
	})
}

func GetSubscription(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	plan, err := stripeService.GetCurrentPlan(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscription"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func CancelSubscription(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var userSub billingModels.UserSubscription
	err := database.DB.Where("user_id = ? AND status = ?", currentUser.ID, "active").First(&userSub).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User has no active subscription"})
		return
	}

	_, err = stripeService.CancelSubscription(userSub.SubscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subscription scheduled for cancellation at period end",
	})
}

func GetBillingPlans(c *gin.Context) {
	// Define available plans
	plans := []billingModels.SubscriptionPlan{
		{
			ID:          "price_free",
			Name:        "Free Plan",
			Description: "Basic plan with limited features",
			Price:       0,
			Currency:    "usd",
			Interval:    "month",
			Features: []billingModels.Feature{
				{Name: "Monthly Events", Quota: 10000, Available: true},
				{Name: "Retention Days", Quota: 7, Available: true},
				{Name: "Alerts", Quota: 100, Available: true},
				{Name: "API Access", Available: true},
				{Name: "Basic Support", Available: true},
			},
		},
		{
			ID:          "price_pro",
			Name:        "Pro Plan",
			Description: "Professional plan with advanced features",
			Price:       2900, // $29 in cents
			Currency:    "usd",
			Interval:    "month",
			Features: []billingModels.Feature{
				{Name: "Monthly Events", Quota: 100000, Available: true},
				{Name: "Retention Days", Quota: 30, Available: true},
				{Name: "Alerts", Quota: 1000, Available: true},
				{Name: "API Access", Available: true},
				{Name: "Advanced Support", Available: true},
				{Name: "Custom Dashboards", Available: true},
				{Name: "Anomaly Prediction", Available: true},
			},
		},
		{
			ID:          "price_business",
			Name:        "Business Plan",
			Description: "Enterprise plan with premium features",
			Price:       9900, // $99 in cents
			Currency:    "usd",
			Interval:    "month",
			Features: []billingModels.Feature{
				{Name: "Monthly Events", Quota: -1, Available: true}, // Unlimited
				{Name: "Retention Days", Quota: 90, Available: true},
				{Name: "Alerts", Quota: -1, Available: true}, // Unlimited
				{Name: "API Access", Available: true},
				{Name: "Priority Support", Available: true},
				{Name: "Custom Dashboards", Available: true},
				{Name: "Anomaly Prediction", Available: true},
				{Name: "SSO Integration", Available: true},
				{Name: "Advanced Analytics", Available: true},
				{Name: "Dedicated Account Manager", Available: true},
			},
		},
	}

	c.JSON(http.StatusOK, plans)
}

func GetUsage(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	// Get query params for date range
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// For now, return sample usage data
	// In a real implementation, this would query the database
	usage := map[string]interface{}{
		"user_id": currentUser.ID,
		"period": map[string]string{
			"start": startDate,
			"end":   endDate,
		},
		"features": []map[string]interface{}{
			{
				"name":      "events_ingested",
				"current":   1500,
				"quota":     10000,
				"remaining": 8500,
			},
			{
				"name":      "alerts_created",
				"current":   45,
				"quota":     100,
				"remaining": 55,
			},
		},
	}

	c.JSON(http.StatusOK, usage)
}

func RecordUsage(c *gin.Context) {
	// This would typically be called from internal services, not from frontend
	// For security, this endpoint should be restricted to internal use
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var req struct {
		FeatureName string `json:"feature_name" binding:"required"`
		Quantity    int64  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := stripeService.RecordUsage(currentUser.ID, req.FeatureName, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usage recorded successfully",
	})
}

// Webhook handler for Stripe events
func HandleStripeWebhook(c *gin.Context) {
	// In a real implementation, this would verify the webhook signature
	// and handle different Stripe events like:
	// - invoice.payment_succeeded
	// - customer.subscription.created
	// - customer.subscription.updated
	// - customer.subscription.deleted

	// For now, we'll just acknowledge the webhook
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}
package routes

import (
	"github.com/gin-gonic/gin"

	"driftlock/productized/api/config"
	"driftlock/productized/api/handlers"
	"driftlock/productized/api/middleware"
	billingConfig "driftlock/productized/billing/config"
	billingHandlers "driftlock/productized/billing/handlers"
	emailConfig "driftlock/productized/email/config"
	emailHandlers "driftlock/productized/email/handlers"
	analyticsMiddleware "driftlock/productized/analytics"
	onboarding "driftlock/productized/api/handlers"
	orgMiddleware "driftlock/productized/api/middleware"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize billing service
	stripeConfig := billingConfig.LoadStripeConfig()
	billingHandlers.InitBillingService(stripeConfig)

	// Initialize email service
	emailConfig := emailConfig.LoadEmailConfig()
	emailHandlers.InitEmailService(emailConfig)

	// Initialize analytics and audit services
	analyticsService, auditService := analyticsMiddleware.InitializeAnalytics()

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/refresh", handlers.RefreshToken)
	}

	// Protected routes with analytics and audit middleware
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	protected.Use(orgMiddleware.OrganizationMiddleware())  // Organization context from API gateway
	protected.Use(analyticsMiddleware.AnalyticsMiddleware(analyticsService))
	{
		// User routes
		protected.GET("/user", handlers.GetProfile)
		protected.PUT("/user", handlers.UpdateProfile)

		// Anomaly routes
		anomalies := protected.Group("/anomalies")
		{
			anomalies.GET("", handlers.GetAnomalies)
			anomalies.GET("/:id", handlers.GetAnomaly)
			anomalies.PUT("/:id/resolve", handlers.ResolveAnomaly)
			anomalies.DELETE("/:id", handlers.DeleteAnomaly)
		}

		// Event routes
		events := protected.Group("/events")
		{
			events.POST("/ingest", handlers.IngestEvent)
			events.GET("", handlers.GetEvents)
		}

		// Dashboard routes
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/stats", handlers.GetDashboardStats)
			dashboard.GET("/recent", handlers.GetRecentAnomalies)
		}

		// Billing routes
		billing := protected.Group("/billing")
		{
			billing.GET("/plans", billingHandlers.GetBillingPlans)
			billing.POST("/checkout", billingHandlers.CreateCheckoutSession)
			billing.GET("/portal", billingHandlers.GetCustomerPortal)
			billing.GET("/subscription", billingHandlers.GetSubscription)
			billing.DELETE("/subscription", billingHandlers.CancelSubscription)
			billing.GET("/usage", billingHandlers.GetUsage)
			billing.POST("/usage", billingHandlers.RecordUsage)
		}

		// Email routes
		email := protected.Group("/email")
		{
			email.POST("/test", emailHandlers.SendTestEmail)
			email.POST("/welcome", emailHandlers.SendWelcomeEmail)
			email.POST("/anomaly-alert", emailHandlers.SendAnomalyAlert)
		}

		// Onboarding routes
		onboarding := protected.Group("/onboarding")
		{
			onboarding.GET("/progress", onboarding.GetOnboardingProgress)
			onboarding.POST("/step/complete", onboarding.MarkStepComplete)
			onboarding.POST("/skip", onboarding.SkipOnboarding)
			onboarding.GET("/resources", onboarding.GetOnboardingResources)
		}
	}
	}

	// Public webhook route (no auth)
	router.POST("/api/v1/webhooks/stripe", billingHandlers.HandleStripeWebhook)
}
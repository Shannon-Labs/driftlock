package routes

import (
	"github.com/gin-gonic/gin"

	"driftlock/productized/api/config"
	"driftlock/productized/api/handlers"
	"driftlock/productized/api/middleware"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Public routes
	public := router.Group("/api/v1")
	{
		public.POST("/auth/register", handlers.Register)
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/refresh", handlers.RefreshToken)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
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
	}
}
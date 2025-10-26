package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"driftlock/productized/analytics"
	analyticsConfig "driftlock/productized/analytics/config"
	audit "driftlock/productized/analytics"
)

func AnalyticsMiddleware(analyticsService *analytics.AnalyticsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(start)

		// Get user info if available
		userID := "anonymous"
		user, exists := c.Get("user")
		if exists && user != nil {
			// Assuming user has an ID field, adjust based on your user model
			// userID = fmt.Sprintf("%v", user.(models.User).ID)
		}

		// Track API usage
		analyticsService.TrackAPIUsage(
			userID,
			c.Request.URL.Path,
			responseTime,
		)
	}
}

func AuditMiddleware(auditService *audit.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request if needed
		c.Next()
		
		// Example: log specific sensitive actions
		if c.Request.URL.Path == "/api/v1/anomalies" && c.Request.Method == "DELETE" {
			// Log anomaly deletion
			// auditService.LogAnomalyAction(userID, "delete", anomalyID, c.ClientIP(), c.GetHeader("User-Agent"), "Anomaly deleted by user")
		}
	}
}

func InitializeAnalytics() (*analytics.AnalyticsService, *audit.AuditService) {
	analyticsConfig := analyticsConfig.LoadAnalyticsConfig()
	
	analyticsService := analytics.NewAnalyticsService(analyticsConfig)
	auditService := audit.NewAuditService()
	
	return analyticsService, auditService
}
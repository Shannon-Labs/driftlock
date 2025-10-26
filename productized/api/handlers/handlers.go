package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"driftlock/productized/api/models"
	"driftlock/productized/api/services"
)

// HealthCheck returns the health status of the API
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   services.GetCurrentTime(),
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(input.Email, input.Name, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create a JWT token for the user
	token, err := services.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user login
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := services.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	// Implementation would verify refresh token and issue new access token
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetProfile returns the authenticated user's profile
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile updates the authenticated user's profile
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UpdateUser(userID.(uint), input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAnomalies returns a list of anomalies for the authenticated user's tenant
func GetAnomalies(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	
	anomalies, err := services.GetAnomalies(tenantID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get anomalies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"anomalies": anomalies})
}

// GetAnomaly returns a specific anomaly
func GetAnomaly(c *gin.Context) {
	anomalyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anomaly ID"})
		return
	}

	anomaly, err := services.GetAnomalyByID(uint(anomalyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anomaly not found"})
		return
	}

	// Check if the anomaly belongs to the user's tenant
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	if anomaly.TenantID != tenantID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"anomaly": anomaly})
}

// ResolveAnomaly marks an anomaly as resolved
func ResolveAnomaly(c *gin.Context) {
	anomalyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anomaly ID"})
		return
	}

	anomaly, err := services.GetAnomalyByID(uint(anomalyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anomaly not found"})
		return
	}

	// Check if the anomaly belongs to the user's tenant
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	if anomaly.TenantID != tenantID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = services.ResolveAnomaly(uint(anomalyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve anomaly"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Anomaly resolved successfully"})
}

// DeleteAnomaly deletes an anomaly
func DeleteAnomaly(c *gin.Context) {
	anomalyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid anomaly ID"})
		return
	}

	anomaly, err := services.GetAnomalyByID(uint(anomalyID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anomaly not found"})
		return
	}

	// Check if the anomaly belongs to the user's tenant
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	if anomaly.TenantID != tenantID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = services.DeleteAnomaly(uint(anomalyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete anomaly"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Anomaly deleted successfully"})
}

// IngestEvent handles incoming events for anomaly detection
func IngestEvent(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	var event map[string]interface{}
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process the event and potentially detect anomalies
	processed, err := services.ProcessEvent(event, tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"processed": processed})
}

// GetEvents returns a list of events
func GetEvents(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	events, err := services.GetEvents(tenantID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	stats, err := services.GetDashboardStats(tenantID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetRecentAnomalies returns recent anomalies for the dashboard
func GetRecentAnomalies(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tenant ID not found in context"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	anomalies, err := services.GetRecentAnomalies(tenantID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent anomalies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"anomalies": anomalies})
}
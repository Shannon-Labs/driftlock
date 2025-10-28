package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"driftlock/productized/api/database"
	"driftlock/productized/api/models"
	"driftlock/productized/api/utils"
)

// CreateUser creates a new user with a hashed password
func CreateUser(email, name, password string) (*models.User, error) {
	db := database.GetDB()

	// Check if user already exists
	var existingUser models.User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Name:     name,
		Password: string(hashedPassword),
		Role:     "user", // Default role
	}

	// Save to database
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	// Create a default tenant for the user
	tenant := &models.Tenant{
		Name:    fmt.Sprintf("%s's Organization", name),
		OwnerID: user.ID,
	}

	if err := db.Create(tenant).Error; err != nil {
		return nil, err
	}

	// Create default tenant settings
	settings := &models.TenantSettings{
		TenantID:              tenant.ID,
		LogAnomalyEnabled:     true,
		MetricAnomalyEnabled:  true,
		AnomalyThresholdMs:    1000,
	}

	if err := db.Create(settings).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser validates user credentials
func AuthenticateUser(email, password string) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func GetUserByID(id uint) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates a user's information
func UpdateUser(id uint, name, email string) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Update fields
	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}

	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAnomalies retrieves anomalies for a tenant
func GetAnomalies(tenantID uint, page, limit int) ([]models.Anomaly, error) {
	db := database.GetDB()

	var anomalies []models.Anomaly
	offset := (page - 1) * limit

	if err := db.Where("tenant_id = ?", tenantID).
		Order("detected_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&anomalies).Error; err != nil {
		return nil, err
	}

	return anomalies, nil
}

// GetAnomalyByID retrieves a specific anomaly by ID
func GetAnomalyByID(id uint) (*models.Anomaly, error) {
	db := database.GetDB()

	var anomaly models.Anomaly
	if err := db.First(&anomaly, id).Error; err != nil {
		return nil, err
	}

	return &anomaly, nil
}

// ResolveAnomaly marks an anomaly as resolved
func ResolveAnomaly(id uint) error {
	db := database.GetDB()

	now := time.Now()
	if err := db.Model(&models.Anomaly{}).Where("id = ?", id).Updates(map[string]interface{}{
		"resolved":    true,
		"resolved_at": &now,
	}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAnomaly deletes an anomaly
func DeleteAnomaly(id uint) error {
	db := database.GetDB()

	if err := db.Delete(&models.Anomaly{}, id).Error; err != nil {
		return err
	}

	return nil
}

// ProcessEvent processes an incoming event and potentially detects anomalies
func ProcessEvent(event map[string]interface{}, tenantID uint) (bool, error) {
	// In a real implementation, you would:
	// 1. Validate the event format
	// 2. Send to Kafka for processing
	// 3. Run anomaly detection algorithms
	// 4. Store results in database
	
	// For now, just return success
	return true, nil
}

// GetEvents retrieves events for a tenant
func GetEvents(tenantID uint, page, limit int) ([]map[string]interface{}, error) {
	// In a real implementation, you would query the database
	// For now, return mock data
	events := []map[string]interface{}{
		{
			"id":        1,
			"tenant_id": tenantID,
			"type":      "log",
			"source":    "service-1",
			"timestamp": time.Now().Add(-5 * time.Minute),
			"data":      `{"level": "INFO", "message": "Request processed successfully"}`,
		},
		{
			"id":        2,
			"tenant_id": tenantID,
			"type":      "metric",
			"source":    "api-gateway",
			"timestamp": time.Now().Add(-3 * time.Minute),
			"data":      `{"metric": "response_time", "value": 250}`,
		},
	}

	return events, nil
}

// GetDashboardStats returns dashboard statistics
func GetDashboardStats(tenantID uint) (map[string]interface{}, error) {
	db := database.GetDB()

	var totalAnomalies int64
	var unresolvedAnomalies int64
	var criticalAnomalies int64
	var last24hAnomalies int64

	db.Model(&models.Anomaly{}).Where("tenant_id = ?", tenantID).Count(&totalAnomalies)
	db.Model(&models.Anomaly{}).Where("tenant_id = ? AND resolved = ?", tenantID, false).Count(&unresolvedAnomalies)
	db.Model(&models.Anomaly{}).Where("tenant_id = ? AND severity = ?", tenantID, "critical").Count(&criticalAnomalies)
	db.Model(&models.Anomaly{}).Where("tenant_id = ? AND detected_at >= ?", tenantID, time.Now().Add(-24*time.Hour)).Count(&last24hAnomalies)

	stats := map[string]interface{}{
		"total_anomalies":      totalAnomalies,
		"unresolved_anomalies": unresolvedAnomalies,
		"critical_anomalies":   criticalAnomalies,
		"last_24h_anomalies":   last24hAnomalies,
		"anomaly_rate":         0.0, // Would need to calculate based on total events
	}

	return stats, nil
}

// GetRecentAnomalies returns recent anomalies
func GetRecentAnomalies(tenantID uint, limit int) ([]models.Anomaly, error) {
	db := database.GetDB()

	var anomalies []models.Anomaly
	if err := db.Where("tenant_id = ?", tenantID).
		Order("detected_at DESC").
		Limit(limit).
		Find(&anomalies).Error; err != nil {
		return nil, err
	}

	return anomalies, nil
}

// GenerateJWT creates a new JWT token
func GenerateJWT(userID uint, email, role string) (string, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}

// GetCurrentTime returns the current time in RFC3339 format
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}
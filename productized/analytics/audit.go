package audit

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"driftlock/productized/api/database"
)

type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	ResourceID  string    `json:"resource_id"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Details     string    `json:"details"`
	CreatedAt   time.Time `json:"created_at"`
}

type AuditService struct {
	db *gorm.DB
}

func NewAuditService() *AuditService {
	return &AuditService{
		db: database.DB,
	}
}

func (a *AuditService) LogAction(userID uint, action, resource, resourceID, ipAddress, userAgent, details string) error {
	auditLog := AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		Details:    details,
	}

	if err := a.db.Create(&auditLog).Error; err != nil {
		return fmt.Errorf("failed to create audit log: %v", err)
	}

	// Also log to standard logger
	log.Printf("AUDIT: User %d performed %s on %s:%s from %s - %s", 
		userID, action, resource, resourceID, ipAddress, details)

	return nil
}

func (a *AuditService) GetLogsByUser(userID uint, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := a.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs for user %d: %v", userID, err)
	}

	return logs, nil
}

func (a *AuditService) GetLogsByResource(resource, resourceID string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := a.db.Where("resource = ? AND resource_id = ?", resource, resourceID).Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs for resource %s:%s: %v", resource, resourceID, err)
	}

	return logs, nil
}

func (a *AuditService) GetLogsByAction(action string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := a.db.Where("action = ?", action).Order("created_at DESC").Limit(limit).Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs for action %s: %v", action, err)
	}

	return logs, nil
}

func (a *AuditService) GetUserActivitySummary(userID uint) (map[string]interface{}, error) {
	// Get total actions by this user
	var totalActions int64
	err := a.db.Model(&AuditLog{}).Where("user_id = ?", userID).Count(&totalActions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get total actions for user %d: %v", userID, err)
	}

	// Get recent actions
	recentLogs, err := a.GetLogsByUser(userID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logs for user %d: %v", userID, err)
	}

	summary := map[string]interface{}{
		"total_actions":   totalActions,
		"recent_actions":  recentLogs,
		"user_id":         userID,
		"last_activity":   time.Now(),
	}

	return summary, nil
}

// LogUserLogin records a user login event
func (a *AuditService) LogUserLogin(userID uint, ipAddress, userAgent string) error {
	return a.LogAction(
		userID,
		"login",
		"user",
		fmt.Sprintf("%d", userID),
		ipAddress,
		userAgent,
		"User successfully logged in",
	)
}

// LogUserLogout records a user logout event
func (a *AuditService) LogUserLogout(userID uint, ipAddress, userAgent string) error {
	return a.LogAction(
		userID,
		"logout",
		"user",
		fmt.Sprintf("%d", userID),
		ipAddress,
		userAgent,
		"User logged out",
	)
}

// LogAnomalyAction records anomaly-related actions
func (a *AuditService) LogAnomalyAction(userID uint, action, anomalyID, ipAddress, userAgent, details string) error {
	return a.LogAction(
		userID,
		action,
		"anomaly",
		anomalyID,
		ipAddress,
		userAgent,
		details,
	)
}

// LogBillingAction records billing-related actions
func (a *AuditService) LogBillingAction(userID uint, action, billingID, ipAddress, userAgent, details string) error {
	return a.LogAction(
		userID,
		action,
		"billing",
		billingID,
		ipAddress,
		userAgent,
		details,
	)
}
package onboarding

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"driftlock/productized/api/models"
)

type OnboardingStep struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	Requires    []string `json:"requires"` // Steps that must be completed before this one
}

// GetOnboardingProgress returns the current onboarding progress for a user
func GetOnboardingProgress(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	// In a real implementation, this would fetch from the database
	// For now, return a sample onboarding flow
	onboardingSteps := []OnboardingStep{
		{
			ID:          1,
			Name:        "profile_setup",
			Title:       "Complete Your Profile",
			Description: "Add your name, company, and contact information",
			Completed:   true, // Assume profile exists if user is logged in
			Requires:    []string{},
		},
		{
			ID:          2,
			Name:        "connect_data_source",
			Title:       "Connect Your Data Source",
			Description: "Connect your application or system to start monitoring",
			Completed:   false,
			Requires:    []string{"profile_setup"},
		},
		{
			ID:          3,
			Name:        "configure_alerts",
			Title:       "Configure Alert Settings",
			Description: "Set up how and when you want to be notified about anomalies",
			Completed:   false,
			Requires:    []string{"profile_setup"},
		},
		{
			ID:          4,
			Name:        "review_insights",
			Title:       "Review Initial Insights",
			Description: "Check your anomaly dashboard and initial insights",
			Completed:   false,
			Requires:    []string{"connect_data_source"},
		},
	}

	totalSteps := len(onboardingSteps)
	completedSteps := 0
	for _, step := range onboardingSteps {
		if step.Completed {
			completedSteps++
		}
	}

	progress := map[string]interface{}{
		"total_steps":     totalSteps,
		"completed_steps": completedSteps,
		"progress":        float64(completedSteps) / float64(totalSteps) * 100,
		"steps":           onboardingSteps,
		"user_id":         currentUser.ID,
	}

	c.JSON(http.StatusOK, progress)
}

// MarkStepComplete marks an onboarding step as completed for the user
func MarkStepComplete(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var req struct {
		StepName string `json:"step_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, this would update the database
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{
		"message":  "Step marked as complete",
		"step":     req.StepName,
		"user_id":  currentUser.ID,
		"status":   "completed",
	})
}

// SkipOnboarding allows a user to skip the onboarding process
func SkipOnboarding(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	// In a real implementation, this would update the database
	// to mark onboarding as skipped

	c.JSON(http.StatusOK, gin.H{
		"message": "Onboarding skipped",
		"user_id": currentUser.ID,
	})
}

// GetOnboardingResources returns helpful resources for new users
func GetOnboardingResources(c *gin.Context) {
	resources := []map[string]interface{}{
		{
			"title":       "Getting Started Guide",
			"description": "Learn how to set up DriftLock in 5 minutes",
			"url":         "https://docs.driftlock.com/getting-started",
			"type":        "guide",
			"duration":    "5 min",
		},
		{
			"title":       "Data Source Integration",
			"description": "Connect your application logs, metrics, and traces",
			"url":         "https://docs.driftlock.com/integrations",
			"type":        "integration",
			"duration":    "10 min",
		},
		{
			"title":       "Alert Configuration",
			"description": "Set up custom alert rules and notification preferences",
			"url":         "https://docs.driftlock.com/alerts",
			"type":        "configuration",
			"duration":    "8 min",
		},
		{
			"title":       "Dashboard Overview",
			"description": "Understand your anomaly detection dashboard",
			"url":         "https://docs.driftlock.com/dashboard",
			"type":        "overview",
			"duration":    "12 min",
		},
		{
			"title":       "Anomaly Resolution",
			"description": "Learn how to investigate and resolve detected anomalies",
			"url":         "https://docs.driftlock.com/anomaly-resolution",
			"type":        "guide",
			"duration":    "15 min",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
	})
}
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"driftlock/productized/api/models"
	"driftlock/productized/email/config"
	"driftlock/productized/email/services"
)

var emailService *services.EmailService

func InitEmailService(emailConfig *config.EmailConfig) {
	emailService = services.NewEmailService(emailConfig)
}

func SendTestEmail(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var req struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := emailService.SendEmail(req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
		"to":      req.To,
	})
}

func SendWelcomeEmail(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	err := emailService.SendWelcomeEmail(currentUser.Email, currentUser.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send welcome email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome email sent successfully",
		"to":      currentUser.Email,
	})
}

func SendAnomalyAlert(c *gin.Context) {
	user, _ := c.Get("user")
	currentUser := user.(*models.User)

	var req struct {
		AnomalyTitle       string `json:"anomaly_title" binding:"required"`
		AnomalyDescription string `json:"anomaly_description" binding:"required"`
		To                 string `json:"to" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := emailService.SendAnomalyAlertEmail(req.To, req.AnomalyTitle, req.AnomalyDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send anomaly alert email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Anomaly alert email sent successfully",
		"to":      req.To,
	})
}
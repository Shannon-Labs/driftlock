package services

import (
	"fmt"
	"net/smtp"
	"bytes"
	"text/template"

	"driftlock/productized/email/config"
)

type EmailService struct {
	config *config.EmailConfig
}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func NewEmailService(config *config.EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	if e.config.Provider == "smtp" {
		return e.sendViaSMTP(to, subject, body)
	} else if e.config.Provider == "sendgrid" {
		return e.sendViaSendGrid(to, subject, body)
	}
	
	return fmt.Errorf("email provider %s not supported", e.config.Provider)
}

func (e *EmailService) sendViaSMTP(to, subject, body string) error {
	from := e.config.FromEmail
	password := e.config.SMTPPassword
	
	// Set up authentication
	auth := smtp.PlainAuth("", e.config.SMTPUsername, password, e.config.SMTPHost)

	// Create the email message
	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s\r\n",
		from, to, subject, body)

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort),
		auth,
		from,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %v", err)
	}

	return nil
}

func (e *EmailService) sendViaSendGrid(to, subject, body string) error {
	// This would be implemented with SendGrid's Go library
	// For the purpose of this implementation, I'll provide the structure
	// but not the full implementation since we don't have the dependency added yet
	
	// In a real implementation, you would use:
	// import "github.com/sendgrid/sendgrid-go"
	// import "github.com/sendgrid/sendgrid-go/helpers/mail"
	
	// mail := mail.NewSingleEmail(
	//     mail.NewEmail("", e.config.FromEmail),
	//     subject,
	//     mail.NewEmail("", to),
	//     "",
	//     body,
	// )
	//
	// response, err := e.sgClient.Send(mail)
	// if err != nil {
	//     return fmt.Errorf("failed to send email via SendGrid: %v", err)
	// }
	
	// For now, we'll return an error to indicate it's not fully implemented
	return fmt.Errorf("SendGrid provider not fully implemented yet")
}

// SendTemplateEmail sends an email with a template
func (e *EmailService) SendTemplateEmail(to, subject, templateName string, data interface{}) error {
	// Load and parse the email template
	templatePath := fmt.Sprintf("email/templates/%s.html", templateName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %v", templateName, err)
	}

	// Execute the template with provided data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %v", templateName, err)
	}

	// Send the templated email
	return e.SendEmail(to, subject, buf.String())
}

// SendWelcomeEmail sends a welcome email to new users
func (e *EmailService) SendWelcomeEmail(to, userName string) error {
	subject := "Welcome to DriftLock!"
	data := map[string]interface{}{
		"UserName": userName,
		"SiteURL":  "https://driftlock.com",
		"SupportEmail": "support@driftlock.com",
	}

	return e.SendTemplateEmail(to, subject, "welcome", data)
}

// SendAnomalyAlertEmail sends an alert when an anomaly is detected
func (e *EmailService) SendAnomalyAlertEmail(to, anomalyTitle, anomalyDescription string) error {
	subject := fmt.Sprintf("Anomaly Alert: %s", anomalyTitle)
	data := map[string]interface{}{
		"AnomalyTitle":       anomalyTitle,
		"AnomalyDescription": anomalyDescription,
		"DashboardURL":       "https://app.driftlock.com/dashboard",
	}

	return e.SendTemplateEmail(to, subject, "anomaly-alert", data)
}

// SendBillingAlertEmail sends a billing-related email
func (e *EmailService) SendBillingAlertEmail(to, planName string, isExceeding bool) error {
	subject := "Billing Alert"
	templateName := "billing-threshold"
	if isExceeding {
		subject = "Billing Threshold Exceeded"
		templateName = "billing-exceeded"
	}
	
	data := map[string]interface{}{
		"PlanName":       planName,
		"IsExceeding":    isExceeding,
		"UsageCenterURL": "https://app.driftlock.com/billing",
	}

	return e.SendTemplateEmail(to, subject, templateName, data)
}
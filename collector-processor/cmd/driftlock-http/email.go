package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

// Email service for SendGrid integration

type emailService struct {
	apiKey   string
	fromAddr string
	fromName string
	enabled  bool
}

func newEmailService() *emailService {
	apiKey := env("SENDGRID_API_KEY", "")
	return &emailService{
		apiKey:   apiKey,
		fromAddr: env("EMAIL_FROM_ADDRESS", "noreply@driftlock.net"),
		fromName: env("EMAIL_FROM_NAME", "Driftlock"),
		enabled:  apiKey != "",
	}
}

// SendGrid API types
type sendgridEmail struct {
	Personalizations []sendgridPersonalization `json:"personalizations"`
	From             sendgridAddress           `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendgridContent         `json:"content"`
}

type sendgridPersonalization struct {
	To []sendgridAddress `json:"to"`
}

type sendgridAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type sendgridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Send email via SendGrid API
func (e *emailService) Send(ctx context.Context, to, subject, htmlContent string) error {
	if !e.enabled {
		log.Printf("[email] SendGrid not configured, skipping email to %s: %s", to, subject)
		return nil
	}

	payload := sendgridEmail{
		Personalizations: []sendgridPersonalization{
			{
				To: []sendgridAddress{{Email: to}},
			},
		},
		From:    sendgridAddress{Email: e.fromAddr, Name: e.fromName},
		Subject: subject,
		Content: []sendgridContent{
			{Type: "text/html", Value: htmlContent},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+e.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid returned status %d", resp.StatusCode)
	}

	log.Printf("[email] Sent email to %s: %s", to, subject)
	return nil
}

// SendWelcome sends welcome email with API key
func (e *emailService) SendWelcome(ctx context.Context, to, companyName, apiKey string) error {
	data := struct {
		CompanyName string
		APIKey      string
		DocsURL     string
	}{
		CompanyName: companyName,
		APIKey:      apiKey,
		DocsURL:     env("DOCS_URL", "https://driftlock.net"),
	}

	html, err := renderTemplate(welcomeEmailTemplate, data)
	if err != nil {
		return fmt.Errorf("render welcome email: %w", err)
	}

	return e.Send(ctx, to, "Welcome to Driftlock - Your API Key", html)
}

// SendVerification sends email verification with token
func (e *emailService) SendVerification(ctx context.Context, to, token string) error {
	verifyURL := fmt.Sprintf("%s/api/v1/onboard/verify?token=%s", env("APP_URL", "https://driftlock.net"), token)
	data := struct {
		VerifyURL string
	}{
		VerifyURL: verifyURL,
	}

	html, err := renderTemplate(verificationEmailTemplate, data)
	if err != nil {
		return fmt.Errorf("render verification email: %w", err)
	}

	return e.Send(ctx, to, "Verify your Driftlock account", html)
}

// SendTrialExpiring sends trial expiration warning
func (e *emailService) SendTrialExpiring(ctx context.Context, to, companyName string, daysRemaining int) error {
	data := struct {
		CompanyName   string
		DaysRemaining int
		UpgradeURL    string
	}{
		CompanyName:   companyName,
		DaysRemaining: daysRemaining,
		UpgradeURL:    fmt.Sprintf("%s/pricing", env("APP_URL", "https://driftlock.net")),
	}

	html, err := renderTemplate(trialExpiringTemplate, data)
	if err != nil {
		return fmt.Errorf("render trial expiring email: %w", err)
	}

	return e.Send(ctx, to, fmt.Sprintf("Your Driftlock trial expires in %d days", daysRemaining), html)
}

// SendAsync sends email asynchronously
func (e *emailService) SendAsync(to, subject, htmlContent string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := e.Send(ctx, to, subject, htmlContent); err != nil {
			log.Printf("[email] async send failed: %v", err)
		}
	}()
}

func renderTemplate(tmpl string, data any) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Email templates
const welcomeEmailTemplate = `<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #1e3a5f; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border: 1px solid #e5e7eb; }
        .api-key { background: #1f2937; color: #10b981; padding: 15px; border-radius: 8px; font-family: monospace; word-break: break-all; margin: 20px 0; }
        .cta { display: inline-block; background: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-top: 20px; }
        .footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Driftlock</h1>
        </div>
        <div class="content">
            <p>Hi {{.CompanyName}},</p>
            <p>Your Driftlock account has been created successfully. Here's your API key:</p>
            <div class="api-key">{{.APIKey}}</div>
            <p><strong>Important:</strong> Save this API key securely. It won't be shown again.</p>
            <p>Your free trial includes:</p>
            <ul>
                <li>10,000 events per month</li>
                <li>14 days of access</li>
                <li>Full API access</li>
            </ul>
            <a href="{{.DocsURL}}" class="cta">View Documentation</a>
        </div>
        <div class="footer">
            <p>&copy; 2024 Shannon Labs. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const verificationEmailTemplate = `<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #1e3a5f; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border: 1px solid #e5e7eb; }
        .cta { display: inline-block; background: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-top: 20px; }
        .footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Verify Your Email</h1>
        </div>
        <div class="content">
            <p>Please click the button below to verify your email address and activate your Driftlock account:</p>
            <a href="{{.VerifyURL}}" class="cta">Verify Email</a>
            <p style="margin-top: 20px; font-size: 12px; color: #6b7280;">
                If the button doesn't work, copy and paste this link into your browser:<br>
                <code>{{.VerifyURL}}</code>
            </p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Shannon Labs. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

const trialExpiringTemplate = `<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #f59e0b; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border: 1px solid #e5e7eb; }
        .cta { display: inline-block; background: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-top: 20px; }
        .footer { text-align: center; padding: 20px; color: #6b7280; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Trial Expiring Soon</h1>
        </div>
        <div class="content">
            <p>Hi {{.CompanyName}},</p>
            <p>Your Driftlock trial expires in <strong>{{.DaysRemaining}} days</strong>.</p>
            <p>To continue using Driftlock's explainable anomaly detection, please upgrade your plan:</p>
            <a href="{{.UpgradeURL}}" class="cta">View Pricing</a>
            <p style="margin-top: 20px;">Questions? Reply to this email and we'll help you find the right plan.</p>
        </div>
        <div class="footer">
            <p>&copy; 2024 Shannon Labs. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

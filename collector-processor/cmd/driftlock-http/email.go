package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type emailService struct {
	client      *sendgrid.Client
	fromAddress string
	fromName    string
}

func newEmailService() *emailService {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: SENDGRID_API_KEY not set, email service disabled")
		return nil
	}
	return &emailService{
		client:      sendgrid.NewSendClient(apiKey),
		fromAddress: env("EMAIL_FROM_ADDRESS", "noreply@driftlock.net"),
		fromName:    env("EMAIL_FROM_NAME", "Driftlock"),
	}
}

func (s *emailService) sendWelcomeEmail(toEmail, companyName, apiKey string) {
	if s == nil {
		log.Printf("MOCK EMAIL: Welcome email to %s (%s) [API key redacted]", toEmail, companyName)
		return
	}

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := "Welcome to Driftlock!"
	to := mail.NewEmail(companyName, toEmail)
	
	plainTextContent := fmt.Sprintf(`Welcome to Driftlock, %s!

Your API Key is: %s

You can use this key to start sending events to our API.
Documentation: https://driftlock.net/docs

Happy Detecting!
The Driftlock Team`, companyName, apiKey)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Welcome to Driftlock!</h2>
			<p>Hi %s,</p>
			<p>Thanks for signing up. Here is your API key to get started:</p>
			<div style="background: #f4f4f4; padding: 15px; border-radius: 5px; font-family: monospace; margin: 20px 0;">
				%s
			</div>
			<p>You can view our <a href="https://driftlock.net/docs">documentation</a> to learn how to integrate.</p>
			<p>Happy Detecting!<br>The Driftlock Team</p>
		</div>
	`, companyName, apiKey)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	
	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send welcome email to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid API returned %d for welcome email to %s: %s", response.StatusCode, toEmail, response.Body)
		} else {
			log.Printf("Sent welcome email to %s", toEmail)
		}
	}()
}

// sendVerificationEmailSync sends a verification email synchronously and returns any error.
// Use this for critical signup flows where we need to ensure delivery before returning success.
func (s *emailService) sendVerificationEmailSync(toEmail, companyName, token string) error {
	if s == nil {
		log.Printf("MOCK EMAIL: Verification email to %s (%s) [token redacted]", toEmail, companyName)
		return nil // Mock mode always succeeds
	}

	verifyLink := fmt.Sprintf("https://driftlock.net/verify?token=%s", token)

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := "Verify your Driftlock account"
	to := mail.NewEmail(companyName, toEmail)

	plainTextContent := fmt.Sprintf(`Welcome to Driftlock!

Please verify your email address by clicking the link below:
%s

If you didn't sign up for Driftlock, you can ignore this email.`, verifyLink)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Verify your email</h2>
			<p>Welcome to Driftlock! Please verify your email address to activate your account.</p>
			<p>
				<a href="%s" style="background: #2563eb; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Verify Email</a>
			</p>
			<p style="font-size: 12px; color: #666;">Or paste this link in your browser: %s</p>
		</div>
	`, verifyLink, verifyLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Synchronous send with retry for transient failures
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*500) * time.Millisecond) // Backoff: 0, 500ms, 1s
		}
		response, err := s.client.Send(message)
		if err != nil {
			lastErr = fmt.Errorf("send failed: %w", err)
			log.Printf("ERROR: Verification email attempt %d failed for %s: %v", attempt+1, toEmail, err)
			continue
		}
		if response.StatusCode >= 500 {
			// Server error - retry
			lastErr = fmt.Errorf("SendGrid returned %d: %s", response.StatusCode, response.Body)
			log.Printf("ERROR: Verification email attempt %d failed for %s: %v", attempt+1, toEmail, lastErr)
			continue
		}
		if response.StatusCode >= 400 {
			// Client error - don't retry
			return fmt.Errorf("SendGrid rejected email: %d - %s", response.StatusCode, response.Body)
		}
		log.Printf("Sent verification email to %s", toEmail)
		return nil
	}
	return lastErr
}

// sendVerificationEmail sends a verification email asynchronously (fire-and-forget).
// Use sendVerificationEmailSync for critical flows where delivery confirmation is needed.
func (s *emailService) sendVerificationEmail(toEmail, companyName, token string) {
	go func() {
		if err := s.sendVerificationEmailSync(toEmail, companyName, token); err != nil {
			log.Printf("ERROR: Async verification email failed for %s: %v", toEmail, err)
		}
	}()
}

func (s *emailService) sendTrialEndingEmail(toEmail, companyName string, daysRemaining int) {
	if s == nil {
		log.Printf("MOCK EMAIL: Trial ending for %s (%s) in %d days", toEmail, companyName, daysRemaining)
		return
	}

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := fmt.Sprintf("Your Driftlock trial ends in %d days", daysRemaining)
	to := mail.NewEmail(companyName, toEmail)

	plainTextContent := fmt.Sprintf(`Hi %s,

Your Driftlock trial ends in %d days. To continue using all features, please add a payment method.

Upgrade now: https://driftlock.net/dashboard

Thanks,
The Driftlock Team`, companyName, daysRemaining)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Your trial is ending soon</h2>
			<p>Hi %s,</p>
			<p>Your Driftlock trial ends in <strong>%d days</strong>. To continue enjoying all features, please add a payment method.</p>
			<p>
				<a href="https://driftlock.net/dashboard" style="background: #2563eb; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Upgrade Now</a>
			</p>
			<p>Thanks,<br>The Driftlock Team</p>
		</div>
	`, companyName, daysRemaining)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send trial ending email to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid returned %d for trial ending email to %s", response.StatusCode, toEmail)
		} else {
			log.Printf("Sent trial ending email to %s", toEmail)
		}
	}()
}

func (s *emailService) sendPaymentFailedEmail(toEmail, companyName string, gracePeriodEnd time.Time) {
	if s == nil {
		log.Printf("MOCK EMAIL: Payment failed for %s (%s), grace period until %s", toEmail, companyName, gracePeriodEnd)
		return
	}

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := "Action required: Payment failed for your Driftlock subscription"
	to := mail.NewEmail(companyName, toEmail)

	formattedDate := gracePeriodEnd.Format("January 2, 2006")

	plainTextContent := fmt.Sprintf(`Hi %s,

We were unable to process your payment. Please update your payment method to avoid service interruption.

Your service will remain active until %s.

Update payment: https://driftlock.net/dashboard

Thanks,
The Driftlock Team`, companyName, formattedDate)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2 style="color: #dc2626;">Payment Failed</h2>
			<p>Hi %s,</p>
			<p>We were unable to process your payment. Please update your payment method to avoid service interruption.</p>
			<p>Your service will remain active until <strong>%s</strong>.</p>
			<p>
				<a href="https://driftlock.net/dashboard" style="background: #dc2626; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Update Payment Method</a>
			</p>
			<p>Thanks,<br>The Driftlock Team</p>
		</div>
	`, companyName, formattedDate)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send payment failed email to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid returned %d for payment failed email to %s", response.StatusCode, toEmail)
		} else {
			log.Printf("Sent payment failed email to %s", toEmail)
		}
	}()
}

func (s *emailService) sendAdminAlert(toEmail, subject, body string) {
	if s == nil {
		log.Printf("MOCK ADMIN ALERT to %s: %s - %s", toEmail, subject, body)
		return
	}

	from := mail.NewEmail(s.fromName, s.fromAddress)
	to := mail.NewEmail("Admin", toEmail)

	plainTextContent := body
	htmlContent := fmt.Sprintf(`
		<div style="font-family: monospace; color: #333; background: #f4f4f4; padding: 20px;">
			<h2 style="color: #dc2626;">%s</h2>
			<pre style="white-space: pre-wrap; background: white; padding: 15px; border-radius: 5px;">%s</pre>
		</div>
	`, subject, body)

	message := mail.NewSingleEmail(from, "[Driftlock Alert] "+subject, to, plainTextContent, htmlContent)

	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send admin alert to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid returned %d for admin alert to %s", response.StatusCode, toEmail)
		} else {
			log.Printf("Sent admin alert to %s: %s", toEmail, subject)
		}
	}()
}

func (s *emailService) sendGraceExpiredEmail(toEmail, companyName string) {
	if s == nil {
		log.Printf("MOCK EMAIL: Grace period expired for %s (%s)", toEmail, companyName)
		return
	}

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := "Your Driftlock subscription has been downgraded"
	to := mail.NewEmail(companyName, toEmail)

	plainTextContent := fmt.Sprintf(`Hi %s,

Your grace period has expired and your subscription has been downgraded to our free Pulse tier.

Your data is still safe, but you'll have reduced feature access and lower usage limits.

To restore your subscription, update your payment method:
https://driftlock.net/dashboard

Thanks,
The Driftlock Team`, companyName)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Subscription Downgraded</h2>
			<p>Hi %s,</p>
			<p>Your grace period has expired and your subscription has been downgraded to our free <strong>Pulse</strong> tier.</p>
			<p>Your data is still safe, but you'll have reduced feature access and lower usage limits.</p>
			<p>
				<a href="https://driftlock.net/dashboard" style="background: #2563eb; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Restore Subscription</a>
			</p>
			<p>Thanks,<br>The Driftlock Team</p>
		</div>
	`, companyName)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send grace expired email to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid returned %d for grace expired email to %s", response.StatusCode, toEmail)
		} else {
			log.Printf("Sent grace expired email to %s", toEmail)
		}
	}()
}

